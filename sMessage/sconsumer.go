package sMessage

import (
	"log"
	"strings"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

/*
	CreateConsumer will create a consumer
		AckPolicy: Default value: none
			value is set using: none, all, explicit
			Sote defaults value: explicit
			Sote immutable: no
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
		ReplayPolicy: Default value: instant
			value is set using: instant, original
			Sote defaults value: instant
			Sote immutable: no
*/
func CreateConsumer(streamName, durableName, deliveryPolicy, deliverySubject, subjectFilter, replayPolicy string, maxDeliveries int, nc *nats.Conn) (pConsumer *jsm.Consumer,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if soteErr = validateStreamName(streamName); soteErr.ErrCode == nil {
		if len(durableName) == 0 && soteErr.ErrCode == nil {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"durableName"}), nil)
		}
		if len(deliverySubject) == 0 && soteErr.ErrCode == nil {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"deliverySubject"}), nil)
		}
		if len(subjectFilter) == 0 && soteErr.ErrCode == nil {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"subjectFilter"}), nil)
		}
		if len(replayPolicy) == 0 && soteErr.ErrCode == nil {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"replayPolicy"}), nil)
		}
		if (maxDeliveries <= 0 || maxDeliveries >= 4) && soteErr.ErrCode == nil {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"replayPolicy"}), nil)
		}
		if soteErr.ErrCode == nil {
			soteErr = validateConnection(nc)
		}
	}

	if soteErr.ErrCode == nil {
		switch strings.ToLower(deliveryPolicy) {
		case "all":
			pConsumer, err = jsm.NewConsumer(streamName, jsm.DurableName(durableName), jsm.DeliverAllAvailable(), jsm.FilterStreamBySubject(subjectFilter),
				jsm.ConsumerConnection(jsm.WithConnection(nc)))
		case "last":
			// TBD
		case "new":
			// TBD
		case "next":
			// TBD
		default:
			pConsumer, err = jsm.NewConsumer(streamName, jsm.DurableName(durableName), jsm.FilterStreamBySubject(subjectFilter),
				jsm.ConsumerConnection(jsm.WithConnection(nc)))
		}
	}

	if err != nil {
		soteErr = sError.GetSError(805599, sError.BuildParams([]string{streamName, durableName}), nil)
		log.Fatal(soteErr.FmtErrMsg)
	}

	return
}

func LoadConsumer(streamName, durableName string) (pConsumer *jsm.Consumer, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if soteErr = validateStreamName(streamName); soteErr.ErrCode == nil {
		pConsumer, err = jsm.LoadConsumer(streamName, durableName)
		if err != nil {
			soteErr = sError.GetSError(805000, nil, nil)
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

func validateConsumer(pConsumer *jsm.Consumer) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if pConsumer == nil {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"NATS.io Consumer"}), nil)
	}

	return
}
