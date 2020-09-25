/*
	sMessage supports the following NATS.io objects
		Stream:
			Limits
			Work
		Consumer:


	NATS.io Jetstream default values and Sotes defaults. This are values that are set and not changed in this package.
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
*/
package sMessage

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/jsm.go/api"
	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

const (
	// F      = "f"
	// FILE   = "file"
	M      = "m"
	MEMORY = "memory"
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

/*
	CreateLimitsStream will create a stream.  If it exists, it will return an error
		Name: Required
			value is set using: no spaces and case sensitive
			Sote defaults value: Required
			Sote immutable: no
		Replicas: Default value: 1
			value is set using: >0
			Sote defaults value: 1
			Sote immutable: no
		Storage: Required
			value is set using: (f)ile, (m)emory
			example: f, file, m, memory
			Sote defaults value: f
			Sote immutable: no
		Subjects: Required
			format of string: <subject>,<subject>,...
			example: images,animals,cars,BILLS
 			Sote defaults value: None
			Sote immutable: no
		--
		Retention: Default value: Limits
			value is set using: Limits, Interest, work queue, workq, work
			Sote defaults value: Limits
			Sote immutable: yes
*/
func CreateLimitsStream(streamName, subjects, storage string, replicas int, nc *nats.Conn) (stream *jsm.Stream, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if soteErr = validateStreamParams(streamName, subjects, nc); soteErr.ErrCode == nil {
		if strings.ToLower(storage) == MEMORY || strings.ToLower(storage) == M {
			stream, err = jsm.NewStream(streamName, jsm.Subjects(subjects), jsm.MemoryStorage(), jsm.Replicas(setReplicas(replicas)), jsm.StreamConnection(jsm.WithConnection(nc)),
				jsm.LimitsRetention())
		} else {
			stream, err = jsm.NewStream(streamName, jsm.Subjects(subjects), jsm.FileStorage(), jsm.Replicas(setReplicas(replicas)), jsm.StreamConnection(jsm.WithConnection(nc)), jsm.LimitsRetention())
		}
		if err != nil {
			soteErr = sError.GetSError(335280, sError.BuildParams([]string{streamName}), nil)
			log.Fatal(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
	CreateWorkStream will create a stream.  If it exists, it will return an error
		Name: Required
			value is set using: no spaces and case sensitive
			Sote defaults value: Required
			Sote immutable: no
		Replicas: Default value: 1
			value is set using: >0
			Sote defaults value: 1
			Sote immutable: no
		Storage: Required
			value is set using: (f)ile, (m)emory
			example: f, file, m, memory
			Sote defaults value: f
			Sote immutable: no
		Subjects: Required
			format of string: <subject>,<subject>,...
			example: images,animals,cars,BILLS
 			Sote defaults value: None
			Sote immutable: no
		--
		Retention: Default value: Limits
			value is set using: Limits, Interest, work queue, workq, work
			Sote defaults value: Limits
			Sote immutable: yes
*/
func CreateWorkStream(streamName, subjects, storage string, replicas int, nc *nats.Conn) (stream *jsm.Stream, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if soteErr = validateStreamParams(streamName, subjects, nc); soteErr.ErrCode == nil {
		if strings.ToLower(storage) == MEMORY || strings.ToLower(storage) == M {
			stream, err = jsm.NewStream(streamName, jsm.Subjects(subjects), jsm.MemoryStorage(), jsm.Replicas(setReplicas(replicas)), jsm.StreamConnection(jsm.WithConnection(nc)),
				jsm.LimitsRetention())
		} else {
			stream, err = jsm.NewStream(streamName, jsm.Subjects(subjects), jsm.FileStorage(), jsm.Replicas(setReplicas(replicas)), jsm.StreamConnection(jsm.WithConnection(nc)), jsm.LimitsRetention())
		}
		if err != nil {
			soteErr = sError.GetSError(335280, sError.BuildParams([]string{streamName}), nil)
			log.Fatal(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
	DeleteStream will destroy the stream
*/
func DeleteStream(pStream *jsm.Stream) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = validateStream(pStream); soteErr.ErrCode == nil {
		err := pStream.Delete()
		if err != nil {
			soteErr = sError.GetSError(805000, nil, nil)
			log.Fatal(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
	LoadStream loads an existing stream by name
*/
func LoadStream(streamName string) (pStream *jsm.Stream, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if soteErr = validateStreamName(streamName); soteErr.ErrCode == nil {
		pStream, err = jsm.LoadStream(streamName)
		if err != nil {
			soteErr = sError.GetSError(805000, nil, nil)
			log.Fatal(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
	StreamInfo loads an existing stream information
*/
func StreamInfo(pStream *jsm.Stream) (info *api.StreamInfo, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if soteErr = validateStream(pStream); soteErr.ErrCode == nil {
		info, err = pStream.LatestInformation()
		if err != nil {
			soteErr = sError.GetSError(805000, nil, nil)
			log.Fatal(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
	PurgeStream will remove all messages from the stream
*/
func PurgeStream(pStream *jsm.Stream) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = validateStream(pStream); soteErr.ErrCode == nil {
		err := pStream.Purge()
		if err != nil {
			soteErr = sError.GetSError(805000, nil, nil)
			log.Fatal(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
	DeleteMessageFromStream will remove a message from the stream
*/
func DeleteMessageFromStream(pStream *jsm.Stream, sequenceNumber int) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if sequenceNumber <= 0 {
		soteErr = sError.GetSError(400005, sError.BuildParams([]string{"sequenceNumber", "0"}), nil)
	} else {
		// TODO Test what happens if there is no message with the sequence number
		if soteErr = validateStream(pStream); soteErr.ErrCode == nil {
			err := pStream.DeleteMessage(sequenceNumber)
			if err != nil {
				if strings.Contains(err.Error(), "not found") {
					soteErr = sError.GetSError(109999, sError.BuildParams([]string{"sequenceNumber(" + strconv.Itoa(sequenceNumber) + ")"}), nil)
				} else {
					soteErr = sError.GetSError(805000, nil, nil)
					log.Fatal(soteErr.FmtErrMsg)
				}
			}
		}
	}

	return
}

/*
	CreateConsumer will load a consumer or create one based on the combination of the stream and durable name
		AckPolicy: Default value: none
			value is set using: none, all, explicit
			Sote defaults value: explicit
			Sote immutable: yes
		AckWait: Default value: -1s (forever)
			value is set using: (s)econds, (m)inutes, (h)ours, (y)ears, (M)onths, (d)ays
			Sote defaults value: 1s
			Sote immutable: yes
		DeliverPolicy: Default value: "" (pull based consumer)
			value is set using: all, last, new or next, DeliverByStartSequence or DeliverByStartTime
			Sote defaults value: all
			Sote immutable: yes
		DeliverySubject: Default value: instant
			value is set using: instant, original
			Sote defaults value: instant
			Sote immutable: yes
		Durable Name: Default value: ""
			value is set using: <durable name>
			example: TEST_CONSUMER_NAME, test_consumer_name, Test_Consumer_Name
			Sote defaults value: Required, not set
			Sote immutable: yes
		FilterSubject: Default value: "" (all)
			value is set using: <stream name>.<subject name>
			example: TEST_STREAM_NAME.* for all messages from the TEST_STREAM_NAME stream, TEST_STREAM_NAME.cat for only cat messages
			Sote defaults value: Required, not set
			Sote immutable: yes
		MaxDeliver: Default value: -1 (unlimited)
			value is set using: >0
			Sote defaults value: 3
			Sote immutable: yes to values of 1,2 or 3
		OptStartSeq: Default value: (Required when using DeliveryPolicy with DeliverByStartSequence
			value is set using: >0
			Sote defaults value: Not Supported
			Sote immutable: no
		ReplayPolicy: Default value: instant
			value is set using: instant, original
			Sote defaults value: instant
			Sote immutable: yes
		SampleFrequency: Default value: -1
			value is set using: 0 to 100
			Sote defaults value: Not Supported
			Sote immutable: no
		OptStartTime: Default value: Now
			value is set using: (s)econds, (m)inutes, (h)ours, (y)ears, (M)onths, (d)ays
			example: 1s, 1h, 1M, 1seconds, 1Months
			Sote defaults value: Not Supported
			Sote immutable: no
		RateLimit: Default value: 0
			value is set using: >0
			Sote defaults value: Not Supported
			Sote immutable: no
*/
func CreateConsumer(streamName, durableName, deliveryPolicy, deliverySubject, subjectFilter, replayPolicy string, maxDeliveries int, nc *nats.Conn) (consumer *jsm.Consumer,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if soteErr = validateStreamName(streamName); soteErr.ErrCode == nil {
		if len(durableName) == 0 && soteErr.ErrCode == nil {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"durableName"}), nil)
		}
		if len(deliveryPolicy) == 0 && soteErr.ErrCode == nil {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"subjectFilter"}), nil)
		}
		if len(deliverySubject) == 0 && soteErr.ErrCode == nil {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"subjectFilter"}), nil)
		}
		if len(subjectFilter) == 0 && soteErr.ErrCode == nil {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"subjectFilter"}), nil)
		}
		if len(replayPolicy) == 0 && soteErr.ErrCode == nil {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"subjectFilter"}), nil)
		}
		if soteErr.ErrCode == nil {
			soteErr = validateConnection(nc)
		}
	}

	if soteErr.ErrCode == nil {
		consumer, err = jsm.LoadOrNewConsumer(streamName, durableName, jsm.FilterStreamBySubject(subjectFilter), jsm.ConsumerConnection(jsm.WithConnection(nc)))
		if err != nil {
			soteErr = sError.GetSError(805599, sError.BuildParams([]string{streamName, durableName}), nil)
			log.Fatal(soteErr.FmtErrMsg)
		}
	}

	return
}

func validateStreamParams(streamName, subjects string, nc *nats.Conn) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = validateStreamName(streamName); soteErr.ErrCode == nil && len(subjects) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"subjects"}), nil)
	}

	if soteErr.ErrCode == nil {
		soteErr = validateConnection(nc)
	}

	return
}

func validateConnection(nc *nats.Conn) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if nc == nil {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"NATS.io Connect"}), nil)
	}

	return

}

func validateStreamName(streamName string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(streamName) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"streamName"}), nil)
	}

	return
}

func validateStream(pStream *jsm.Stream) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if pStream == nil {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"NATS.io Stream"}), nil)
	}

	return
}

func setReplicas(tReplicas int) (replicas int) {
	sLogger.DebugMethod()

	if tReplicas <= 0 {
		replicas = 1
	} else {
		replicas = tReplicas
	}

	return
}
