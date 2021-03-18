package sDocument

import (
	"testing"
)

func TestNew(t *testing.T) {
	if _, soteErr := New(GetTessdataPrefix()); soteErr.ErrCode != nil {
		t.Errorf("New failed: Expected error code to be nil")
	}

	if _, soteErr := New(""); soteErr.ErrCode != 209100 {
		t.Errorf("TestNewExpectErrorCode209100 failed: Expected error code to be 209100")
	}
}
