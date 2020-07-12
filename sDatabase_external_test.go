package packages

import (
	"testing"

	"gitlab.com/soteapps/packages/v2020/sDatabase"
)
//
// sconnection
//
func TestVerifyConnection(t *testing.T) {
	var tConnInfo sDatabase.ConnInfo
	soteErr := sDatabase.VerifyConnection(tConnInfo)
	if soteErr.ErrCode != 602999 {
		t.Errorf("VerifyConnection Failed: Expected 602999 error code.")
		t.Fail()
	}

	tConnInfo, soteErr = sDatabase.GetConnection("sote_development", "sote", "password", "localhost", "disable", 5442, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("setConnectionValues Failed: Expected a nil error code.")
		t.Fail()
	}

	soteErr = sDatabase.VerifyConnection(tConnInfo)
	if soteErr.ErrCode != nil {
		t.Errorf("VerifyConnection Failed: Expected a nil error code.")
		t.Fail()
	}
}
func TestToJSONString(t *testing.T) {
	dbConnInfo, soteErr := sDatabase.GetConnection("sote_development", "sote", "password", "localhost", "disable", 5442, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("GetConnection Failed: Please Investigate")
		t.Fail()
	}

	var dbConnJSONString string
	if dbConnJSONString, soteErr = sDatabase.ToJSONString(dbConnInfo.DSConnValues); soteErr.ErrCode != nil {
		t.Errorf("ToJSONString Failed: Please Investigate")
		t.Fail()
	}

	if len(dbConnJSONString) == 0 {
		t.Errorf("ToJSONString Failed: Please Investigate")
		t.Fail()
	}
}
//
// sconstraintinfo
//
func TestGetSingleColumnConstraintInfo(t *testing.T) {
	tConnInfo, soteErr := sDatabase.GetConnection("sote_development", "sote", "password", "localhost", "disable", 5442, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("GetConnection Failed: Please investigate")
		t.Fail()
	}

	var myConstraints []sDatabase.SConstraint
	if myConstraints, soteErr = sDatabase.GetSingleColumnConstraintInfo("sote", tConnInfo); len(myConstraints) == 0 {
		t.Errorf("GetSingleColumnConstraintInfo Failed: myContraints should not be empty")
	}
}
//
// stableinfo
//
func TestGetTables(t *testing.T) {
	var tConnInfo sDatabase.ConnInfo
	if _, soteErr := sDatabase.GetTableList("sote", tConnInfo); soteErr.ErrCode != 602999 {
		t.Errorf("Get Tables Failed: Expected error code of 602999")
		t.Fail()
	}

	tConnInfo, soteErr := sDatabase.GetConnection("sote_development", "sote", "password", "localhost", "disable", 5442, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("Get Connection Failed: Please Investigate")
		t.Fail()
	}

	var tableList []string
	if tableList, soteErr = sDatabase.GetTableList("sote", tConnInfo); soteErr.ErrCode != nil {
		t.Errorf("Get Tables Failed: Expected error code to be nil")
		t.Fail()
	}

	if len(tableList) == 0 {
		t.Errorf("Get Tables Failed: Expected at least one table name to be returned")
		t.Fail()
	}
}