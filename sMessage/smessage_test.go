package sMessage

import (
	"testing"
	"time"

	"gitlab.com/soteapps/packages/v2020/sConfigParams"
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

func TestNew(t *testing.T) {
	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds", TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("New Failed: Expected error code to be nil")
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL,1, 250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("New Failed: Expected error code to be nil")
	}

	if _, soteErr := New("", sConfigParams.STAGING, "", TESTSYNADIAURL,1, 250*time.Millisecond); soteErr.ErrCode != 200513 {
		t.Errorf("New Failed: Expected error code of 200513")
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, "", "", TESTSYNADIAURL,1, 250*time.Millisecond); soteErr.ErrCode != 601010 {
		t.Errorf("New Failed: Expected error code of 601010")
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", "",1, 250*time.Millisecond); soteErr.ErrCode != 609990 {
		t.Errorf("New Failed: Expected error code of 609990")
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL,0, 250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("New Failed: Expected error code to be nil")
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL,200, 250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("New Failed: Expected error code to be nil")
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL,1, 1*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("New Failed: Expected error code to be nil")
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL,1, 2*time.Minute); soteErr.ErrCode != nil {
		t.Errorf("New Failed: Expected error code to be nil")
	}
}
