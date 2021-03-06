package sDatabase

import (
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

func init() {
	sLogger.SetLogMessagePrefix("sconstraininfo_test.go")
}

// func TestGetSingleColumnConstraintInfo(t *testing.T) {
// 	if soteErr := GetAWSParams(); soteErr.ErrCode != nil {
// 		t.Errorf("getAWSParams Failed: Expected error code to be nil.")
// 		t.Fatal()
// 	}
//
// 	tConnInfo, soteErr := GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
// 	if soteErr.ErrCode != nil {
// 		t.Errorf("GetConnection Failed: Please investigate")
// 		t.Fail()
// 	}
//
// 	var myConstraints []SConstraint
// 	if myConstraints, soteErr = GetSingleColumnConstraintInfo(SOTETESTSCHEMA, tConnInfo); len(myConstraints) == 0 {
// 		t.Errorf("GetSingleColumnConstraintInfo Failed: myContraints should not be empty")
// 	}
// 	if myConstraints, soteErr = GetSingleColumnConstraintInfo(EMPTYVALUE, tConnInfo); soteErr.ErrCode != 200513 {
// 		t.Errorf("pkLookup Failed: Expected error code to be 200513.")
//
// 	}
// }
// func TestGetSingleColumnConstraintInfoNoConn(t *testing.T) {
// 	var tConnInfo = ConnInfo{nil, ConnValues{}}
// 	if _, soteErr := GetSingleColumnConstraintInfo(SOTETESTSCHEMA, tConnInfo); soteErr.ErrCode != 209299 {
// 		t.Errorf("pkLookup Failed: Expected error code to be 209299.")
//
// 	}
// }
