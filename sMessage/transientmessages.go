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
func (mmPtr *MessageManager) Publish(subject string, message string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Subject: "] = subject
	if err := mmPtr.NatsConnectionPtr.Publish(subject, []byte(message)); err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}
	return
}

/*
	Subscribe will express interest in the given subject. The subject can have wildcards (partial:*, full:>).
	Messages will be delivered to the associated MsgHandler.
*/
func (mmPtr *MessageManager) Subscribe(subject string) (msg *nats.Msg, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Subject: "] = subject
	if _, err := mmPtr.NatsConnectionPtr.Subscribe(subject, func(msgIn *nats.Msg) {
		msg = msgIn
	}); err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

/*
	PublishRequest will perform a Publish() expecting a response on the reply subject.
*/
func (mmPtr *MessageManager) PublishRequest(subject string, reply string, message string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	params := make(map[string]string)
	params["Subject: "] = subject
	params["Reply: "] = reply
	if err := mmPtr.NatsConnectionPtr.PublishRequest(subject, reply, []byte(message)); err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
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

	params := make(map[string]string)
	params["Subscription Name: "] = subscriptionName
	params["Subject: "] = subject
	if subscriptionName == "" {
		soteErr = mmPtr.natsErrorHandle(errors.New("SubscribeSync name must be populated"), params)
	} else {
		if mmPtr.SyncSubscriptions[subscriptionName], err = mmPtr.NatsConnectionPtr.SubscribeSync(subject); err != nil {
			soteErr = mmPtr.natsErrorHandle(err, params)
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

	params := make(map[string]string)
	params["Subscription Name: "] = subscriptionName
	if msg, err = mmPtr.SyncSubscriptions[subscriptionName].NextMsg(1 * time.Second); err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

/*
	Request will send a request payload and deliver the response message, or an error,
*/
func (mmPtr *MessageManager) Request(subject string, message string, time time.Duration) (msg *nats.Msg, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	params := make(map[string]string)
	params["Subject: "] = subject
	params["Time: "] = time.String()
	if msg, err = mmPtr.NatsConnectionPtr.Request(subject, []byte(message), time); err != nil {
		soteErr = mmPtr.natsErrorHandle(err, params)
	}

	return
}

/*
	RequestReply listens to a subject argument and sends data argument as reply to a request.
*/
func (mmPtr *MessageManager) RequestReply(subject string, message string) (s *nats.Subscription, soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var err error

	params := make(map[string]string)
	params["Subject: "] = subject
	_, err = mmPtr.NatsConnectionPtr.Subscribe(subject, func(msg *nats.Msg) {
		if err = msg.Respond([]byte(message)); err != nil {
			soteErr = mmPtr.natsErrorHandle(err, params)
		}
	})

	return
}
