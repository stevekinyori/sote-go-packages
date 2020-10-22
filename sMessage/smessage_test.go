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
}

// func TestSetStreamName(t *testing.T) {
// 	if _, soteErr := SetStreamName(TESTSTREAMNAME); soteErr.ErrCode != nil {
// 		t.Errorf("SetStreamName Failed: Expected error code to be nil")
// 	}
//
// 	if _, soteErr := SetStreamName(""); soteErr.ErrCode != 200513 {
// 		t.Errorf("SetStreamName Failed: Expected error code of 200513")
// 	}
// }
// func TestSetCredentials(t *testing.T) {
// 	if _, soteErr := SetCredentialsFile(TESTSTREAMCREDITIALFILE); soteErr.ErrCode != nil {
// 		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
// 	}
//
// 	if _, soteErr := SetCredentialsFile(""); soteErr.ErrCode != 200513 {
// 		t.Errorf("SetCredentialsFile Failed: Expected error code of 200513")
// 	}
//
// 	if _, soteErr := SetCredentialsFile(TESTNONEXISTINGFILE); soteErr.ErrCode != 600010 {
// 		t.Errorf("SetCredentialsFile Failed: Expected error code of 600010")
// 	}
// }
// func TestSetReconnectOptions(t *testing.T) {
// 	if _, soteErr := SetReconnectOptions(TESTMAXRECONNECT, TESTRECONNECTWAIT); soteErr.ErrCode != nil {
// 		t.Errorf("SetReconnectOptions Failed: Expected error code to be nil")
// 	}
//
// 	if _, soteErr := SetReconnectOptions(0, 0); soteErr.ErrCode != 200512 {
// 		t.Errorf("SetReconnectOptions Failed: Expected error code of 200512")
// 	}
//
// 	if _, soteErr := SetReconnectOptions(0, TESTRECONNECTWAIT); soteErr.ErrCode != 200513 {
// 		t.Errorf("SetReconnectOptions Failed: Expected error code of 200513")
// 	}
//
// 	if _, soteErr := SetReconnectOptions(TESTMAXRECONNECT, 0); soteErr.ErrCode != 200513 {
// 		t.Errorf("SetReconnectOptions Failed: Expected error code of 200513")
// 	}
// }
// func TestSetTimeOut(t *testing.T) {
// 	if _, soteErr := SetTimeOut(TESTTIMEOUT); soteErr.ErrCode != nil {
// 		t.Errorf("SetTimeOut Failed: Expected error code to be nil")
// 	}
//
// 	if _, soteErr := SetTimeOut(0); soteErr.ErrCode != 200513 {
// 		t.Errorf("SetTimeOut Failed: Expected error code of 200513")
// 	}
// }
// func TestSetAllOptions(t *testing.T) {
// 	if _, soteErr := SetAllOptions(TESTSTREAMNAME, TESTSTREAMCREDITIALFILE, TESTMAXRECONNECT, TESTRECONNECTWAIT, TESTTIMEOUT); soteErr.ErrCode != nil {
// 		t.Errorf("SetAllOptions Failed: Expected error code to be nil")
// 	}
//
// 	if _, soteErr := SetAllOptions("", TESTSTREAMCREDITIALFILE, TESTMAXRECONNECT, TESTRECONNECTWAIT, TESTTIMEOUT); soteErr.ErrCode != 200513 {
// 		t.Errorf("SetAllOptions Failed: Expected error code of 200513")
// 	}
//
// 	if _, soteErr := SetAllOptions(TESTSTREAMNAME, "", TESTMAXRECONNECT, TESTRECONNECTWAIT, TESTTIMEOUT); soteErr.ErrCode != nil {
// 		t.Errorf("SetAllOptions Failed: Expected error code to be nil")
// 	}
//
// 	if _, soteErr := SetAllOptions(TESTSTREAMNAME, TESTSTREAMCREDITIALFILE, 0, TESTRECONNECTWAIT, TESTTIMEOUT); soteErr.ErrCode != 200513 {
// 		t.Errorf("SetAllOptions Failed: Expected error code of 200513")
// 	}
//
// 	if _, soteErr := SetAllOptions(TESTSTREAMNAME, TESTSTREAMCREDITIALFILE, TESTMAXRECONNECT, 0, TESTTIMEOUT); soteErr.ErrCode != 200513 {
// 		t.Errorf("SetAllOptions Failed: Expected error code of 200513")
// 	}
//
// 	if _, soteErr := SetAllOptions(TESTSTREAMNAME, TESTSTREAMCREDITIALFILE, TESTMAXRECONNECT, TESTRECONNECTWAIT, 0); soteErr.ErrCode != 200513 {
// 		t.Errorf("SetAllOptions Failed: Expected error code of 200513")
// 	}
// }
// func TestConnect(t *testing.T) {
// 	var (
// 		opts    []nats.Option
// 		soteErr sError.SoteError
// 	)
//
// 	if opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds"); soteErr.ErrCode != nil {
// 		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
// 	} else {
// 		if _, soteErr = Connect(TESTSYNADIAURL, opts); soteErr.ErrCode != nil {
// 			t.Errorf("Connect Failed: Expected error code to be nil")
// 		}
// 	}
//
// 	if _, soteErr = Connect(TESTSYNADIAURL, nil); soteErr.ErrCode != 200513 {
// 		t.Errorf("Connect Failed: Expected error code of 200513")
// 	}
//
// 	if _, soteErr = Connect(TESTBADURL, opts); soteErr.ErrCode != 603999 {
// 		t.Errorf("Connect Failed: Expected error code of 603999")
// 	}
// 	if _, soteErr = Connect("", opts); soteErr.ErrCode != 200513 {
// 		t.Errorf("Connect Failed: Expected error code of 200513")
// 	}
//
// 	if _, soteErr = Connect(TESTSYNADIAURL, opts); soteErr.ErrCode != nil {
// 		t.Errorf("Connect Failed: Expected error code to be nil")
// 	}
// }
// func TestGetManager(t *testing.T) {
// 	var (
// 		nc      *nats.Conn
// 		opts    []nats.Option
// 		soteErr sError.SoteError
// 	)
//
// 	if opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds"); soteErr.ErrCode != nil {
// 		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
// 	} else {
// 		if nc, soteErr = Connect(TESTSYNADIAURL, opts); soteErr.ErrCode != nil {
// 			t.Errorf("Connect Failed: Expected error code to be nil")
// 		}
// 	}
//
// 	if _, soteErr = GetJSMManager(nc); soteErr.ErrCode != nil {
// 		t.Errorf("GetManager Failed: Expected error code to be nil")
// 	}
// }
// func TestGetJSMManagerWithConnOptions(t *testing.T) {
// 	var (
// 		opts    []nats.Option
// 		soteErr sError.SoteError
// 	)
//
// 	opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds")
//
// 	if soteErr.ErrCode == nil {
// 		if _, soteErr = GetJSMManagerWithConnOptions(TESTSYNADIAURL, opts); soteErr.ErrCode != nil {
// 			t.Errorf("GetManager Failed: Expected error code to be nil")
// 		}
// 	}
// }
