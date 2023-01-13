package sDatabase

import (
	"runtime"
	"testing"

	"gitlab.com/soteapps/packages/v2023/sError"
)

func TestConnInfo_ExecDBStmt(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	tPtr.Run("valid query", func(tPtr *testing.T) {
		if soteErr := testConnInfo.ExecDBStmt(parentCtx, "SELECT 1", "query test"); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil but got %v", testName, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("invalid query", func(tPtr *testing.T) {
		if soteErr := testConnInfo.ExecDBStmt(parentCtx, "FAKE SQL STATEMENT", "query test"); soteErr.ErrCode != sError.ErrSQLError {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, sError.ErrSQLError, soteErr.FmtErrMsg)
		}
	})
}

func TestConnInfo_ExecDBStmts(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	tPtr.Run("valid queries", func(tPtr *testing.T) {
		if soteErr := testConnInfo.ExecDBStmts(parentCtx, []Query{{
			Statement: "SELECT 1",
			ErrorKey:  "query test",
		}}); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil but got %v", testName, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("invalid queries", func(tPtr *testing.T) {
		if soteErr := testConnInfo.ExecDBStmts(parentCtx, []Query{{
			Statement: "FAKE SQL STATEMENT",
			ErrorKey:  "query test",
		}}); soteErr.ErrCode != sError.ErrSQLError {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, sError.ErrSQLError, soteErr.FmtErrMsg)
		}
	})
}

func TestConnInfo_QueryDBStmt(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	tPtr.Run("valid query", func(tPtr *testing.T) {
		if _, soteErr := testConnInfo.QueryDBStmt(parentCtx, "SELECT 1", "query test"); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil but got %v", testName, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("invalid query", func(tPtr *testing.T) {
		if _, soteErr := testConnInfo.QueryDBStmt(parentCtx, "FAKE SQL STATEMENT", "query test"); soteErr.ErrCode != sError.ErrSQLError {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, sError.ErrSQLError, soteErr.FmtErrMsg)
		}
	})
}

func TestConnInfo_QueryOneColumn(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	tPtr.Run("valid query", func(tPtr *testing.T) {
		var resp int
		if soteErr := testConnInfo.QueryOneColumn(parentCtx, "SELECT 1", &resp, "query test"); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil but got %v", testName, soteErr.FmtErrMsg)
		}

		if resp != 1 {
			tPtr.Errorf("%v Failed: Expected column value to be %v got %v", testName, 1, resp)
		}
	})

	tPtr.Run("invalid query", func(tPtr *testing.T) {
		var resp int
		if soteErr := testConnInfo.QueryOneColumn(parentCtx, "FAKE SQL STATEMENT", &resp, "query test"); soteErr.ErrCode != sError.ErrSQLError {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, sError.ErrSQLError, soteErr.FmtErrMsg)
		}

		if resp == 1 {
			tPtr.Errorf("%v Failed: Expected column value not to be %v got %v", testName, 1, resp)
		}
	})
}

func TestConnInfo_QueryOneRow(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	tPtr.Run("valid query", func(tPtr *testing.T) {
		var (
			row   []interface{}
			count int
		)
		if row, soteErr = testConnInfo.QueryOneRow(parentCtx, "SELECT 1", "query test"); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil but got %v", testName, soteErr.FmtErrMsg)
		}

		count = len(row)
		if count != 1 {
			tPtr.Errorf("%v Failed: Expected column value to be %v got %v", testName, 1, count)
		}
	})

	tPtr.Run("invalid query", func(tPtr *testing.T) {
		var (
			row   []interface{}
			count int
		)
		if row, soteErr = testConnInfo.QueryOneRow(parentCtx, "FAKE SQL STATEMENT", "query test"); soteErr.ErrCode != sError.ErrSQLError {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, sError.ErrSQLError, soteErr.FmtErrMsg)
		}

		count = len(row)
		if count != 0 {
			tPtr.Errorf("%v Failed: Expected column value to be %v got %v", testName, 0, count)
		}
	})
}

func TestConnInfo_QueryOneRowWithDest(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	tPtr.Run("valid query", func(tPtr *testing.T) {
		var (
			resp  int
			row   = []interface{}{&resp}
			count int
		)
		if soteErr = testConnInfo.QueryOneRowWithDest(parentCtx, "SELECT 1", row, "query test"); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil but got %v", testName, soteErr.FmtErrMsg)
		}

		count = len(row)
		if count != 1 {
			tPtr.Errorf("%v Failed: Expected column value to be %v got %v", testName, 1, count)
		}
	})

	tPtr.Run("invalid query", func(tPtr *testing.T) {
		var (
			resp int
			row  = []interface{}{&resp}
		)
		if soteErr = testConnInfo.QueryOneRowWithDest(parentCtx, "FAKE SQL STATEMENT", row, "query test"); soteErr.ErrCode != sError.ErrSQLError {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, sError.ErrSQLError, soteErr.FmtErrMsg)
		}

		if resp == 1 {
			tPtr.Errorf("%v Failed: Expected column value not to be %v got %v", testName, 1, resp)
		}
	})
}
