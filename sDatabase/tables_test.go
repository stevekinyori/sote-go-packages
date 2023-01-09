package sDatabase

import (
	"runtime"
	"testing"

	"gitlab.com/soteapps/packages/v2023/sError"
)

func TestConnInfo_CreateTable(tPtr *testing.T) {
	var (
		soteErr           sError.SoteError
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		tableName         = "test_table2"
	)

	tPtr.Cleanup(func() {
		testConnInfo.DropTable(parentCtx, tableName, MigrationType)
	})

	testConnInfo.DropTable(parentCtx, tableName, MigrationType)
	tableInfo := TableInfo{
		Name: "test_table2",
		PrimaryKey: &PrimaryKeyInfo{
			Columns: []string{"test_id"},
			AutoIncrementInfo: &AutoIncrementInfo{
				IsAutoIncrement:       true,
				AutoIncrementStartBy:  0,
				AutoIncrementInterval: 0,
			},
			Description: "this is the unique identifier",
		},
		Description: "test table",
	}

	columnInfo := []ColumnInfo{
		{
			Name:       "name",
			Default:    "mary",
			DataType:   STRING,
			Length:     20,
			IsNullable: false,
		},
		{
			Name:       "age",
			Default:    11,
			DataType:   INTEGER,
			IsNullable: true,
		},
	}

	if soteErr = testConnInfo.CreateTable(parentCtx, tableInfo, columnInfo, MigrationType); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code %v got %v", testName, "nil", soteErr.FmtErrMsg)
	}
}

func TestConnInfo_DropTable(tPtr *testing.T) {
	var (
		soteErr           sError.SoteError
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		tableName         = "test_table2"
		primaryKeyInfo    = &PrimaryKeyInfo{
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

	createTable(tPtr, tableName, nil, primaryKeyInfo)
	if soteErr = testConnInfo.DropTable(parentCtx, tableName, MigrationType); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code %v got %v", testName, "nil", soteErr.FmtErrMsg)
	}

}

func TestConnInfo_RenameTable(tPtr *testing.T) {
	var (
		soteErr           sError.SoteError
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		tableName         = "test_table2"
		primaryKeyInfo    = &PrimaryKeyInfo{
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
		testConnInfo.DropTable(parentCtx, tableName+"_new", MigrationType)
	})

	createTable(tPtr, tableName, nil, primaryKeyInfo)
	if soteErr = testConnInfo.RenameTable(parentCtx, ObjRename{
		OldName: tableName,
		NewName: tableName + "_new",
	}, MigrationType); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code %v got %v", testName, "nil", soteErr.FmtErrMsg)
	}

}

func TestConnInfo_HasTable(tPtr *testing.T) {
	var (
		soteErr           sError.SoteError
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		tableName         = "test_table2"
	)

	if _, soteErr = testConnInfo.HasTable(parentCtx, tableName, MigrationType); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code %v got %v", testName, "nil", soteErr.FmtErrMsg)
	}
}
