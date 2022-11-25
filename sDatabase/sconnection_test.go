package sDatabase

import (
	"context"
	"runtime"
	"testing"

	"gitlab.com/soteapps/packages/v2022/sConfigParams"
	"gitlab.com/soteapps/packages/v2022/sError"
	"gitlab.com/soteapps/packages/v2022/sLogger"
)

var (
	parentCtx    = context.Background()
	testConnInfo = ConnInfo{}
)

func init() {
	sLogger.SetLogMessagePrefix("sconnection_test.go")
	testConnInfo, _ = New(parentCtx, sConfigParams.DEVELOPMENT)
}

func TestSetConnectionValue(tPtr *testing.T) {
	_, soteErr := setConnectionValues("dbName", "User", "Password", "Host", "INVALID", 1, 1)
	if soteErr.ErrCode != 209220 {
		tPtr.Errorf("setConnectionValues Failed: Error code is not for an invalid sslMode.")
		tPtr.Fail()
	}
	_, soteErr = setConnectionValues("dbName", "User", "Password", "Host", "disable", 1, 1)
	if soteErr.ErrCode != nil {
		tPtr.Errorf("setConnectionValues Failed: Expected a nil error code.")
		tPtr.Fail()
	}
}
func TestSetConnectionValues(tPtr *testing.T) {
	_, soteErr := setConnectionValues("dbName", "User", "Password", "Host", "disable", 1, 1)
	if soteErr.ErrCode != nil {
		tPtr.Errorf("setConnectionValues Failed: Expected a nil error code.")
		tPtr.Fail()
	}
}
func TestVerifyConnection(tPtr *testing.T) {
	var (
		tConnInfo ConnInfo
		config    = &sConfigParams.Database{}
	)
	soteErr := VerifyConnection(tConnInfo)
	if soteErr.ErrCode != 209299 {
		tPtr.Errorf("VerifyConnection Failed: Expected 209299 error code.")
		tPtr.Fail()
	}

	tPtr.Run("test", func(tPtr *testing.T) {
		if config, soteErr = sConfigParams.GetAWSParams(parentCtx, sConfigParams.API, sConfigParams.DEMO); soteErr.ErrCode != nil {
			tPtr.Errorf("GetAWSParams Failed: Expected error code to be nil.")
			tPtr.Fatal()
		}
	})

	if config, soteErr = sConfigParams.GetAWSParams(parentCtx, sConfigParams.API, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetAWSParams Failed: Expected error code to be nil.")
		tPtr.Fatal()
	}

	tConnInfo, soteErr = GetConnection(config.Name, config.User, config.Password, config.Host, config.SSLMode, config.Port, 3)
	if soteErr.ErrCode != nil {
		tPtr.Errorf("setConnectionValues Failed: Expected a nil error code.")
		tPtr.Fail()
	}

	soteErr = VerifyConnection(tConnInfo)
	if soteErr.ErrCode != nil {
		tPtr.Errorf("VerifyConnection Failed: Expected a nil error code.")
		tPtr.Fail()
	}

	// This will test the condition that no database is available to connect
	tConnInfo, soteErr = GetConnection(config.Name, config.User, config.Password, config.Host, config.SSLMode, 65000, 3)
	if soteErr.ErrCode != 209299 {
		tPtr.Errorf("setConnectionValues Failed: Expected 209299 error code.")
		tPtr.Fail()
	}

}
func TestToJSONString(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		config  = &sConfigParams.Database{}
	)

	if config, soteErr = sConfigParams.GetAWSParams(parentCtx, sConfigParams.API, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {

		tPtr.Errorf("GetAWSParams Failed: Expected error code to be nil.")
		tPtr.Fatal()
	}

	tConnInfo, soteErr := GetConnection(config.Name, config.User, config.Password, config.Host, config.SSLMode, config.Port, 3)
	if soteErr.ErrCode != nil {
		tPtr.Errorf("GetConnection Failed: Please Investigate")
		tPtr.Fail()
	}

	var dbConnJSONString string
	if dbConnJSONString, soteErr = ToJSONString(tConnInfo.DSConnValues); soteErr.ErrCode != nil {
		tPtr.Errorf("ToJSONString Failed: Please Investigate")
		tPtr.Fail()
	}

	if len(dbConnJSONString) == 0 {
		tPtr.Errorf("ToJSONString Failed: Please Investigate")
		tPtr.Fail()
	}
}
func TestContext(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		config  = &sConfigParams.Database{}
	)

	if config, soteErr = sConfigParams.GetAWSParams(parentCtx, sConfigParams.API, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {

		tPtr.Errorf("GetAWSParams Failed: Expected error code to be nil.")
		tPtr.Fatal()
	}

	tConnInfo, soteErr := GetConnection(config.Name, config.User, config.Password, config.Host, config.SSLMode, config.Port, 3)
	if soteErr.ErrCode != nil {
		tPtr.Errorf("setConnectionValues Failed: Expected a nil error code.")
		tPtr.Fail()
	}

	if tConnInfo.DBContext == nil {
		tPtr.Errorf("TestContext testing DBContext Failed: Expected a non-nil error code.")
		tPtr.Fail()
	}
}
func TestSRow(tPtr *testing.T) {
	tRow := SRow(nil)
	if tRow != nil {
		tPtr.Errorf("TestSRow testing creation of SRow variable Failed: Expected error code to be nil.")
		tPtr.Fail()
	}
}
func TestSRows(tPtr *testing.T) {
	tRows := SRows(nil)
	if tRows != nil {
		tPtr.Errorf("TestSRows testing creation of SRows variable Failed: Expected error code to be nil.")
		tPtr.Fail()
	}
}
func getMyDBConn(tPtr *testing.T) (myDBConn ConnInfo, soteErr sError.SoteError) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		config            = &sConfigParams.Database{}
	)

	if config, soteErr = sConfigParams.GetAWSParams(parentCtx, sConfigParams.API, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {

		tPtr.Errorf("%v Failed: Expected error code to be nil.", testName)
		tPtr.Fatal()
	}

	myDBConn, soteErr = GetConnection(config.Name, config.User, config.Password, config.Host, config.SSLMode, config.Port, 3)
	if soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected a nil error code.", testName)
		tPtr.Fail()
	}

	return
}
