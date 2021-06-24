package sDatabase

import (
	"runtime"
	"strconv"
	"testing"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

const (
	TESTINFOSCHEMA  = "information_schema"
	INFOSCHEMATABLE = "columns"
)

func init() {
	sLogger.SetLogMessagePrefix("scolumninfo_test.go")
}

func TestGetColumnInfo(tPtr *testing.T) {
	var (
		soteErr           sError.SoteError
		tConnInfo         ConnInfo
		columnInfo        []SColumnInfo
		columnInfoJSON    []byte
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	if _, soteErr = GetColumnInfo(TESTINFOSCHEMA, INFOSCHEMATABLE, tConnInfo); soteErr.ErrCode != 209299 {
		tPtr.Errorf("%v Failed: Expected error code of 209299", testName)
		tPtr.Fail()
	}

	if tConnInfo, soteErr = getMyDBConn(tPtr); soteErr.ErrCode == nil {
		if _, soteErr = GetColumnInfo("", INFOSCHEMATABLE, tConnInfo); soteErr.ErrCode != 200513 {
			tPtr.Errorf("%v Failed: Expected error code of 200513", testName)
			tPtr.Fail()
		}
		if _, soteErr = GetColumnInfo(TESTINFOSCHEMA, "", tConnInfo); soteErr.ErrCode != 200513 {
			tPtr.Errorf("%v Failed: Expected error code of 200513", testName)
			tPtr.Fail()
		}
		if columnInfo, soteErr = GetColumnInfo(TESTINFOSCHEMA, INFOSCHEMATABLE, tConnInfo); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil [%v]", testName, strconv.Itoa(soteErr.ErrCode.(int)))
			tPtr.Fail()
		}
		if len(columnInfo) == 0 {
			tPtr.Errorf("%v Failed: Expected at least one column's info to be returned", testName)
			tPtr.Fail()
		} else {
			if columnInfo[0].ColName == "" {
				tPtr.Errorf("%v Failed: Expected the column name to be returned", testName)
				tPtr.Fail()
			}
		}
		if columnInfoJSON, soteErr = GetColumnInfoJSONFormat(TESTINFOSCHEMA, INFOSCHEMATABLE, tConnInfo); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil [%v]", testName, strconv.Itoa(soteErr.ErrCode.(int)))
			tPtr.Fail()
		}
		if string(columnInfoJSON) == "" {
			tPtr.Errorf("%v Failed: Expected at least one column's info to be returned", testName)
			tPtr.Fail()
		}
	}

	tConnInfo.DBPoolPtr.Close()
}
