package tests

import (
	"context"
	"runtime"
	"testing"

	"gitlab.com/soteapps/packages/v2023/sConfigParams"
	"gitlab.com/soteapps/packages/v2023/sDatabase"
	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
)

const (
	TESTINFOSCHEMA   = "information_schema"
	INFOSCHEMATABLE  = "columns"
	SOTETESTSCHEMA   = "sotetest"
	REFERENCETABLE   = "referencetable"
	REFTBLCOLUMNNAME = "reference_name"
	PARENTCHILDTABLE = "parentchildtable"
	PCTBLCOLUMNNAME  = "reference_name"
	EMPTYVALUE       = ""
	SELF             = "self"
)

func init() {
	sLogger.SetLogMessagePrefix("packages")
}

//
// sconnection
//

func TestVerifyConnection(tPtr *testing.T) {
	var (
		tConnInfo sDatabase.ConnInfo
		config    = &sConfigParams.Database{}
		soteErr   sError.SoteError
	)

	soteErr = sDatabase.VerifyConnection(tConnInfo)
	if soteErr.ErrCode != sError.ErrDBConnectionError {
		tPtr.Errorf("VerifyConnection Failed: Expected %v error code.", sError.ErrDBConnectionError)
		tPtr.Fail()
	}

	if config, soteErr = sConfigParams.GetAWSParams(context.Background(), sConfigParams.API, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetAWSParams Failed: Expected error code to be nil.")
	}

	tConnInfo, soteErr = sDatabase.GetConnection(config.Name, config.Schema, config.User, config.Password, config.Host, config.SSLMode,
		config.Port, 3)
	if soteErr.ErrCode != nil {
		tPtr.Errorf("setConnectionValues Failed: Expected a nil error code.")
	}

	soteErr = sDatabase.VerifyConnection(tConnInfo)
	if soteErr.ErrCode != nil {
		tPtr.Errorf("VerifyConnection Failed: Expected a nil error code.")
	}
}
func TestToJSONString(tPtr *testing.T) {
	var (
		config  = &sConfigParams.Database{}
		soteErr sError.SoteError
	)

	if config, soteErr = sConfigParams.GetAWSParams(context.Background(), sConfigParams.API, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetAWSParams Failed: Expected error code to be nil.")
	}

	tConnInfo, soteErr := sDatabase.GetConnection(config.Name, config.Schema, config.User, config.Password, config.Host, config.SSLMode,
		config.Port, 3)
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
	var (
		config  = &sConfigParams.Database{}
		soteErr sError.SoteError
	)

	if config, soteErr = sConfigParams.GetAWSParams(context.Background(), sConfigParams.API, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetAWSParams Failed: Expected error code to be nil.")
	}

	tConnInfo, soteErr := sDatabase.GetConnection(config.Name, config.Schema, config.User, config.Password, config.Host, config.SSLMode,
		config.Port, 3)
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

// filters
func TestFormatArrayFilterCondition(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	tPtr.Run("multiple prefixes", func(tPtr *testing.T) {
		if _, soteErr := sDatabase.FormatListQueryConditions(context.Background(), &sDatabase.FormatConditionParams{
			InitialParamCount: 0,
			RecordLimitCount:  0,
			TblPrefixes:       []string{"table1.", "table2.", "table3."},
			Filters: map[string][]sDatabase.FilterFields{
				"AND": {
					sDatabase.FilterFields{
						FilterCommon: sDatabase.FilterCommon{
							Operator: "IN",
							Value:    []string{"name1", "name2"},
						},
						FieldName: "column-name",
					},
					sDatabase.FilterFields{
						FilterCommon: sDatabase.FilterCommon{
							Operator: "=",
							Value:    1,
						},
						FieldName: "column-id",
					},
				},
			},
			SortOrderKeysMap: map[string]sDatabase.TableColumn{
				"column-name": {
					ColumnName:      "column_name",
					CaseInsensitive: true,
				},
				"column-id": {
					ColumnName:      "column_id",
					CaseInsensitive: false,
				},
			},
			SortOrder: sDatabase.SortOrder{
				TblPrefix: "table.",
				Fields:    map[string]string{"column-name": "DESC", "column-id": "DESC"},
			},
		}); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
		}
	})
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
// tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.Name, sDatabase.User, sDatabase.Password, sDatabase.Host, sDatabase.SSLMode, sDatabase.Port, 3)
// if soteErr.ErrCode != nil {
// 	tPtr.Errorf("GetConnection Failed: Please investigate")
// 	tPtr.Fail()
// }
//
// var myConstraints []sDatabase.SConstraint
// if myConstraints, soteErr = sDatabase.GetSingleColumnConstraintInfo(SOTETESTSCHEMA, tConnInfo); len(myConstraints) == 0 {
// 	tPtr.Errorf("GetSingleColumnConstraintInfo Failed: myContraints should not be empty")
// }
// if myConstraints, soteErr = sDatabase.GetSingleColumnConstraintInfo(EMPTYVALUE, tConnInfo); soteErr.ErrCode != sError.ErrMissingParameters {
// 	tPtr.Errorf("pkLookup Failed: Expected error code to be %v.",sError.ErrMissingParameters)
//
// }
// }
// func TestGetSingleColumnConstraintInfoNoConn(tPtr *testing.T) {
// 	var tConnInfo = sDatabase.ConnInfo{nil, sDatabase.ConnValues{}}
// 	if _, soteErr := sDatabase.GetSingleColumnConstraintInfo(SOTETESTSCHEMA, tConnInfo); soteErr.ErrCode != sError.ErrDBConnectionError {
// 		tPtr.Errorf("pkLookup Failed: Expected error code to be %v.",sError.ErrDBConnectionError)
//
// 	}
// }

//
// stableinfo
//
// func TestGetTables(tPtr *testing.T) {
// var tConnInfo sDatabase.ConnInfo
// if _, soteErr := sDatabase.GetTableList("sote", tConnInfo); soteErr.ErrCode != sError.ErrDBConnectionError {
// 	tPtr.Errorf("Get Tables Failed: Expected error code of %v",sError.ErrDBConnectionError)
// 	tPtr.Fail()
// }
//
// if soteErr := sDatabase.GetAWSParams(); soteErr.ErrCode != nil {
// 	tPtr.Errorf("GetAWSParams Failed: Expected error code to be nil.")
// 	t.Fatal()
// }
//
// tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.Name, sDatabase.User, sDatabase.Password, sDatabase.Host, sDatabase.SSLMode, sDatabase.Port, 1)
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
// if _, soteErr := sDatabase.GetColumnInfo(SOTETESTSCHEMA, REFERENCETABLE, tConnInfo); soteErr.ErrCode != sError.ErrDBConnectionError {
// 	tPtr.Errorf("GetColumnInfo Failed: Expected error code of %v",sError.ErrDBConnectionError)
// 	tPtr.Fail()
// }
//
// if soteErr := sDatabase.GetAWSParams(); soteErr.ErrCode != nil {
// 	tPtr.Errorf("getAWSParams Failed: Expected error code to be nil.")
// 	t.Fatal()
// }
//
// tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.Name, sDatabase.User, sDatabase.Password, sDatabase.Host, sDatabase.SSLMode, sDatabase.Port, 1)
// if soteErr.ErrCode != nil {
// 	tPtr.Errorf("GetConnection Failed: Please Investigate")
// 	tPtr.Fail()
// }
//
// var columnInfo []sDatabase.ColumnInfo
// if columnInfo, soteErr = sDatabase.GetColumnInfo(SOTETESTSCHEMA, REFERENCETABLE, tConnInfo); soteErr.ErrCode != nil {
// 	tPtr.Errorf("GetTableList Failed: Expected error code to be nil")
// 	tPtr.Fail()
// }
//
// if len(columnInfo) == 0 {
// 	tPtr.Errorf("GetColumnInfo Failed: Expected at least one column's info to be returned")
// 	tPtr.Fail()
// } else {
// 	if columnInfo[0].Name == "" {
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
// 	tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.Name, sDatabase.User, sDatabase.Password, sDatabase.Host, sDatabase.SSLMode, sDatabase.Port, 3)
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
// 	tConnInfo, soteErr := sDatabase.GetConnection(sDatabase.Name, sDatabase.User, sDatabase.Password, sDatabase.Host, sDatabase.SSLMode, sDatabase.Port, 3)
// 	if soteErr.ErrCode != nil {
// 		tPtr.Errorf("GetConnection Failed: Please investigate")
// 		tPtr.Fail()
// 	}
//
// 	if _, soteErr := sDatabase.PKLookup(EMPTYVALUE, PARENTCHILDTABLE, PCTBLCOLUMNNAME, tConnInfo); soteErr.ErrCode != sError.ErrMissingParameters {
// 		tPtr.Errorf("pkLookup Failed: Expected error code to be %v.",sError.ErrMissingParameters)
// 	}
// 	if _, soteErr := sDatabase.PKLookup(SOTETESTSCHEMA, EMPTYVALUE, PCTBLCOLUMNNAME, tConnInfo); soteErr.ErrCode != sError.ErrMissingParameters {
// 		tPtr.Errorf("pkLookup Failed: Expected error code to be %v.",sError.ErrMissingParameters)
// 	}
// 	if _, soteErr := sDatabase.PKLookup(SOTETESTSCHEMA, PARENTCHILDTABLE, EMPTYVALUE, tConnInfo); soteErr.ErrCode != sError.ErrMissingParameters {
// 		tPtr.Errorf("pkLookup Failed: Expected error code to be %v.",sError.ErrMissingParameters)
// 	}
// }
