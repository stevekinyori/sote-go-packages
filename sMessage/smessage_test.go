package sMessage

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2020/sError"
)

const (
	TESTSTREAMNAME            = "TEST_STREAM_NAME"
	TESTSTREAMSUBJECT         = "TEST_STREAM_NAME"
	TESTSTREAMSUBJECTWILDCARD = "TEST_STREAM_NAME.*"
	TESTSTREAMCREDITIALFILE   = "./test_fake.creds"
	TESTSTREAMFAKETOKEN       = "eyJ0eXAiFAKEd3QiLCJhbGciOiJlZDI1NTE5In0.eyJqdGkiOiIyCREDSkQzWjZUUUxFQ1FKUkNSNEEzS1M2SzZRM1hOWDcyTTJLRUFUQjRJFAKEUVJFVVpRIiwiaWF0IjoxNTk4OTE1Njc5LCJpc3MiOiJBQVZOVzU0UFJCREDSEUzWEM2SEVIQUNTTEM2RkRQWFIyV1g2M1NGTjJVS0k1R0VLUE9WMktPTiIsIm5hbWUiOiJiZSIsInN1YiI6IlVER0VZWUVJWUFMTzdTUFpPWlFPRDVDUUlLWVZXRlNMNFVQRUFTNDc0NkpTRkFSMk5QS0hUSEdPIiwidHlwZSI6InVzZXIiLCJuYXRzIjp7InB1YiI6e30sInN1YiI6e319fQ.qKhcZTO8P1ZxM28B5fZ-5NOrKB3GUv_60jVMohzR2p1PqJR5rmXLA22IpOYnApJXmkI8z1UJNJ7CHj6wpottCg\n"
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
)

func TestSetStreamName(t *testing.T) {
	if _, soteErr := SetStreamName(TESTSTREAMNAME); soteErr.ErrCode != nil {
		t.Errorf("SetStreamName Failed: Expected error code to be nil")
	}

	if _, soteErr := SetStreamName(""); soteErr.ErrCode != 200513 {
		t.Errorf("SetStreamName Failed: Expected error code of 200513")
	}
}
func TestSetCredentials(t *testing.T) {
	if _, soteErr := SetCredentialsFile(TESTSTREAMCREDITIALFILE); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
	}

	if _, soteErr := SetCredentialsFile(""); soteErr.ErrCode != 200513 {
		t.Errorf("SetCredentialsFile Failed: Expected error code of 200513")
	}

	if _, soteErr := SetCredentialsFile(TESTNONEXISTINGFILE); soteErr.ErrCode != 600010 {
		t.Errorf("SetCredentialsFile Failed: Expected error code of 600010")
	}

	if _, soteErr := SetCredentialsToken(TESTSTREAMFAKETOKEN); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsToken Failed: Expected error code to be nil")
	}

	if _, soteErr := SetCredentialsToken(""); soteErr.ErrCode != 200513 {
		t.Errorf("SetCredentialsToken Failed: Expected error code of 200513")
	}
}
func TestSetReconnectOptions(t *testing.T) {
	if _, soteErr := SetReconnectOptions(TESTMAXRECONNECT, TESTRECONNECTWAIT); soteErr.ErrCode != nil {
		t.Errorf("SetReconnectOptions Failed: Expected error code to be nil")
	}

	if _, soteErr := SetReconnectOptions(0, 0); soteErr.ErrCode != 200512 {
		t.Errorf("SetReconnectOptions Failed: Expected error code of 200512")
	}

	if _, soteErr := SetReconnectOptions(0, TESTRECONNECTWAIT); soteErr.ErrCode != 200513 {
		t.Errorf("SetReconnectOptions Failed: Expected error code of 200513")
	}

	if _, soteErr := SetReconnectOptions(TESTMAXRECONNECT, 0); soteErr.ErrCode != 200513 {
		t.Errorf("SetReconnectOptions Failed: Expected error code of 200513")
	}
}
func TestSetTimeOut(t *testing.T) {
	if _, soteErr := SetTimeOut(TESTTIMEOUT); soteErr.ErrCode != nil {
		t.Errorf("SetTimeOut Failed: Expected error code to be nil")
	}

	if _, soteErr := SetTimeOut(0); soteErr.ErrCode != 200513 {
		t.Errorf("SetTimeOut Failed: Expected error code of 200513")
	}
}
func TestSetAllOptions(t *testing.T) {
	if _, soteErr := SetAllOptions(TESTSTREAMNAME, TESTSTREAMCREDITIALFILE, TESTSTREAMFAKETOKEN, TESTMAXRECONNECT, TESTRECONNECTWAIT, TESTTIMEOUT); soteErr.ErrCode != nil {
		t.Errorf("SetAllOptions Failed: Expected error code to be nil")
	}

	if _, soteErr := SetAllOptions("", TESTSTREAMCREDITIALFILE, TESTSTREAMFAKETOKEN, TESTMAXRECONNECT, TESTRECONNECTWAIT, TESTTIMEOUT); soteErr.ErrCode != 200513 {
		t.Errorf("SetAllOptions Failed: Expected error code of 200513")
	}

	if _, soteErr := SetAllOptions(TESTSTREAMNAME, "", TESTSTREAMFAKETOKEN, TESTMAXRECONNECT, TESTRECONNECTWAIT, TESTTIMEOUT); soteErr.ErrCode != 200513 {
		t.Errorf("SetAllOptions Failed: Expected error code of 200513")
	}

	if _, soteErr := SetAllOptions(TESTSTREAMNAME, TESTSTREAMCREDITIALFILE, "", TESTMAXRECONNECT, TESTRECONNECTWAIT, TESTTIMEOUT); soteErr.ErrCode != 200513 {
		t.Errorf("SetAllOptions Failed: Expected error code of 200513")
	}

	if _, soteErr := SetAllOptions(TESTSTREAMNAME, TESTSTREAMCREDITIALFILE, TESTSTREAMFAKETOKEN, 0, TESTRECONNECTWAIT, TESTTIMEOUT); soteErr.ErrCode != 200513 {
		t.Errorf("SetAllOptions Failed: Expected error code of 200513")
	}

	if _, soteErr := SetAllOptions(TESTSTREAMNAME, TESTSTREAMCREDITIALFILE, TESTSTREAMFAKETOKEN, TESTMAXRECONNECT, 0, TESTTIMEOUT); soteErr.ErrCode != 200513 {
		t.Errorf("SetAllOptions Failed: Expected error code of 200513")
	}

	if _, soteErr := SetAllOptions(TESTSTREAMNAME, TESTSTREAMCREDITIALFILE, TESTSTREAMFAKETOKEN, TESTMAXRECONNECT, TESTRECONNECTWAIT, 0); soteErr.ErrCode != 200513 {
		t.Errorf("SetAllOptions Failed: Expected error code of 200513")
	}
}
func TestConnect(t *testing.T) {
	var (
		opts    []nats.Option
		soteErr sError.SoteError
	)

	if opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds"); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
	} else {
		if _, soteErr = Connect(TESTSYNADIAURL, opts); soteErr.ErrCode != nil {
			t.Errorf("Connect Failed: Expected error code to be nil")
		}
	}

	if _, soteErr = Connect(TESTSYNADIAURL, nil); soteErr.ErrCode != 200513 {
		t.Errorf("Connect Failed: Expected error code of 200513")
	}

	if _, soteErr = Connect(TESTBADURL, opts); soteErr.ErrCode != 603999 {
		t.Errorf("Connect Failed: Expected error code of 603999")
	}
	if _, soteErr = Connect("", opts); soteErr.ErrCode != 200513 {
		t.Errorf("Connect Failed: Expected error code of 200513")
	}
}
