package sDatabase

import (
	"testing"

	"gitlab.com/soteapps/packages/v2021/sLogger"
)

func init() {
	sLogger.SetLogMessagePrefix("sconnection_test.go")
}

func TestSetConnectionValue(t *testing.T) {
	_, soteErr := setConnectionValues("dbName", "User", "Password", "Host", "INVALID", 1, 1)
	if soteErr.ErrCode != 209220 {
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
	var tConnInfo ConnInfo
	soteErr := VerifyConnection(tConnInfo)
	if soteErr.ErrCode != 209299 {
		t.Errorf("VerifyConnection Failed: Expected 209299 error code.")
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
	if soteErr.ErrCode != 209299 {
		t.Errorf("setConnectionValues Failed: Expected 209299 error code.")
		t.Fail()
	}

}
func TestToJSONString(t *testing.T) {
	if soteErr := GetAWSParams(); soteErr.ErrCode != nil {
		t.Errorf("GetAWSParams Failed: Expected error code to be nil.")
		t.Fatal()
	}

	tConnInfo, soteErr := GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
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
func TestContext(t *testing.T) {
	if soteErr := GetAWSParams(); soteErr.ErrCode != nil {
		t.Errorf("GetAWSParams Failed: Expected error code to be nil.")
		t.Fatal()
	}

	tConnInfo, soteErr := GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("setConnectionValues Failed: Expected a nil error code.")
		t.Fail()
	}

	if tConnInfo.DBContext == nil {
		t.Errorf("TestContext testing DBContext Failed: Expected a non-nil error code.")
		t.Fail()
	}
}
func TestSRow(t *testing.T) {
	tRow := SRow(nil)
	if tRow != nil {
		t.Errorf("TestSRow testing creation of SRow variable Failed: Expected error code to be nil.")
		t.Fail()
	}
}
func TestSRows(t *testing.T) {
	tRows := SRows(nil)
	if tRows != nil {
		t.Errorf("TestSRows testing creation of SRows variable Failed: Expected error code to be nil.")
		t.Fail()
	}
}
