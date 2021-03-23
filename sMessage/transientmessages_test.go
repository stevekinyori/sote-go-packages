package sMessage

import (
	"testing"
	"time"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
)

// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPublish(t *testing.T) {
	var (
		nm *MessageManager
		soteErr sError.SoteError
	)

	if nm, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if soteErr = nm.Publish("greeting", []byte("Hello world")); soteErr.ErrCode != nil {
			t.Errorf("TestPublish Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	nm = nm.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestSubscribe(t *testing.T) {
	var (
		nm *MessageManager
		soteErr sError.SoteError
	)

	if nm, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if _, soteErr = nm.Subscribe("greeting"); soteErr.ErrCode != nil {
			t.Errorf("TestSubscribe Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	nm = nm.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPublishRequest(t *testing.T) {
	var (
		nm *MessageManager
		soteErr sError.SoteError
	)

	if nm, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if soteErr = nm.PublishRequest("greeting", "greeting-reply", []byte("Back At You!")); soteErr.ErrCode != nil {
			t.Errorf("TestPublishRequest Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	nm = nm.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestSubscribeSync(t *testing.T) {
	var (
		nm *MessageManager
		soteErr sError.SoteError
	)

	if nm, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if soteErr = nm.SubscribeSync("greeting", "greeting-reply"); soteErr.ErrCode != nil {
			t.Errorf("TestSubscribeSync Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	nm = nm.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestNextMsg(t *testing.T) {
	var (
		nm *MessageManager
		soteErr sError.SoteError
	)

	if nm, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if soteErr = nm.SubscribeSync("greeting", "greeting-reply"); soteErr.ErrCode != nil {
			t.Errorf("TestSubscribeSync Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
		if soteErr = nm.Publish("greeting-reply", []byte("Hello world")); soteErr.ErrCode != nil {
			t.Errorf("TestPublish Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = nm.NextMsg("greeting"); soteErr.ErrCode != nil {
			t.Errorf("TestNextMsg Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	nm = nm.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestRequest(t *testing.T) {
	// TODO This code is not being tested at this time. How to test this needs to be investigated.
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestRequestReply(t *testing.T) {
	var (
		nm *MessageManager
		soteErr sError.SoteError
	)

	if nm, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if _, soteErr = nm.RequestReply("greeting", []byte("Hello World")); soteErr.ErrCode != nil {
			t.Errorf("TestRequestReply Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	nm = nm.Close()
}
