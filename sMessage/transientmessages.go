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
func (mm *MessageManager) Publish(subject string, data []byte) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if err := mm.NatsConnection.Publish(subject, data); err != nil {
		soteErr = mm.natsErrorHandle(err, subject, "", "", data)
	}
	return
}

/*
	Subscribe will express interest in the given subject. The subject can have wildcards (partial:*, full:>).
	Messages will be delivered to the associated MsgHandler.
*/
func (mm *MessageManager) Subscribe(subject string) (msg *nats.Msg, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if _, err := mm.NatsConnection.Subscribe(subject, func(msgIn *nats.Msg) {
		msg = msgIn
	}); err != nil {
		soteErr = mm.natsErrorHandle(err, subject, "", "", msg.Data)
	}

	return
}

/*
	PublishRequest will perform a Publish() expecting a response on the reply subject.
*/
func (mm *MessageManager) PublishRequest(subject string, reply string, data []byte) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if err := mm.NatsConnection.PublishRequest(subject, reply, data); err != nil {
		soteErr = mm.natsErrorHandle(err, subject, reply, "", data)
	}

	return
}

/*
	Subscribe will express interest in the given subject. The subject can have wildcards (partial:*, full:>).
	Messages will be delivered to the associated MsgHandler. Returns an error and the subscription.
*/
func (mm *MessageManager) SubscribeSync(subscriptionName, subject string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if subscriptionName == "" {
		soteErr = mm.natsErrorHandle(errors.New("SubscribeSync name must be populated"), subject, "", "", nil)
	} else {
		if mm.SyncSubscriptions[subscriptionName], err = mm.NatsConnection.SubscribeSync(subject); err != nil {
			soteErr = mm.natsErrorHandle(err, subject, "", "", nil)
		}
	}

	return
}

/*
	NextMsg will return the next message available to a synchronous subscriber or block until one is available.
*/
func (mm *MessageManager) NextMsg(subscriptionName string) (msg *nats.Msg, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if msg, err = mm.SyncSubscriptions[subscriptionName].NextMsg(1 * time.Second); err != nil {
		soteErr = mm.natsErrorHandle(err, "", "", subscriptionName, nil)
	}

	return
}

/*
	Request will send a request payload and deliver the response message, or an error,
*/
func (mm *MessageManager) Request(subject string, data []byte, time time.Duration) (msg *nats.Msg, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if msg, err = mm.NatsConnection.Request(subject, data, time); err != nil {
		soteErr = mm.natsErrorHandle(err, "", "", "", nil)
	}

	return
}

/*
	RequestReply listens to a subject argument and sends data argument as reply to a request.
*/
func (mm *MessageManager) RequestReply(subject string, data []byte) (s *nats.Subscription, soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var err error

	_, err = mm.NatsConnection.Subscribe(subject, func(msg *nats.Msg) {
		if err = msg.Respond(data); err != nil {
			soteErr = mm.natsErrorHandle(err, "", "", "", nil)
		}
	})

	return
}
