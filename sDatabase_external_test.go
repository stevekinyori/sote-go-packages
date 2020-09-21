package packages

import (
	"testing"

	"gitlab.com/soteapps/packages/v2020/sDatabase"
)

const (
	TESTINFOSCHEMA  = "information_schema"
	INFOSCHEMATABLE = "columns"
	SOTETESTSCHEMA   = "sotetest"
	REFERENCETABLE   = "referencetable"
	REFTBLCOLUMNNAME = "reference_name"
	PARENTCHILDTABLE = "parentchildtable"
	PCTBLCOLUMNNAME  = "reference_name"
	EMPTYVALUE       = ""
	SELF             = "self"
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

	if soteErr = sDatabase.GetAWSParams(); soteErr.ErrCode != nil {
		t.Errorf("GetAWSParams Failed: Expected error code to be nil.")
		t.Fatal()
	}

	tConnInfo, soteErr = sDatabase.GetConnection(sDatabase.DBName, sDatabase.DBUser, sDatabase.DBPassword, sDatabase.DBHost, sDatabase.DBSSLMode, sDatabase.DBPort, 3)
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
	if soteErr := sDatabase.GetAWSParams(); soteErr.ErrCode != nil {
		t.Errorf("GetAWSParams Failed: Expected error code to be nil.")
		t.Fatal()
	}

	tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.DBName, sDatabase.DBUser, sDatabase.DBPassword, sDatabase.DBHost, sDatabase.DBSSLMode, sDatabase.DBPort, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("GetConnection Failed: Please Investigate")
		t.Fail()
	}

	var dbConnJSONString string
	if dbConnJSONString, soteErr = sDatabase.ToJSONString(tConnInfo.DSConnValues); soteErr.ErrCode != nil {
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
	// if soteErr := sDatabase.GetAWSParams(); soteErr.ErrCode != nil {
	// 	t.Errorf("GetAWSParams Failed: Expected error code to be nil.")
	// 	t.Fatal()
	// }
	//
	// tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.DBName, sDatabase.DBUser, sDatabase.DBPassword, sDatabase.DBHost, sDatabase.DBSSLMode, sDatabase.DBPort, 3)
	// if soteErr.ErrCode != nil {
	// 	t.Errorf("GetConnection Failed: Please investigate")
	// 	t.Fail()
	// }
	//
	// var myConstraints []sDatabase.SConstraint
	// if myConstraints, soteErr = sDatabase.GetSingleColumnConstraintInfo(SOTETESTSCHEMA, tConnInfo); len(myConstraints) == 0 {
	// 	t.Errorf("GetSingleColumnConstraintInfo Failed: myContraints should not be empty")
	// }
	// if myConstraints, soteErr = sDatabase.GetSingleColumnConstraintInfo(EMPTYVALUE, tConnInfo); soteErr.ErrCode != 200513 {
	// 	t.Errorf("pkLookup Failed: Expected error code to be 200513.")
	//
	// }
}
func TestGetSingleColumnConstraintInfoNoConn(t *testing.T) {
	var tConnInfo = sDatabase.ConnInfo{nil, sDatabase.ConnValues{}}
	if _, soteErr := sDatabase.GetSingleColumnConstraintInfo(SOTETESTSCHEMA, tConnInfo); soteErr.ErrCode != 602999 {
		t.Errorf("pkLookup Failed: Expected error code to be 602999.")

	}
}

//
// stableinfo
//
func TestGetTables(t *testing.T) {
	// var tConnInfo sDatabase.ConnInfo
	// if _, soteErr := sDatabase.GetTableList("sote", tConnInfo); soteErr.ErrCode != 602999 {
	// 	t.Errorf("Get Tables Failed: Expected error code of 602999")
	// 	t.Fail()
	// }
	//
	// if soteErr := sDatabase.GetAWSParams(); soteErr.ErrCode != nil {
	// 	t.Errorf("GetAWSParams Failed: Expected error code to be nil.")
	// 	t.Fatal()
	// }
	//
	// tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.DBName, sDatabase.DBUser, sDatabase.DBPassword, sDatabase.DBHost, sDatabase.DBSSLMode, sDatabase.DBPort, 1)
	// if soteErr.ErrCode != nil {
	// 	t.Errorf("Get Connection Failed: Please Investigate")
	// 	t.Fail()
	// }
	//
	// var tableList []string
	// if tableList, soteErr = sDatabase.GetTableList(SOTETESTSCHEMA, tConnInfo); soteErr.ErrCode != nil {
	// 	t.Errorf("Get Tables Failed: Expected error code to be nil")
	// 	t.Fail()
	// }
	//
	// if len(tableList) == 0 {
	// 	t.Errorf("Get Tables Failed: Expected at least one table name to be returned")
	// 	t.Fail()
	// }
}

//
// scolumninfo
//
func TestGetColumnInfo(t *testing.T) {
	// var tConnInfo sDatabase.ConnInfo
	// if _, soteErr := sDatabase.GetColumnInfo(SOTETESTSCHEMA, REFERENCETABLE, tConnInfo); soteErr.ErrCode != 602999 {
	// 	t.Errorf("GetColumnInfo Failed: Expected error code of 602999")
	// 	t.Fail()
	// }
	//
	// if soteErr := sDatabase.GetAWSParams(); soteErr.ErrCode != nil {
	// 	t.Errorf("getAWSParams Failed: Expected error code to be nil.")
	// 	t.Fatal()
	// }
	//
	// tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.DBName, sDatabase.DBUser, sDatabase.DBPassword, sDatabase.DBHost, sDatabase.DBSSLMode, sDatabase.DBPort, 1)
	// if soteErr.ErrCode != nil {
	// 	t.Errorf("GetConnection Failed: Please Investigate")
	// 	t.Fail()
	// }
	//
	// var columnInfo []sDatabase.SColumnInfo
	// if columnInfo, soteErr = sDatabase.GetColumnInfo(SOTETESTSCHEMA, REFERENCETABLE, tConnInfo); soteErr.ErrCode != nil {
	// 	t.Errorf("GetTableList Failed: Expected error code to be nil")
	// 	t.Fail()
	// }
	//
	// if len(columnInfo) == 0 {
	// 	t.Errorf("GetColumnInfo Failed: Expected at least one column's info to be returned")
	// 	t.Fail()
	// } else {
	// 	if columnInfo[0].ColName == "" {
	// 		t.Errorf("GetColumnInfo Failed: Expected the column name to be returned")
	// 		t.Fail()
	// 	}
	// }
}

//
// sprimarykeyinfo
//
func TestPKLookup(t *testing.T) {
	if soteErr := sDatabase.GetAWSParams(); soteErr.ErrCode != nil {
		t.Errorf("getAWSParams Failed: Expected error code to be nil.")
		t.Fatal()
	}

	tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.DBName, sDatabase.DBUser, sDatabase.DBPassword, sDatabase.DBHost, sDatabase.DBSSLMode, sDatabase.DBPort, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("GetConnection Failed: Please investigate")
		t.Fail()
	}

	if tbName, soteErr := sDatabase.PKLookup(SOTETESTSCHEMA, REFERENCETABLE, REFTBLCOLUMNNAME, tConnInfo); soteErr.ErrCode != nil && tbName != SELF {
		t.Errorf("pkLookup Failed: Expected error code to be nil and tbName should be self.")
	}

	if tbName, soteErr := sDatabase.PKLookup(SOTETESTSCHEMA, PARENTCHILDTABLE, PCTBLCOLUMNNAME, tConnInfo); soteErr.ErrCode != nil && tbName != REFERENCETABLE {
		t.Errorf("pkLookup Failed: Expected error code to be nil and tbName should be referencetable.")
	}
}
func TestPKLookupEmptyValues(t *testing.T) {
	if soteErr := sDatabase.GetAWSParams(); soteErr.ErrCode != nil {
		t.Errorf("getAWSParams Failed: Expected error code to be nil.")
		t.Fatal()
	}

	tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.DBName, sDatabase.DBUser, sDatabase.DBPassword, sDatabase.DBHost, sDatabase.DBSSLMode, sDatabase.DBPort, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("GetConnection Failed: Please investigate")
		t.Fail()
	}

	if _, soteErr := sDatabase.PKLookup(EMPTYVALUE, PARENTCHILDTABLE, PCTBLCOLUMNNAME, tConnInfo); soteErr.ErrCode != 200513 {
		t.Errorf("pkLookup Failed: Expected error code to be 200513.")
	}
	if _, soteErr := sDatabase.PKLookup(SOTETESTSCHEMA, EMPTYVALUE, PCTBLCOLUMNNAME, tConnInfo); soteErr.ErrCode != 200513 {
		t.Errorf("pkLookup Failed: Expected error code to be 200513.")
	}
	if _, soteErr := sDatabase.PKLookup(SOTETESTSCHEMA, PARENTCHILDTABLE, EMPTYVALUE, tConnInfo); soteErr.ErrCode != 200513 {
		t.Errorf("pkLookup Failed: Expected error code to be 200513.")
	}
}
