package packages

import (
	"testing"

	"gitlab.com/soteapps/packages/v2021/sDatabase"
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
func TestVerifyConnection(tPtr *testing.T) {
	var tConnInfo sDatabase.ConnInfo
	soteErr := sDatabase.VerifyConnection(tConnInfo)
	if soteErr.ErrCode != 209299 {
		tPtr.Errorf("VerifyConnection Failed: Expected 209299 error code.")
		tPtr.Fail()
	}

	if soteErr = sDatabase.GetAWSParams(); soteErr.ErrCode != nil {
		tPtr.Errorf("GetAWSParams Failed: Expected error code to be nil.")
	}

	tConnInfo, soteErr = sDatabase.GetConnection(sDatabase.DBName, sDatabase.DBUser, sDatabase.DBPassword, sDatabase.DBHost, sDatabase.DBSSLMode, sDatabase.DBPort, 3)
	if soteErr.ErrCode != nil {
		tPtr.Errorf("setConnectionValues Failed: Expected a nil error code.")
	}

	soteErr = sDatabase.VerifyConnection(tConnInfo)
	if soteErr.ErrCode != nil {
		tPtr.Errorf("VerifyConnection Failed: Expected a nil error code.")
	}
}
func TestToJSONString(tPtr *testing.T) {
	if soteErr := sDatabase.GetAWSParams(); soteErr.ErrCode != nil {
		tPtr.Errorf("GetAWSParams Failed: Expected error code to be nil.")
	}

	tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.DBName, sDatabase.DBUser, sDatabase.DBPassword, sDatabase.DBHost, sDatabase.DBSSLMode, sDatabase.DBPort, 3)
	if soteErr.ErrCode != nil {
		tPtr.Errorf("GetConnection Failed: Please Investigate")
	}

	var dbConnJSONString string
	if dbConnJSONString, soteErr = sDatabase.ToJSONString(tConnInfo.DSConnValues); soteErr.ErrCode != nil {
		tPtr.Errorf("ToJSONString Failed: Please Investigate")
	}

	if len(dbConnJSONString) == 0 {
		tPtr.Errorf("ToJSONString Failed: Please Investigate")
	}
}
func TestContext(tPtr *testing.T) {
	if soteErr := sDatabase.GetAWSParams(); soteErr.ErrCode != nil {
		tPtr.Errorf("GetAWSParams Failed: Expected error code to be nil.")
	}

	tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.DBName, sDatabase.DBUser, sDatabase.DBPassword, sDatabase.DBHost, sDatabase.DBSSLMode, sDatabase.DBPort, 3)
	if soteErr.ErrCode != nil {
		tPtr.Errorf("setConnectionValues Failed: Expected a nil error code.")
	}

	if tConnInfo.DBContext == nil {
		tPtr.Errorf("TestContext testing DBContext Failed: Expected a non-nil error code.")
	}
}
func TestSRow(tPtr *testing.T) {
	tRow := sDatabase.SRow(nil)
	if tRow != nil {
		tPtr.Errorf("TestSRow testing creation of SRow variable Failed: Expected error code to be nil.")
	}
}
func TestSRows(tPtr *testing.T) {
	tRows := sDatabase.SRows(nil)
	if tRows != nil {
		tPtr.Errorf("TestSRows testing creation of SRows variable Failed: Expected error code to be nil.")
	}
}

//
// sconstraintinfo
//
// func TestGetSingleColumnConstraintInfo(tPtr *testing.T) {
	// if soteErr := sDatabase.GetAWSParams(); soteErr.ErrCode != nil {
	// 	tPtr.Errorf("GetAWSParams Failed: Expected error code to be nil.")
	// 	t.Fatal()
	// }
	//
	// tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.DBName, sDatabase.DBUser, sDatabase.DBPassword, sDatabase.DBHost, sDatabase.DBSSLMode, sDatabase.DBPort, 3)
	// if soteErr.ErrCode != nil {
	// 	tPtr.Errorf("GetConnection Failed: Please investigate")
	// 	tPtr.Fail()
	// }
	//
	// var myConstraints []sDatabase.SConstraint
	// if myConstraints, soteErr = sDatabase.GetSingleColumnConstraintInfo(SOTETESTSCHEMA, tConnInfo); len(myConstraints) == 0 {
	// 	tPtr.Errorf("GetSingleColumnConstraintInfo Failed: myContraints should not be empty")
	// }
	// if myConstraints, soteErr = sDatabase.GetSingleColumnConstraintInfo(EMPTYVALUE, tConnInfo); soteErr.ErrCode != 200513 {
	// 	tPtr.Errorf("pkLookup Failed: Expected error code to be 200513.")
	//
	// }
// }
// func TestGetSingleColumnConstraintInfoNoConn(tPtr *testing.T) {
// 	var tConnInfo = sDatabase.ConnInfo{nil, sDatabase.ConnValues{}}
// 	if _, soteErr := sDatabase.GetSingleColumnConstraintInfo(SOTETESTSCHEMA, tConnInfo); soteErr.ErrCode != 209299 {
// 		tPtr.Errorf("pkLookup Failed: Expected error code to be 209299.")
//
// 	}
// }

//
// stableinfo
//
// func TestGetTables(tPtr *testing.T) {
	// var tConnInfo sDatabase.ConnInfo
	// if _, soteErr := sDatabase.GetTableList("sote", tConnInfo); soteErr.ErrCode != 209299 {
	// 	tPtr.Errorf("Get Tables Failed: Expected error code of 209299")
	// 	tPtr.Fail()
	// }
	//
	// if soteErr := sDatabase.GetAWSParams(); soteErr.ErrCode != nil {
	// 	tPtr.Errorf("GetAWSParams Failed: Expected error code to be nil.")
	// 	t.Fatal()
	// }
	//
	// tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.DBName, sDatabase.DBUser, sDatabase.DBPassword, sDatabase.DBHost, sDatabase.DBSSLMode, sDatabase.DBPort, 1)
	// if soteErr.ErrCode != nil {
	// 	tPtr.Errorf("Get Connection Failed: Please Investigate")
	// 	tPtr.Fail()
	// }
	//
	// var tableList []string
	// if tableList, soteErr = sDatabase.GetTableList(SOTETESTSCHEMA, tConnInfo); soteErr.ErrCode != nil {
	// 	tPtr.Errorf("Get Tables Failed: Expected error code to be nil")
	// 	tPtr.Fail()
	// }
	//
	// if len(tableList) == 0 {
	// 	tPtr.Errorf("Get Tables Failed: Expected at least one table name to be returned")
	// 	tPtr.Fail()
	// }
// }

//
// scolumninfo
//
// func TestGetColumnInfo(tPtr *testing.T) {
	// var tConnInfo sDatabase.ConnInfo
	// if _, soteErr := sDatabase.GetColumnInfo(SOTETESTSCHEMA, REFERENCETABLE, tConnInfo); soteErr.ErrCode != 209299 {
	// 	tPtr.Errorf("GetColumnInfo Failed: Expected error code of 209299")
	// 	tPtr.Fail()
	// }
	//
	// if soteErr := sDatabase.GetAWSParams(); soteErr.ErrCode != nil {
	// 	tPtr.Errorf("getAWSParams Failed: Expected error code to be nil.")
	// 	t.Fatal()
	// }
	//
	// tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.DBName, sDatabase.DBUser, sDatabase.DBPassword, sDatabase.DBHost, sDatabase.DBSSLMode, sDatabase.DBPort, 1)
	// if soteErr.ErrCode != nil {
	// 	tPtr.Errorf("GetConnection Failed: Please Investigate")
	// 	tPtr.Fail()
	// }
	//
	// var columnInfo []sDatabase.SColumnInfo
	// if columnInfo, soteErr = sDatabase.GetColumnInfo(SOTETESTSCHEMA, REFERENCETABLE, tConnInfo); soteErr.ErrCode != nil {
	// 	tPtr.Errorf("GetTableList Failed: Expected error code to be nil")
	// 	tPtr.Fail()
	// }
	//
	// if len(columnInfo) == 0 {
	// 	tPtr.Errorf("GetColumnInfo Failed: Expected at least one column's info to be returned")
	// 	tPtr.Fail()
	// } else {
	// 	if columnInfo[0].ColName == "" {
	// 		tPtr.Errorf("GetColumnInfo Failed: Expected the column name to be returned")
	// 		tPtr.Fail()
	// 	}
	// }
// }

//
// sprimarykeyinfo
//
// func TestPKLookup(tPtr *testing.T) {
// 	if soteErr := sDatabase.GetAWSParams(); soteErr.ErrCode != nil {
// 		tPtr.Errorf("getAWSParams Failed: Expected error code to be nil.")
// 		t.Fatal()
// 	}
//
// 	tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.DBName, sDatabase.DBUser, sDatabase.DBPassword, sDatabase.DBHost, sDatabase.DBSSLMode, sDatabase.DBPort, 3)
// 	if soteErr.ErrCode != nil {
// 		tPtr.Errorf("GetConnection Failed: Please investigate")
// 		tPtr.Fail()
// 	}
//
// 	if tbName, soteErr := sDatabase.PKLookup(SOTETESTSCHEMA, REFERENCETABLE, REFTBLCOLUMNNAME, tConnInfo); soteErr.ErrCode != nil && tbName != SELF {
// 		tPtr.Errorf("pkLookup Failed: Expected error code to be nil and tbName should be self.")
// 	}
//
// 	if tbName, soteErr := sDatabase.PKLookup(SOTETESTSCHEMA, PARENTCHILDTABLE, PCTBLCOLUMNNAME, tConnInfo); soteErr.ErrCode != nil && tbName != REFERENCETABLE {
// 		tPtr.Errorf("pkLookup Failed: Expected error code to be nil and tbName should be referencetable.")
// 	}
// }
// func TestPKLookupEmptyValues(tPtr *testing.T) {
// 	if soteErr := sDatabase.GetAWSParams(); soteErr.ErrCode != nil {
// 		tPtr.Errorf("getAWSParams Failed: Expected error code to be nil.")
// 		t.Fatal()
// 	}
//
// 	tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.DBName, sDatabase.DBUser, sDatabase.DBPassword, sDatabase.DBHost, sDatabase.DBSSLMode, sDatabase.DBPort, 3)
// 	if soteErr.ErrCode != nil {
// 		tPtr.Errorf("GetConnection Failed: Please investigate")
// 		tPtr.Fail()
// 	}
//
// 	if _, soteErr := sDatabase.PKLookup(EMPTYVALUE, PARENTCHILDTABLE, PCTBLCOLUMNNAME, tConnInfo); soteErr.ErrCode != 200513 {
// 		tPtr.Errorf("pkLookup Failed: Expected error code to be 200513.")
// 	}
// 	if _, soteErr := sDatabase.PKLookup(SOTETESTSCHEMA, EMPTYVALUE, PCTBLCOLUMNNAME, tConnInfo); soteErr.ErrCode != 200513 {
// 		tPtr.Errorf("pkLookup Failed: Expected error code to be 200513.")
// 	}
// 	if _, soteErr := sDatabase.PKLookup(SOTETESTSCHEMA, PARENTCHILDTABLE, EMPTYVALUE, tConnInfo); soteErr.ErrCode != 200513 {
// 		tPtr.Errorf("pkLookup Failed: Expected error code to be 200513.")
// 	}
// }
