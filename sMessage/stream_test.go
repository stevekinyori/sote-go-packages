package sMessage

import (
	"testing"
	"time"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
)

const (
	TESTSTREAMNAME    = "test-stream-delete-me"
)

var (
	testSubjects = []string{"test-subject-delete-me-1","test-subject-delete-me-2"}
)

func TestSetReplicas(tPtr *testing.T) {
	if replicas := setReplicas(-1); replicas != 1 {
		tPtr.Errorf("TestSetReplicas Failed: Expected replicase value to be 1")
	}
	if replicas := setReplicas(0); replicas != 1 {
		tPtr.Errorf("TestSetReplicas Failed: Expected replicase value to be 1")
	}
	if replicas := setReplicas(1); replicas != 1 {
		tPtr.Errorf("TestSetReplicas Failed: Expected replicase value to be 1")
	}
	if replicas := setReplicas(11); replicas != 10 {
		tPtr.Errorf("TestSetReplicas Failed: Expected replicase value to be 10")
	}
}
func TestCreateLimitsStreamWithFileStorage(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		mmPtr   *MessageManager
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStreamWithFileStorage failed: Expected soteErr to be nil got %v", soteErr.FmtErrMsg)
	}

	if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
		tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
	}

	if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1,false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}

	if soteErr := mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}
}
func TestCreateLimitsStreamWithMemoryStorage(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		mmPtr   *MessageManager
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStreamWithMemoryStorage failed: Expected soteErr to be nil got %v", soteErr.FmtErrMsg)
	}

	if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
		tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
	}

	if _, soteErr = mmPtr.CreateLimitsStreamWithMemoryStorage(TESTSTREAMNAME, testSubjects, 1, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStreamWithMemoryStorage Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}

	if soteErr := mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStreamWithMemoryStorage Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}
}
func TestCreateWorkQueueStreamWithFileStorage(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		mmPtr   *MessageManager
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateWorkQueueStreamWithFileStorage failed: Expected soteErr to be nil got %v", soteErr.FmtErrMsg)
	}

	if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
		tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
	}

	if _, soteErr = mmPtr.CreateWorkQueueStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateWorkQueueStreamWithFileStorage Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}

	if soteErr := mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateWorkQueueStreamWithFileStorage Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}
}
func TestCreateWorkQueueStreamWithMemoryStorage(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		mmPtr   *MessageManager
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateWorkQueueStreamWithMemoryStorage failed: Expected soteErr to be nil got %v", soteErr.FmtErrMsg)
	}

	if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
		tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
	}

	if _, soteErr = mmPtr.CreateWorkQueueStreamWithMemoryStorage(TESTSTREAMNAME, testSubjects, 1,false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateWorkQueueStreamWithMemoryStorage Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}

	if soteErr := mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateWorkQueueStreamWithMemoryStorage Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}
}
