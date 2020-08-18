package sDatabase

import (
	"testing"

	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

var (
	tConnInfo                                                ConnInfo
	soteErr                                                  sError.SoteError
)

func init() {
	sLogger.SetLogMessagePrefix("sconnection_test.go")
}

func TestSetConnectionValue(t *testing.T) {
	_, soteErr := setConnectionValues("dbName", "User", "Password", "Host", "INVALID", 1, 1)
	if soteErr.ErrCode != 602020 {
		t.Errorf("setConnectionValues Failed: Error code is not for an invalid sslMode.")
		t.Fail()
	}
	_, soteErr = setConnectionValues("dbName", "User", "Password", "Host", "disable", 1, 1)
	if soteErr.ErrCode != nil {
		t.Errorf("setConnectionValues Failed: Expected a nil error code.")
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
func TestVerifyConnection(t *testing.T) {
	soteErr = VerifyConnection(tConnInfo)
	if soteErr.ErrCode != 602999 {
		t.Errorf("VerifyConnection Failed: Expected 602999 error code.")
		t.Fail()
	}

	if soteErr = GetAWSParams(); soteErr.ErrCode != nil {
		t.Errorf("GetAWSParams Failed: Expected error code to be nil.")
		t.Fatal()
	}

	tConnInfo, soteErr = GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("setConnectionValues Failed: Expected a nil error code.")
		t.Fail()
	}

	soteErr = VerifyConnection(tConnInfo)
	if soteErr.ErrCode != nil {
		t.Errorf("VerifyConnection Failed: Expected a nil error code.")
		t.Fail()
	}

	// This will test the condition that no database is available to connect
	tConnInfo, soteErr = GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, 65000, 3)
	if soteErr.ErrCode != 602999 {
		t.Errorf("setConnectionValues Failed: Expected 602999 error code.")
		t.Fail()
	}

}
func TestToJSONString(t *testing.T) {
	if soteErr = GetAWSParams(); soteErr.ErrCode != nil {
		t.Errorf("GetAWSParams Failed: Expected error code to be nil.")
		t.Fatal()
	}

	tConnInfo, soteErr = GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("GetConnection Failed: Please Investigate")
		t.Fail()
	}

	var dbConnJSONString string
	if dbConnJSONString, soteErr = ToJSONString(tConnInfo.DSConnValues); soteErr.ErrCode != nil {
		t.Errorf("ToJSONString Failed: Please Investigate")
		t.Fail()
	}

	if len(dbConnJSONString) == 0 {
		t.Errorf("ToJSONString Failed: Please Investigate")
		t.Fail()
	}
}
