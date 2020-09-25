package sMessage

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/jsm.go/api"
	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2020/sError"
)

const (
	STREAMNAME          = "TEST_STREAM_NAME"
	STREAMSUBJECT       = "TEST_STREAM_NAME"
	STREAMCREDITIALFILE = "./test_fake.creds"
	STREAMFAKETOKEN     = "eyJ0eXAiFAKEd3QiLCJhbGciOiJlZDI1NTE5In0.eyJqdGkiOiIyCREDSkQzWjZUUUxFQ1FKUkNSNEEzS1M2SzZRM1hOWDcyTTJLRUFUQjRJFAKEUVJFVVpRIiwiaWF0IjoxNTk4OTE1Njc5LCJpc3MiOiJBQVZOVzU0UFJCREDSEUzWEM2SEVIQUNTTEM2RkRQWFIyV1g2M1NGTjJVS0k1R0VLUE9WMktPTiIsIm5hbWUiOiJiZSIsInN1YiI6IlVER0VZWUVJWUFMTzdTUFpPWlFPRDVDUUlLWVZXRlNMNFVQRUFTNDc0NkpTRkFSMk5QS0hUSEdPIiwidHlwZSI6InVzZXIiLCJuYXRzIjp7InB1YiI6e30sInN1YiI6e319fQ.qKhcZTO8P1ZxM28B5fZ-5NOrKB3GUv_60jVMohzR2p1PqJR5rmXLA22IpOYnApJXmkI8z1UJNJ7CHj6wpottCg\n"
	CONSUMERDURABLENAME = "TEST_CONSUMER_NAME"
	SUBJECTFILETER      = "TEST_STREAM_NAME.TEST1"
	NONEXISTINGFILE     = "./file_does_not.exist"
	MAXRECONNECT        = 5
	RECONNECTWAIT       = 1 * time.Second
	SYNADIAURL          = "euwest1.aws.ngs.global"
	TIMEOUT             = 1 * time.Second
	BADURL              = "google.com"
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
func TestConnect(t *testing.T) {
	var (
		opts    []nats.Option
		soteErr sError.SoteError
	)

	if opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds"); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
	} else {
		if _, soteErr = Connect(SYNADIAURL, opts); soteErr.ErrCode != nil {
			t.Errorf("Connect Failed: Expected error code to be nil")
		}
	}

	if _, soteErr = Connect(SYNADIAURL, nil); soteErr.ErrCode != 200513 {
		t.Errorf("Connect Failed: Expected error code of 200513")
	}

	if _, soteErr = Connect(BADURL, opts); soteErr.ErrCode != 603999 {
		t.Errorf("Connect Failed: Expected error code of 603999")
	}
}
func TestValidateStream(t *testing.T) {
	var (
		opts    []nats.Option
		soteErr sError.SoteError
		nc      *nats.Conn
	)

	if soteErr = validateStream(nil); soteErr.ErrCode != 200513 {
		t.Errorf("validateStream Failed: Expected error code of 200513")
	}

	if opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds"); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
	} else {
		if nc, soteErr = Connect(SYNADIAURL, opts); soteErr.ErrCode != nil {
			t.Errorf("Connect Failed: Expected error code to be nil")
		}
	}

	if soteErr = validateStreamParams(STREAMNAME, STREAMSUBJECT, nc); soteErr.ErrCode != nil {
		t.Errorf("validateStreamParams Failed: Expected error code to be nil")
	}
}
func TestValidateStreamName(t *testing.T) {
	var (
		soteErr sError.SoteError
	)

	if soteErr = validateStreamName(""); soteErr.ErrCode != 200513 {
		t.Errorf("validateStreamName Failed: Expected error code of 200513")
	}

	if soteErr = validateStreamName(STREAMNAME); soteErr.ErrCode != nil {
		t.Errorf("validateStreamName Failed: Expected error code to be nil")
	}
}
func TestValidateConnection(t *testing.T) {
	var (
		opts    []nats.Option
		soteErr sError.SoteError
		nc      *nats.Conn
	)

	if soteErr = validateConnection(nil); soteErr.ErrCode != 200513 {
		t.Errorf("validateConnection Failed: Expected error code of 200513")
	}

	if opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds"); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
	} else {
		if nc, soteErr = Connect(SYNADIAURL, opts); soteErr.ErrCode != nil {
			t.Errorf("Connect Failed: Expected error code to be nil")
		}
	}

	if soteErr = validateConnection(nc); soteErr.ErrCode != nil {
		t.Errorf("validateConnection Failed: Expected error code to be nil")
	}
}
func TestValidateStreamParams(t *testing.T) {
	var (
		opts    []nats.Option
		soteErr sError.SoteError
		nc      *nats.Conn
	)

	if opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds"); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
	} else {
		if nc, soteErr = Connect(SYNADIAURL, opts); soteErr.ErrCode != nil {
			t.Errorf("Connect Failed: Expected error code to be nil")
		}
	}

	if soteErr = validateStreamParams(STREAMNAME, STREAMSUBJECT, nc); soteErr.ErrCode != nil {
		t.Errorf("validateStreamParams Failed: Expected error code to be nil")
	}
	if soteErr = validateStreamParams("", STREAMSUBJECT, nc); soteErr.ErrCode != 200513 {
		t.Errorf("validateStreamParams Failed: Expected error code of 200513")
	}
	if soteErr = validateStreamParams(STREAMNAME, "", nc); soteErr.ErrCode != 200513 {
		t.Errorf("validateStreamParams Failed: Expected error code of 200513")
	}
	if soteErr = validateStreamParams(STREAMNAME, STREAMSUBJECT, nil); soteErr.ErrCode != 200513 {
		t.Errorf("validateStreamParams Failed: Expected error code of 200513")
	}
}
func TestSetReplicas(t *testing.T) {
	if replicas := setReplicas(-1); replicas != 1 {
		t.Errorf("setReplicas Failed: Expected replicase value to be 1")
	}
	if replicas := setReplicas(0); replicas != 1 {
		t.Errorf("setReplicas Failed: Expected replicase value to be 1")
	}
	if replicas := setReplicas(1); replicas != 1 {
		t.Errorf("setReplicas Failed: Expected replicase value to be 1")
	}
	if replicas := setReplicas(10); replicas != 10 {
		t.Errorf("setReplicas Failed: Expected replicase value to be 10")
	}
}
func TestDeleteStream(t *testing.T) {
	// Testing that Stream pointer is not nil
	if soteErr := DeleteStream(nil); soteErr.ErrCode != 335260 {
		t.Errorf("DeleteStream Failed: Expected error code of 335260")
	}
}
func TestCreateLimitsStreamDeleteStream(t *testing.T) {
	var (
		opts    []nats.Option
		soteErr sError.SoteError
		nc      *nats.Conn
		pStream *jsm.Stream
	)

	if opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds"); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
	} else {
		if nc, soteErr = Connect(SYNADIAURL, opts); soteErr.ErrCode != nil {
			t.Errorf("Connect Failed: Expected error code to be nil")
		}
	}

	// Test the default storage setting when creating a stream
	if pStream, soteErr = CreateLimitsStream(STREAMNAME, STREAMSUBJECT, "", 1, nc); soteErr.ErrCode != nil {
		t.Errorf("CreateLimitsStream Failed: Expected error code to be nil")
	}

	if soteErr := DeleteStream(pStream); soteErr.ErrCode != nil {
		t.Errorf("DeleteStream Failed: Expected error code to be nil")
	}
}
func TestCreateWorkStreamDeleteStream(t *testing.T) {
	var (
		opts    []nats.Option
		soteErr sError.SoteError
		nc      *nats.Conn
		pStream *jsm.Stream
	)

	if opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds"); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
	} else {
		if nc, soteErr = Connect(SYNADIAURL, opts); soteErr.ErrCode != nil {
			t.Errorf("Connect Failed: Expected error code to be nil")
		}
	}

	// Test the default storage setting when creating a stream
	if pStream, soteErr = CreateWorkStream(STREAMNAME, STREAMSUBJECT, "", 1, nc); soteErr.ErrCode != nil {
		t.Errorf("CreateWorkStream Failed: Expected error code to be nil")
	}

	if soteErr := DeleteStream(pStream); soteErr.ErrCode != nil {
		t.Errorf("DeleteStream Failed: Expected error code to be nil")
	}
}
func TestStreamInfo(t *testing.T) {
	var (
		opts    []nats.Option
		soteErr sError.SoteError
		nc      *nats.Conn
		pStream *jsm.Stream
		info    *api.StreamInfo
	)

	// Testing when stream doesn't exist
	if info, soteErr = StreamInfo(nil); soteErr.ErrCode != 200513 {
		t.Errorf("StreamInfo Failed: Expected error code of 200513")
	}

	if opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds"); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
	} else {
		if nc, soteErr = Connect(SYNADIAURL, opts); soteErr.ErrCode != nil {
			t.Errorf("Connect Failed: Expected error code to be nil")
		}
	}

	// Test the default storage setting when creating a stream
	if pStream, soteErr = CreateWorkStream(STREAMNAME, STREAMSUBJECT, "", 1, nc); soteErr.ErrCode != nil {
		t.Errorf("CreateWorkStream Failed: Expected error code to be nil")
	}

	if info, soteErr = StreamInfo(pStream); soteErr.ErrCode != nil {
		x, _ := json.Marshal(info)
		t.Errorf("Stream Info: " + string(x))
		t.Errorf("StreamInfo Failed: Expected error code to be nil")
	}
}
func TestPurgeStream(t *testing.T) {
	var (
		opts    []nats.Option
		soteErr sError.SoteError
		nc      *nats.Conn
		pStream *jsm.Stream
	)

	// Testing when stream doesn't exist
	if soteErr = PurgeStream(nil); soteErr.ErrCode != 200513 {
		t.Errorf("PurgeStream Failed: Expected error code of 200513")
	}

	if opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds"); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
	} else {
		if nc, soteErr = Connect(SYNADIAURL, opts); soteErr.ErrCode != nil {
			t.Errorf("Connect Failed: Expected error code to be nil")
		}
	}

	// Test the default storage setting when creating a stream
	if pStream, soteErr = CreateWorkStream(STREAMNAME, STREAMSUBJECT, "", 1, nc); soteErr.ErrCode != nil {
		t.Errorf("CreateWorkStream Failed: Expected error code to be nil")
	}

	if soteErr = PurgeStream(pStream); soteErr.ErrCode != nil {
		t.Errorf("PurgeStream Failed: Expected error code to be nil")
	}
}
func TestDeleteMessageFromStream(t *testing.T) {
	var (
		opts    []nats.Option
		soteErr sError.SoteError
		nc      *nats.Conn
		pStream *jsm.Stream
	)

	// Testing when stream doesn't exist
	if soteErr = DeleteMessageFromStream(nil, 0); soteErr.ErrCode != 400005 {
		t.Errorf("DeleteMessageFromStream Failed: Expected error code of 400005")
	}

	if opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds"); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
	} else {
		if nc, soteErr = Connect(SYNADIAURL, opts); soteErr.ErrCode != nil {
			t.Errorf("Connect Failed: Expected error code to be nil")
		}
	}

	// Test the default storage setting when creating a stream
	if pStream, soteErr = CreateWorkStream(STREAMNAME, STREAMSUBJECT, "", 1, nc); soteErr.ErrCode != nil {
		t.Errorf("CreateWorkStream Failed: Expected error code to be nil")
	}

	if soteErr = DeleteMessageFromStream(pStream, 1); soteErr.ErrCode != 109999 {
		t.Errorf("DeleteMessageFromStream Failed: Expected error code of 109999")
	}
}
// func TestCreateConsumer(t *testing.T) {
// 	var (
// 		opts    []nats.Option
// 		soteErr sError.SoteError
// 		nc      *nats.Conn
// 	)
//
// 	if opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds"); soteErr.ErrCode != nil {
// 		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
// 	} else {
// 		if nc, soteErr = Connect(SYNADIAURL, opts); soteErr.ErrCode != nil {
// 			t.Errorf("Connect Failed: Expected error code to be nil")
// 		}
// 	}
//
// 	_, soteErr = CreateConsumer(STREAMNAME, CONSUMERDURABLENAME, SUBJECTFILETER, nc)
// }
