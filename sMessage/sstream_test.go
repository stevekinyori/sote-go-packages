package sMessage

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/jsm.go/api"
	"gitlab.com/soteapps/packages/v2020/sConfigParams"
	"gitlab.com/soteapps/packages/v2020/sError"
)

func TestValidateStreamWhenNil(t *testing.T) {
	if soteErr := validateStream(nil); soteErr.ErrCode != 335260 {
		t.Errorf("validateStream Failed: Expected error code of 335260")
	}
}
func TestValidateStreamName(t *testing.T) {
	if soteErr := validateStreamName(""); soteErr.ErrCode != 200513 {
		t.Errorf("validateStreamName Failed: Expected error code of 200513")
	}

	if soteErr := validateStreamName(TESTSTREAMNAME); soteErr.ErrCode != nil {
		t.Errorf("validateStreamName Failed: Expected error code to be nil")
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
	if replicas := setReplicas(11); replicas != 10 {
		t.Errorf("setReplicas Failed: Expected replicase value to be 10")
	}
}
func TestValidateStreamParams(t *testing.T) {
	var (
		soteErr     sError.SoteError
		pJSMManager *JSMManager
	)

	// Setup
	if pJSMManager, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("New Failed: Expected error code to be nil")
	}

	// Stream Params
	if soteErr = validateStreamParams(TESTSTREAMNAME, TESTSTREAMSUBJECT, pJSMManager); soteErr.ErrCode != nil {
		t.Errorf("validateStreamParams Failed: Expected error code to be nil")
	}
	if soteErr = validateStreamParams("", TESTSTREAMSUBJECT, pJSMManager); soteErr.ErrCode != 200513 {
		t.Errorf("validateStreamParams Failed: Expected error code of 200513")
	}
	if soteErr = validateStreamParams(TESTSTREAMNAME, "", pJSMManager); soteErr.ErrCode != 200513 {
		t.Errorf("validateStreamParams Failed: Expected error code of 200513")
	}
	if soteErr = validateStreamParams(TESTSTREAMNAME, TESTSTREAMSUBJECT, nil); soteErr.ErrCode != 200513 {
		t.Errorf("validateStreamParams Failed: Expected error code of 200513")
	}
}
func TestDeleteStream(t *testing.T) {
	// Testing when the Stream pointer is nil
	if soteErr := DeleteStream(nil); soteErr.ErrCode != 335260 {
		t.Errorf("DeleteStream Failed: Expected error code of 335260")
	}
}
func TestCreateOrLoadLimitsStreamDeleteStream(t *testing.T) {
	var (
		soteErr     sError.SoteError
		pJSMManager *JSMManager
		pStream     *jsm.Stream
	)

	// Setup
	if pJSMManager, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("New Failed: Expected error code to be nil")
	}

	// Test the creation of the stream
	if _, soteErr = CreateOrLoadLimitsStream(TESTSTREAMNAME, TESTSTREAMSUBJECT, "", 1, pJSMManager); soteErr.ErrCode != nil {
		t.Errorf("CreateLimitsStream Failed: Expected error code to be nil")
	}

	// Test the loading of the stream
	if pStream, soteErr = CreateOrLoadLimitsStream(TESTSTREAMNAME, TESTSTREAMSUBJECT, "", 1, pJSMManager); soteErr.ErrCode != nil {
		t.Errorf("CreateLimitsStream Failed: Expected error code to be nil")
	}

	// Testing when the Stream pointer exists
	if soteErr := DeleteStream(pStream); soteErr.ErrCode != nil {
		t.Errorf("DeleteStream Failed: Expected error code to be nil")
	}
}
func TestCreateOrLoadWorkStreamDeleteStream(t *testing.T) {
	var (
		soteErr     sError.SoteError
		pJSMManager *JSMManager
		pStream     *jsm.Stream
	)

	// Setup
	if pJSMManager, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("New Failed: Expected error code to be nil")
	}

	// Test the creation of the stream
	if pStream, soteErr = CreateOrLoadWorkStream(TESTSTREAMNAME, TESTSTREAMSUBJECT, "", 1, pJSMManager); soteErr.ErrCode != nil {
		t.Errorf("CreateWorkStream Failed: Expected error code to be nil")
	}

	// Test the loading of the stream
	if pStream, soteErr = CreateOrLoadWorkStream(TESTSTREAMNAME, TESTSTREAMSUBJECT, "", 1, pJSMManager); soteErr.ErrCode != nil {
		t.Errorf("CreateWorkStream Failed: Expected error code to be nil")
	}

	// Clean up
	if soteErr := DeleteStream(pStream); soteErr.ErrCode != nil {
		t.Errorf("DeleteStream Failed: Expected error code to be nil")
	}
}
func TestStreamInfo(t *testing.T) {
	var (
		soteErr     sError.SoteError
		pJSMManager *JSMManager
		pStream     *jsm.Stream
		info        *api.StreamInfo
	)

	// Testing when stream doesn't exist
	if info, soteErr = StreamInfo(nil); soteErr.ErrCode != 335260 {
		t.Errorf("StreamInfo Failed: Expected error code of 335260")
	}

	// Setup
	if pJSMManager, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("New Failed: Expected error code to be nil")
	}

	// Test the default storage setting when creating a stream
	if pStream, soteErr = CreateOrLoadWorkStream(TESTSTREAMNAME, TESTSTREAMSUBJECT, "", 1, pJSMManager); soteErr.ErrCode != nil {
		t.Errorf("CreateWorkStream Failed: Expected error code to be nil")
	}

	if info, soteErr = StreamInfo(pStream); soteErr.ErrCode != nil {
		x, _ := json.Marshal(info)
		t.Errorf("Stream Info: " + string(x))
		t.Errorf("StreamInfo Failed: Expected error code to be nil")
	}

	// Clean up
	if soteErr := DeleteStream(pStream); soteErr.ErrCode != nil {
		t.Errorf("DeleteStream Failed: Expected error code to be nil")
	}
}
func TestPurgeStream(t *testing.T) {
	var (
		soteErr    sError.SoteError
		pJSMManager *JSMManager
		pStream    *jsm.Stream
	)

	// Testing when stream doesn't exist
	if soteErr = PurgeStream(nil); soteErr.ErrCode != 335260 {
		t.Errorf("PurgeStream Failed: Expected error code of 335260")
	}

	// Setup
	if pJSMManager, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("New Failed: Expected error code to be nil")
	}

	// Test the default storage setting when creating a stream
	if pStream, soteErr = CreateOrLoadWorkStream(TESTSTREAMNAME, TESTSTREAMSUBJECT, "", 1, pJSMManager); soteErr.ErrCode != nil {
		t.Errorf("CreateWorkStream Failed: Expected error code to be nil")
	}

	// Testing that the call works, not the jsm code
	if soteErr = PurgeStream(pStream); soteErr.ErrCode != nil {
		t.Errorf("PurgeStream Failed: Expected error code to be nil")
	}

	// Clean up
	if soteErr := DeleteStream(pStream); soteErr.ErrCode != nil {
		t.Errorf("DeleteStream Failed: Expected error code to be nil")
	}
}
func TestDeleteMessageFromStream(t *testing.T) {
	var (
		soteErr    sError.SoteError
		pJSMManager *JSMManager
		pStream    *jsm.Stream
	)

	// Testing when stream doesn't exist
	if soteErr = DeleteMessageFromStream(nil, 0); soteErr.ErrCode != 400005 {
		t.Errorf("DeleteMessageFromStream Failed: Expected error code of 400005")
	}

	// Setup
	if pJSMManager, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, 1, 250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("New Failed: Expected error code to be nil")
	}

	// Test the default storage setting when creating a stream
	if pStream, soteErr = CreateOrLoadWorkStream(TESTSTREAMNAME, TESTSTREAMSUBJECT, "", 1, pJSMManager); soteErr.ErrCode != nil {
		t.Errorf("CreateWorkStream Failed: Expected error code to be nil")
	}

	// Testing that the call works, not the jsm code
	if soteErr = DeleteMessageFromStream(pStream, 1); soteErr.ErrCode != 109999 {
		t.Errorf("DeleteMessageFromStream Failed: Expected error code of 109999")
	}

	// Clean up
	if soteErr := DeleteStream(pStream); soteErr.ErrCode != nil {
		t.Errorf("DeleteStream Failed: Expected error code to be nil")
	}
}
