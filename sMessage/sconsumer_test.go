package sMessage

import (
	"testing"

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
func TestValidateDurableName(t *testing.T) {
	if soteErr := validateDurableName(""); soteErr.ErrCode != 200513 {
		t.Errorf("validateDurableName Failed: Expected error code of 200513")
	}
	if soteErr := validateDurableName(TESTDURABLENAME); soteErr.ErrCode != nil {
		t.Errorf("validateDurableName Failed: Expected error code to be nil")
	}
}
func TestValidateDeliverySubject(t *testing.T) {
	if soteErr := validateDeliverySubject(""); soteErr.ErrCode != 200513 {
		t.Errorf("validateDeliverySubject Failed: Expected error code of 200513")
	}
	if soteErr := validateDeliverySubject(TESTDELIVERSUBJECT); soteErr.ErrCode != nil {
		t.Errorf("validateDeliverySubject Failed: Expected error code to be nil")
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
func TestValidateConsumerParams(t *testing.T) {
	var (
		opts    []nats.Option
		soteErr sError.SoteError
		nc      *nats.Conn
	)

	// Setup
	if opts, soteErr = SetCredentialsFile("/Users/syacko/.nkeys/creds/synadia/NATS_CONNECT/NATS_CONNECT.creds"); soteErr.ErrCode != nil {
		t.Errorf("SetCredentialsFile Failed: Expected error code to be nil")
	} else {
		if nc, soteErr = Connect(TESTSYNADIAURL, opts); soteErr.ErrCode != nil {
			t.Errorf("Connect Failed: Expected error code to be nil")
		}
	}

	if soteErr = validateConsumerParams(TESTSTREAMNAME, TESTDURABLENAME, TESTDELIVERSUBJECT, TESTSUBJECTFILETER, nc); soteErr.ErrCode != nil {
		t.Errorf("validateConsumerParams Failed: Expected error code to be nil")
	}
	if soteErr = validateConsumerParams("", TESTDURABLENAME, TESTDELIVERSUBJECT, TESTSUBJECTFILETER, nc); soteErr.ErrCode != 200513 {
		t.Errorf("validateConsumerParams Failed: Expected error code of 200513")
	}
	if soteErr = validateConsumerParams(TESTSTREAMNAME, "", TESTDELIVERSUBJECT, TESTSUBJECTFILETER, nc); soteErr.ErrCode != 200513 {
		t.Errorf("validateConsumerParams Failed: Expected error code of 200513")
	}
	if soteErr = validateConsumerParams(TESTSTREAMNAME, TESTDURABLENAME, "", TESTSUBJECTFILETER, nc); soteErr.ErrCode != 200513 {
		t.Errorf("validateConsumerParams Failed: Expected error code of 200513")
	}
	if soteErr = validateConsumerParams(TESTSTREAMNAME, TESTDURABLENAME, TESTDELIVERSUBJECT, "", nc); soteErr.ErrCode != 200513 {
		t.Errorf("validateConsumerParams Failed: Expected error code of 200513")
	}
	if soteErr = validateConsumerParams(TESTSTREAMNAME, TESTDURABLENAME, TESTDELIVERSUBJECT, TESTSUBJECTFILETER, nil); soteErr.ErrCode != 200513 {
		t.Errorf("validateConsumerParams Failed: Expected error code of 200513")
	}
}
func TestDeleteConsumer(t *testing.T) {
	// Testing when the Stream pointer is nil
	if soteErr := DeleteConsumer(nil); soteErr.ErrCode != 200513 {
		t.Errorf("DeleteStream Failed: Expected error code of 200513")
	}
}
