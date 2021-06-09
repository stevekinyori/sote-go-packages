package sMessage

import (
	"runtime"
	"testing"
	"time"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
)

const (
	TESTSTREAMNAME    = "test-stream-delete-me"
)

var (
	testPullSubjects = []string{"test-subject-delete-me-1","test-subject-delete-me-2"}
	testPushSubjects = []string{"test-subject-delete-me-push-1","test-subject-delete-me-push-2"}
)

func TestSetReplicas(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)
	if replicas := setReplicas(-1); replicas != 1 {
		tPtr.Errorf("%v Failed: Expected replicase value to be 1", testName)
	}
	if replicas := setReplicas(0); replicas != 1 {
		tPtr.Errorf("%v Failed: Expected replicase value to be 1", testName)
	}
	if replicas := setReplicas(1); replicas != 1 {
		tPtr.Errorf("%v Failed: Expected replicase value to be 1", testName)
	}
	if replicas := setReplicas(11); replicas != 10 {
		tPtr.Errorf("%v Failed: Expected replicase value to be 10", testName)
	}
}
func TestCreateLimitsStreamWithFileStorage(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr sError.SoteError
		mmPtr   *MessageManager
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v failed: Expected soteErr to be nil got %v", testName, soteErr.FmtErrMsg)
	}

	if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
		tPtr.Errorf("%v Failed: Expected error code to be nil or 109999 got %v", testName, soteErr.FmtErrMsg)
	}

	if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testPullSubjects, 1,false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
	}

	if soteErr := mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
	}
}
func TestCreateLimitsStreamWithMemoryStorage(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr sError.SoteError
		mmPtr   *MessageManager
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v failed: Expected soteErr to be nil got %v", testName, soteErr.FmtErrMsg)
	}

	if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
		tPtr.Errorf("%v Failed: Expected error code to be nil or 109999 got %v", testName, soteErr.FmtErrMsg)
	}

	if _, soteErr = mmPtr.CreateLimitsStreamWithMemoryStorage(TESTSTREAMNAME, testPullSubjects, 1, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
	}

	if soteErr := mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
	}
}
func TestCreateWorkQueueStreamWithFileStorage(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr sError.SoteError
		mmPtr   *MessageManager
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v failed: Expected soteErr to be nil got %v", testName, soteErr.FmtErrMsg)
	}

	if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
		tPtr.Errorf("%v Failed: Expected error code to be nil or 109999 got %v", testName, soteErr.FmtErrMsg)
	}

	if _, soteErr = mmPtr.CreateWorkQueueStreamWithFileStorage(TESTSTREAMNAME, testPullSubjects, 1, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
	}

	if soteErr := mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
	}
}
func TestCreateWorkQueueStreamWithMemoryStorage(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr sError.SoteError
		mmPtr   *MessageManager
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v failed: Expected soteErr to be nil got %v", testName, soteErr.FmtErrMsg)
	}

	if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
		tPtr.Errorf("%v Failed: Expected error code to be nil or 109999 got %v", testName, soteErr.FmtErrMsg)
	}

	if _, soteErr = mmPtr.CreateWorkQueueStreamWithMemoryStorage(TESTSTREAMNAME, testPullSubjects, 1,false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
	}

	if soteErr := mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
	}
}
