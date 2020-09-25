/*
	sMessage supports the following NATS.io objects
		Stream:
			Limits
			Work
		Consumer:


	NATS.io Jetstream default values and Sotes defaults. This are values that are set and not changed in this package.
	STREAM SETTINGS:
		Ack: Required for published
			value is set using: true, false
			Sote defaults value: true
			Sote immutable: yes
		Discard: Default value: old
			value is set using: new, old
			Sote defaults value: old
			Sote immutable: yes
		Duplicates: Default value: ""
			value is set using: (s)econds, (m)inutes, (h)ours, (y)ears, (M)onths, (d)ays
			example: 1s, 1h, 1M, 1seconds, 1Months
			Sote defaults value: ""
			Sote immutable: yes
		MaxAge: Default value: -1 (unlimited)
			value is set using: (s)econds, (m)inutes, (h)ours, (y)ears, (M)onths, (d)ays
			example: 1s, 1h, 1M, 1seconds, 1Months
			Sote defaults value: -1 (unlimited)
			Sote immutable: yes
		MaxBytes: Default value: -1 (unlimited)
			value is set using: (B)ytes, (k)ilobytes, (m)egabytes
			example: 512B, 1k, 1m
			Sote defaults value: -1 (unlimited)
			Sote immutable: yes
		MaxConsumers: Default value: -1 (unlimited)
			value is set using: >0
			Sote defaults value: -1 (unlimited)
			Sote immutable: yes
		MaxMsgs: Required
			value is set using: -1 (unlimited) or >0
			Sote defaults value: -1 (unlimited)
			Sote immutable: yes
		MaxMsgSize: Default value: -1 (unlimited)
			value is set using: (B)ytes, (k)ilobytes, (m)egabytes
			example: 512B, 1k, 1m
			Sote defaults value: -1 (unlimited)
			Sote immutable: yes
	CONSUMER SETTINGS:
		AckPolicy: Default value: none
			value is set using: none, all, explicit
			Sote defaults value: explicit
			Sote immutable: no
		AckWait: Default value: -1s (forever)
			value is set using: (s)econds, (m)inutes, (h)ours, (y)ears, (M)onths, (d)ays
			Sote defaults value: 1s
			Sote immutable: yes
		DeliverPolicy: Default value: "" (pull based consumer)
			value is set using: all, last, new or next, DeliverByStartSequence or DeliverByStartTime
			Sote defaults value: all
			Sote immutable: no
		DeliverySubject: Default value: instant
			value is set using: instant, original
			Sote defaults value: instant
			Sote immutable: no
		Durable Name: Default value: ""
			value is set using: <durable name>
			example: TEST_CONSUMER_NAME, test_consumer_name, Test_Consumer_Name
			Sote defaults value: Required, not set
			Sote immutable: no
		FilterSubject: Default value: "" (all)
			value is set using: <stream name>.<subject name>
			example: TEST_STREAM_NAME.* for all messages from the TEST_STREAM_NAME stream, TEST_STREAM_NAME.cat for only cat messages
			Sote defaults value: Required, not set
			Sote immutable: no
		MaxDeliver: Default value: -1 (unlimited)
			value is set using: >0
			Sote defaults value: 3
			Sote immutable: no to values of 1,2 or 3
		OptStartSeq: Default value: (Required when using DeliveryPolicy with DeliverByStartSequence
			value is set using: >0
			Sote defaults value: Not Supported
			Sote immutable: yes
		ReplayPolicy: Default value: instant
			value is set using: instant, original
			Sote defaults value: instant
			Sote immutable: no
		SampleFrequency: Default value: -1
			value is set using: 0 to 100
			Sote defaults value: Not Supported
			Sote immutable: yes
		OptStartTime: Default value: Now
			value is set using: (s)econds, (m)inutes, (h)ours, (y)ears, (M)onths, (d)ays
			example: 1s, 1h, 1M, 1seconds, 1Months
			Sote defaults value: Not Supported
			Sote immutable: yes
		RateLimit: Default value: 0
			value is set using: >0
			Sote defaults value: Not Supported
			Sote immutable: yes
*/
package sMessage

import (
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

func SetAllOptions(streamName, streamCredentialFile, streamCredentialToken string, maxReconnect int,
	reconnectWait, timeOut time.Duration) (opts []nats.Option, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if opts, soteErr = SetStreamName(streamName); soteErr.ErrCode == nil {
		if opts, soteErr = SetCredentialsFile(streamCredentialFile); soteErr.ErrCode == nil {
			if opts, soteErr = SetCredentialsToken(streamCredentialToken); soteErr.ErrCode == nil {
				if opts, soteErr = SetReconnectOptions(maxReconnect, reconnectWait); soteErr.ErrCode == nil {
					opts, soteErr = SetTimeOut(timeOut)
				}
			}
		}
	}

	return
}

func SetStreamName(streamName string) (opts []nats.Option, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = validateStreamName(streamName); soteErr.ErrCode == nil {
		opts = []nats.Option{nats.Name(streamName)}
	}

	return
}

/*
	streamCredentialToken will take precedence over streamCredentialFile
*/
func SetCredentialsFile(streamCredentialFile string) (opts []nats.Option, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(streamCredentialFile) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"streamCredentialFile"}), nil)
	} else if _, err := os.Stat(streamCredentialFile); err != nil {
		soteErr = sError.GetSError(600010, sError.BuildParams([]string{streamCredentialFile, err.Error()}), nil)
	} else {
		opts = []nats.Option{nats.UserCredentials(streamCredentialFile)}
	}

	return
}

/*
	streamCredentialToken will take precedence over streamCredentialFile
*/
func SetCredentialsToken(streamCredentialToken string) (opts []nats.Option, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(streamCredentialToken) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"streamCredentialToken"}), nil)
	} else {
		opts = []nats.Option{nats.Token(streamCredentialToken)}
	}

	return
}

func SetReconnectOptions(maxReconnect int, reconnectWait time.Duration) (opts []nats.Option, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if maxReconnect == 0 && reconnectWait == 0 {
		soteErr = sError.GetSError(200512, sError.BuildParams([]string{"maxReconnect", "reconnectWait"}), nil)
	} else {
		if reconnectWait == 0 {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"reconnectWait"}), nil)
		} else {
			opts = []nats.Option{nats.ReconnectWait(reconnectWait)}
		}
		if maxReconnect == 0 {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"reconnectWait"}), nil)
		} else {
			opts = []nats.Option{nats.MaxReconnects(maxReconnect)}
		}
	}

	return
}

func SetTimeOut(timeOut time.Duration) (opts []nats.Option, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if timeOut == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"timeOut"}), nil)
	} else {
		opts = []nats.Option{nats.Timeout(timeOut)}
	}

	return
}

// This will connect to the NATS network
func Connect(url string, opts []nats.Option) (nc *nats.Conn, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(url) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"url"}), nil)
	}

	if opts == nil {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"Options (opts)"}), nil)
	}

	if soteErr.ErrCode == nil {
		var err error
		// Connect to NATS  --  default "euwest1.aws.ngs.global"
		nc, err = nats.Connect(url, opts...)
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
				panic("sDatabase.sconnection.GetConnection Failed")
			}
		}
		// defer nc.Close()
	}

	return
}
