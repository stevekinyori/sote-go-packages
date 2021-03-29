package sMessage

import (
	"testing"
	"time"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
)

const (
	TESTCONSUMERNAME = "test-consumer-delete-me"
)

func TestSetMaxDeliver(tPtr *testing.T) {
	if maxDeliver := setMaxDeliver(-1); maxDeliver != 1 {
		tPtr.Errorf("setMaxDeliver Failed: Expected maxDeliver value to be 1")
	}
	if maxDeliver := setMaxDeliver(0); maxDeliver != 1 {
		tPtr.Errorf("setMaxDeliver Failed: Expected maxDeliver value to be 1")
	}
	if maxDeliver := setMaxDeliver(1); maxDeliver != 1 {
		tPtr.Errorf("setMaxDeliver Failed: Expected maxDeliver value to be 1")
	}
	if maxDeliver := setMaxDeliver(11); maxDeliver != 3 {
		tPtr.Errorf("setMaxDeliver Failed: Expected maxDeliver value to be 10")
	}
}
func TestPullReplayInstantConsumer(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		mmPtr   *MessageManager
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1); soteErr.ErrCode == nil {
			if soteErr = mmPtr.PullReplayInstantConsumer(TESTSTREAMNAME, TESTCONSUMERNAME, testSubjects[0], 1); soteErr.ErrCode != nil {
				tPtr.Errorf("TestPullReplayInstantConsumer Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
			}
		}
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
	}
}
