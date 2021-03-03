package sMessage

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

type NCManager struct {
	Manager     *jsm.Manager
	NC          *nats.Conn
	Application string
	Environment string
	sURL        string
	connOpts    []nats.Option

	sync.Mutex
}

/*
	NewNC will create a Sote Jetstream Manager for NATS.  The application and environment are required.
	credentialFileName is optional and should not be used except in development.
*/
func NewNC(application, environment, credentialFileName, sURL string, maxReconnect int, reconnectWait time.Duration) (pNCManager *NCManager, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err      error
		tmpCreds interface{}
	)

	if soteErr = sConfigParams.ValidateApplication(application); soteErr.ErrCode == nil {
		pNCManager = &NCManager{Application: application}
	}

	if soteErr.ErrCode == nil {
		if soteErr = sConfigParams.ValidateEnvironment(environment); soteErr.ErrCode == nil {
			pNCManager.Environment = environment
		}
	}

	if soteErr.ErrCode == nil {
		if len(sURL) > 0 {
			soteErr = pNCManager.setURL(sURL)
		} else {
			var getURL string
			getURL, soteErr = sConfigParams.GetNATSURL(pNCManager.Application, pNCManager.Environment)
			soteErr = pNCManager.setURL(getURL)
		}

		soteErr = pNCManager.setURL(sURL)
	}

	// Setting connection options
	if soteErr.ErrCode == nil {
		soteErr = pNCManager.setReconnectOptions(maxReconnect, reconnectWait)
	}

	if soteErr.ErrCode == nil {
		if len(credentialFileName) > 0 {
			soteErr = pNCManager.setCredentialsFile(credentialFileName)
		} else {
			getCreds := sConfigParams.GetNATSCredentials()
			tmpCreds, soteErr = getCreds(pNCManager.Application, pNCManager.Environment)
			pNCManager.connOpts = append(pNCManager.connOpts, pNCManager.userCredsFromRaw([]byte(tmpCreds.(string))))
		}
		// Making connection to server
		if soteErr.ErrCode == nil {
			if pNCManager.NC, soteErr = pNCManager.connect(); soteErr.ErrCode == nil {
				// Creating the JSM Manager
				pNCManager.Manager, err = jsm.New(pNCManager.NC)
				if err != nil {
					log.Panic(err.Error())
				}
			}
		}
	}

	return
}

/*
	setCredentialsFile will pull the credentials from the file system.
	** THIS IS NOT THE RECOMMENDED APPROACH! YOU SHOULD USE setCredentialFromSystemParameters **
*/
func (ncm *NCManager) setCredentialsFile(streamCredentialFile string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(streamCredentialFile) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"streamCredentialFile"}), sError.EmptyMap)
	} else if _, err := os.Stat(streamCredentialFile); err != nil {
		soteErr = sError.GetSError(209010, sError.BuildParams([]string{streamCredentialFile, err.Error()}), sError.EmptyMap)
	} else {
		ncm.connOpts = append(ncm.connOpts, nats.UserCredentials(streamCredentialFile))
	}

	return
}

/*
	UserCredsFromRaw will take a credential file content that is not stored on the file system, such as AWS System Manager Parameters
*/
func (ncm *NCManager) userCredsFromRaw(rawCredentials []byte) nats.Option {
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
func (ncm *NCManager) setReconnectOptions(maxReconnect int, reconnectWait time.Duration) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if reconnectWait < 250*time.Millisecond || reconnectWait > 1*time.Minute {
		ncm.connOpts = append(ncm.connOpts, nats.ReconnectWait(1*time.Second))
		sLogger.Info("Resetting reconnectWait to 1 second")
	} else {
		ncm.connOpts = append(ncm.connOpts, nats.ReconnectWait(reconnectWait))
	}
	if maxReconnect < 1 || maxReconnect > 5 {
		ncm.connOpts = append(ncm.connOpts, nats.MaxReconnects(1))
		sLogger.Info("Resetting maxReconnect to 1 try")
	} else {
		ncm.connOpts = append(ncm.connOpts, nats.MaxReconnects(maxReconnect))
	}

	return
}

func (ncm *NCManager) setURL(sURL string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if _, err := url.Parse(sURL); err != nil || sURL == "" {
		soteErr = sError.GetSError(210090, sError.BuildParams([]string{sURL}), nil)
	} else {
		ncm.sURL = sURL
	}

	return
}

/*
	This will connect to the NATS network using the values set in the NCManager
*/
func (ncm *NCManager) connect() (nc *nats.Conn, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	// Connect to NATS  --  default "euwest1.aws.ngs.global"
	nc, err = nats.Connect(ncm.sURL, ncm.connOpts...)
	if err != nil {
		if strings.Contains(err.Error(), "no servers") {
			soteErr = sError.GetSError(209499, nil, sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		} else if strings.Contains(err.Error(), "no nkey seed found") {
			soteErr = sError.GetSError(209398, nil, sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		} else {
			var errDetails = make(map[string]string)
			errDetails, soteErr = sError.ConvertErr(err)
			if soteErr.ErrCode != nil {
				sLogger.Info(soteErr.FmtErrMsg)
				panic("sError.ConvertErr Failed")
			}
			sLogger.Info(sError.GetSError(210400, nil, errDetails).FmtErrMsg)
			panic("sMessages.connect Failed")
		}
	}

	return
}

/*
	This will close the connection to the NATS network using the NCManager
*/
func (ncm *NCManager) Close() {
	sLogger.DebugMethod()

	// Closing connect to NATS  --  default "euwest1.aws.ngs.global"
	ncm.NC.Close()

	ncm = nil

	return
}

func (ncm *NCManager) natsErrorHandle(err error, subject string, data []byte) (soteErr sError.SoteError) {
	if err == nil {
		ncm.NC.Flush()

		if err = ncm.NC.LastError(); err != nil {
			soteErr = sError.GetSError(210400, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
			panic("An unexpected error has occurred.")
		} else {
			sLogger.Info(fmt.Sprintf("Published [%v] : '%v'\n", subject, string(data)))
		}
	} else {

		if strings.Contains(err.Error(), "invalid subject") {
			soteErr = sError.GetSError(208310, nil, sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
	SPublish will publish the data argument to the given subject. The data argument is left untouched and needs
	to be correctly interpreted on the receiver
*/
func (ncm *NCManager) SPublish(subject string, data []byte) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var err error

	err = ncm.NC.Publish(subject, data)

	soteErr = ncm.natsErrorHandle(err, subject, data)

	return
}

/*
	SPublishMsg publishes the Msg structure, which includes the Subject, an optional Reply and an optional Data field.
*/
func (ncm *NCManager) SPublishMsg(m *nats.Msg) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var err error

	err = ncm.NC.PublishMsg(m)

	soteErr = ncm.natsErrorHandle(err, m.Subject, m.Data)

	return
}

/*
	SPublishRequest will perform a Publish() expecting a response on the reply subject.
*/
func (ncm *NCManager) SPublishRequest(subject string, reply string, data []byte) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var err error

	if err = ncm.NC.PublishRequest(subject, reply, data); err == nil {
		ncm.NC.Flush()
		if err = ncm.NC.LastError(); err != nil {
			soteErr = sError.GetSError(210400, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
			panic("An unexpected error has occurred.")
		} else {
			sLogger.Info(fmt.Sprintf("Published [%s] : '%s'\n", subject, string(data)))
		}
	} else {
		soteErr = sError.GetSError(210400, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
		sLogger.Info(soteErr.FmtErrMsg)
	}

	return
}

/*
	SSubscribe will express interest in the given subject. The subject can have wildcards (partial:*, full:>).
	Messages will be delivered to the associated MsgHandler. Returns an error and the subscription.
*/
func (ncm *NCManager) SSubscribe(subject string) (s *nats.Subscription, soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var err error

	if s, err = ncm.NC.Subscribe(subject, func(msg *nats.Msg) {
		sLogger.Info(fmt.Sprintf("Received message [%s]:  %s\n", msg.Subject, string(msg.Data)))
	}); err != nil {
		if strings.Contains(err.Error(), "invalid subject") {
			soteErr = sError.GetSError(208310, nil, sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
	SNextMsg will return the next message available to a synchronous subscriber or block until one is available.
	An error is returned if the subscription is invalid (ErrBadSubscription), the connection is closed (ErrConnectionClosed),
	or the timeout is reached (ErrTimeout).
*/
func (ncm *NCManager) SNextMsg(s *nats.Subscription, t time.Duration) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var (
		m   *nats.Msg
		err error
	)

	if m, err = s.NextMsg(t); err == nil {
		sLogger.Info(fmt.Sprintf("Received message [%s]: %s\n", m.Subject, string(m.Data)))
	} else {
		if strings.Contains(err.Error(), "nats: connection closed") {
			soteErr = sError.GetSError(209499, nil, sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
	SRequest will send a request payload and deliver the response message, or an error,
	including a timeout if no message was received properly
*/
func (ncm *NCManager) SRequest(subject string, data []byte, time time.Duration) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		msg *nats.Msg
		err error
	)

	if msg, err = ncm.NC.Request(subject, data, time); err == nil {
		sLogger.Info(fmt.Sprintf("Published [%s] : '%s'", subject, data))
		sLogger.Info(fmt.Sprintf("Received Reply [%v] : '%s'", msg.Subject, string(msg.Data)))
	} else {
		if ncm.NC.LastError() != nil {
			soteErr = sError.GetSError(210400, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
			panic("An unexpected error has occurred.")
		}
		sLogger.Info(fmt.Sprintf("%v", err))
	}

	return
}

/*
	SRequestReply listens  to subject argument and sends data argument as reply to a request. Returns the subscription
	and error.
*/
func (ncm *NCManager) SRequestReply(subject string, data []byte) (s *nats.Subscription, soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var err error

	s, err = ncm.NC.Subscribe(subject, func(msg *nats.Msg) {
		sLogger.Info(fmt.Sprintf("Received message [%s]:  %s\n", msg.Subject, string(msg.Data)))
		if err = msg.Respond(data); err != nil {
			sLogger.Info(fmt.Sprintf("%v", err))
			panic("Response to request failed")
		}
	})

	ncm.NC.Flush()

	if err = ncm.NC.LastError(); err != nil {
		soteErr = sError.GetSError(210400, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
		sLogger.Info(soteErr.FmtErrMsg)
		panic("An unexpected error has occurred.")
	}

	sLogger.Info(fmt.Sprintf("Listening on [%s]", s.Subject))

	return
}

/*
	SRequestMsg will send a request payload including optional headers and deliver the response message, or an error,
	including a timeout if no message was received properly.
*/
func (ncm *NCManager) SRequestMsg(m *nats.Msg, time time.Duration) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var (
		msg *nats.Msg
		err error
	)

	if msg, err = ncm.NC.RequestMsg(m, time); err == nil {
		sLogger.Info(fmt.Sprintf("Published [%s] : '%s'", m.Subject, string(m.Data)))
		sLogger.Info(fmt.Sprintf("Received Reply [%v] : '%s'", msg.Subject, string(msg.Data)))
	} else {
		if ncm.NC.LastError() != nil {
			soteErr = sError.GetSError(210400, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
			panic("An unexpected error has occurred.")
		}
		panic(fmt.Sprintf("%v", err))
	}

	return
}
