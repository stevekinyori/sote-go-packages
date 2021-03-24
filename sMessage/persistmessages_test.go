package sMessage

import (
	"testing"
	"time"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
)

// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPPublish(tPtr *testing.T) {
	var (
		mmPtr   *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if _, soteErr = mmPtr.PPublish("greeting", "Hello world"); soteErr.ErrCode != nil {
			tPtr.Errorf("TestPPublish Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr = mmPtr.Close()
}

// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPSubscribe(tPtr *testing.T) {
	var (
		mmPtr   *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if soteErr = mmPtr.PSubscribe("greeting", "mySub", nil); soteErr.ErrCode != nil {
			tPtr.Errorf("TestPSubscribe Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr = mmPtr.Close()
}

// // We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPSubscribeSync(tPtr *testing.T) {
	var (
		mmPtr   *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if soteErr = mmPtr.PSubscribeSync("greeting", "mySub"); soteErr.ErrCode != nil {
			tPtr.Errorf("TestPSubscribeSync Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr = mmPtr.Close()
}

// // We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPFetch(tPtr *testing.T) {
	var (
		mmPtr   *MessageManager
		soteErr sError.SoteError
		err     error
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if soteErr = mmPtr.PSubscribeSync("greeting", "mySub"); soteErr.ErrCode != nil {
			tPtr.Errorf("TestPFetch Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.PPublish("greeting-reply", "Hello world"); soteErr.ErrCode != nil {
			tPtr.Errorf("TestPFetch Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
		if _, err = mmPtr.Subscriptions["mySub"].Fetch(1); err != nil {
			tPtr.Errorf("TestPFetch Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr = mmPtr.Close()
}
