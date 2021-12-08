package sHelper

import (
	"testing"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sMessage"
)

func newRun() *Run {
	env, _ := NewEnvironment(ENVDEFAULTAPPNAME, ENVDEFAULTTARGET, ENVDEFAULTTARGET)
	return NewRun(env)
}

func testSubscribeListener(s *Subscriber, m *Msg) sError.SoteError {
	return sError.SoteError{}
}

func TestRunInitApp(t *testing.T) {
	run := newRun()
	run.ValidateEnvironment = func(environment string) sError.SoteError {
		AssertEqual(t, environment, ENVDEFAULTTARGET)
		return sError.SoteError{}
	}
	run.GetNATSURL = func(application, environment string) (string, sError.SoteError) {
		AssertEqual(t, application, ENVDEFAULTAPPNAME)
		AssertEqual(t, environment, ENVDEFAULTTARGET)
		return "west.eu.geo.ngs.global", sError.SoteError{}
	}
	run.NewMessage = func(env Environment, natsURL string) (*sMessage.MessageManager, sError.SoteError) {
		AssertEqual(t, natsURL, "west.eu.geo.ngs.global")
		return nil, sError.SoteError{}
	}
	soteErr := run.InitApp()
	AssertEqual(t, soteErr.ErrCode, nil)
}

func TestRunPanic(t *testing.T) {
	defer func() {
		r := recover()
		AssertEqual(t, r, "210599: Business Service error has occurred that is not expected.")
	}()
	run := newRun()
	run.PanicService(NewError().InternalError())
}

func TestRunError(t *testing.T) {
	defer func() {
		r := recover()
		AssertEqual(t, r, "210599: Business Service error has occurred that is not expected.")
	}()
	run := newRun()
	run.Error(NewError().InternalError(), &Msg{Subject: "test-subject"})
}

func TestRunListen(t *testing.T) {
	index := 0
	run := newRun()
	defer func() {
		r := recover()
		AssertEqual(t, r, "210599: Business Service error has occurred that is not expected.")
	}()
	run.AddSubscriber(&Subscriber{
		Run:          run,
		StreamName:   BSLSTREAMNAME,
		ConsumerName: "test-cosumer",
		Subject:      "test-subject",
	}, testSubscribeListener)
	run.Listen(func(*Subscriber) sError.SoteError {
		if index == 0 {
			index += 1
			return sError.SoteError{}
		} else {
			return NewError().InternalError() //panic
		}
	})
}

func TestRunReturnChainError(t *testing.T) {
	env, _ := NewEnvironment(ENVDEFAULTAPPNAME, "production", "production")
	run := NewRun(env)
	defer func() {
		r := recover()
		AssertEqual(t, r, "210599: Business Service error has occurred that is not expected.")
	}()
	run.AddSubscriber(&Subscriber{
		Run:          run,
		StreamName:   BSLSTREAMNAME,
		ConsumerName: "test-cosumer",
		Subject:      "test-subject",
	}, testSubscribeListener)
	run.Listen(func(s *Subscriber) (soteError sError.SoteError) {
		s.End(&Msg{Subject: "Test Subject"}, NewError().InvalidJson("Body"))
		return NewError().InternalError() //panic
	})
}
