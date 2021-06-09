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
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v failed: Expected soteErr to be nil got %v", testName, soteErr.FmtErrMsg)
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v failed: Expected soteErr to be nil got %v", testName, soteErr.FmtErrMsg)
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
		tPtr.Errorf("%v failed: Expected soteErr to be nil got %v", testName, soteErr.FmtErrMsg)
	}
	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, homedir+"/.nkeys/creds/synadia/sote-staging/staging-soteadmin."+
		"creds",
		TESTSYNADIAURL, "test", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v failed: Expected soteErr to be nil got %v", testName, soteErr.FmtErrMsg)
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "/XXXX/syacko/.nkeys/creds/synadia/sote-staging/staging-soteadmin.creds",
		TESTSYNADIAURL, "test", true, 1, 250*time.Millisecond, false); soteErr.ErrCode != 209010 {
		tPtr.Errorf("%v failed: Expected error code of 209010 got %v", testName, soteErr.FmtErrMsg)
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, "INVALID", "", TESTSYNADIAURL, "test", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != 209110 {
		tPtr.Errorf("%v failed: Expected error code of 209110 got %v", testName, soteErr.FmtErrMsg)
	}

	if _, soteErr := New("XXXX", sConfigParams.STAGING, "", TESTSYNADIAURL, "test", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != 109999 {
		tPtr.Errorf("%v failed: Expected error code of 109999 got %v", testName, soteErr.FmtErrMsg)
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", true, -1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
	}

	if _, soteErr := New(TESTAPPLICATIONSYNADIA, sConfigParams.STAGING, "", TESTSYNADIAURL, "test", true, 6,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
	}
}