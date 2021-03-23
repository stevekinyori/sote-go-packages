package sMessage

import (
	"testing"
	"time"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
)

const (
	TESTAPPLICATIONSYNADIA = "synadia"
	TESTSYNADIAURL         = "euwest1.aws.ngs.global"
)

func TestNew(t *testing.T) {
	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("TestNew failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", true, 1,
		250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("TestNew failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "/Users/syacko/.nkeys/creds/synadia/sote-staging/staging-soteadmin.creds",
		TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("TestNew failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "/Users/syacko/.nkeys/creds/synadia/sote-staging/staging-soteadmin.creds",
		TESTSYNADIAURL, "test", true, 1,
		250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("TestNew failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "/XXXX/syacko/.nkeys/creds/synadia/sote-staging/staging-soteadmin.creds",
		TESTSYNADIAURL, "test", true, 1, 250*time.Millisecond); soteErr.ErrCode != 209010 {
		t.Errorf("TestNew failed: Expected error code of 209010: %v", soteErr.FmtErrMsg)
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, "INVALID", "", TESTSYNADIAURL, "test", true, 1,
		250*time.Millisecond); soteErr.ErrCode != 209110 {
		t.Errorf("TestNew failed: Expected error code of 209110: %v", soteErr.FmtErrMsg)
	}

	if _, soteErr := New("XXXX", sConfigParams.STAGING, "", TESTSYNADIAURL, "test", true, 1,
		250*time.Millisecond); soteErr.ErrCode != 109999 {
		t.Errorf("TestNew failed: Expected error code of 109999: %v", soteErr.FmtErrMsg)
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", true, -1,
		250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("TestNew Failed: Expected error code to be nil")
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", true, 6,
		250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("TestNew Failed: Expected error code to be nil")
	}
}
func TestClose(t *testing.T) {
	var (
		nm      *MessageManager
		soteErr sError.SoteError
	)

	if nm, soteErr = New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", true, 1,
		250*time.Millisecond); soteErr.ErrCode != nil {
		t.Errorf("New Failed: Expected error code to be nil")
	}

	if nm = nm.Close(); nm != nil {
		t.Errorf("TestClose Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}
}
