package sDatabase

import (
	"testing"

	"gitlab.com/soteapps/packages/v2020/sConfigParams"
	"gitlab.com/soteapps/packages/v2020/sError"
)

var (
	tConnInfo                                                ConnInfo
	awsRegion, dbName, dbUser, dbPassword, dbHost, dbSSLMode string
	dbPort                                                   int
	soteErr                                                  sError.SoteError
)

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

	if soteErr = getAWSParams(); soteErr.ErrCode != nil {
		t.Errorf("getAWSParams Failed: Expected error code to be nil.")
		t.Fatal()
	}

	tConnInfo, soteErr = GetConnection(dbName, dbUser, dbPassword, dbHost, dbSSLMode, dbPort, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("setConnectionValues Failed: Expected a nil error code.")
		t.Fail()
	}

	soteErr = VerifyConnection(tConnInfo)
	if soteErr.ErrCode != nil {
		t.Errorf("VerifyConnection Failed: Expected a nil error code.")
		t.Fail()
	}
}
func TestToJSONString(t *testing.T) {
	if soteErr = getAWSParams(); soteErr.ErrCode != nil {
		t.Errorf("getAWSParams Failed: Expected error code to be nil.")
		t.Fatal()
	}

	tConnInfo, soteErr = GetConnection(dbName, dbUser, dbPassword, dbHost, dbSSLMode, dbPort, 3)
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
func getAWSParams() (soteErr sError.SoteError) {
	awsRegion, soteErr = sConfigParams.GetEnvironmentAWSRegion()
	if soteErr.ErrCode == nil {
		dbName, soteErr = sConfigParams.GetDBName("API", awsRegion)
		if soteErr.ErrCode == nil {
			dbUser, soteErr = sConfigParams.GetDBUser("API", awsRegion)
			if soteErr.ErrCode == nil {
				dbPassword, soteErr = sConfigParams.GetDBPassword("API", awsRegion)
				if soteErr.ErrCode == nil {
					dbHost, soteErr = sConfigParams.GetDBHost("API", awsRegion)
					if soteErr.ErrCode == nil {
						dbSSLMode, soteErr = sConfigParams.GetDBSSLMode("API", awsRegion)
						if soteErr.ErrCode == nil {
							dbPort, soteErr = sConfigParams.GetDBPort("API", awsRegion)
						}
					}
				}
			}
		}
	}

	return
}
