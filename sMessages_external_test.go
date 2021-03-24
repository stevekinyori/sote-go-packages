package packages

import (
	"testing"
	"time"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sMessage/sJetStream"
)

const (
	TESTSTREAMNAME            = "TEST_STREAM_NAME"
	TESTSTREAMSUBJECT         = "TEST_STREAM_NAME"
	TESTSTREAMSUBJECTWILDCARD = "TEST_STREAM_NAME.*"
	TESTSTREAMCREDITIALFILE   = "./test_fake.creds"
	//
	TESTCONSUMERNAME       = "TEST_CONSUMER_NAME"
	TESTDURABLENAME        = "TEST_DURABLE_NAME"
	TESTDELIVERSUBJECT     = "TEST_STREAM_NAME.TEST1"
	TESTDELIVERSUBJECTPULL = ""
	TESTSUBJECTFILETER     = "TEST_STREAM_NAME.TEST1"
	//
	TESTNONEXISTINGFILE = "./file_does_not.exist"
	TESTMAXRECONNECT    = 5
	TESTRECONNECTWAIT   = 1 * time.Second
	TESTSYNADIAURL      = "euwest1.aws.ngs.global"
	TESTTIMEOUT         = 1 * time.Second
	TESTBADURL          = "google.com"
	//
	TESTAPPLICATIONSYNADIA = "synadia"
)

func TestNew(tPtr *testing.T) {
	if _, soteErr := sJetStream.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL,1, 250*time.Millisecond); soteErr.ErrCode != nil {
		tPtr.Errorf("New Failed: Expected error code to be nil")
	}

	if _, soteErr := sJetStream.New("", sConfigParams.STAGING, "", TESTSYNADIAURL,1, 250*time.Millisecond); soteErr.ErrCode != 200513 {
		tPtr.Errorf("New Failed: Expected error code of 200513")
	}

	if _, soteErr := sJetStream.New(TESTAPPLICATIONSYNADIA, "", "", TESTSYNADIAURL,1, 250*time.Millisecond); soteErr.ErrCode != 209110 {
		tPtr.Errorf("New Failed: Expected error code of 209110")
	}

	if _, soteErr := sJetStream.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", "",1, 250*time.Millisecond); soteErr.ErrCode != 210090 {
		tPtr.Errorf("New Failed: Expected error code of 210090")
	}

	if _, soteErr := sJetStream.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL,0, 250*time.Millisecond); soteErr.ErrCode != nil {
		tPtr.Errorf("New Failed: Expected error code to be nil")
	}

	if _, soteErr := sJetStream.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL,200, 250*time.Millisecond); soteErr.ErrCode != nil {
		tPtr.Errorf("New Failed: Expected error code to be nil")
	}

	if _, soteErr := sJetStream.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL,1, 1*time.Millisecond); soteErr.ErrCode != nil {
		tPtr.Errorf("New Failed: Expected error code to be nil")
	}

	if _, soteErr := sJetStream.New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL,1, 2*time.Minute); soteErr.ErrCode != nil {
		tPtr.Errorf("New Failed: Expected error code to be nil")
	}
}
