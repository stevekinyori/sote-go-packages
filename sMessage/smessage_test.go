package sMessage

import (
	"testing"
	"time"
)

const (
	STREAMNAME = "TEST_STREAM_NAME"
	STREAMCREDITIALFILE = "./test_fake.creds"
	NONEXISTINGFILE = "./file_does_not.exist"
	STREAMFAKETOKEN = "eyJ0eXAiFAKEd3QiLCJhbGciOiJlZDI1NTE5In0.eyJqdGkiOiIyCREDSkQzWjZUUUxFQ1FKUkNSNEEzS1M2SzZRM1hOWDcyTTJLRUFUQjRJFAKEUVJFVVpRIiwiaWF0IjoxNTk4OTE1Njc5LCJpc3MiOiJBQVZOVzU0UFJCREDSEUzWEM2SEVIQUNTTEM2RkRQWFIyV1g2M1NGTjJVS0k1R0VLUE9WMktPTiIsIm5hbWUiOiJiZSIsInN1YiI6IlVER0VZWUVJWUFMTzdTUFpPWlFPRDVDUUlLWVZXRlNMNFVQRUFTNDc0NkpTRkFSMk5QS0hUSEdPIiwidHlwZSI6InVzZXIiLCJuYXRzIjp7InB1YiI6e30sInN1YiI6e319fQ.qKhcZTO8P1ZxM28B5fZ-5NOrKB3GUv_60jVMohzR2p1PqJR5rmXLA22IpOYnApJXmkI8z1UJNJ7CHj6wpottCg\n"
	MAXRECONNECT = 5
	RECONNECTWAIT = 1 * time.Second
	TIMEOUT = 1 * time.Second
)

func TestSetStreamName(t *testing.T) {
	if _, soteErr := SetStreamName(STREAMNAME); soteErr.ErrCode != nil {
		t.Errorf("SetStreamName Failed: Expected error code to be nil")
	}

	if _, soteErr := SetStreamName(""); soteErr.ErrCode != 200513 {
		t.Errorf("SetStreamName Failed: Expected error code of 200513")
	}
}
func TestSetCredentials(t *testing.T) {
	if _, soteErr := SetCredentialsFile(STREAMCREDITIALFILE); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
	}

	if _, soteErr := SetCredentialsFile(""); soteErr.ErrCode != 200513 {
		t.Errorf("SetCredentialsFile Failed: Expected error code of 200513")
	}

	if _, soteErr := SetCredentialsFile(NONEXISTINGFILE); soteErr.ErrCode != 600010 {
		t.Errorf("SetCredentialsFile Failed: Expected error code of 600010")
	}

	if _, soteErr := SetCredentialsToken(STREAMFAKETOKEN); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsToken Failed: Expected error code to be nil")
	}

	if _, soteErr := SetCredentialsToken(""); soteErr.ErrCode != 200513 {
		t.Errorf("SetCredentialsToken Failed: Expected error code of 200513")
	}
}
func TestSetReconnectOptions(t *testing.T) {
	if _, soteErr := SetReconnectOptions(MAXRECONNECT, RECONNECTWAIT); soteErr.ErrCode != nil {
		t.Errorf("SetReconnectOptions Failed: Expected error code to be nil")
	}

	if _, soteErr := SetReconnectOptions(0, 0); soteErr.ErrCode != 200512 {
		t.Errorf("SetReconnectOptions Failed: Expected error code of 200512")
	}

	if _, soteErr := SetReconnectOptions(0, RECONNECTWAIT); soteErr.ErrCode != 200513 {
		t.Errorf("SetReconnectOptions Failed: Expected error code of 200513")
	}

	if _, soteErr := SetReconnectOptions(MAXRECONNECT, 0); soteErr.ErrCode != 200513 {
		t.Errorf("SetReconnectOptions Failed: Expected error code of 200513")
	}
}
func TestSetTimeOut(t *testing.T) {
	if _, soteErr := SetTimeOut(TIMEOUT); soteErr.ErrCode != nil {
		t.Errorf("SetTimeOut Failed: Expected error code to be nil")
	}

	if _, soteErr := SetTimeOut(0); soteErr.ErrCode != 200513 {
		t.Errorf("SetTimeOut Failed: Expected error code of 200513")
	}
}
func TestSetAllOptions(t *testing.T) {
	if _, soteErr := SetAllOptions(STREAMNAME, STREAMCREDITIALFILE, STREAMFAKETOKEN, MAXRECONNECT, RECONNECTWAIT, TIMEOUT); soteErr.ErrCode != nil {
		t.Errorf("SetAllOptions Failed: Expected error code to be nil")
	}

	if _, soteErr := SetAllOptions("", STREAMCREDITIALFILE, STREAMFAKETOKEN, MAXRECONNECT, RECONNECTWAIT, TIMEOUT); soteErr.ErrCode != 200513 {
		t.Errorf("SetAllOptions Failed: Expected error code of 200513")
	}

	if _, soteErr := SetAllOptions(STREAMNAME, "", STREAMFAKETOKEN, MAXRECONNECT, RECONNECTWAIT, TIMEOUT); soteErr.ErrCode != 200513 {
		t.Errorf("SetAllOptions Failed: Expected error code of 200513")
	}

	if _, soteErr := SetAllOptions(STREAMNAME, STREAMCREDITIALFILE, "", MAXRECONNECT, RECONNECTWAIT, TIMEOUT); soteErr.ErrCode != 200513 {
		t.Errorf("SetAllOptions Failed: Expected error code of 200513")
	}

	if _, soteErr := SetAllOptions(STREAMNAME, STREAMCREDITIALFILE, STREAMFAKETOKEN, 0, RECONNECTWAIT, TIMEOUT); soteErr.ErrCode != 200513 {
		t.Errorf("SetAllOptions Failed: Expected error code of 200513")
	}

	if _, soteErr := SetAllOptions(STREAMNAME, STREAMCREDITIALFILE, STREAMFAKETOKEN, MAXRECONNECT, 0, TIMEOUT); soteErr.ErrCode != 200513 {
		t.Errorf("SetAllOptions Failed: Expected error code of 200513")
	}

	if _, soteErr := SetAllOptions(STREAMNAME, STREAMCREDITIALFILE, STREAMFAKETOKEN, MAXRECONNECT, RECONNECTWAIT, 0); soteErr.ErrCode != 200513 {
		t.Errorf("SetAllOptions Failed: Expected error code of 200513")
	}
}