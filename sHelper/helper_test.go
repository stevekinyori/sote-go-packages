package sHelper

import (
	"testing"

	"gitlab.com/soteapps/packages/v2021/sError"
)

func testListener(s *Subscriber, m *Msg) sError.SoteError {
	return sError.SoteError{}
}

func testNewHelper(t *testing.T) *Helper {
	var s Subscriber
	env, _ := NewEnvironment(ENVDEFAULTAPPNAME, ENVDEFAULTTARGET, ENVDEFAULTTARGET)
	helper := NewHelper(env)
	helper.CreateSubscriber = func(consumerName, subject string) *Subscriber {
		s = Subscriber{
			Run: helper.r,
		}
		s.Run.returnChain = make(chan *ReturnChain, 1)
		s.Run.Listen = func(listener func(*Subscriber) sError.SoteError) {
			listener(&s)
		}
		s.PullSubscribe = func() sError.SoteError {
			return sError.SoteError{}
		}
		s.Fetch = func() ([]Msg, sError.SoteError) {
			messages := make([]Msg, 1)
			messages[0] = Msg{
				Subject: "Test-subject",
				Data:    []byte("Data"),
			}
			return messages, sError.SoteError{}
		}
		return &s
	}
	helper.CreateDatabase = func() sError.SoteError {
		return sError.SoteError{}
	}
	helper.InitApp = func() sError.SoteError {
		return sError.SoteError{}
	}
	return helper
}

func TestNewHelperAddSubscriber(t *testing.T) {
	helper := testNewHelper(t)
	soteErr := helper.AddSubscriber("bsl-notification-wildcard", "bsl.notification.add", testListener, nil)
	AssertEqual(t, soteErr.FmtErrMsg, "")
}

func TestHelperMutipleSubscribers(t *testing.T) {
	helper := testNewHelper(t)
	soteErr := helper.AddSubscriber("bsl-notification-wildcard", "bsl.notification.add", testListener, nil)
	AssertEqual(t, soteErr.ErrCode, nil)
	soteErr = helper.AddSubscriber("bsl-notification-wildcard", "bsl.notification.remove", testListener, nil)
	AssertEqual(t, soteErr.ErrCode, nil)
}

func TestHelperWithSchema(t *testing.T) {
	helper := testNewHelper(t)
	schema := Schema{
		FileName:  "schema_test.json",
		StructRef: &TestSchema{},
	}
	soteErr := helper.AddSubscriber("bsl-notification-wildcard", "bsl.notification.add", testListener, &schema)
	AssertEqual(t, soteErr.ErrCode, nil)
}

func TestHelperRun(t *testing.T) {
	helper := testNewHelper(t)
	testListener := func(s *Subscriber, m *Msg) sError.SoteError {
		return sError.SoteError{}
	}
	soteErr := helper.AddSubscriber("bsl-notification-wildcard", "bsl.notification.add", testListener, nil)
	AssertEqual(t, soteErr.ErrCode, nil)
	helper.Run(false)
}

func TestHelperRunAsync(t *testing.T) {
	helper := testNewHelper(t)
	testListener := func(s *Subscriber, m *Msg) sError.SoteError {
		go func() {
			for rc := range s.Run.returnChain {
				AssertEqual(t, rc.soteErr.FmtErrMsg, "")
				AssertEqual(t, rc.msg.Subject, "Test-subject")
			}
		}()
		return sError.SoteError{}
	}
	soteErr := helper.AddSubscriber("bsl-notification-wildcard", "bsl.notification.add", testListener, nil)
	AssertEqual(t, soteErr.ErrCode, nil)
	helper.Run(true)
}

func TestHelperCreateSubscriber(t *testing.T) {
	helper := testNewHelper(t)
	s := helper.createSubscriber("bsl-notification-wildcard", "bsl.notification.add")
	AssertEqual(t, s != nil, true)
}

func TestHelperCreateDatabase(t *testing.T) {
	helper := testNewHelper(t)
	soteErr := helper.createDatabase()
	if soteErr.FmtErrMsg != "" &&
		soteErr.FmtErrMsg != "209299: No database connection has been established" &&
		soteErr.FmtErrMsg != "109999: /sote/api/staging/DB_NAME was/were not found" {
		AssertEqual(t, soteErr.FmtErrMsg, "")
	}
}

func TestHelperInitApp(t *testing.T) {
	helper := testNewHelper(t)
	helper.r.GetNATSURL = func(application, environment string) (string, sError.SoteError) {
		AssertEqual(t, application, ENVDEFAULTAPPNAME)
		AssertEqual(t, environment, ENVDEFAULTTARGET)
		return "west.eu.geo.ngs.global", sError.SoteError{}
	}
	soteErr := helper.initApp()
	if soteErr.FmtErrMsg != "" && soteErr.FmtErrMsg != "109999: /sote/synadia/tls-urlmask was/were not found" {
		AssertEqual(t, soteErr.FmtErrMsg, "")
	}
}
