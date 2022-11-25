package sMigration

import (
	"runtime"
	"testing"

	"gitlab.com/soteapps/packages/v2022/sConfigParams"
	"gitlab.com/soteapps/packages/v2022/sDatabase"
	"gitlab.com/soteapps/packages/v2022/sError"
)

func TestSeed(tPtr *testing.T) {
	var (
		soteErr           sError.SoteError
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		testConnInfo, _   = sDatabase.New(parentCtx, sConfigParams.DEVELOPMENT)
		tableName         = "sote_student_test"
	)

	tPtr.Cleanup(func() {
		testConnInfo.DropTable(parentCtx, tableName, SeedingType)
	})

	testConnInfo.DropTable(parentCtx, tableName, MigrationType)
	tableInfo := sDatabase.TableInfo{
		Name: tableName,
		PrimaryKey: &sDatabase.PrimaryKeyInfo{
			Columns: []string{"sote_student_test_id"},
			AutoIncrementInfo: &sDatabase.AutoIncrementInfo{
				IsAutoIncrement:       true,
				AutoIncrementInterval: 3,
			},
		},
	}

	columnInfo := []sDatabase.ColumnInfo{
		{
			Name:       "name",
			DataType:   sDatabase.STRING,
			IsNullable: false,
		},
		{
			Name:       "age",
			DataType:   sDatabase.INTEGER,
			IsNullable: false,
		},
		{
			Name:       "class",
			DataType:   sDatabase.STRING,
			IsNullable: false,
		},
	}

	if soteErr = testConnInfo.CreateTable(parentCtx, tableInfo, columnInfo, MigrationType); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code %v got %v", testName, "nil", soteErr.FmtErrMsg)
	}

	if soteErr = Seed(parentCtx, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code %v got %v", testName, "nil", soteErr.FmtErrMsg)
	}
}
