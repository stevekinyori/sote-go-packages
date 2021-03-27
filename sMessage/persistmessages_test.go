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
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1); soteErr.ErrCode == nil {
			if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world"); soteErr.ErrCode != nil {
				tPtr.Errorf("TestPPublish Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
			}
		}
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME); soteErr.ErrCode != nil {
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
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1); soteErr.ErrCode == nil {
			if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world"); soteErr.ErrCode == nil {
				if soteErr = mmPtr.PSubscribe(testSubjects[0], "mySub", nil); soteErr.ErrCode != 200513 {
					tPtr.Errorf("TestPSubscribe Failed: Expected error code to be 206050 got %v", soteErr.FmtErrMsg)
				}
			}
		}
	}

	mmPtr.DeleteStream(TESTSTREAMNAME)

	mmPtr = mmPtr.Close()
}

// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPSubscribeSync(tPtr *testing.T) {
	var (
		mmPtr   *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1); soteErr.ErrCode == nil {
			if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world"); soteErr.ErrCode == nil {
				if soteErr = mmPtr.PSubscribeSync(testSubjects[0], "mySub"); soteErr.ErrCode != nil {
					tPtr.Errorf("TestPSubscribeSync Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
			}
		}
	}

	mmPtr.DeleteStream(TESTSTREAMNAME)

	mmPtr = mmPtr.Close()
}

// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPDeleteMsg(tPtr *testing.T) {
	var (
		mmPtr   *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1); soteErr.ErrCode == nil {
			if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world"); soteErr.ErrCode == nil {
				if soteErr = mmPtr.PDeleteMsg(TESTSTREAMNAME, 1); soteErr.ErrCode != nil {
					tPtr.Errorf("TestPSubscribeSync Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
			}
		}
	}

	mmPtr.DeleteStream(TESTSTREAMNAME)

	mmPtr = mmPtr.Close()
}

// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPGetMsg(tPtr *testing.T) {
	var (
		mmPtr   *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1); soteErr.ErrCode == nil {
			if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world"); soteErr.ErrCode == nil {
				if _, soteErr = mmPtr.PGetMsg(TESTSTREAMNAME, 1); soteErr.ErrCode != nil {
					tPtr.Errorf("TestPSubscribeSync Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
			}
		}
	}

	mmPtr.DeleteStream(TESTSTREAMNAME)

	mmPtr = mmPtr.Close()
}

// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPPullSubscribe(tPtr *testing.T) {
	var (
		mmPtr   *MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if _, soteErr = mmPtr.CreateWorkQueueStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1); soteErr.ErrCode == nil {
			if soteErr = mmPtr.PullReplayInstantConsumer(TESTSTREAMNAME, TESTCONSUMERNAME, testSubjects[0], 1); soteErr.ErrCode == nil {
				if soteErr = mmPtr.PPullSubscribe(testSubjects[0], "mySub"); soteErr.ErrCode != nil {
					tPtr.Errorf("TestPSubscribeSync Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
			}
		}
	}

	mmPtr.DeleteStream(TESTSTREAMNAME)

	mmPtr = mmPtr.Close()
}

// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPFetch(tPtr *testing.T) {
	var (
		mmPtr   *MessageManager
		soteErr sError.SoteError
		err     error
	)

	if mmPtr, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode == nil {
		if _, soteErr = mmPtr.CreateWorkQueueStreamWithFileStorage(TESTSTREAMNAME, testSubjects, 1); soteErr.ErrCode == nil {
			if soteErr = mmPtr.PSubscribeSync(testSubjects[0], "mySub"); soteErr.ErrCode != nil {
				tPtr.Errorf("TestPFetch Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
			}
			if _, soteErr = mmPtr.PPublish(testSubjects[0], "Hello world"); soteErr.ErrCode != nil {
				tPtr.Errorf("TestPFetch Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
			}
			if _, err = mmPtr.Subscriptions["mySub"].Fetch(1); err != nil {
				tPtr.Errorf("TestPFetch Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
			}
		}
	}

	mmPtr = mmPtr.Close()
}
