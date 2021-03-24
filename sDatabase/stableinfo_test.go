package sDatabase

import (
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

func init() {
	sLogger.SetLogMessagePrefix("stableinfo_test.go")
}

// func TestGetTables(tPtr *testing.T) {
// 	var tConnInfo ConnInfo
// 	if _, soteErr := GetTableList("sote", tConnInfo); soteErr.ErrCode != 209299 {
// 		tPtr.Errorf("GetTableList Failed: Expected error code of 209299")
// 		tPtr.Fail()
// 	}
//
// 	if soteErr := GetAWSParams(); soteErr.ErrCode != nil {
// 		tPtr.Errorf("getAWSParams Failed: Expected error code to be nil.")
// 		t.Fatal()
// 	}
//
// 	tConnInfo, soteErr := GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
// 	if soteErr.ErrCode != nil {
// 		tPtr.Errorf("GetConnection Failed: Please Investigate")
// 		tPtr.Fail()
// 	}
//
// 	var tableList []string
// 	if tableList, soteErr = GetTableList("sote", tConnInfo); soteErr.ErrCode != nil {
// 		tPtr.Errorf("GetTableList Failed: Expected error code to be nil")
// 		tPtr.Fail()
// 	}
//
// 	if len(tableList) == 0 {
// 		tPtr.Errorf("GetTableList Failed: Expected at least one table name to be returned")
// 		tPtr.Fail()
// 	}
// }
