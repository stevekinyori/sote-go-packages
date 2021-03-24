/*
	This is a wrapper for Sote Golang developers to access services from NATS. This does not support JetStream.
*/
package sMessage

import (
	"errors"
	"time"

	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

/*
	Publish will push a message to NATS.
*/
func (mmPtr *MessageManager) Publish(subject string, data string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if err := mmPtr.NatsConnectionPtr.Publish(subject, []byte(data)); err != nil {
		soteErr = mmPtr.natsErrorHandle(err, subject, "", "", data)
	}
	return
}

/*
	Subscribe will express interest in the given subject. The subject can have wildcards (partial:*, full:>).
	Messages will be delivered to the associated MsgHandler.
*/
func (mmPtr *MessageManager) Subscribe(subject string) (msg *nats.Msg, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if _, err := mmPtr.NatsConnectionPtr.Subscribe(subject, func(msgIn *nats.Msg) {
		msg = msgIn
	}); err != nil {
		soteErr = mmPtr.natsErrorHandle(err, subject, "", "", string(msg.Data))
	}

	return
}

/*
	PublishRequest will perform a Publish() expecting a response on the reply subject.
*/
func (mmPtr *MessageManager) PublishRequest(subject string, reply string, data []byte) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if err := mmPtr.NatsConnectionPtr.PublishRequest(subject, reply, data); err != nil {
		soteErr = mmPtr.natsErrorHandle(err, subject, reply, "", string(data))
	}

	return
}

/*
	Subscribe will express interest in the given subject. The subject can have wildcards (partial:*, full:>).
	Messages will be delivered to the associated MsgHandler. Returns an error and the subscription.
*/
func (mmPtr *MessageManager) SubscribeSync(subscriptionName, subject string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if subscriptionName == "" {
		soteErr = mmPtr.natsErrorHandle(errors.New("SubscribeSync name must be populated"), subject, "", "", "")
	} else {
		if mmPtr.SyncSubscriptions[subscriptionName], err = mmPtr.NatsConnectionPtr.SubscribeSync(subject); err != nil {
			soteErr = mmPtr.natsErrorHandle(err, subject, "", "", "")
		}
	}

	return
}

/*
	NextMsg will return the next message available to a synchronous subscriber or block until one is available.
*/
func (mmPtr *MessageManager) NextMsg(subscriptionName string) (msg *nats.Msg, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if msg, err = mmPtr.SyncSubscriptions[subscriptionName].NextMsg(1 * time.Second); err != nil {
		soteErr = mmPtr.natsErrorHandle(err, "", "", subscriptionName, "")
	}

	return
}

/*
	Request will send a request payload and deliver the response message, or an error,
*/
func (mmPtr *MessageManager) Request(subject string, data []byte, time time.Duration) (msg *nats.Msg, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if msg, err = mmPtr.NatsConnectionPtr.Request(subject, data, time); err != nil {
		soteErr = mmPtr.natsErrorHandle(err, "", "", "", "")
	}

	return
}

/*
	RequestReply listens to a subject argument and sends data argument as reply to a request.
*/
func (mmPtr *MessageManager) RequestReply(subject string, data []byte) (s *nats.Subscription, soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var err error

	_, err = mmPtr.NatsConnectionPtr.Subscribe(subject, func(msg *nats.Msg) {
		if err = msg.Respond(data); err != nil {
			soteErr = mmPtr.natsErrorHandle(err, "", "", "", "")
		}
	})

	return
}
