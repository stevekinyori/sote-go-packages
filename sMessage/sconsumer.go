/*
	Consumer setting supported
		AckPolicy: Default value: none
			value is set using: none, all, explicit
			Sote defaults value: explicit
			Sote immutable: no
		DeliverPolicy: Default value: "" (pull based consumer)
			value is set using: all, last, next, DeliverByStartSequence or DeliverByStartTime
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
		ReplayPolicy: Default value: instant
			value is set using: instant, original
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
	CreateConsumerWithDeliverAllReplayInstantMax3 will create a consumer. If the consumer exists, it will load.
		Set values for this function
			AckPolicy: explicit
			DeliverPolicy: all
			ReplayPolicy: instant
*/
func CreateDeliverAllReplayInstantConsumer(streamName, durableName, deliverySubject, subjectFilter string, maxDeliveries int, nc *nats.Conn) (pConsumer *jsm.Consumer,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	pConsumer, err = jsm.LoadOrNewConsumer(streamName, durableName, jsm.DeliverAllAvailable(), jsm.FilterStreamBySubject(subjectFilter), jsm.ConsumerConnection(jsm.WithConnection(nc)),
		jsm.ReplayInstantly(), jsm.MaxDeliveryAttempts(3))
	if err != nil {
		soteErr = sError.GetSError(805599, sError.BuildParams([]string{streamName, durableName}), nil)
		log.Fatal(soteErr.FmtErrMsg)
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
		if validateDurableName(durableName); soteErr.ErrCode != nil {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"durableName"}), nil)
		}
	}

	if soteErr.ErrCode == nil {
		if validateDeliverySubject(deliverySubject); soteErr.ErrCode != nil {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"deliverySubject"}), nil)
		}
	}

	if soteErr.ErrCode == nil {
		if validateSubjectFilter(subjectFilter); soteErr.ErrCode != nil {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"subjectFilter"}), nil)
		}
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
