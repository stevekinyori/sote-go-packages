package sDatabase

import (
	"testing"
)

func TestSetConnectionStringInvalidSSLMode(t *testing.T) {
	_, soteErr := setConnectionValues("dbName", "User", "Password", "Host", "INVALID", 1, 1)
	if soteErr.ErrCode != 602020 {
		t.Errorf("setConnectionValues Failed: Error code is not for an invalid sslMode.")
		t.Fail()
	}
}
func TestSetConnectionValues(t *testing.T) {
	_, soteErr := setConnectionValues("dbName", "User", "Password", "Host", "disable", 1, 1)
	if soteErr.ErrCode != nil {
		t.Errorf("setConnectionValues Failed: Expected a nil error code.")
		t.Fail()
	}
}
func TestGetConnection(t *testing.T) {
	_, soteErr := GetConnection("sote_development", "sote", "password", "localhost", "disable", 5442, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("GetConnection Failed: Please Investigate")
		t.Fail()
	}

	_, soteErr = setConnectionValues("dbName", "User", "Password", "Host", "disable", 1, 1)
	if soteErr.ErrCode != nil {
		t.Errorf("setConnectionValues Failed: Expected a nil error code.")
		t.Fail()
	}
}
func TestToJSONString(t *testing.T) {
	dbConnInfo, soteErr := GetConnection("sote_development", "sote", "password", "localhost", "disable", 5442, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("GetConnection Failed: Please Investigate")
		t.Fail()
	}

	var dbConnJSONString string
	if dbConnJSONString, soteErr = ToJSONString(dbConnInfo.dsConnValues); soteErr.ErrCode != nil {
		t.Errorf("ToJSONString Failed: Please Investigate")
		t.Fail()
	}

	if len(dbConnJSONString) == 0 {
		t.Errorf("ToJSONString Failed: Please Investigate")
		t.Fail()
	}
}