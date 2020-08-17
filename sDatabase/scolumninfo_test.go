package sDatabase

import (
	"testing"
)

const (
	REFERENCETABLE = "referencetable"
)

func TestGetColumnInfo(t *testing.T) {
	var tConnInfo ConnInfo
	if _, soteErr := GetColumnInfo("sote", REFERENCETABLE, tConnInfo); soteErr.ErrCode != 602999 {
		t.Errorf("GetColumnInfo Failed: Expected error code of 602999")
		t.Fail()
	}

	if soteErr = GetAWSParams(); soteErr.ErrCode != nil {
		t.Errorf("getAWSParams Failed: Expected error code to be nil.")
		t.Fatal()
	}

	tConnInfo, soteErr = GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("GetConnection Failed: Please Investigate")
		t.Fail()
	}

	var tableList []string
	if tableList, soteErr = GetTableList("sote", tConnInfo); soteErr.ErrCode != nil {
		t.Errorf("GetTableList Failed: Expected error code to be nil")
		t.Fail()
	}

	if len(tableList) == 0 {
		t.Errorf("GetTableList Failed: Expected at least one table name to be returned")
		t.Fail()
	}
}
