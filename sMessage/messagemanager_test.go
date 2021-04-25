package sMessage

import (
	"os"
	"runtime"
	"testing"
	"time"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
)

const (
	TESTAPPLICATIONSYNADIA = "synadia"
	TESTSYNADIAURL         = "west.eu.geo.ngs.global"
)

func TestNew(tPtr *testing.T) {
	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestNew failed: Expected soteErr to be nil got %v", soteErr.FmtErrMsg)
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestNew failed: Expected soteErr to be nil got %v", soteErr.FmtErrMsg)
	}
	var homedir string

	if runtime.GOOS == "windows" {
		homedir = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if homedir == "" {
			homedir = os.Getenv("USERPROFILE")
		}
	} else {
		homedir = os.Getenv("HOME")
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, homedir+"/.nkeys/creds/synadia/sote-staging/staging-soteadmin.creds",
		TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestNew failed: Expected soteErr to be nil got %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, homedir+"/.nkeys/creds/synadia/sote-staging/staging-soteadmin."+
		"creds",
		TESTSYNADIAURL, "test", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestNew failed: Expected soteErr to be nil got %v", soteErr.FmtErrMsg)
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "/XXXX/syacko/.nkeys/creds/synadia/sote-staging/staging-soteadmin.creds",
		TESTSYNADIAURL, "test", true, 1, 250*time.Millisecond, false); soteErr.ErrCode != 209010 {
		tPtr.Errorf("TestNew failed: Expected error code of 209010 got %v", soteErr.FmtErrMsg)
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, "INVALID", "", TESTSYNADIAURL, "test", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != 209110 {
		tPtr.Errorf("TestNew failed: Expected error code of 209110 got %v", soteErr.FmtErrMsg)
	}

	if _, soteErr := New("XXXX", sConfigParams.STAGING, "", TESTSYNADIAURL, "test", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != 109999 {
		tPtr.Errorf("TestNew failed: Expected error code of 109999 got %v", soteErr.FmtErrMsg)
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", true, -1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestNew Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", true, 6,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("TestNew Failed: Expected error code to be nil got %v", soteErr.FmtErrMsg)
	}
}