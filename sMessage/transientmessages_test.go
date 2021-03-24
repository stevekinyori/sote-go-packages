package sMessage

import (
	"testing"
	"time"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
)

// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPublish(tPtr *testing.T) {
	var (
		mmPtr *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if soteErr = mmPtr.Publish("greeting", "Hello world"); soteErr.ErrCode != nil {
			tPtr.Errorf("TestPublish Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr = mmPtr.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestSubscribe(tPtr *testing.T) {
	var (
		mmPtr *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if _, soteErr = mmPtr.Subscribe("greeting"); soteErr.ErrCode != nil {
			tPtr.Errorf("TestSubscribe Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr = mmPtr.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPublishRequest(tPtr *testing.T) {
	var (
		mmPtr *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if soteErr = mmPtr.PublishRequest("greeting", "greeting-reply", []byte("Back At You!")); soteErr.ErrCode != nil {
			tPtr.Errorf("TestPublishRequest Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr = mmPtr.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestSubscribeSync(tPtr *testing.T) {
	var (
		mmPtr *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if soteErr = mmPtr.SubscribeSync("greeting", "greeting-reply"); soteErr.ErrCode != nil {
			tPtr.Errorf("TestSubscribeSync Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr = mmPtr.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestNextMsg(tPtr *testing.T) {
	var (
		mmPtr *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if soteErr = mmPtr.SubscribeSync("greeting", "greeting-reply"); soteErr.ErrCode != nil {
			tPtr.Errorf("TestSubscribeSync Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
		if soteErr = mmPtr.Publish("greeting-reply", "Hello world"); soteErr.ErrCode != nil {
			tPtr.Errorf("TestPublish Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.NextMsg("greeting"); soteErr.ErrCode != nil {
			tPtr.Errorf("TestNextMsg Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr = mmPtr.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestRequest(tPtr *testing.T) {
	// TODO This code is not being tested at this time. How to test this needs to be investigated.
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestRequestReply(tPtr *testing.T) {
	var (
		mmPtr *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if _, soteErr = mmPtr.RequestReply("greeting", []byte("Hello World")); soteErr.ErrCode != nil {
			tPtr.Errorf("TestRequestReply Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr = mmPtr.Close()
}
