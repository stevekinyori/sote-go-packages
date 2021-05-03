package sHelper

import (
	"testing"

	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sMessage"
)

func newSubscriber() *Subscriber {
	env, _ := NewEnvironment(ENVDEFAULTAPPNAME, ENVDEFAULTTARGET, ENVDEFAULTTARGET)
	run := NewRun(env)
	run.myMMPtr = &sMessage.MessageManager{}
	return NewSubscriber(run, "test-consumer", "test-subject")
}

func TestNewSubscriber(t *testing.T) {
	s := newSubscriber()
	AssertEqual(t, s.ConsumerName, "test-consumer")
	AssertEqual(t, s.Subject, "test-subject")
	AssertEqual(t, s.StreamName, STREAMNAME)

	s = NewSubscriber(s.Run, "test-consumer", "test-subject", "test-stream")
	AssertEqual(t, s.ConsumerName, "test-consumer")
	AssertEqual(t, s.Subject, "test-subject")
	AssertEqual(t, s.StreamName, "test-stream")
}

func TestSubscribe(t *testing.T) {
	defer func() {
		recover()
	}()
	s := newSubscriber()
	s.subscribe() //Expect to get an error NatsConnectionPtr is nil
}

func TestSubscribeConsumerInfo(t *testing.T) {
	defer func() {
		recover()
	}()
	s := newSubscriber()
	s.getConsumerInfo() //Expect to get an error NatsConnectionPtr is nil
}

func TestSubscribeFetch(t *testing.T) {
	s := newSubscriber()
	s.GetConsumerInfo = func() (*nats.ConsumerInfo, sError.SoteError) {
		return &nats.ConsumerInfo{NumPending: 1}, sError.SoteError{}
	}
	s.DoFetch = func(consumerInfo *ConsumerInfo) sError.SoteError {
		s.Run.myMMPtr.Messages = make([]*nats.Msg, consumerInfo.NumPending)
		s.Run.myMMPtr.Messages[0] = &nats.Msg{
			Subject: "Subject",
			Data:    []byte("Data"),
		}
		return sError.SoteError{}
	}
	s.fetch()
}

func TestSubscribeStart(t *testing.T) {
	s := newSubscriber()
	data := []byte("Test Data")
	s.Start(&Msg{
		uuid:    UUID(UUIDKind.Short),
		Subject: "Test Subject",
		Data:    data,
	})
}

func TestSubscribeEnd(t *testing.T) {
	s := newSubscriber()
	s.Run.returnChain = make(chan *ReturnChain, 1)
	go func() {
		for rc := range s.Run.returnChain {
			AssertEqual(t, rc.soteErr.FmtErrMsg, "210599: Business Service error has occurred that is not expected.")
			AssertEqual(t, rc.msg.Subject, "Test-subject")
		}
	}()
	s.End(&Msg{Subject: "Test-subject"}, NewError().InternalError())
}

func TestSubscribeConsumerError(t *testing.T) {
	defer func() {
		recover()
	}()
	s := newSubscriber()
	s.GetConsumerInfo = func() (*nats.ConsumerInfo, sError.SoteError) {
		return &nats.ConsumerInfo{NumPending: 1}, sError.SoteError{}
	}
	s.fetch() //Expect to get an error NatsConnectionPtr is nil
}

func TestSubscribeFetchError(t *testing.T) {
	defer func() {
		recover()
	}()
	s := newSubscriber()
	s.GetConsumerInfo = func() (*nats.ConsumerInfo, sError.SoteError) {
		return nil, NewError().InternalError()
	}
	s.fetch() //Expect to get an error NatsConnectionPtr is nil
}

func TestSubscribePublish(t *testing.T) {
	defer func() {
		recover()
	}()
	s := newSubscriber()
	s.publish("Hello") //Expect to get an error NatsConnectionPtr is nil
}
