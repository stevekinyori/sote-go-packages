package sMessage

import (
	"testing"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2020/sError"
)

func TestValidateConsumerWhenNil(t *testing.T) {
	var (
		soteErr sError.SoteError
	)

	if soteErr = validateConsumer(nil); soteErr.ErrCode != 200513 {
		t.Errorf("validateConsumer Failed: Expected error code of 200513")
	}
}
func TestValidateConsumerName(t *testing.T) {
	if soteErr := validateConsumerName(""); soteErr.ErrCode != 200513 {
		t.Errorf("validateDConsumerName Failed: Expected error code of 200513")
	}
	if soteErr := validateConsumerName(TESTDURABLENAME); soteErr.ErrCode != nil {
		t.Errorf("validateConsumerName Failed: Expected error code to be nil")
	}
}
func TestValidateDurableName(t *testing.T) {
	// Test missing durable name when a pull consumer
	if soteErr := validateDurableName("", TESTDELIVERSUBJECTPULL); soteErr.ErrCode != 200513 {
		t.Errorf("validateDurableName Failed: Expected error code of 200513")
	}
	// Test durable name when a pull consumer
	if soteErr := validateDurableName(TESTDURABLENAME, TESTDELIVERSUBJECTPULL); soteErr.ErrCode != nil {
		t.Errorf("validateDurableName Failed: Expected error code to be nil")
	}
	// Test optional durable name when a push consumer with deliver subject
	if soteErr := validateDurableName("", TESTDELIVERSUBJECT); soteErr.ErrCode != nil {
		t.Errorf("validateDurableName Failed: Expected error code to be nil")
	}
	// Test optional durable name when a push consumer without deliver subject
	if soteErr := validateDurableName(TESTDURABLENAME, TESTDELIVERSUBJECT); soteErr.ErrCode != nil {
		t.Errorf("validateDurableName Failed: Expected error code to be nil")
	}
}
func TestValidateDeliverySubject(t *testing.T) {
	if soteErr := validateDeliverySubject(TESTSUBJECTFILETER, TESTDELIVERSUBJECTPULL); soteErr.ErrCode != 200513 {
		t.Errorf("validateDeliverySubject Failed: Expected error code of 200513")
	}
	if soteErr := validateDeliverySubject("", TESTDELIVERSUBJECTPULL); soteErr.ErrCode != nil {
		t.Errorf("validateDeliverySubject Failed: Expected error code to be nil")
	}
	if soteErr := validateDeliverySubject(TESTSUBJECTFILETER, TESTDELIVERSUBJECT); soteErr.ErrCode != nil {
		t.Errorf("validateDeliverySubject Failed: Expected error code to be nil")
	}
	if soteErr := validateDeliverySubject("", TESTDELIVERSUBJECT); soteErr.ErrCode != 200513 {
		t.Errorf("validateDeliverySubject Failed: Expected error code to of 200513")
	}
}
func TestValidateSubjectFilter(t *testing.T) {
	if soteErr := validateSubjectFilter(""); soteErr.ErrCode != 200513 {
		t.Errorf("validateSubjectFilter Failed: Expected error code of 200513")
	}
	if soteErr := validateSubjectFilter(TESTSUBJECTFILETER); soteErr.ErrCode != nil {
		t.Errorf("validateSubjectFilter Failed: Expected error code to be nil")
	}
}
func TestSetMaxDeliver(t *testing.T) {
	if maxDeliver := setMaxDeliver(-1); maxDeliver != 1 {
		t.Errorf("setMaxDeliver Failed: Expected maxDeliver value to be 1")
	}
	if maxDeliver := setMaxDeliver(0); maxDeliver != 1 {
		t.Errorf("setMaxDeliver Failed: Expected maxDeliver value to be 1")
	}
	if maxDeliver := setMaxDeliver(1); maxDeliver != 1 {
		t.Errorf("setMaxDeliver Failed: Expected maxDeliver value to be 1")
	}
	if maxDeliver := setMaxDeliver(4); maxDeliver != 3 {
		t.Errorf("setMaxDeliver Failed: Expected maxDeliver value to be 10")
	}
}
func TestValidateConsumerParams(t *testing.T) {
	var (
		opts       []nats.Option
		soteErr    sError.SoteError
		jsmManager *jsm.Manager
	)

	// Setup
	if opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds"); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
	} else {
		if jsmManager, soteErr = GetJSMManagerWithConnOptions(TESTSYNADIAURL, opts); soteErr.ErrCode != nil {
			t.Errorf("Connect Failed: Expected error code to be nil")
		}
	}

	if soteErr = validateConsumerParams(TESTSTREAMNAME, TESTCONSUMERNAME, TESTDURABLENAME, TESTDELIVERSUBJECT, TESTSUBJECTFILETER, jsmManager); soteErr.ErrCode != nil {
		t.Errorf("validateConsumerParams Failed: Expected error code to be nil")
	}
	if soteErr = validateConsumerParams("", TESTCONSUMERNAME, TESTDURABLENAME, TESTDELIVERSUBJECT, TESTSUBJECTFILETER, jsmManager); soteErr.ErrCode != 200513 {
		t.Errorf("validateConsumerParams Failed: Expected error code of 200513")
	}
	if soteErr = validateConsumerParams(TESTSTREAMNAME, "", TESTDURABLENAME, TESTDELIVERSUBJECT, TESTSUBJECTFILETER, jsmManager); soteErr.ErrCode != 200513 {
		t.Errorf("validateConsumerParams Failed: Expected error code of 200513")
	}
	if soteErr = validateConsumerParams(TESTSTREAMNAME, TESTCONSUMERNAME, TESTDURABLENAME, "", TESTSUBJECTFILETER, jsmManager); soteErr.ErrCode != 200513 {
		t.Errorf("validateConsumerParams Failed: Expected error code of 200513")
	}
	if soteErr = validateConsumerParams(TESTSTREAMNAME, TESTCONSUMERNAME, TESTDURABLENAME, TESTDELIVERSUBJECT, "", jsmManager); soteErr.ErrCode != 200513 {
		t.Errorf("validateConsumerParams Failed: Expected error code of 200513")
	}
	if soteErr = validateConsumerParams(TESTSTREAMNAME, TESTCONSUMERNAME, TESTDURABLENAME, TESTDELIVERSUBJECT, TESTSUBJECTFILETER, nil); soteErr.ErrCode != 200513 {
		t.Errorf("validateConsumerParams Failed: Expected error code of 200513")
	}
}
func TestDeleteConsumer(t *testing.T) {
	// Testing when the Stream pointer is nil
	if soteErr := DeleteConsumer(nil); soteErr.ErrCode != 200513 {
		t.Errorf("DeleteStream Failed: Expected error code of 200513")
	}
}
func TestCreateDeliverAllReplayInstantConsumer(t *testing.T) {
	var (
		opts       []nats.Option
		soteErr    sError.SoteError
		jsmManager *jsm.Manager
		pStream    *jsm.Stream
	)

	// Setup
	if opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds"); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
	} else {
		if jsmManager, soteErr = GetJSMManagerWithConnOptions(TESTSYNADIAURL, opts); soteErr.ErrCode != nil {
			t.Errorf("Connect Failed: Expected error code to be nil")
		}
	}

	if pStream, soteErr = CreateOrLoadLimitsStream(TESTSTREAMNAME, TESTSTREAMSUBJECT, "", 1, jsmManager); soteErr.ErrCode != nil {
		t.Errorf("CreateLimitsStream Failed: Expected error code to be nil")
	}

	// Test error code 336100 - Consumer subject filter is not subset of stream subject
	if _, soteErr := CreatePullConsumerWithReplayInstant(TESTSTREAMNAME, TESTCONSUMERNAME, TESTDURABLENAME, TESTSUBJECTFILETER, 1, jsmManager); soteErr.ErrCode != 336100 {
		t.Errorf("CreateDeliverAllReplayInstantConsumer Failed: Expected error code of 336100")
	}
	// Clean Up
	if soteErr := DeleteStream(pStream); soteErr.ErrCode != nil {
		t.Errorf("DeleteStream Failed: Expected error code to be nil")
	}

	// Test when consumer subject filter is a subset of the stream subject
	if pStream, soteErr = CreateOrLoadLimitsStream(TESTSTREAMNAME, TESTSTREAMSUBJECTWILDCARD, "", 1, jsmManager); soteErr.ErrCode != nil {
		t.Errorf("CreateLimitsStream Failed: Expected error code to be nil")
	}

	// Test that the consumer is loaded without error
	if _, soteErr := CreatePullConsumerWithReplayInstant(TESTSTREAMNAME, TESTCONSUMERNAME, TESTDURABLENAME, TESTSUBJECTFILETER, 1, jsmManager); soteErr.ErrCode != nil {
		t.Errorf("CreateDeliverAllReplayInstantConsumer Failed: Expected error code to be nil")
	}

	// Clean Up
	if soteErr := DeleteStream(pStream); soteErr.ErrCode != nil {
		t.Errorf("DeleteStream Failed: Expected error code to be nil")
	}
}
func TestCreatePullConsumerWithReplayOriginal(t *testing.T) {
	var (
		opts       []nats.Option
		soteErr    sError.SoteError
		jsmManager *jsm.Manager
		pStream    *jsm.Stream
	)

	// Setup
	if opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds"); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
	} else {
		if jsmManager, soteErr = GetJSMManagerWithConnOptions(TESTSYNADIAURL, opts); soteErr.ErrCode != nil {
			t.Errorf("Connect Failed: Expected error code to be nil")
		}
	}

	if pStream, soteErr = CreateOrLoadLimitsStream(TESTSTREAMNAME, TESTSTREAMSUBJECT, "", 1, jsmManager); soteErr.ErrCode != nil {
		t.Errorf("CreateLimitsStream Failed: Expected error code to be nil")
	}

	// Test error code 336100 - Consumer subject filter is not subset of stream subject
	if _, soteErr := CreatePullConsumerWithReplayOriginal(TESTSTREAMNAME, TESTCONSUMERNAME, TESTDURABLENAME, TESTSUBJECTFILETER, 1, jsmManager); soteErr.ErrCode != 336100 {
		t.Errorf("CreateDeliverAllReplayInstantConsumer Failed: Expected error code of 336100")
	}
	// Clean Up
	if soteErr := DeleteStream(pStream); soteErr.ErrCode != nil {
		t.Errorf("DeleteStream Failed: Expected error code to be nil")
	}

	// Test when consumer subject filter is a subset of the stream subject
	if pStream, soteErr = CreateOrLoadLimitsStream(TESTSTREAMNAME, TESTSTREAMSUBJECTWILDCARD, "", 1, jsmManager); soteErr.ErrCode != nil {
		t.Errorf("CreateLimitsStream Failed: Expected error code to be nil")
	}

	// Test that the consumer is loaded without error
	if _, soteErr := CreatePullConsumerWithReplayOriginal(TESTSTREAMNAME, TESTCONSUMERNAME, TESTDURABLENAME, TESTSUBJECTFILETER, 1, jsmManager); soteErr.ErrCode != nil {
		t.Errorf("CreateDeliverAllReplayInstantConsumer Failed: Expected error code to be nil")
	}

	// Clean Up
	if soteErr := DeleteStream(pStream); soteErr.ErrCode != nil {
		t.Errorf("DeleteStream Failed: Expected error code to be nil")
	}
}
