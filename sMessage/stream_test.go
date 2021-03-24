package sMessage

import (
	"testing"
	"time"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
)

const (
	TESTSTREAMNAME    = "test-stream-delete-me"
	TESTMEMORYSTORAGE = "memory"
)

var (
	testSubjects = []string{"test-subject-delete-me1","test-subject-delete-me2"}
)

func TestValidateStreamName(tPtr *testing.T) {
	if soteErr := validateStreamName(""); soteErr.ErrCode != 206300 {
		tPtr.Errorf("TestValidateStreamName Failed: Expected error code of 206300")
	}

	if soteErr := validateStreamName(""); soteErr.ErrCode != 200513 {
		tPtr.Errorf("TestValidateStreamName Failed: Expected error code of 200513")
	}

	if soteErr := validateStreamName(TESTSTREAMNAME); soteErr.ErrCode != nil {
		tPtr.Errorf("TestValidateStreamName Failed: Expected error code to be nil")
	}
}
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
func TestDeleteStream(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		mmPtr   *MessageManager
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode != nil {
		tPtr.Errorf("TestDeleteStream failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}

	if soteErr := mmPtr.DeleteStream(""); soteErr.ErrCode != 206300 {
		tPtr.Errorf("TestDeleteStream Failed: Expected error code of 206300")
	}
}
func TestCreateLimitsStream(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		mmPtr   *MessageManager
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStream failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}

	if _, soteErr = mmPtr.CreateLimitsStream(TESTSTREAMNAME, TESTMEMORYSTORAGE, testSubjects, 1); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStream Failed: Expected error code to be nil")
	}

	if soteErr := mmPtr.DeleteStream(TESTSTREAMNAME); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStream Failed: Expected error code to be nil")
	}
}
func TestCreateWorkQueueStream(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		mmPtr   *MessageManager
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStream failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}

	if _, soteErr = mmPtr.CreateWorkQueueStream(TESTSTREAMNAME, TESTMEMORYSTORAGE, testSubjects, 1); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStream Failed: Expected error code to be nil")
	}

	if soteErr := mmPtr.DeleteStream(TESTSTREAMNAME); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStream Failed: Expected error code to be nil")
	}
}
