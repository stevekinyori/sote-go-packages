/*
	This is a wrapper for Sote Golang developers to access services from NATS. This does not support JetStream.
*/
package sMessage

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"gitlab.com/soteapps/packages/v2022/sConfigParams"
	"gitlab.com/soteapps/packages/v2022/sError"
	"gitlab.com/soteapps/packages/v2022/sLogger"
)

type MessageManager struct {
	NatsConnectionPtr *nats.Conn
	application       string
	environment       string
	connectionURL     string
	connectionOptions []nats.Option
	Subscriptions     map[string]*nats.Subscription
	SyncSubscriptions map[string]*nats.Subscription
	PullSubscriptions map[string]*nats.Subscription
	Messages          []*nats.Msg
	RawMessage        *nats.RawStreamMsg
}

type SStreamInfo *nats.StreamInfo
type SConsumerInfo *nats.ConsumerInfo
type SMsg *nats.Msg

/*
	New will create a Sote Message Manager and a connection to the NATS network.
*/
func New(ctx context.Context, application, environment, credentialFileName, connectionURL, connectionName string, secure bool, maxReconnect int,
	reconnectWait time.Duration, testMode bool) (MessageManagerPtr *MessageManager, soteErr sError.SoteError) {
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
			connectionURL, soteErr = sConfigParams.GetNATSURL(ctx, MessageManagerPtr.application, MessageManagerPtr.environment)
		}
		soteErr = MessageManagerPtr.setURL(ctx, connectionURL, secure)
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
			getCreds := sConfigParams.GetNATSCredentials(ctx)
			tmpCreds, soteErr = getCreds(MessageManagerPtr.application, MessageManagerPtr.environment)
			MessageManagerPtr.connectionOptions = append(MessageManagerPtr.connectionOptions,
				MessageManagerPtr.setCredentialFromSystemParameters([]byte(tmpCreds.(string))))
		}
		// Making connection to server
		if soteErr.ErrCode == nil {
			soteErr = MessageManagerPtr.connect(testMode)
			MessageManagerPtr.Subscriptions = make(map[string]*nats.Subscription)
			MessageManagerPtr.SyncSubscriptions = make(map[string]*nats.Subscription)
			MessageManagerPtr.PullSubscriptions = make(map[string]*nats.Subscription)
		}
	}

	return
}

/*
	Close will terminate the connection to the NATS network
*/
func (mmPtr *MessageManager) Close() {
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

func (mmPtr *MessageManager) setURL(ctx context.Context, connectionURL string, secure bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		mask string
	)

	if _, err := url.Parse(connectionURL); err != nil || connectionURL == "" {
		soteErr = sError.GetSError(210090, sError.BuildParams([]string{connectionURL}), nil)
	} else {
		if secure {
			mask, soteErr = sConfigParams.GetNATSTLSURLMask(ctx, mmPtr.application)
			connectionURL = mask + connectionURL
		}
		mmPtr.connectionURL = connectionURL
	}

	return
}

/*
	This will connect to the NATS network using the values set in the MessageManager
*/
func (mmPtr *MessageManager) connect(testMode bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	params := make(map[string]string)
	params["Connection URL"] = mmPtr.connectionURL
	params["testMode"] = strconv.FormatBool(testMode)

	// Connect to NATS
	mmPtr.NatsConnectionPtr, err = nats.Connect(mmPtr.connectionURL, mmPtr.connectionOptions...)
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

func (mmPtr *MessageManager) natsErrorHandle(err error, params map[string]string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		errorDetail = make(map[string]string)
	)
	switch err.Error() {
	case nats.ErrInvalidConnection.Error():
		soteErr = sError.GetSError(210499, nil, sError.EmptyMap)
	case nats.ErrBadSubject.Error():
		soteErr = sError.GetSError(208310, sError.BuildParams([]string{params["Subject"]}), sError.EmptyMap)
	case nats.ErrNoServers.Error():
		soteErr = sError.GetSError(209499, nil, sError.EmptyMap)
	case "no nkey seed found":
		soteErr = sError.GetSError(209398, nil, sError.EmptyMap)
	case nats.ErrNoMatchingStream.Error():
		soteErr = sError.GetSError(210599, nil, sError.EmptyMap)
	case nats.ErrTimeout.Error():
		errorDetail["raw_message"] = "nats: timeout"
		soteErr = sError.GetSError(101010, sError.BuildParams([]string{"nats"}), errorDetail)
	case "context deadline exceeded":
		errorDetail["raw_message"] = "context deadline exceeded"
		soteErr = sError.GetSError(101010, sError.BuildParams([]string{"nats"}), errorDetail)
	case nats.ErrConnectionClosed.Error():
		errorDetail["raw_message"] = "nats: connection closed"
		soteErr = sError.GetSError(209499, nil, errorDetail)
	case nats.ErrBadSubscription.Error():
		errorDetail["raw_message"] = "nats: invalid subscription"
		soteErr = sError.GetSError(206050, sError.BuildParams([]string{params["Subscription Name"], params["Subject"]}), errorDetail)
	case nats.ErrStreamNotFound.Error():
		errorDetail["raw_message"] = "stream not found"
		soteErr = sError.GetSError(109999, sError.BuildParams([]string{params["Stream Name"]}), errorDetail)
	// 	TODO This should be removed once the NATS bug is resolved.
	case "too many open files":
		soteErr = sError.GetSError(109999, sError.BuildParams([]string{params["Stream Name"]}), sError.EmptyMap)
	case "no message found":
		soteErr = sError.GetSError(109999, sError.BuildParams([]string{params["Stream Name"], params["Message Sequence"]}), sError.EmptyMap)
	default:
		soteErr = sError.GetSError(199999, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
	}
	sLogger.Info(fmt.Sprintf("ERROR IN: messagemanager.go err: %v | %v", err.Error(), dumpParams(params)))
	sLogger.Info(soteErr.FmtErrMsg)

	return
}

func dumpParams(params map[string]string) (paramString string) {
	sLogger.DebugMethod()

	for key, value := range params {
		paramString += fmt.Sprintf("%v: %v | ", key, value)
	}

	return
}

// NewMessage creates a message for publishing that will use headers.
func NewMessage(subject string) (msg *nats.Msg) {
	sLogger.DebugMethod()

	msg = nats.NewMsg(subject)

	return
}
