package sDatabase

import (
	"runtime"
	"strconv"
	"testing"

	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
)

const (
	TESTINFOSCHEMA  = "information_schema"
	INFOSCHEMATABLE = "columns"
	VALIDDATATYPE   = STRING
	INVALIDDATATYPE = "invalid"
)

func init() {
	sLogger.SetLogMessagePrefix("columns_test.go")
}

func TestGetColumnInfo(tPtr *testing.T) {
	var (
		soteErr           sError.SoteError
		tConnInfo         ConnInfo
		columnInfo        []ColumnInfo
		columnInfoJSON    []byte
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	if _, soteErr = GetColumnInfo(TESTINFOSCHEMA, INFOSCHEMATABLE, tConnInfo); soteErr.ErrCode != sError.ErrDBConnectionError {
		tPtr.Errorf("%v Failed: Expected error code of %v", testName, sError.ErrDBConnectionError)
		tPtr.Fail()
	}

	if tConnInfo, soteErr = getMyDBConn(tPtr); soteErr.ErrCode == nil {
		if _, soteErr = GetColumnInfo("", INFOSCHEMATABLE, tConnInfo); soteErr.ErrCode != sError.ErrMissingParameters {
			tPtr.Errorf("%v Failed: Expected error code of %v", testName, sError.ErrMissingParameters)
			tPtr.Fail()
		}
		if _, soteErr = GetColumnInfo(TESTINFOSCHEMA, "", tConnInfo); soteErr.ErrCode != sError.ErrMissingParameters {
			tPtr.Errorf("%v Failed: Expected error code of %v", testName, sError.ErrMissingParameters)
			tPtr.Fail()
		}
		if columnInfo, soteErr = GetColumnInfo(TESTINFOSCHEMA, INFOSCHEMATABLE, tConnInfo); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil [%v]", testName, strconv.Itoa(soteErr.ErrCode.(int)))
			tPtr.Fail()
		}
		if len(columnInfo) == 0 {
			tPtr.Errorf("%v Failed: Expected at least one column's info to be returned", testName)
			tPtr.Fail()
		} else {
			if columnInfo[0].Name == "" {
				tPtr.Errorf("%v Failed: Expected the column name to be returned", testName)
				tPtr.Fail()
			}
		}
		if columnInfoJSON, soteErr = GetColumnInfoJSONFormat(TESTINFOSCHEMA, INFOSCHEMATABLE, tConnInfo); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil [%v]", testName, strconv.Itoa(soteErr.ErrCode.(int)))
			tPtr.Fail()
		}
		if string(columnInfoJSON) == "" {
			tPtr.Errorf("%v Failed: Expected at least one column's info to be returned", testName)
			tPtr.Fail()
		}
	}

	tConnInfo.DBPoolPtr.Close()
}
func TestValidateDataType(tPtr *testing.T) {
	var (
		soteErr           sError.SoteError
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	tPtr.Run("Valid type", func(tPtr *testing.T) {
		if soteErr = validateDataType(VALIDDATATYPE); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.ErrCode)
		}
	})

	tPtr.Run("Invalid type", func(tPtr *testing.T) {
		if soteErr = validateDataType(INVALIDDATATYPE); soteErr.ErrCode != sError.ErrInvalidSQLDataType {
			tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.ErrCode)
		}
	})
}

func TestConnInfo_AddColumns(tPtr *testing.T) {
	var (
		soteErr           sError.SoteError
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		tableName         = "test_table2"
		columnInfo        = []ColumnInfo{
			{
				Name:       "name",
				DataType:   STRING,
				Length:     20,
				IsNullable: false,
			},
			{
				Name:       "age",
				DataType:   INTEGER,
				IsNullable: false,
			},
		}
		primaryKeyInfo = &PrimaryKeyInfo{
			Columns: []string{"test_id"},
			AutoIncrementInfo: &AutoIncrementInfo{
				IsAutoIncrement:       true,
				AutoIncrementStartBy:  0,
				AutoIncrementInterval: 0,
			},
			Description: "this is the unique identifier",
		}
	)

	tPtr.Cleanup(func() {
		testConnInfo.DropTable(parentCtx, tableName, MigrationType)
	})

	createTable(tPtr, tableName, columnInfo, primaryKeyInfo)
	if soteErr = testConnInfo.AddColumns(parentCtx, tableName, []ColumnInfo{
		{
			Name:        "school",
			DataType:    CHARACTERVARYING,
			Description: "this describes school name",
			IsNullable:  true,
		},
		{
			Name:        "home",
			DataType:    CHARACTERVARYING,
			Description: "this describes school name",
			IsNullable:  true,
		},
	}, MigrationType); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code %v got %v", testName, "nil", soteErr.FmtErrMsg)
	}
}

func TestConnInfo_DropColumns(tPtr *testing.T) {
	var (
		soteErr           sError.SoteError
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		tableName         = "test_table2"
		columnInfo        = []ColumnInfo{
			{
				Name:       "name",
				DataType:   STRING,
				Length:     20,
				IsNullable: false,
			},
			{
				Name:       "age",
				DataType:   INTEGER,
				IsNullable: false,
			},
		}
		primaryKeyInfo = &PrimaryKeyInfo{
			Columns: []string{"test_id"},
			AutoIncrementInfo: &AutoIncrementInfo{
				IsAutoIncrement:       true,
				AutoIncrementStartBy:  0,
				AutoIncrementInterval: 0,
			},
			Description: "this is the unique identifier",
		}
	)

	tPtr.Cleanup(func() {
		testConnInfo.DropTable(parentCtx, tableName, MigrationType)
	})

	createTable(tPtr, tableName, columnInfo, primaryKeyInfo)
	if soteErr = testConnInfo.DropColumns(parentCtx, tableName, []string{"Name"}, MigrationType); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code %v got %v", testName, "nil", soteErr.FmtErrMsg)
	}

}

func TestConnInfo_RenameColumns(tPtr *testing.T) {
	var (
		soteErr           sError.SoteError
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		tableName         = "test_table2"
		columnInfo        = []ColumnInfo{
			{
				Name:       "name",
				DataType:   STRING,
				Length:     20,
				IsNullable: false,
			},
			{
				Name:       "age",
				DataType:   INTEGER,
				IsNullable: false,
			},
		}
		primaryKeyInfo = &PrimaryKeyInfo{
			Columns: []string{"test_id"},
			AutoIncrementInfo: &AutoIncrementInfo{
				IsAutoIncrement:       true,
				AutoIncrementStartBy:  0,
				AutoIncrementInterval: 0,
			},
			Description: "this is the unique identifier",
		}
	)

	tPtr.Cleanup(func() {
		testConnInfo.DropTable(parentCtx, tableName, MigrationType)
	})

	createTable(tPtr, tableName, columnInfo, primaryKeyInfo)
	if soteErr = testConnInfo.RenameColumns(parentCtx, tableName, []ObjRename{
		{
			OldName: "name",
			NewName: "new_name",
		},
	}, MigrationType); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code %v got %v", testName, "nil", soteErr.FmtErrMsg)
	}
}

// HELPER FUNCTIONS

func createTable(tPtr *testing.T, tableName string, columnInfo []ColumnInfo, primaryKeyInfo *PrimaryKeyInfo) {
	tPtr.Helper()
	var (
		soteErr           sError.SoteError
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	testConnInfo.DropTable(parentCtx, tableName, MigrationType)
	tableInfo := TableInfo{
		Name:        tableName,
		PrimaryKey:  primaryKeyInfo,
		Description: "test table",
	}

	if soteErr = testConnInfo.CreateTable(parentCtx, tableInfo, columnInfo, MigrationType); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code %v got %v", testName, "nil", soteErr.FmtErrMsg)
	}
}
