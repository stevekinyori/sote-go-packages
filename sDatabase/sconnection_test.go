package sDatabase

import (
	"testing"
)

const EMPTYCONNECTIONJSON = `{"dbName":"","user":"","password":"","host":"","port":0,"timeout":0,"sslMode":""}`
const POPULATEDCONNECTIONJSON = `{"dbName":"dbName","user":"User","password":"Password","host":"Host","port":1,"timeout":1,"sslMode":"disable"}`

func TestGetConnectionStringEmpty(t *testing.T) {
	var dsConnValues ConnValues
	if s := GetConnectionValuesJSON(dsConnValues); s != EMPTYCONNECTIONJSON {
		t.Errorf("Get Connection String Empty: Returned value should have been empty.")
		t.Fail()
	}
}
func TestSetConnectionStringInvalidSSLMode(t *testing.T) {
	_, soteErr := setConnectionValues("dbName", "User", "Password", "Host", "INVALID", 1, 1)
	if soteErr.ErrCode != 602020 {
		t.Errorf("Set Connection String Invalid SSL Mode: Error code is not for an invalid sslMode.")
		t.Fail()
	}
}
func TestSetGetConnectionString(t *testing.T) {
	connValues, soteErr := setConnectionValues("dbName", "User", "Password", "Host", "disable", 1, 1)
	if soteErr.ErrCode != nil {
		t.Errorf("Set Connection String: Expected a nil error code.")
		t.Fail()
	}
	s := GetConnectionValuesJSON(connValues)
	if s == EMPTYCONNECTIONJSON {
		t.Errorf("Get Connection Values JSON: Expected a nil error code.")
		t.Fail()
	}
	if s != POPULATEDCONNECTIONJSON {
		t.Errorf("Get Connection Values JSON: Should have looked like %v.", POPULATEDCONNECTIONJSON)
		t.Fail()
	}

}
