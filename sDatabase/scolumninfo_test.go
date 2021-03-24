package sDatabase

import (
	"strconv"
	"testing"

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
	var tConnInfo ConnInfo
	if _, soteErr := GetColumnInfo(TESTINFOSCHEMA, INFOSCHEMATABLE, tConnInfo); soteErr.ErrCode != 209299 {
		tPtr.Errorf("GetColumnInfo Failed: Expected error code of 209299")
		tPtr.Fail()
	}

	if soteErr := GetAWSParams(); soteErr.ErrCode != nil {
		tPtr.Errorf("getAWSParams Failed: Expected error code to be nil.")
		tPtr.Fatal()
	}

	tConnInfo, soteErr := GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
	if soteErr.ErrCode != nil {
		tPtr.Errorf("GetConnection Failed: Please Investigate")
		tPtr.Fail()
	}

	var columnInfo []SColumnInfo
	if columnInfo, soteErr = GetColumnInfo(TESTINFOSCHEMA, INFOSCHEMATABLE, tConnInfo); soteErr.ErrCode != nil {
		tPtr.Errorf("GetColumnInfo Failed: Expected error code to be nil [" + strconv.Itoa(soteErr.ErrCode.(int)) + "]")
		tPtr.Fail()
	}

	if len(columnInfo) == 0 {
		tPtr.Errorf("GetColumnInfo Failed: Expected at least one column's info to be returned")
		tPtr.Fail()
	} else {
		if columnInfo[0].ColName == "" {
			tPtr.Errorf("GetColumnInfo Failed: Expected the column name to be returned")
			tPtr.Fail()
		}
	}
}
