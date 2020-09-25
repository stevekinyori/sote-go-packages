package sMessage
import (
	"encoding/json"
	"testing"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/jsm.go/api"
	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2020/sError"
)

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
	if soteErr := DeleteStream(nil); soteErr.ErrCode != 200513 {
		t.Errorf("DeleteStream Failed: Expected error code of 200513")
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
