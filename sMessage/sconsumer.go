/*
General information about the consumers. (Gathered from https://github.com/nats-io/jetstream)
CONSUMERS:
	Consumers come in two flavors, push and pull.  Push is expecting that the service is available at all times.  Pull holds the
	messages until the process is active. There are many setting for these two types,which effect the behavior of the
	consumer.  Only a limited set of push and pull consumers are supported below. Consumers can either be push based where
	JetStream will deliver the messages as fast as possible to a subject of your choice or pull based for typical work queue like
	behavior.

	The Consumer subject filter must be a subset of the Stream subject. Error Code: 206700
	Here are some examples:
		Stream subject: IMAGES  Consumer subject filter: IMAGES -> works
		Stream subject: IMAGES  Consumer subject filter: IMAGES.cat -> DOES NOT works
		--
		Stream subject: IMAGES.*  Consumer subject filter: IMAGES -> works
		Stream subject: IMAGES.*  Consumer subject filter: IMAGES.cat -> works
		--
		Stream subject: IMAGES.CATS  Consumer subject filter: IMAGES -> DOES NOT works
		Stream subject: IMAGES.cat  Consumer subject filter: IMAGES.* -> works
*/
package sMessage

import (
	"strconv"

	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

const (
	PULLREPLAYINSTANTCONSUMER = "pull-replay-instant"
	PUSHREPLAYINSTANTCONSUMER = "push-replay-instant"
)

/*
CreatePullConsumerWithReplayInstant will create a consumer. If the consumer exists, it will load.
	This will read all the message in the stream without concern of the order.  This is should be used when transaction
	order doesn't matter.
	Required parameters:
		streamName
		durableName
		subjectFilter
		maxDeliver
			Sote defaults value: 1 (Sote Max value is 10, if great it is set to 3)

	Set values:
		AckPolicy: explicit (explicit is required for a pull consumer)
		DeliverPolicy: all
		DeliverySubject: "" (nil string is required for a pull consumer)
		ReplayPolicy: instant
*/
func (mmPtr *MessageManager) CreatePullReplayInstantConsumer(streamName, durableName, subjectFilter string,
	maxDeliveries int, testMode bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	_, soteErr = mmPtr.createConsumer(PULLREPLAYINSTANTCONSUMER, streamName, durableName, "", subjectFilter, maxDeliveries, testMode)

	return
}

/*
CreatePushReplayInstantConsumer will create a consumer. If the consumer exists, it will load.
	This will read all the message in the stream without concern of the order.  CreatePushReplayInstantConsumer should be used when transaction
	order doesn't matter.
	Required parameters:
		streamName
		durableName
		deliversubject
		subjectFilter
		maxDeliver
			Sote defaults value: 1 (Sote Max value is 10, if great it is set to 3)

	Set values:
		AckPolicy: explicit (explicit is required for a pull consumer)
		DeliverPolicy: all
		DeliverySubject: "" (nil string is required for a pull consumer)
		ReplayPolicy: instant
*/
func (mmPtr *MessageManager) CreatePushReplayInstantConsumer(streamName, durableName, deliverySubject, subjectFilter string,
	maxDeliveries int, testMode bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	_, soteErr = mmPtr.createConsumer(PUSHREPLAYINSTANTCONSUMER, streamName, durableName, deliverySubject, subjectFilter, maxDeliveries, testMode)

	return
}

func (mmPtr *MessageManager) DeleteConsumer(streamName, durableName string, testMode bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Stream Name"] = streamName
	params["Durable Name"] = durableName
	params["testMode"] = strconv.FormatBool(testMode)

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	err = js.DeleteConsumer(streamName, durableName)
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

func (mmPtr *MessageManager) createConsumer(consumerType, streamName, durableName, deliverySubject, subjectFilter string, maxDeliveries int,
	testMode bool) (sConsumer *nats.ConsumerInfo, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		sConsumerConfig *nats.ConsumerConfig
	)

	// The default sote setting will change over time, so they are called out here.
	switch consumerType {
	case PULLREPLAYINSTANTCONSUMER:
		sConsumerConfig = &nats.ConsumerConfig{
			Durable:         durableName,
			DeliverSubject:  "",
			DeliverPolicy:   nats.DeliverAllPolicy,
			OptStartSeq:     0,
			OptStartTime:    nil,
			AckPolicy:       nats.AckExplicitPolicy,
			AckWait:         0,
			MaxDeliver:      setMaxDeliver(maxDeliveries),
			FilterSubject:   subjectFilter,
			ReplayPolicy:    nats.ReplayInstantPolicy,
			RateLimit:       0,
			SampleFrequency: "",
			MaxWaiting:      0,
			MaxAckPending:   0,
		}
	case PUSHREPLAYINSTANTCONSUMER:
		// SAMPLE: nats con add ORDERS MONITOR --filter '' --ack none --target monitor.ORDERS --deliver last --replay instant
		sConsumerConfig = &nats.ConsumerConfig{
			Durable:         durableName,
			DeliverSubject:  deliverySubject,
			DeliverPolicy:   nats.DeliverAllPolicy,
			OptStartSeq:     0,
			OptStartTime:    nil,
			AckPolicy:       nats.AckExplicitPolicy,
			AckWait:         0,
			MaxDeliver:      setMaxDeliver(maxDeliveries),
			FilterSubject:   subjectFilter,
			ReplayPolicy:    nats.ReplayInstantPolicy,
			RateLimit:       0,
			SampleFrequency: "",
			MaxWaiting:      0,
			MaxAckPending:   0,
		}
	}

	params := make(map[string]string)
	params["Stream Name"] = streamName
	params["Consumer Type"] = consumerType
	params["Durable Name"] = durableName
	params["Delivery Subject"] = deliverySubject
	params["Filter Subject"] = subjectFilter
	params["Max Deliveries"] = strconv.Itoa(maxDeliveries)
	params["testMode"] = strconv.FormatBool(testMode)

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	sConsumer, err = js.AddConsumer(streamName, sConsumerConfig)
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

func (mmPtr *MessageManager) GetConsumerInfo(streamName, durableName string, testMode bool) (sConsumer *nats.ConsumerInfo, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Stream Name"] = streamName
	params["Durable Name"] = durableName
	params["testMode"] = strconv.FormatBool(testMode)

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	sConsumer, err = js.ConsumerInfo(streamName, durableName)
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

/*
	setMaxDeliver forces the value be between 1 and 10.  If it below 1, then 1 and if greater than 10, its set to 3.
*/
func setMaxDeliver(tMaxDeliver int) (maxDeliver int) {
	sLogger.DebugMethod()

	if tMaxDeliver <= 0 {
		maxDeliver = 1
	} else if tMaxDeliver > 10 {
		maxDeliver = 3
	} else {
		maxDeliver = tMaxDeliver
	}

	return
}
