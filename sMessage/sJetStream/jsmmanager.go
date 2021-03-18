/*
NATS.io Jetstream default values are listed below. (Gathered from https://github.com/nats-io/jetstream)
These are values that can be set natively.  sstream and consumer place limitation on the configurations available from NATS.

	NATS STREAM SETTINGS:
		Ack: Required for published
			value is set using: true, false
		Discard: Default value: old
			value is set using: jsm.DiscardNew(), jsm.DiscardOld()
		Duplicates: Default value: ""
			value is set using: (s)econds, (m)inutes, (h)ours, (y)ears, (M)onths, (d)ays
			example: 1s, 1h, 1M, 1seconds, 1Months
		MaxAge: Default value: -1 (unlimited)
			value is set using: (s)econds, (m)inutes, (h)ours, (y)ears, (M)onths, (d)ays
			example: 1s, 1h, 1M, 1seconds, 1Months
		MaxBytes: Default value: -1 (unlimited)
			value is set using: (B)ytes, (k)ilobytes, (m)egabytes
			example: 512B, 1k, 1m
		MaxConsumers: Default value: -1 (unlimited)
			value is set using: >0
		MaxMsgs: Required
			value is set using: -1 (unlimited) or >0
		MaxMsgSize: Default value: -1 (unlimited)
			value is set using: (B)ytes, (k)ilobytes, (m)egabytes
			example: 512B, 1k, 1m
		Name: Required
			The name of the stream and the stream subjects can be different
		NoAck: Default is false
			values is set using 'true' or 'false'
		Replicas: No default
			value is from 1 to n
		Retention: No default
			value is set using jsm.LimitsRetention(), jsm.WorkQueueRetention() jsm.InterestPolicy()
		Storage: No default
			value is set using jsm..MemoryStorage() or jsm.FileStorage()
		Subjects: No default
			value is a string of values that are comma separated.  Dot can be used to create sub-subjects and '*' is a wildcard
	NATS CONSUMER SETTINGS:
		AckPolicy: Default value: none
			value is set using: none, all, explicit (explicit required for pull consumers)
		AckWait: Default value: -1s (forever)
			value is set using: (s)econds, (m)inutes, (h)ours, (y)ears, (M)onths, (d)ays
		DeliverPolicy: Default value: ""
			value is set using: all, last, new or next, DeliverByStartSequence or DeliverByStartTime
		DeliverySubject: Default value: "" ("" required for pull consumers)
			value is set using: <target subject>
			example: TEST_CONSUMER_NAME, test_consumer_name, Test_Consumer_Name
		Durable Name: Default value: ""
			value is set using: <durable name>
			example: TEST_CONSUMER_NAME, test_consumer_name, Test_Consumer_Name
		FilterSubject: Default value: "" (all)
			value is set using: <stream name>.<subject name>
			example: TEST_STREAM_NAME.* for all messages from the TEST_STREAM_NAME stream, TEST_STREAM_NAME.cat for only cat messages
		MaxDeliver: Default value: -1 (unlimited)
			value is set using: >0
		OptStartSeq: Default value: (Required when using DeliveryPolicy with DeliverByStartSequence)
			value is set using: >0
		ReplayPolicy: Default value: instant
			value is set using: instant, original
		SampleFrequency: Default value: -1
			value is set using: 0 to 100
		OptStartTime: Default value: Now
			value is set using: (s)econds, (m)inutes, (h)ours, (y)ears, (M)onths, (d)ays
			example: 1s, 1h, 1M, 1seconds, 1Months
		RateLimit: Default value: 0
			value is set using: >0
*/
package sJetStream

import (
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

type JSMManager struct {
	Manager     *jsm.Manager
	nc          *nats.Conn
	Application string
	Environment string
	sURL        string
	connOpts    []nats.Option

	sync.Mutex
}

/*
	New will create a Sote Jetstream Manager for NATS.  The application and environment are required.
	credentialFileName is optional and should not be used except in development.
*/
func New(application, environment, credentialFileName, sURL string, maxReconnect int, reconnectWait time.Duration) (pJSMManager *JSMManager, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err      error
		tmpCreds interface{}
	)

	if soteErr = sConfigParams.ValidateApplication(application); soteErr.ErrCode == nil {
		pJSMManager = &JSMManager{Application: application}
	}

	if soteErr.ErrCode == nil {
		if soteErr = sConfigParams.ValidateEnvironment(environment); soteErr.ErrCode == nil {
			pJSMManager.Environment = environment
		}
	}

	if soteErr.ErrCode == nil {
		if len(sURL) > 0 {
			soteErr = pJSMManager.setURL(sURL)
		} else {
			var getURL string
			getURL, soteErr = sConfigParams.GetNATSURL(pJSMManager.Application, pJSMManager.Environment)
			soteErr = pJSMManager.setURL(getURL)
		}

		soteErr = pJSMManager.setURL(sURL)
	}

	// Setting connection options
	if soteErr.ErrCode == nil {
		soteErr = pJSMManager.setReconnectOptions(maxReconnect, reconnectWait)
	}

	if soteErr.ErrCode == nil {
		if len(credentialFileName) > 0 {
			soteErr = pJSMManager.setCredentialsFile(credentialFileName)
		} else {
			getCreds := sConfigParams.GetNATSCredentials()
			tmpCreds, soteErr = getCreds(pJSMManager.Application, pJSMManager.Environment)
			pJSMManager.connOpts = append(pJSMManager.connOpts, pJSMManager.userCredsFromRaw([]byte(tmpCreds.(string))))
		}
		// Making connection to server
		if soteErr.ErrCode == nil {
			if pJSMManager.nc, soteErr = pJSMManager.connect(); soteErr.ErrCode == nil {
				// Creating the JSM Manager
				pJSMManager.Manager, err = jsm.New(pJSMManager.nc)
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
func (jsmm *JSMManager) setCredentialsFile(streamCredentialFile string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(streamCredentialFile) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"streamCredentialFile"}), sError.EmptyMap)
	} else if _, err := os.Stat(streamCredentialFile); err != nil {
		soteErr = sError.GetSError(209010, sError.BuildParams([]string{streamCredentialFile, err.Error()}), sError.EmptyMap)
	} else {
		jsmm.connOpts = append(jsmm.connOpts, nats.UserCredentials(streamCredentialFile))
	}

	return
}

/*
	UserCredsFromRaw will take a credential file content that is not stored on the file system, such as AWS System Manager Parameters
*/
func (jsmm *JSMManager) userCredsFromRaw(rawCredentials []byte) nats.Option {
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
func (jsmm *JSMManager) setReconnectOptions(maxReconnect int, reconnectWait time.Duration) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if reconnectWait < 250*time.Millisecond || reconnectWait > 1*time.Minute {
		jsmm.connOpts = append(jsmm.connOpts, nats.ReconnectWait(1*time.Second))
		sLogger.Info("Resetting reconnectWait to 1 second")
	} else {
		jsmm.connOpts = append(jsmm.connOpts, nats.ReconnectWait(reconnectWait))
	}
	if maxReconnect < 1 || maxReconnect > 5 {
		jsmm.connOpts = append(jsmm.connOpts, nats.MaxReconnects(1))
		sLogger.Info("Resetting maxReconnect to 1 try")
	} else {
		jsmm.connOpts = append(jsmm.connOpts, nats.MaxReconnects(maxReconnect))
	}

	return
}

func (jsmm *JSMManager) setURL(sURL string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if _, err := url.Parse(sURL); err != nil || sURL == "" {
		soteErr = sError.GetSError(210090, sError.BuildParams([]string{sURL}), nil)
	} else {
		jsmm.sURL = sURL
	}

	return
}

/*
	This will connect to the NATS network using the values set in the JSMManager
*/
func (jsmm *JSMManager) connect() (nc *nats.Conn, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	// Connect to NATS  --  default "euwest1.aws.ngs.global"
	nc, err = nats.Connect(jsmm.sURL, jsmm.connOpts...)
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
	This will close the connection to the NATS network using the JSMManager
*/
func (jsmm *JSMManager) Close() {
	sLogger.DebugMethod()

	// Closing connect to NATS  --  default "euwest1.aws.ngs.global"
	jsmm.nc.Close()

	jsmm = nil

	return
}
