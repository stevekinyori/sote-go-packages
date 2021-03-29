/*
	This is a wrapper for Sote Golang developers to access services from NATS JetStream.
*/
package sMessage

import (
	"strconv"

	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

/*
	PPublish will send a persist message to the stream that owns the subject
*/
func (mmPtr *MessageManager) PPublish(subject, message string) (acknowledgement *nats.PubAck, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Subject: "] = subject

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
	PSubscribe will listen for message from the stream that owns the subject
*/
func (mmPtr *MessageManager) PSubscribe(subject, durableName string, callback nats.MsgHandler) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Subject"] = subject
	params["Durable Name"] = durableName

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
	PSubscribeSync will listen synchronously for message from the stream that owns the subject
*/
func (mmPtr *MessageManager) PSubscribeSync(subject, durableName string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Subject"] = subject
	params["Durable Name"] = durableName

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
	PPullSubscribe creates a subscription that can be used to fetch messages
*/
func (mmPtr *MessageManager) PullSubscribe(subject, durableName string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Subject"] = subject
	params["Durable Name"] = durableName

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}
	mmPtr.PullSubscriptions[durableName], err = js.PullSubscribe(subject, nats.SubOpt(nats.Durable(durableName)))
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

/*
	PDeleteMsg will remove a message from the stream
*/
func (mmPtr *MessageManager) DeleteMsg(streamName string, messageSequence int) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Stream Name"] = streamName
	params["Message Sequence"] = strconv.Itoa(messageSequence)

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
	PGetMsg retrieves a message using the sequence number directly from the stream
*/
func (mmPtr *MessageManager) GetMsg(streamName string, messageSequence int) (message *nats.RawStreamMsg, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Stream Name"] = streamName
	params["Message Sequence"] = strconv.Itoa(messageSequence)

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}
	message, err = js.GetMsg(streamName, uint64(messageSequence))
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

/*
	PFetch creates a subscription that can be used to fetch messages
*/
func (mmPtr *MessageManager) Fetch(durableName string, messageCount int) (messages []*nats.Msg, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	params := make(map[string]string)
	params["Durable Name"] = durableName
	params["Message Count"] = strconv.Itoa(messageCount)

	messages, err = mmPtr.Subscriptions[durableName].Fetch(messageCount)
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

/*
	PAck acknowledges a message
*/
func (mmPtr *MessageManager) Ack(message *nats.Msg) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)

	if err := message.Ack(); err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}
