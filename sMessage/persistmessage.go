/*
	This is a wrapper for Sote Golang developers to access services from NATS JetStream.
*/
package sMessage

import (
	"context"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

/*
	PPublish will send a persist message to the stream that owns the subject
*/
func (mmPtr *MessageManager) PPublish(subject, message string, testMode bool) (acknowledgement *nats.PubAck, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Subject: "] = subject
	params["testMode"] = strconv.FormatBool(testMode)

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}
	acknowledgement, err = js.Publish(subject, []byte(message))
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

/*
	PPublishMsg is not supported at this time
*/

/*
	PSubscribe will listen for message from the stream that owns the subject.
The subscription is saved in the an map of pull subscriptions in the Message Manager structure. The durable name is the index to the subscription.
*/
func (mmPtr *MessageManager) PSubscribe(subject, durableName string, callback nats.MsgHandler, testMode bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Subject"] = subject
	params["Durable Name"] = durableName
	params["testMode"] = strconv.FormatBool(testMode)

	if callback == nil {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"callback"}), nil)
	} else {
		js, err := mmPtr.NatsConnectionPtr.JetStream()
		if err != nil {
			soteErr = mmPtr.natsErrorHandle(err, params)
		}
		mmPtr.Subscriptions[durableName], err = js.Subscribe(subject, callback, nats.Durable(durableName))
		if err != nil {
			soteErr = mmPtr.natsErrorHandle(err, params)
		}
	}

	return
}

/*
	PSubscribeSync will listen synchronously for message from the stream that owns the subject.
The subscription is saved in the an map of pull subscriptions in the Message Manager structure. The durable name is the index to the subscription.
*/
func (mmPtr *MessageManager) PSubscribeSync(subject, durableName string, testMode bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Subject"] = subject
	params["Durable Name"] = durableName
	params["testMode"] = strconv.FormatBool(testMode)

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}
	mmPtr.Subscriptions[durableName], err = js.SubscribeSync(subject, nats.Durable(durableName))
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

/*
	PullSubscribe creates a subscription that can be used to fetch messages.
The subscription is saved in the an map of pull subscriptions in the Message Manager structure. The durable name is the index to the subscription.
*/
func (mmPtr *MessageManager) PullSubscribe(subject, durableName string, testMode bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Subject"] = subject
	params["Durable Name"] = durableName
	params["testMode"] = strconv.FormatBool(testMode)

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}
	mmPtr.PullSubscriptions[durableName], err = js.PullSubscribe(subject, durableName)
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

/*
	DeleteMsg will remove a message from the stream
*/
func (mmPtr *MessageManager) DeleteMsg(streamName string, messageSequence int, testMode bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Stream Name"] = streamName
	params["Message Sequence"] = strconv.Itoa(messageSequence)
	params["testMode"] = strconv.FormatBool(testMode)

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}
	err = js.DeleteMsg(streamName, uint64(messageSequence))
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

/*
	GetMsg retrieves a message using the sequence number directly from the stream
*/
func (mmPtr *MessageManager) GetMsg(streamName string, messageSequence int, testMode bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Stream Name"] = streamName
	params["Message Sequence"] = strconv.Itoa(messageSequence)
	params["testMode"] = strconv.FormatBool(testMode)

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}
	mmPtr.RawMessage, err = js.GetMsg(streamName, uint64(messageSequence))
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

/*
	Fetch creates a pull subscription that can be used to fetch messages.
With autoAck set to true each message fetched will be acknowledged before the method returns call to the caller.
Acknowledgement is needed for messages otherwise the consumer will choke with the max acknowledge limit is reached.
Acknowledgement can be manually done at anything so long as the pointer to the message still exists.
*/
func (mmPtr *MessageManager) Fetch(durableName string, messageCount int, autoAck, testMode bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	params := make(map[string]string)
	params["Durable Name"] = durableName
	params["Message Count"] = strconv.Itoa(messageCount)
	params["autoAck"] = strconv.FormatBool(autoAck)
	params["testMode"] = strconv.FormatBool(testMode)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	mmPtr.Messages, err = mmPtr.PullSubscriptions[durableName].Fetch(messageCount, nats.Context(ctx))
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}
	if autoAck {
		for _, message := range mmPtr.Messages {
			mmPtr.Ack(message, false)
			// 	TODO Review if err needs to be handled for failed Acks
		}
	}

	return
}

/*
	Ack acknowledges a message
*/
func (mmPtr *MessageManager) Ack(message *nats.Msg, testMode bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["testMode"] = strconv.FormatBool(testMode)

	if err := message.Ack(); err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}
