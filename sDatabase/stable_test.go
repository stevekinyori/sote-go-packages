package sDatabase

import (
	"testing"
)

func TestGetTables(t *testing.T) {
	if _, soteErr := getTables("sote"); soteErr.ErrCode != 602999 {
		t.Errorf("Get Tables Failed: Expected error code of 602999")
		t.Fail()
	}

	soteErr := GetConnection("single", "sote_development", "sote", "password", "localhost", "disable", 5442, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("Get Connection Failed: Please Investigate")
		t.Fail()
	}

	var tableList []string
	if tableList, soteErr = getTables("sote"); soteErr.ErrCode != nil {
		t.Errorf("Get Tables Failed: Expected error code to be nil")
		t.Fail()
	}
	if len(tableList) == 0 {
		t.Errorf("Get Tables Failed: Expected at least one table name to be returned")
		t.Fail()
	}

	soteErr = GetConnection("pool", "sote_development", "sote", "password", "localhost", "disable", 5442, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("Get Connection Failed: Please Investigate")
		t.Fail()
	}

	if tableList, soteErr = getTables("sote"); soteErr.ErrCode != nil {
		t.Errorf("Get Tables Failed: Expected error code to be nil")
		t.Fail()
	}
	if len(tableList) == 0 {
		t.Errorf("Get Tables Failed: Expected at least one table name to be returned")
		t.Fail()
	}
}
