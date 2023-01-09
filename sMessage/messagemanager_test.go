package sMessage

import (
	"context"
	"os"
	"runtime"
	"testing"
	"time"

	"gitlab.com/soteapps/packages/v2023/sConfigParams"
	"gitlab.com/soteapps/packages/v2023/sError"
)

const (
	TESTAPPLICATIONSYNADIA = "synadia"
	TESTSYNADIAURL         = "west.eu.geo.ngs.global"
)

var parentCtx = context.Background()

func TestNew(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	if _, soteErr := New(parentCtx, TESTAPPLICATIONSYNADIA, sConfigParams.DEVELOPMENT, "", TESTSYNADIAURL, "test", false, 1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v failed: Expected soteErr to be nil got %v", testName, soteErr.FmtErrMsg)
	}

	if _, soteErr := New(parentCtx, TESTAPPLICATIONSYNADIA, sConfigParams.DEVELOPMENT, "", TESTSYNADIAURL, "test", true, 1,
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

	if _, soteErr := New(parentCtx, TESTAPPLICATIONSYNADIA, sConfigParams.DEVELOPMENT,
		homedir+"/.nkeys/creds/synadia/sote-staging/staging-soteadmin."+
			"creds", TESTSYNADIAURL, "test", false, 1, 250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v failed: Expected soteErr to be nil got %v", testName, soteErr.FmtErrMsg)
	}
	if _, soteErr := New(parentCtx, TESTAPPLICATIONSYNADIA, sConfigParams.DEVELOPMENT,
		homedir+"/.nkeys/creds/synadia/sote-staging/staging-soteadmin."+
			"creds", TESTSYNADIAURL, "test", true, 1, 250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v failed: Expected soteErr to be nil got %v", testName, soteErr.FmtErrMsg)
	}

	if _, soteErr := New(parentCtx, TESTAPPLICATIONSYNADIA, sConfigParams.DEVELOPMENT,
		"/XXXX/syacko/.nkeys/creds/synadia/sote-staging/staging-soteadmin.creds",
		TESTSYNADIAURL, "test", true, 1, 250*time.Millisecond, false); soteErr.ErrCode != sError.ErrMissingFile {
		tPtr.Errorf("%v failed: Expected error code of %v got %v", testName, sError.ErrMissingFile, soteErr.FmtErrMsg)
	}

	if _, soteErr := New(parentCtx, TESTAPPLICATIONSYNADIA, "INVALID", "", TESTSYNADIAURL, "test", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != sError.ErrInvalidEnvValue {
		tPtr.Errorf("%v failed: Expected error code of %v got %v", testName, sError.ErrInvalidEnvValue,
			soteErr.FmtErrMsg)
	}

	if _, soteErr := New(parentCtx, "XXXX", sConfigParams.DEVELOPMENT, "", TESTSYNADIAURL, "test", true, 1,
		250*time.Millisecond, false); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("%v failed: Expected error code of %v got %v", testName, sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}

	if _, soteErr := New(parentCtx, TESTAPPLICATIONSYNADIA, sConfigParams.DEVELOPMENT, "", TESTSYNADIAURL, "test", true, -1,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
	}

	if _, soteErr := New(parentCtx, TESTAPPLICATIONSYNADIA, sConfigParams.DEVELOPMENT, "", TESTSYNADIAURL, "test", true, 6,
		250*time.Millisecond, false); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
	}
}
