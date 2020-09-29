/*
General information about the consumers. (Gathered from https://github.com/nats-io/jetstream)
CONSUMERS:
	Consumers come in two flavors, push and pull.  Push is expecting that the service is available at all times.  Pull holds the
	messages until the process is active. There are many setting for these two types,which effect the behavior of the
	consumer.  Only a limited set of push and pull consumers will are supported below. Consumers can either be push based where
	JetStream will deliver the messages as fast as possible to a subject of your choice or pull based for typical work queue like
	behavior.

	Consumer setting supported
		AckPolicy: Default value: none
			value is set using: none, all, explicit
			NOTE: Is only for push streams. "Explicit" is for each message while "all" will cover all message in the stream or if a sequence is provided,
				all the message up to that sequence number.  Example using 200 message and ack'ing message 100 will ack 1 to 99 also. Good for
				performance boost when high value of messages are being processed. "None" means no ack is needed.
			Sote defaults value: explicit
			Sote immutable: no
		DeliverPolicy: Default value: "" (pull based consumer)
			value is set using: all, last, next, DeliverByStartSequence or DeliverByStartTime
			NOTE: This is where in the stream will start sending messages to the consumer.
			Sote defaults value: all
			Sote immutable: no
		DeliverySubject: Default value: instant
			value is set using: instant, original
			NOTE:
			Sote defaults value: instant
			Sote immutable: no
		Durable Name: Default value: ""
			value is set using: <durable name>
			example: TEST_CONSUMER_NAME, test_consumer_name, Test_Consumer_Name
			NOTE:
			Sote defaults value: Required, not set
			Sote immutable: no
		FilterSubject: Default value: "" (all)
			value is set using: <stream name>.<subject name>
			example: TEST_STREAM_NAME.* for all messages from the TEST_STREAM_NAME stream, TEST_STREAM_NAME.cat for only cat messages
			NOTE:
			Sote defaults value: Required, not set
			Sote immutable: no
		MaxDeliver: Default value: -1 (unlimited)
			value is set using: >0
			Sote defaults value: 3
			Sote immutable: no Allowed values are 1,2 or 3
		ReplayPolicy: Default value: instant
			value is set using: instant, original
			NOTE: Instant will send the message to the consumer as fast as possible. Original will replay as received.
			Sote defaults value: instant
			Sote immutable: no
*/
package sMessage

import (
	"log"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

const (
	ACKPOLICYNONE     = "none"
	ACKPOLICYALL      = "all"
	ACKPOLICYEXPLICIT = "explicit"
	//
	DELIVERYPOLICYALL  = "all"
	DELIVERYPOLICYLAST = "last"
	DELIVERYPOLICYNEXT = "next"
	DELIVERYPOLICYPULL = ""
	//
	REPLAYPOLICYINSTANT  = "instant"
	REPLAYPOLICYORIGINAL = "original"
)

/*
	CreatePullConsumerWithReplayInstantMax3 will create a consumer. If the consumer exists, it will load.
		Set values for this function
			AckPolicy: explicit (required for pull consumer)
			DeliverPolicy: explicit
			DeliverySubject: ""
			ReplayPolicy: instant
*/
func CreatePulConsumerWithReplayInstantMax3(streamName, durableName, subjectFilter string, maxDeliveries int, nc *nats.Conn) (pConsumer *jsm.Consumer,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if soteErr = validateConsumerParams(streamName, durableName, DELIVERYPOLICYPULL, subjectFilter, nc); soteErr.ErrCode == nil {
		pConsumer, err = jsm.LoadOrNewConsumer(streamName, durableName, jsm.DeliverAllAvailable(), jsm.FilterStreamBySubject(subjectFilter), jsm.ConsumerConnection(jsm.WithConnection(nc)),
			jsm.ReplayInstantly(), jsm.MaxDeliveryAttempts(3))
		if err != nil {
			soteErr = sError.GetSError(805000, sError.BuildParams([]string{streamName, durableName}), nil)
			log.Fatal(soteErr.FmtErrMsg)
		}
	}

	return
}

func DeleteConsumer(pConsumer *jsm.Consumer) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = validateConsumer(pConsumer); soteErr.ErrCode == nil {
		err := pConsumer.Delete()
		if err != nil {
			soteErr = sError.GetSError(805000, nil, nil)
			log.Fatal(soteErr.FmtErrMsg)
		}
	}

	return
}

func validateConsumerParams(streamName, durableName, deliverySubject, subjectFilter string, nc *nats.Conn) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	soteErr = validateStreamName(streamName)

	if soteErr.ErrCode == nil {
		soteErr = validateDurableName(durableName)
	}

	if soteErr.ErrCode == nil && deliverySubject != DELIVERYPOLICYPULL {
		soteErr = validateDeliverySubject(deliverySubject)
	}

	if soteErr.ErrCode == nil {
		soteErr = validateSubjectFilter(subjectFilter)
	}

	if soteErr.ErrCode == nil {
		soteErr = validateConnection(nc)
	}

	return
}

func validateDurableName(durableName string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(durableName) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"durableName"}), nil)
	}

	return
}

func validateDeliverySubject(subjectFilter string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(subjectFilter) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"subjectFilter"}), nil)
	}

	return
}

func validateSubjectFilter(deliverySubject string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(deliverySubject) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"deliverySubject"}), nil)
	}

	return
}

func validateConsumer(pConsumer *jsm.Consumer) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if pConsumer == nil {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"NATS.io Consumer"}), nil)
	}

	return
}
