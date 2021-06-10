package packages

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sMessage"
)

const (
	TESTSTREAMNAME         = "TEST_STREAM_NAME"
	TESTCONSUMERNAMEPULL       = "TEST_CONSUMER_NAME"
	TESTSYNADIAURL         = "west.eu.geo.ngs.global"
	TESTAPPLICATIONSYNADIA = "synadia"
)

var (
	testPullSubjects = []string{"test-subject-delete-me-1", "test-subject-delete-me-2"}
)

func TestNew(tPtr *testing.T) {
	if _, soteErr := sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "myConnection", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("New Failed: Expected error code to be nil")
	}

	if _, soteErr := sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "myConnection", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("New Failed: Expected error code of 200513")
	}

	if _, soteErr := sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "myConnection", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("New Failed: Expected error code of 209110")
	}

	if _, soteErr := sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "myConnection", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("New Failed: Expected error code of 210090")
	}

	if _, soteErr := sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "myConnection", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("New Failed: Expected error code to be nil")
	}

	if _, soteErr := sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "myConnection", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("New Failed: Expected error code to be nil")
	}

	if _, soteErr := sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "myConnection", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("New Failed: Expected error code to be nil")
	}

	if _, soteErr := sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "myConnection", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("New Failed: Expected error code to be nil")
	}
}
func TestPPublish(tPtr *testing.T) {
	var (
		mmPtr   *sMessage.MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testPullSubjects, 1, false); soteErr.ErrCode == nil {
			if _, soteErr = mmPtr.PPublish(testPullSubjects[0], "Hello world", false); soteErr.ErrCode != nil {
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
		mmPtr   *sMessage.MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testPullSubjects, 1, false); soteErr.ErrCode == nil {
			if _, soteErr = mmPtr.PPublish(testPullSubjects[0], "Hello world", false); soteErr.ErrCode == nil {
				if soteErr = mmPtr.PSubscribe(testPullSubjects[0], TESTCONSUMERNAMEPULL, nil, false); soteErr.ErrCode != 200513 {
					tPtr.Errorf("TestPSubscribe Failed: Expected error code to be 200513 got %v", soteErr.FmtErrMsg)
				}
				if soteErr = mmPtr.PSubscribe(testPullSubjects[0], TESTCONSUMERNAMEPULL, func(msgIn *nats.Msg) {
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
		mmPtr   *sMessage.MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testPullSubjects, 1, false); soteErr.ErrCode == nil {
			if _, soteErr = mmPtr.PPublish(testPullSubjects[0], "Hello world", false); soteErr.ErrCode == nil {
				if soteErr = mmPtr.PSubscribeSync(testPullSubjects[0], TESTCONSUMERNAMEPULL, false); soteErr.ErrCode != nil {
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
func TestPPullSubscribe(tPtr *testing.T) {
	var (
		mmPtr   *sMessage.MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testPullSubjects, 1, false); soteErr.ErrCode == nil {
			if soteErr = mmPtr.CreatePullReplayInstantConsumer(TESTSTREAMNAME, TESTCONSUMERNAMEPULL, testPullSubjects[0], 1, false); soteErr.ErrCode == nil {
				if soteErr = mmPtr.PullSubscribe(testPullSubjects[0], TESTCONSUMERNAMEPULL, false); soteErr.ErrCode != nil {
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
func TestPDeleteMsg(tPtr *testing.T) {
	var (
		mmPtr   *sMessage.MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testPullSubjects, 1, false); soteErr.ErrCode == nil {
			if _, soteErr = mmPtr.PPublish(testPullSubjects[0], "Hello world", false); soteErr.ErrCode == nil {
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
func TestPGetMsg(tPtr *testing.T) {
	var (
		mmPtr   *sMessage.MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testPullSubjects, 1, false); soteErr.ErrCode == nil {
			if _, soteErr = mmPtr.PPublish(testPullSubjects[0], "Hello world", false); soteErr.ErrCode == nil {
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
func TestPFetch(tPtr *testing.T) {
	var (
		mmPtr   *sMessage.MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testPullSubjects, 1, false); soteErr.ErrCode == nil {
			if soteErr = mmPtr.CreatePullReplayInstantConsumer(TESTSTREAMNAME, TESTCONSUMERNAMEPULL, testPullSubjects[0], 1, false); soteErr.ErrCode == nil {
				if _, soteErr = mmPtr.PPublish(testPullSubjects[0], "Hello world", false); soteErr.ErrCode != nil {
					tPtr.Errorf("TestPFetch Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
				if soteErr = mmPtr.PullSubscribe(testPullSubjects[0], TESTCONSUMERNAMEPULL, false); soteErr.ErrCode != nil {
					tPtr.Errorf("TestPSubscribeSync Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
				if soteErr = mmPtr.Fetch(TESTCONSUMERNAMEPULL, 1, true, false); soteErr.ErrCode != nil && soteErr.ErrCode != 101010 {
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
func TestPullReplayInstantConsumer(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		mmPtr   *sMessage.MessageManager
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestPullReplayInstantConsumer Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testPullSubjects, 1, false); soteErr.ErrCode == nil {
			if soteErr = mmPtr.CreatePullReplayInstantConsumer(TESTSTREAMNAME, TESTCONSUMERNAMEPULL, testPullSubjects[0], 1, false); soteErr.ErrCode != nil {
				tPtr.Errorf("TestPullReplayInstantConsumer Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
			}
		}
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestPullReplayInstantConsumer Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
	}
}
func TestGetConsumerInfo(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		mmPtr   *sMessage.MessageManager
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
			tPtr.Errorf("TestGetConsumerInfo Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testPullSubjects, 1, false); soteErr.ErrCode == nil {
			if soteErr = mmPtr.CreatePullReplayInstantConsumer(TESTSTREAMNAME, TESTCONSUMERNAMEPULL, testPullSubjects[0], 1, false); soteErr.ErrCode == nil {
				if _, soteErr = mmPtr.GetConsumerInfo(TESTSTREAMNAME, TESTCONSUMERNAMEPULL, false); soteErr.ErrCode != nil {
					tPtr.Errorf("TestGetConsumerInfo Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
				}
			}
			if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
				tPtr.Errorf("TestGetConsumerInfo Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
			}
		}
	}
}
func TestCreateLimitsStreamWithFileStorage(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		mmPtr   *sMessage.MessageManager
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStreamWithFileStorage failed: Expected soteErr to be nil got %v", soteErr.FmtErrMsg)
	}

	if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
		tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
	}

	if _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage(TESTSTREAMNAME, testPullSubjects, 1,false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}

	if soteErr := mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}
}
func TestCreateLimitsStreamWithMemoryStorage(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		mmPtr   *sMessage.MessageManager
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStreamWithMemoryStorage failed: Expected soteErr to be nil got %v", soteErr.FmtErrMsg)
	}

	if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
		tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
	}

	if _, soteErr = mmPtr.CreateLimitsStreamWithMemoryStorage(TESTSTREAMNAME, testPullSubjects, 1, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStreamWithMemoryStorage Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}

	if soteErr := mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateLimitsStreamWithMemoryStorage Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}
}
func TestCreateWorkQueueStreamWithFileStorage(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		mmPtr   *sMessage.MessageManager
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateWorkQueueStreamWithFileStorage failed: Expected soteErr to be nil got %v", soteErr.FmtErrMsg)
	}

	if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
		tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
	}

	if _, soteErr = mmPtr.CreateWorkQueueStreamWithFileStorage(TESTSTREAMNAME, testPullSubjects, 1, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateWorkQueueStreamWithFileStorage Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}

	if soteErr := mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateWorkQueueStreamWithFileStorage Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}
}
func TestCreateWorkQueueStreamWithMemoryStorage(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		mmPtr   *sMessage.MessageManager
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateWorkQueueStreamWithMemoryStorage failed: Expected soteErr to be nil got %v", soteErr.FmtErrMsg)
	}

	if soteErr = mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil && soteErr.ErrCode != 109999 {
		tPtr.Errorf("TestCreateLimitsStreamWithFileStorage Failed: Expected error code to be nil or 109999 got %v", soteErr.FmtErrMsg)
	}

	if _, soteErr = mmPtr.CreateWorkQueueStreamWithMemoryStorage(TESTSTREAMNAME, testPullSubjects, 1,false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateWorkQueueStreamWithMemoryStorage Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}

	if soteErr := mmPtr.DeleteStream(TESTSTREAMNAME, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestCreateWorkQueueStreamWithMemoryStorage Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}
}
func TestPublish(tPtr *testing.T) {
	var (
		mmPtr *sMessage.MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.Publish("greeting", "Hello world", false); soteErr.ErrCode != nil {
			tPtr.Errorf("TestPublish Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestSubscribe(tPtr *testing.T) {
	var (
		mmPtr *sMessage.MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if _, soteErr = mmPtr.Subscribe("greeting", false); soteErr.ErrCode != nil {
			tPtr.Errorf("TestSubscribe Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestPublishRequest(tPtr *testing.T) {
	var (
		mmPtr *sMessage.MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.PublishRequest("greeting", "greeting-reply", "Back At You!", false); soteErr.ErrCode != nil {
			tPtr.Errorf("TestPublishRequest Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestSubscribeSync(tPtr *testing.T) {
	var (
		mmPtr *sMessage.MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.SubscribeSync("greeting", "greeting-reply", false); soteErr.ErrCode != nil {
			tPtr.Errorf("TestSubscribeSync Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestNextMsg(tPtr *testing.T) {
	var (
		mmPtr *sMessage.MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if soteErr = mmPtr.SubscribeSync("greeting", "greeting-reply", false); soteErr.ErrCode != nil {
			tPtr.Errorf("TestSubscribeSync Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
		if soteErr = mmPtr.Publish("greeting-reply", "Hello world", false); soteErr.ErrCode != nil {
			tPtr.Errorf("TestPublish Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
		if _, soteErr = mmPtr.NextMsg("greeting", false); soteErr.ErrCode != nil {
			tPtr.Errorf("TestNextMsg Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr.Close()
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestRequest(tPtr *testing.T) {
	// TODO This code is not being tested at this time. How to test this needs to be investigated.
}
// We are not testing to see if NATS messaging works. We are only testing if the code works.
func TestRequestReply(tPtr *testing.T) {
	var (
		mmPtr *sMessage.MessageManager
		soteErr sError.SoteError
	)

	if mmPtr, soteErr = sMessage.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode == nil {
		if _, soteErr = mmPtr.RequestReply("greeting", "Hello World", false); soteErr.ErrCode != nil {
			tPtr.Errorf("TestRequestReply Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
		}
	}

	mmPtr.Close()
}
