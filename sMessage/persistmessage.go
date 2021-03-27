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
func (mmPtr *MessageManager) PSubscribe(subject, subscriptionName string, callback nats.MsgHandler) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Subject"] = subject
	params["Subscription Name"] = subscriptionName

	if callback == nil {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"callback"}), nil)
	} else {
		js, err := mmPtr.NatsConnectionPtr.JetStream()
		if err != nil {
			soteErr = mmPtr.natsErrorHandle(err, params)
		}
		mmPtr.Subscriptions[subscriptionName], err = js.Subscribe(subject, callback)
		if err != nil {
			soteErr = mmPtr.natsErrorHandle(err, params)
		}
	}

	return
}

/*
	PSubscribeSync will listen synchronously for message from the stream that owns the subject
*/
func (mmPtr *MessageManager) PSubscribeSync(subject, subscriptionName string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Subject"] = subject
	params["Subscription Name"] = subscriptionName
	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}
	mmPtr.Subscriptions[subscriptionName], err = js.SubscribeSync(subject)
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

/*
	PDeleteMsg will remove a message from the stream
*/
func (mmPtr *MessageManager) PDeleteMsg(streamName string, messageSequence int) (soteErr sError.SoteError) {
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
func (mmPtr *MessageManager) PGetMsg(streamName string, messageSequence int) (message *nats.RawStreamMsg, soteErr sError.SoteError) {
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
	PPullSubscribe creates a subscription that can be used to fetch messages
*/
func (mmPtr *MessageManager) PPullSubscribe(subject, subscriptionName string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Subject"] = subject
	params["Subscription Name"] = subscriptionName
	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}
	mmPtr.Subscriptions[subscriptionName], err = js.PullSubscribe(subject)
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

/*
	PFetch creates a subscription that can be used to fetch messages
*/
func (mmPtr *MessageManager) PFetch(subscriptionName string, messageCount int) (messages []*nats.Msg, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	params := make(map[string]string)
	params["Subscription Name"] = subscriptionName
	params["Message Count"] = strconv.Itoa(messageCount)
	messages, err = mmPtr.Subscriptions[subscriptionName].Fetch(messageCount)
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

/*
	PChanSubscribe creates a chan based subscription
*/
func (mmPtr *MessageManager) PChanSubscribe(subject, subscriptionName string, channel chan *nats.Msg) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Subject"] = subject
	params["Subscription Name"] = subscriptionName
	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}
	mmPtr.Subscriptions[subscriptionName], err = js.ChanSubscribe(subject, channel)
	if err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}
