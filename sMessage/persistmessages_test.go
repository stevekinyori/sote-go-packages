package sMessage

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
)

var (
	mmPtr       *MessageManager
	initTestRun bool
)

// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPPublish(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	if !initTestRun {
		soteErr = initTest()
	}
	if soteErr.ErrCode == nil {
		if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world", false); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
		}
	}

}

// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPSubscribe(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	if !initTestRun {
		soteErr = initTest()
	}
	if soteErr.ErrCode == nil {
		if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world", false); soteErr.ErrCode == nil {
			if soteErr = mmPtr.PSubscribe(testSubjects[0], TESTCONSUMERNAME, nil, false); soteErr.ErrCode != 200513 {
				tPtr.Errorf("%v Failed: Expected error code to be 200513 got %v", testName, soteErr.FmtErrMsg)
			}
			if soteErr = mmPtr.PSubscribe(testSubjects[0], TESTCONSUMERNAME, func(msgIn *nats.Msg) {
				return
			}, false); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
			}
		}
	}
}

// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPSubscribeSync(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	if !initTestRun {
		soteErr = initTest()
	}
	if soteErr.ErrCode == nil {
		if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world", false); soteErr.ErrCode == nil {
			if soteErr = mmPtr.PSubscribeSync(testSubjects[0], TESTCONSUMERNAME, false); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
			}
		}
	}
}

// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPullSubscribe(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	if !initTestRun {
		soteErr = initTest()
	}
	if soteErr.ErrCode == nil {
		if soteErr = mmPtr.PullSubscribe(testSubjects[0], TESTCONSUMERNAME, false); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
		}
	}
}

// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestDeleteMsg(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	if !initTestRun {
		soteErr = initTest()
	}
	if soteErr.ErrCode == nil {
		if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world", false); soteErr.ErrCode == nil {
			if soteErr = mmPtr.DeleteMsg(TESTSTREAMNAME, 1, false); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
			}
		}
	}
}

// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestGetMsg(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	if !initTestRun {
		soteErr = initTest()
	}
	if soteErr.ErrCode == nil {
		if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world", false); soteErr.ErrCode == nil {
			if soteErr = mmPtr.GetMsg(TESTSTREAMNAME, 1, false); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
			}
		}
	}
}

// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestFetch(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	if !initTestRun {
		soteErr = initTest()
	}
	if soteErr.ErrCode == nil {
		if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world", false); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
		}
		if soteErr = mmPtr.PullSubscribe(testSubjects[0], TESTCONSUMERNAME, false); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
		}
		if soteErr = mmPtr.Fetch(TESTCONSUMERNAME, 1, true, false); soteErr.ErrCode != nil && soteErr.ErrCode != 101010 {
			tPtr.Errorf("%vFetch Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
		}
	}
}

// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestAck(tPtr *testing.T) {
	var (
		soteErr           sError.SoteError
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	if !initTestRun {
		soteErr = initTest()
	}
	if soteErr.ErrCode == nil {
		if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world", false); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
		}
		if soteErr = mmPtr.PullSubscribe(testSubjects[0], TESTCONSUMERNAME, false); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
		}
		// time.Sleep(5 * time.Millisecond)
		if soteErr = mmPtr.Fetch(TESTCONSUMERNAME, 1, false, false); soteErr.ErrCode != nil {
			fmt.Println("Msg: " + soteErr.FmtErrMsg)
			tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
		}
		if mmPtr.Messages != nil {
			if soteErr = mmPtr.Ack(mmPtr.Messages[0], true); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
			}
		} else {
			tPtr.Errorf("%v Failed: Expected at least one message, got Nil", testName)
		}
	}
}
func TestFINALTESTCLEANUP(t *testing.T) {
	cleanUpTest()
}
func initTest() (soteErr sError.SoteError) {
	if !initTestRun {
		if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1, 250*time.Millisecond,
			false); soteErr.ErrCode == nil {
			if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode == nil {
				if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1, false); soteErr.ErrCode == nil {
					initTestRun = true
				}
			}
		}
	}

	return
}
func cleanUpTest() (soteErr sError.SoteError) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	if initTestRun {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode == nil {
			mmPtr.Close()
		} else {
			soteErr = sError.GetSError(199999, sError.BuildParams([]string{testName}), sError.EmptyMap)
		}
	}
	initTestRun = false

	return
}
