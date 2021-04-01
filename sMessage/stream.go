/*
	General information about the streams. (Gathered from https://github.com/nats-io/jetstream)
*/
package sMessage

import (
	"strconv"
	"strings"

	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

const (
	LIMITSFILESTREAM = "limits-file"
	LIMITSMEMORYSTREAM = "limits-memory"
	WORKQUEUEFILESTREAM = "workqueue-file"
	WORKQUEUEMEMORYSTREAM = "workqueue-memory"
)

/*
	CreateLimitsStreamWithFileStorage will create a file based stream that has limited retention of messages.
*/
func (mmPtr *MessageManager) CreateLimitsStreamWithFileStorage(streamName string, subjects []string, replicas int,
	testMode bool) (sStream *nats.StreamInfo,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	sStream, soteErr = mmPtr.createStream(LIMITSFILESTREAM, streamName, subjects, replicas, testMode)

	return
}

/*
	CreateLimitsStreamWithMemoryStorage will create a memory based stream that has limited retention of messages.
*/
func (mmPtr *MessageManager) CreateLimitsStreamWithMemoryStorage(streamName string, subjects []string, replicas int,
	testMode bool) (sStream *nats.StreamInfo,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	sStream, soteErr = mmPtr.createStream(LIMITSMEMORYSTREAM, streamName, subjects, replicas, testMode)

	return
}

/*
	CreateWorkQueueStreamWithFileStorage will create a stream that once the message is pulled, it will be removed from the stream.
*/
func (mmPtr *MessageManager) CreateWorkQueueStreamWithFileStorage(streamName string, subjects []string, replicas int,
	testMode bool) (sStream *nats.StreamInfo,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	sStream, soteErr = mmPtr.createStream(WORKQUEUEFILESTREAM, streamName, subjects, replicas, testMode)

	return
}

/*
	CreateWorkQueueStreamWithFileStorage will create a stream that once the message is pulled, it will be removed from the stream.
*/
func (mmPtr *MessageManager) CreateWorkQueueStreamWithMemoryStorage(streamName string, subjects []string, replicas int,
	testMode bool) (sStream *nats.StreamInfo,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	sStream, soteErr = mmPtr.createStream(WORKQUEUEMEMORYSTREAM, streamName, subjects, replicas, testMode)

	return
}

/*
	DeleteStream will destroy the stream
*/
func (mmPtr *MessageManager) DeleteStream(streamName string, testMode bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Stream Name"] = streamName
	params["testMode"] = strconv.FormatBool(testMode)

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	err = js.DeleteStream(streamName)
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

func (mmPtr *MessageManager) createStream(streamType, streamName string, subjects []string, replicas int, testMode bool) (sStream *nats.StreamInfo,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		sStreamConfig *nats.StreamConfig
	)

	params := make(map[string]string)
	params["Stream Name"] = streamName
	params["Subjects"] = strings.Join(subjects, ", ")
	params["Replicas"] = strconv.Itoa(replicas)
	params["testMode"] = strconv.FormatBool(testMode)

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}
	// The default sote setting will change over time, so they are called out here.
	switch streamType {
	case LIMITSFILESTREAM:
		sStreamConfig = &nats.StreamConfig{
			Name:         streamName,
			Subjects:     subjects,
			Retention:    nats.LimitsPolicy,
			MaxConsumers: 0,
			MaxMsgs:      10000,
			MaxBytes:     0,
			Discard:      nats.DiscardOld,
			MaxAge:       1209600000000000,
			MaxMsgSize:   102400,
			Storage:      nats.FileStorage,
			Replicas:     replicas,
			NoAck:        false,
			Template:     "",
			Duplicates:   0,
			Placement:    nil,
			Mirror:       nil,
			Sources:      nil,
		}
	case LIMITSMEMORYSTREAM:
		sStreamConfig = &nats.StreamConfig{
			Name:         streamName,
			Subjects:     subjects,
			Retention:    nats.LimitsPolicy,
			MaxConsumers: 0,
			MaxMsgs:      10000,
			MaxBytes:     0,
			Discard:      nats.DiscardOld,
			MaxAge:       1209600000000000,
			MaxMsgSize:   102400,
			Storage:      nats.MemoryStorage,
			Replicas:     replicas,
			NoAck:        false,
			Template:     "",
			Duplicates:   0,
			Placement:    nil,
			Mirror:       nil,
			Sources:      nil,
		}
	case WORKQUEUEFILESTREAM:
		sStreamConfig = &nats.StreamConfig{
			Name:         streamName,
			Subjects:     subjects,
			Retention:    nats.WorkQueuePolicy,
			MaxConsumers: 0,
			MaxMsgs:      10000,
			MaxBytes:     0,
			Discard:      nats.DiscardOld,
			MaxAge:       1209600000000000,
			MaxMsgSize:   102400,
			Storage:      nats.FileStorage,
			Replicas:     replicas,
			NoAck:        false,
			Template:     "",
			Duplicates:   0,
			Placement:    nil,
			Mirror:       nil,
			Sources:      nil,
		}
	case WORKQUEUEMEMORYSTREAM:
		sStreamConfig = &nats.StreamConfig{
			Name:         streamName,
			Subjects:     subjects,
			Retention:    nats.WorkQueuePolicy,
			MaxConsumers: 0,
			MaxMsgs:      10000,
			MaxBytes:     0,
			Discard:      nats.DiscardOld,
			MaxAge:       1209600000000000,
			MaxMsgSize:   102400,
			Storage:      nats.MemoryStorage,
			Replicas:     replicas,
			NoAck:        false,
			Template:     "",
			Duplicates:   0,
			Placement:    nil,
			Mirror:       nil,
			Sources:      nil,
		}
	}

	sStream, err = js.AddStream(sStreamConfig)
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
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
