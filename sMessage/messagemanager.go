/*
	This is a wrapper for Sote Golang developers to access services from NATS. This does not support JetStream.
*/
package sMessage

import (
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

type MessageManager struct {
	NatsConnectionPtr *nats.Conn
	application       string
	environment       string
	connectionURL     string
	connectionOptions []nats.Option
	SyncSubscriptions map[string]*nats.Subscription
	Subscriptions     map[string]*nats.Subscription

	sync.Mutex // TODO Do we need this?
}

/*
	New will create a Sote Message Manager and a connection to the NATS network.
*/
func New(application, environment, credentialFileName, connectionURL, connectionName string, secure bool, maxReconnect int,
	reconnectWait time.Duration) (MessageManagerPtr *MessageManager, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		tmpCreds interface{}
	)

	// Initialize the values for Nats Manager
	MessageManagerPtr = &MessageManager{
		NatsConnectionPtr: nil,
		application:       "synadia",
		environment:       "staging",
		connectionURL:     "localhost",
		connectionOptions: []nats.Option{nats.Name(connectionName)},
		Mutex:             sync.Mutex{},
	}

	if soteErr = sConfigParams.ValidateApplication(application); soteErr.ErrCode == nil {
		MessageManagerPtr.application = application
	}

	if soteErr.ErrCode == nil {
		if soteErr = sConfigParams.ValidateEnvironment(environment); soteErr.ErrCode == nil {
			MessageManagerPtr.environment = environment
		}
	}

	if soteErr.ErrCode == nil {
		if len(connectionURL) == 0 {
			connectionURL, soteErr = sConfigParams.GetNATSURL(MessageManagerPtr.application, MessageManagerPtr.environment)
		}
		soteErr = MessageManagerPtr.setURL(connectionURL, secure)
	}

	// Setting connection options
	if soteErr.ErrCode == nil {
		soteErr = MessageManagerPtr.setReconnectOptions(maxReconnect, reconnectWait)
	}

	// Setting credentials
	if soteErr.ErrCode == nil {
		if len(credentialFileName) > 0 {
			soteErr = MessageManagerPtr.setCredentialsFile(credentialFileName)
		} else {
			// This will retrieve value from AWS System Manager Parameter Store
			getCreds := sConfigParams.GetNATSCredentials()
			tmpCreds, soteErr = getCreds(MessageManagerPtr.application, MessageManagerPtr.environment)
			MessageManagerPtr.connectionOptions = append(MessageManagerPtr.connectionOptions,
				MessageManagerPtr.setCredentialFromSystemParameters([]byte(tmpCreds.(string))))
		}
		// Making connection to server
		if soteErr.ErrCode == nil {
			soteErr = MessageManagerPtr.connect()
			MessageManagerPtr.SyncSubscriptions = make(map[string]*nats.Subscription)
		}
	}

	return
}

/*
	Close will terminate the connection to the NATS network
*/
func (mmPtr *MessageManager) Close() (nmOutPtr *MessageManager) {
	sLogger.DebugMethod()

	// Closing connect to NATS
	mmPtr.NatsConnectionPtr.Close()

	return
}

/*
	setCredentialsFile will pull the credentials from the file system.
	** THIS IS NOT THE RECOMMENDED APPROACH! YOU SHOULD USE setCredentialFromSystemParameters **
*/
func (mmPtr *MessageManager) setCredentialsFile(streamCredentialFile string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(streamCredentialFile) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"streamCredentialFile"}), sError.EmptyMap)
	} else if _, err := os.Stat(streamCredentialFile); err != nil {
		soteErr = sError.GetSError(209010, sError.BuildParams([]string{streamCredentialFile, err.Error()}), sError.EmptyMap)
	} else {
		mmPtr.connectionOptions = append(mmPtr.connectionOptions, nats.UserCredentials(streamCredentialFile))
	}

	return
}

/*
	setCredentialFromSystemParameters will take a credential file content that is not stored on the file system, such as AWS System Manager Parameters
*/
func (mmPtr *MessageManager) setCredentialFromSystemParameters(rawCredentials []byte) nats.Option {
	return nats.UserJWT(
		func() (string, error) { return nkeys.ParseDecoratedJWT(rawCredentials) },
		func(nonce []byte) ([]byte, error) {
			kp, err := nkeys.ParseDecoratedNKey(rawCredentials)
			if err != nil {
				return nil, err
			}
			defer kp.Wipe()
			sig, _ := kp.Sign(nonce)
			return sig, nil
		})
}

/*
	setReconnectOptions expects a maxReconnect value between 1 and 5; if not, it is set to 1. The reconnectWait value
	between 250 milliseconds and 1 minute; if not, it is set to 1 second.
*/
func (mmPtr *MessageManager) setReconnectOptions(maxReconnect int, reconnectWait time.Duration) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if reconnectWait < 250*time.Millisecond || reconnectWait > 1*time.Minute {
		mmPtr.connectionOptions = append(mmPtr.connectionOptions, nats.ReconnectWait(1*time.Second))
		sLogger.Info("Resetting reconnectWait to 1 second")
	} else {
		mmPtr.connectionOptions = append(mmPtr.connectionOptions, nats.ReconnectWait(reconnectWait))
	}
	if maxReconnect < 1 || maxReconnect > 5 {
		mmPtr.connectionOptions = append(mmPtr.connectionOptions, nats.MaxReconnects(1))
		sLogger.Info("Resetting maxReconnect to 1 try")
	} else {
		mmPtr.connectionOptions = append(mmPtr.connectionOptions, nats.MaxReconnects(maxReconnect))
	}

	return
}

func (mmPtr *MessageManager) setURL(connectionURL string, secure bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		mask string
	)

	if _, err := url.Parse(connectionURL); err != nil || connectionURL == "" {
		soteErr = sError.GetSError(210090, sError.BuildParams([]string{connectionURL}), nil)
	} else {
		if secure {
			mask, soteErr = sConfigParams.GetNATSTLSURLMask(mmPtr.application)
			connectionURL = mask + connectionURL
		}
		mmPtr.connectionURL = connectionURL
	}

	return
}

/*
	This will connect to the NATS network using the values set in the MessageManager
*/
func (mmPtr *MessageManager) connect() (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	// Connect to NATS
	mmPtr.NatsConnectionPtr, err = nats.Connect(mmPtr.connectionURL, mmPtr.connectionOptions...)
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, "", "", "", "")
	}

	return
}

func (mmPtr *MessageManager) natsErrorHandle(err error, subject, reply, subscriptionName string, data string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		panicError = true
	)

	switch err.Error() {
	case "nats: invalid connection":
		soteErr = sError.GetSError(210499, nil, sError.EmptyMap)
	case "nats: invalid subject":
		soteErr = sError.GetSError(208310, sError.BuildParams([]string{subject}), sError.EmptyMap)
	case "nats: no servers available for connection":
		soteErr = sError.GetSError(209499, nil, sError.EmptyMap)
	case "no nkey seed found":
		soteErr = sError.GetSError(209398, nil, sError.EmptyMap)
	case "nats: timeout":
		soteErr = sError.GetSError(101010, nil, sError.EmptyMap)
		panicError = false
	case "nats: connection closed":
		soteErr = sError.GetSError(209499, nil, sError.EmptyMap)
		panicError = false
	default:
		soteErr = sError.GetSError(199999, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
	}
	sLogger.Info(fmt.Sprintf("ERROR IN: messagemanager.go err: %v subject: %v reply: %v subscription name: %v data: %v", err.Error(), subject,
		reply, subscriptionName, data))
	sLogger.Info(soteErr.FmtErrMsg)

	if panicError {
		panic(soteErr.FmtErrMsg)
	}

	return
}
