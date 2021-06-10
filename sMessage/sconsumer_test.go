package sMessage

import (
	"runtime"
	"testing"
	"time"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
)

const (
	TESTCONSUMERNAMEPULL = "test-consumer-delete-me-pull"
	TESTCONSUMERNAMEPUSH = "test-consumer-delete-me-push"
	TESTDELIVERYSUBJECT = "TEST-ME"
)

func TestSetMaxDeliver(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	if maxDeliver := setMaxDeliver(-1); maxDeliver != 1 {
		tPtr.Errorf("%v Failed: Expected maxDeliver value to be 1", testName)
	}
	if maxDeliver := setMaxDeliver(0); maxDeliver != 1 {
		tPtr.Errorf("%v Failed: Expected maxDeliver value to be 1", testName)
	}
	if maxDeliver := setMaxDeliver(1); maxDeliver != 1 {
		tPtr.Errorf("%v Failed: Expected maxDeliver value to be 1", testName)
	}
	if maxDeliver := setMaxDeliver(11); maxDeliver != 3 {
		tPtr.Errorf("%v Failed: Expected maxDeliver value to be 10", testName)
	}
}
func TestPullReplayInstantConsumer(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		mmPtr             *MessageManager
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("%v Failed: Expected error code to be nil or 109999 got %v", testName, soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testPullSubjects, 1, false); soteErr.ErrCode == nil {
			if soteErr = mmPtr.CreatePullReplayInstantConsumer(TESTSTREAMNAME, TESTCONSUMERNAMEPULL, testPullSubjects[0], 1,
				false); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
			}
		}
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("%v Failed: Expected error code to be nil or 109999 got %v", testName, soteErr.FmtErrMsg)
		}
	}
}
func TestPushReplayInstantConsumer(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		mmPtr             *MessageManager
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("%v Failed: Expected error code to be nil or 109999 got %v", testName, soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testPushSubjects, 1, false); soteErr.ErrCode == nil {
			if soteErr = mmPtr.CreatePushReplayInstantConsumer(TESTSTREAMNAME, TESTCONSUMERNAMEPUSH, TESTDELIVERYSUBJECT, testPushSubjects[0], 1,
				false); soteErr.ErrCode != nil {
				tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
			}
		}
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("%v Failed: Expected error code to be nil or 109999 got %v", testName, soteErr.FmtErrMsg)
		}
	}
}
func TestGetConsumerInfo(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		mmPtr             *MessageManager
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("%v Failed: Expected error code to be nil or 109999 got %v", testName, soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testPullSubjects, 1, false); soteErr.ErrCode == nil {
			if soteErr = mmPtr.CreatePullReplayInstantConsumer(TESTSTREAMNAME, TESTCONSUMERNAMEPULL, testPullSubjects[0], 1,
				false); soteErr.ErrCode == nil {
				if _, soteErr = mmPtr.GetConsumerInfo(TESTSTREAMNAME, TESTCONSUMERNAMEPULL, false); soteErr.ErrCode != nil {
					tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
				}
			}
			if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
				tPtr.Errorf("%v Failed: Expected error code to be nil or 109999 got %v", testName, soteErr.FmtErrMsg)
			}
		}
	}
}
