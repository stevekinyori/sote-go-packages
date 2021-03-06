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

func TestGetColumnInfo(t *testing.T) {
	var tConnInfo ConnInfo
	if _, soteErr := GetColumnInfo(TESTINFOSCHEMA, INFOSCHEMATABLE, tConnInfo); soteErr.ErrCode != 209299 {
		t.Errorf("GetColumnInfo Failed: Expected error code of 209299")
		t.Fail()
	}

	if soteErr := GetAWSParams(); soteErr.ErrCode != nil {
		t.Errorf("getAWSParams Failed: Expected error code to be nil.")
		t.Fatal()
	}

	tConnInfo, soteErr := GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("GetConnection Failed: Please Investigate")
		t.Fail()
	}

	var columnInfo []SColumnInfo
	if columnInfo, soteErr = GetColumnInfo(TESTINFOSCHEMA, INFOSCHEMATABLE, tConnInfo); soteErr.ErrCode != nil {
		t.Errorf("GetColumnInfo Failed: Expected error code to be nil [" + strconv.Itoa(soteErr.ErrCode.(int)) + "]")
		t.Fail()
	}

	if len(columnInfo) == 0 {
		t.Errorf("GetColumnInfo Failed: Expected at least one column's info to be returned")
		t.Fail()
	} else {
		if columnInfo[0].ColName == "" {
			t.Errorf("GetColumnInfo Failed: Expected the column name to be returned")
			t.Fail()
		}
	}
}
