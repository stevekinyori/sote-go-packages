/*
General information about the streams. (Gathered from https://github.com/nats-io/jetstream)
STREAMS:
	Limit streams are used to control the size of the stream.  Limited streams can be limited by the following parameters:
		MaxAge of the message,
		MaxBytes of the stream,
		MaxMsgs that can be in the stream.
	When one of these limits are reached the messages in the stream will be discarded based on the discard setting.  Discard
	can be set to old or new.  Oldest record the newest record is removed.

	Interest streams retain messages so long as there is a consumer active for the subject.  At this time, this is not support by
	Sote sMessage wrapper. Interest stream limits using age, size and count still apply as upper bounds.

	Work or Work Queue streams will retain the messages until the message is consumed by any one consumer. The message is then
	removed by the stream. Work stream limits using age, size and count still apply as upper bounds.
*/
package sMessage

import (
	"log"
	"strconv"
	"strings"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/jsm.go/api"
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

const (
	// F      = "f"
	// FILE   = "file"
	M      = "m"
	MEMORY = "memory"
)

/*
CreateLimitsStream will create a limits stream.  If it exists, it will load the stream
	Required parameters:
		streamName
		subjects (Multiple subject are allowed with a comma separating values)
			example: images,animals,cars,BILLS
		storage
			Sote defaults value: f
		replicas
			Sote defaults value: 1 (Sote Max value is 10)
		nc (pointer to a Jetstream connection)
*/
func CreateOrLoadLimitsStream(streamName, subjects, storage string, replicas int, jsmManager *jsm.Manager) (stream *jsm.Stream, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if soteErr = validateStreamParams(streamName, subjects, jsmManager); soteErr.ErrCode == nil {
		if strings.ToLower(storage) == MEMORY || strings.ToLower(storage) == M {
			stream, err = jsmManager.LoadOrNewStream(streamName, jsm.Subjects(subjects), jsm.MemoryStorage(), jsm.Replicas(setReplicas(replicas)), jsm.LimitsRetention())
		} else {
			stream, err = jsmManager.LoadOrNewStream(streamName, jsm.Subjects(subjects), jsm.FileStorage(), jsm.Replicas(setReplicas(replicas)), jsm.LimitsRetention())
		}
		if err != nil {
			soteErr = sError.GetSError(805000, sError.BuildParams([]string{streamName}), nil)
			log.Fatal(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
CreateOrLoadWorkStream will create a work stream.  If it exists, it will load the stream
	Required parameters:
		streamName
		subjects (Multiple subject are allowed with a comma separating values)
			example: images,animals,cars,BILLS
		storage
			Sote defaults value: f
		replicas
			Sote defaults value: 1 (Sote Max value is 10)
		nc (pointer to a Jetstream connection)
*/
func CreateOrLoadWorkStream(streamName, subjects, storage string, replicas int, jsmManager *jsm.Manager) (stream *jsm.Stream, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if soteErr = validateStreamParams(streamName, subjects, jsmManager); soteErr.ErrCode == nil {
		if strings.ToLower(storage) == MEMORY || strings.ToLower(storage) == M {
			stream, err = jsmManager.LoadOrNewStream(streamName, jsm.Subjects(subjects), jsm.MemoryStorage(), jsm.Replicas(setReplicas(replicas)), jsm.WorkQueueRetention())
		} else {
			stream, err = jsmManager.LoadOrNewStream(streamName, jsm.Subjects(subjects), jsm.FileStorage(), jsm.Replicas(setReplicas(replicas)), jsm.WorkQueueRetention())
		}
		if err != nil {
			soteErr = sError.GetSError(805000, sError.BuildParams([]string{streamName}), nil)
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
	StreamInfo loads an existing stream's information
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

func validateStreamParams(streamName, subjects string, jsmManager *jsm.Manager) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = validateStreamName(streamName); soteErr.ErrCode == nil && len(subjects) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"subjects"}), nil)
	}

	if soteErr.ErrCode == nil {
		soteErr = validateJSMManager(jsmManager)
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
		soteErr = sError.GetSError(335260, sError.BuildParams([]string{"NATS.io Stream"}), nil)
	}

	return
}

func setReplicas(tReplicas int) (replicas int) {
	sLogger.DebugMethod()

	if tReplicas <= 0 {
		replicas = 1
	} else if tReplicas > 10 {
		replicas = 10
	} else {
		replicas = tReplicas
	}

	return
}
