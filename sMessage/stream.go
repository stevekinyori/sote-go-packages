/*
General information about the streams. (Gathered from https://github.com/nats-io/jetstream)
STREAMS:
	Limit streams are used to control the size of the stream.  Limited streams can be limited by the following parameters:
		MaxAge of the message, (Sote supported)
		MaxBytes of the stream,
		MaxMsgs that can be in the stream. (Sote supported)
	When one of these limits are reached the messages in the stream will be discarded based on the discard setting.  Discard
	can be set to old or new.  Oldest record the newest record is removed.

	Interest streams retain messages so long as there is a consumer active for the subject.  At this time, this is not support by
	Sote sJetStream wrapper. Interest stream limits using age, size and count still apply as upper bounds.

	Work or Work Queue streams will retain the messages until the message is consumed by any one consumer. The message is then
	removed by the stream. Work stream limits using age, size and count still apply as upper bounds.
*/
package sMessage

import (
	"strings"

	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

const (
	M              = "m"
	MEMORY         = "memory"
)

var (
	sStreamConfigPtr *nats.StreamConfig
)

/*
	CreateLimitsStream will create a stream that has limited retention of messages.
*/
func (mmPtr *MessageManager) CreateLimitsStream(streamName, storage string, subjects []string, replicas int) (sStream *nats.StreamInfo,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = validateStreamName(streamName); soteErr.ErrCode == nil {
		if strings.ToLower(storage) == MEMORY || strings.ToLower(storage) == M {
			sStreamConfigPtr.Storage = 1
		}
		sStreamConfigPtr.Name = streamName
		sStreamConfigPtr.Subjects = subjects
		sStreamConfigPtr.Replicas = replicas
		js, err := mmPtr.NatsConnectionPtr.JetStream()
		if err != nil {
			mmPtr.natsErrorHandle(err, "", "", "", "")
		}
		sStream, err = js.AddStream(sStreamConfigPtr)
		if err != nil {
			mmPtr.natsErrorHandle(err, "", "", "", "")
		}
	}

	return
}

/*
	CreateWorkQueueStream will create a stream that once the message is pulled, it will be removed from the stream.
*/
func (mmPtr *MessageManager) CreateWorkQueueStream(streamName, storage string, subjects []string, replicas int) (sStream *nats.StreamInfo,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = validateStreamName(streamName); soteErr.ErrCode == nil {
		if strings.ToLower(storage) == MEMORY || strings.ToLower(storage) == M {
			sStreamConfigPtr.Storage = 1
		}
		sStreamConfigPtr.Name = streamName
		sStreamConfigPtr.Subjects = subjects
		sStreamConfigPtr.Replicas = replicas
		sStreamConfigPtr.Retention = nats.WorkQueuePolicy
		js, err := mmPtr.NatsConnectionPtr.JetStream()
		if err != nil {
			mmPtr.natsErrorHandle(err, "", "", "", "")
		}
		sStream, err = js.AddStream(sStreamConfigPtr)
		if err != nil {
			mmPtr.natsErrorHandle(err, "", "", "", "")
		}
	}

	return
}

/*
	DeleteStream will destroy the stream
*/
func (mmPtr *MessageManager) DeleteStream(streamName string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		mmPtr.natsErrorHandle(err, "", "", "", "")
	}

	err = js.DeleteStream(streamName)
	if err != nil {
		mmPtr.natsErrorHandle(err, "", "", "", "")
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
