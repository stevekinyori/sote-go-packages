/*
	This is a wrapper for Sote Golang developers to access services from NATS JetStream.
*/
package sMessage

import (
	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

/*
	PPublish will send a persist message to the stream that owns the subject
*/
func (mmPtr *MessageManager) PPublish(subject, message string) (acknowledgement *nats.PubAck, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		mmPtr.natsErrorHandle(err, subject, "", "", message)
	}
	acknowledgement, err = js.Publish(subject, []byte(message))
	if err != nil {
		mmPtr.natsErrorHandle(err, subject, "", "", message)
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

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		mmPtr.natsErrorHandle(err, subject, "", subscriptionName, "")
	}
	mmPtr.Subscriptions[subscriptionName], err = js.Subscribe(subject, callback)
	if err != nil {
		mmPtr.natsErrorHandle(err, subject, "", subscriptionName, "")
	}

	return
}

/*
	PSubscribeSync will listen synchronously for message from the stream that owns the subject
*/
func (mmPtr *MessageManager) PSubscribeSync(subject, subscriptionName string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		mmPtr.natsErrorHandle(err, subject, "", subscriptionName, "")
	}
	mmPtr.Subscriptions[subscriptionName], err = js.SubscribeSync(subject)
	if err != nil {
		mmPtr.natsErrorHandle(err, subject, "", subscriptionName, "")
	}

	return
}

/*
	PDeleteMsg will remove a message from the stream
*/
func (mmPtr *MessageManager) PDeleteMsg(streamName string, messageSequence uint64) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		mmPtr.natsErrorHandle(err, "", "", "", "")
	}
	err = js.DeleteMsg(streamName, messageSequence)
	if err != nil {
		mmPtr.natsErrorHandle(err, "", "", "", "")
	}

	return
}

/*
	PChanSubscribe creates a chan based subscription
*/
func (mmPtr *MessageManager) PChanSubscribe(subject, subscriptionName string, channel chan *nats.Msg) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		mmPtr.natsErrorHandle(err, subject, "", subscriptionName, "")
	}
	mmPtr.Subscriptions[subscriptionName], err = js.ChanSubscribe(subject, channel)
	if err != nil {
		mmPtr.natsErrorHandle(err, subject, "", subscriptionName, "")
	}

	return
}

/*
	PGetMsg retrieves a message using the sequence number directly from the stream
*/
func (mmPtr *MessageManager) PGetMsg(streamName string, messageSequence uint64) (message *nats.RawStreamMsg, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		mmPtr.natsErrorHandle(err, "", "", "", "")
	}
	message, err = js.GetMsg(streamName, messageSequence)
	if err != nil {
		mmPtr.natsErrorHandle(err, "", "", "", "")
	}

	return
}

/*
	PPullSubscribe creates a subscription that can be used to fetch messages
*/
func (mmPtr *MessageManager) PPullSubscribe(subject, subscriptionName string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	js, err := mmPtr.NatsConnectionPtr.JetStream()
	if err != nil {
		mmPtr.natsErrorHandle(err, subject, "", subscriptionName, "")
	}
	mmPtr.Subscriptions[subscriptionName], err = js.PullSubscribe(subject)
	if err != nil {
		mmPtr.natsErrorHandle(err, subject, "", subscriptionName, "")
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

	messages, err = mmPtr.Subscriptions[subscriptionName].Fetch(messageCount)
	if err != nil {
		mmPtr.natsErrorHandle(err, "", "", subscriptionName, "")
	}

	return
}
