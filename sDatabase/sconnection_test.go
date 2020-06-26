package sDatabase

import (
	"testing"
)

const EMPTYCONNECTIONJSON = `{"connType":"","dbName":"","user":"","password":"","host":"","port":0,"timeout":0,"sslMode":""}`
const POPULATEDCONNECTIONJSON = `{"connType":"CONNTYPE","dbName":"dbName","user":"User","password":"Password","host":"Host","port":1,"timeout":1,"sslMode":"disable"}`

func TestGetConnectionStringEmpty(t *testing.T) {
	if s := GetConnectionValuesJSON(); s != EMPTYCONNECTIONJSON {
		t.Errorf("Get Connection String Empty: Returned value should have been empty.")
		t.Fail()
	}
}
func TestSetConnectionStringInvalidSSLMode(t *testing.T) {
	soteErr := setConnectionValues("connType", "dbName", "User", "Password", "Host", "INVALID", 1, 1)
	if soteErr.ErrCode != 602020 {
		t.Errorf("Set Connection String Invalid SSL Mode: Error code is not for an invalid sslMode.")
		t.Fail()
	}
}
func TestSetGetConnectionString(t *testing.T) {
	soteErr := setConnectionValues("connType", "dbName", "User", "Password", "Host", "disable", 1, 1)
	if soteErr.ErrCode != nil {
		t.Errorf("Set Connection String: Expected a nil error code.")
		t.Fail()
	}
	s := GetConnectionValuesJSON()
	if s == EMPTYCONNECTIONJSON {
		t.Errorf("Get Connection Values JSON: Expected a nil error code.")
		t.Fail()
	}
	if s != POPULATEDCONNECTIONJSON {
		t.Errorf("Get Connection Values JSON: Should have looked like %v.", POPULATEDCONNECTIONJSON)
		t.Fail()
	}

}
func TestGetConnection(t *testing.T) {
	soteErr := GetConnection("single", "sote_development", "sote", "password", "localhost", "disable", 5442, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("Get Connection Failed: Should have returned a pointer to the single database connection")
		t.Fail()
	}

	soteErr = GetConnection("pool", "sote_development", "sote", "password", "localhost", "disable", 5442, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("Get Connection Failed: Should have returned a pointer to the pool database connection")
		t.Fail()
	}

	soteErr = setConnectionValues("connType", "dbName", "User", "Password", "Host", "disable", 1, 1)
	if soteErr.ErrCode != nil {
		t.Errorf("Set Connection String: Expected a nil error code.")
		t.Fail()
	}
}
func TestGetConnectionValues(t *testing.T) {
	soteErr := GetConnection("single", "sote_development", "sote", "password", "localhost", "disable", 5442, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("Get Connection Failed: Should have returned a pointer to the single database connection")
		t.Fail()
	}

	soteErr = setConnectionValues("connType", "dbName", "User", "Password", "Host", "disable", 1, 1)
	if soteErr.ErrCode != nil {
		t.Errorf("Set Connection String: Expected a nil error code.")
		t.Fail()
	}

	s := GetConnectionValuesJSON()
	if s != POPULATEDCONNECTIONJSON {
		t.Errorf("Get Connection Values JSON: Expected JSON string with populated values")
		t.Fail()
	}
}
