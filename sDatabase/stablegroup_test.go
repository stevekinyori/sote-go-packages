package sDatabase

import (
	"runtime"
	"testing"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

func init() {
	sLogger.SetLogMessagePrefix("stableinfo_test.go")
}

func TestGetTableGroupInfo(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		tTableGroupInfo   TableGroup
		soteErr           sError.SoteError
	)

	if _, soteErr = GetTableGroupInfo("", true); soteErr.ErrCode != 200513 {
		tPtr.Errorf("%v Failed: Expected error code to be 200513 got %v", testName, soteErr.FmtErrMsg)
	}
	if _, soteErr = GetTableGroupInfo("file-does-not-exist.json", true); soteErr.ErrCode != 109999 {
		tPtr.Errorf("%v Failed: Expected error code to be 109999 got %v", testName, soteErr.FmtErrMsg)
	}
	if tTableGroupInfo, soteErr = GetTableGroupInfo("organizations-table-group.json", true); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
	} else {
		if tTableGroupInfo.TableGroupName == "" {
			tPtr.Errorf("%v Failed: Expected Table Group Nme to be populated", testName)
		}
		if tTableGroupInfo.Schema == "" {
			tPtr.Errorf("%v Failed: Expected Schema to be populated", testName)
		}
		if len(tTableGroupInfo.Tables) == 0 {
			tPtr.Errorf("%v Failed: Expected table count to be greater than zero", testName)
		}
	}
}
