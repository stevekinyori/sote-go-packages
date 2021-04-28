package sMessage

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
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
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1, false); soteErr.ErrCode == nil {
			if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world", false); soteErr.ErrCode != nil {
				tPtr.Errorf("TestPPublish Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
			}
		}
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPSubscribe(tPtr *testing.T) {
	var (
		mmPtr   *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1, false); soteErr.ErrCode == nil {
			if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world", false); soteErr.ErrCode == nil {
				if soteErr = mmPtr.PSubscribe(testSubjects[0], TESTCONSUMERNAME, nil, false); soteErr.ErrCode != 200513 {
					tPtr.Errorf("TestPSubscribe Failed: Expected error code to be 200513 got %v", soteErr.FmtErrMsg)
				}
				if soteErr = mmPtr.PSubscribe(testSubjects[0], TESTCONSUMERNAME, func(msgIn *nats.Msg) {
					return
				}, false); soteErr.ErrCode != nil {
					tPtr.Errorf("TestPSubscribe Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
			}
		}
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPSubscribeSync(tPtr *testing.T) {
	var (
		mmPtr   *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1, false); soteErr.ErrCode == nil {
			if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world", false); soteErr.ErrCode == nil {
				if soteErr = mmPtr.PSubscribeSync(testSubjects[0], TESTCONSUMERNAME, false); soteErr.ErrCode != nil {
					tPtr.Errorf("TestPSubscribe Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
			}
		}
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPullSubscribe(tPtr *testing.T) {
	var (
		mmPtr   *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1, false); soteErr.ErrCode == nil {
			if soteErr = mmPtr.CreatePullReplayInstantConsumer(TESTSTREAMNAME, TESTCONSUMERNAME, testSubjects[0], 1, false); soteErr.ErrCode == nil {
				if soteErr = mmPtr.PullSubscribe(testSubjects[0], TESTCONSUMERNAME, false); soteErr.ErrCode != nil {
					tPtr.Errorf("TestPSubscribeSync Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
			}
		}
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestDeleteMsg(tPtr *testing.T) {
	var (
		mmPtr   *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1, false); soteErr.ErrCode == nil {
			if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world", false); soteErr.ErrCode == nil {
				if soteErr = mmPtr.DeleteMsg(TESTSTREAMNAME, 1, false); soteErr.ErrCode != nil {
					tPtr.Errorf("TestPSubscribeSync Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
			}
		}
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestGetMsg(tPtr *testing.T) {
	var (
		mmPtr   *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1, false); soteErr.ErrCode == nil {
			if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world", false); soteErr.ErrCode == nil {
				if soteErr = mmPtr.GetMsg(TESTSTREAMNAME, 1, false); soteErr.ErrCode != nil {
					tPtr.Errorf("TestPSubscribeSync Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
			}
		}
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestFetch(tPtr *testing.T) {
	var (
		mmPtr   *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1, false); soteErr.ErrCode == nil {
			if soteErr = mmPtr.CreatePullReplayInstantConsumer(TESTSTREAMNAME, TESTCONSUMERNAME, testSubjects[0], 1, false); soteErr.ErrCode == nil {
				if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world", false); soteErr.ErrCode != nil {
					tPtr.Errorf("TestPFetch Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
				if soteErr = mmPtr.PullSubscribe(testSubjects[0], TESTCONSUMERNAME, false); soteErr.ErrCode != nil {
					tPtr.Errorf("TestPSubscribeSync Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
				if soteErr = mmPtr.Fetch(TESTCONSUMERNAME, 1, true, false); soteErr.ErrCode != nil && soteErr.ErrCode != 101010 {
					tPtr.Errorf("TestPFetch Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
			}
		}
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestAck(tPtr *testing.T) {
	var (
		mmPtr   *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1, false); soteErr.ErrCode == nil {
			if soteErr = mmPtr.CreatePullReplayInstantConsumer(TESTSTREAMNAME, TESTCONSUMERNAME, testSubjects[0], 1, false); soteErr.ErrCode == nil {
				if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world", false); soteErr.ErrCode != nil {
					tPtr.Errorf("TestPFetch Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
				if soteErr = mmPtr.PullSubscribe(testSubjects[0], TESTCONSUMERNAME, false); soteErr.ErrCode != nil {
					tPtr.Errorf("TestPSubscribeSync Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
				if soteErr = mmPtr.Fetch(TESTCONSUMERNAME, 1, false, false); soteErr.ErrCode != nil {
					tPtr.Errorf("TestPFetch Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
				time.Sleep(5 * time.Second)
				if soteErr = mmPtr.Ack(mmPtr.Messages[0], true); soteErr.ErrCode != nil {
					tPtr.Errorf("TestPFetch Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
			}
		}
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr.Close()
}
