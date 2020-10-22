/*
NATS.io Jetstream default values are listed below. (Gathered from https://github.com/nats-io/jetstream)
These are values that can be set natively.  sstream and consumer place limitation on the configurations available from NATS.

	NATS STREAM SETTINGS:
		Ack: Required for published
			value is set using: true, false
		Discard: Default value: old
			value is set using: new, old
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
package sMessage

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
	"gitlab.com/soteapps/packages/v2020/sConfigParams"
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

type JSMManager struct {
	Manager     jsm.Manager
	Application string
	Environment string
	sURL        string
	opts        []nats.Option

	sync.Mutex
}

/*
	New will create a Sote Jetstream Manager for NATS.  The application and environment are required.
	credentialFileName is optional and should not be used except in development.
*/
func New(application, environment, credentialFileName, sURL string, maxReconnect int, reconnectWait time.Duration) (pJSMManager *JSMManager, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		// tmpFileName string
		nc          *nats.Conn
		tmpCreds    interface{}
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
			pJSMManager.opts = append(pJSMManager.opts, pJSMManager.UserCredsFromRaw([]byte(tmpCreds.(string))))
		}
		// Making connection to server
		if soteErr.ErrCode == nil {
			nc, soteErr = pJSMManager.connect()
		}
	}

	log.Println(nc)

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
		soteErr = sError.GetSError(600010, sError.BuildParams([]string{streamCredentialFile, err.Error()}), sError.EmptyMap)
	} else {
		jsmm.opts = append(jsmm.opts, nats.UserCredentials(streamCredentialFile))
	}

	return
}

func (jsmm *JSMManager) UserCredsFromRaw(rawdata []byte) nats.Option {
	return nats.UserJWT(
		func() (string, error) { return nkeys.ParseDecoratedJWT(rawdata) },
		func(nonce []byte) ([]byte, error) {
			kp, err := nkeys.ParseDecoratedNKey(rawdata)
			if err != nil {
				return nil, err
			}
			defer kp.Wipe()
			sig, _ := kp.Sign(nonce)
			return sig, nil
		})
}

/*
	setReconnectOptions limits the maxReconnect to 5 and the highest reconnectWait of 1 second.
*/
func (jsmm *JSMManager) setReconnectOptions(maxReconnect int, reconnectWait time.Duration) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if maxReconnect == 0 && reconnectWait == 0 {
		soteErr = sError.GetSError(200512, sError.BuildParams([]string{"maxReconnect", "reconnectWait"}), sError.EmptyMap)
	} else {
		if reconnectWait == 0 {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"reconnectWait"}), sError.EmptyMap)
		} else {
			if reconnectWait > 1*time.Minute {
				jsmm.opts = append(jsmm.opts, nats.ReconnectWait(1*time.Second))
			} else {
				jsmm.opts = append(jsmm.opts, nats.ReconnectWait(reconnectWait))
			}
		}
		if maxReconnect == 0 {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"reconnectWait"}), sError.EmptyMap)
		} else {
			if maxReconnect > 5 {
				jsmm.opts = append(jsmm.opts, nats.MaxReconnects(5))
			} else {
				jsmm.opts = append(jsmm.opts, nats.MaxReconnects(maxReconnect))
			}
		}
	}

	return
}

func (jsmm *JSMManager) setURL(sURL string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if _, err := url.Parse(sURL); err != nil {
		// TODO Create Error Code
		soteErr = sError.GetSError(100000, nil, nil)
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
	nc, err = nats.Connect(jsmm.sURL, jsmm.opts...)
	if err != nil {
		if strings.Contains(err.Error(), "no servers") {
			soteErr = sError.GetSError(603999, nil, sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		} else {
			var errDetails = make(map[string]string)
			errDetails, soteErr = sError.ConvertErr(err)
			if soteErr.ErrCode != nil {
				sLogger.Info(soteErr.FmtErrMsg)
				panic("sError.ConvertErr Failed")
			}
			sLogger.Info(sError.GetSError(805000, nil, errDetails).FmtErrMsg)
			panic("sMessages.connect Failed")
		}
	}
	defer nc.Close()

	return
}

func SetStreamName(streamName string) (opts []nats.Option, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = validateStreamName(streamName); soteErr.ErrCode == nil {
		opts = []nats.Option{nats.Name(streamName)}
	}

	return
}

func GetJSMManager(nc *nats.Conn) (jsmManager *jsm.Manager, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if nc == nil {
		soteErr = sError.GetSError(603999, nil, sError.EmptyMap)
	}

	// TODO Change Connect to non-exported and call from GetJSMManager

	jsmManager, err = jsm.New(nc, jsm.WithTimeout(2*time.Second))
	if err != nil {
		log.Fatal(err.Error())
	}

	return
}

func GetJSMManagerWithConnOptions(url string, opts []nats.Option) (jsmManager *jsm.Manager, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		nc *nats.Conn
	)

	// nc, soteErr = Connect(url, opts)
	if soteErr.ErrCode == nil {
		jsmManager, soteErr = GetJSMManager(nc)
	}

	return
}

func validateConnection(nc *nats.Conn) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if nc == nil {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"NATS.io Connect"}), sError.EmptyMap)
	}

	return

}

func validateJSMManager(jsmManager *jsm.Manager) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if jsmManager == nil {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"Jetstream Manager"}), sError.EmptyMap)
	}

	if !jsmManager.IsJetStreamEnabled() {
		soteErr = sError.GetSError(300000, nil, sError.EmptyMap)
	}
	return

}
