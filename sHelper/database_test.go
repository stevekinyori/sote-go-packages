package sHelper

import (
	"errors"
	"fmt"
	"os"

	"testing"

	"gitlab.com/soteapps/packages/v2021/sDatabase"
	"gitlab.com/soteapps/packages/v2021/sError"
)

func newDbRun() *Run {
	env, _ := NewEnvironment(ENVDEFAULTAPPNAME, ENVDEFAULTTARGET, ENVDEFAULTTARGET)
	os.Setenv("APP_ENVIRONMENT", env.AppEnvironment)
	return NewRun(env)
}

func createDatabaseHelper(r *Run, result *Result) sError.SoteError {
	r.GetConnection = func(dbName, user, password, host, sslMode string, port, timeout int) (dbConnInfo sDatabase.ConnInfo, soteErr sError.SoteError) {
		if result != nil {
			dbConnInfo = sDatabase.ConnInfo{}
		} else {
			soteErr = NewError().NoDbConnection()
		}
		return
	}
	soteErr := NewDatabase(r)
	if r.dbHelper != nil {
		r.dbHelper.query = func(sql string) (sDatabase.SRows, error) {
			return result.rows, result.err
		}
	}
	return soteErr
}

type Result struct {
	rows sDatabase.SRows
	err  error
}

func TestDatabaseInvalidAppEnv(t *testing.T) {
	run := newDbRun()
	os.Setenv("APP_ENVIRONMENT", "")
	soteErr := NewDatabase(run)
	AssertEqual(t, soteErr.FmtErrMsg, "109999: APP_ENVIRONMENT was/were not found")
	os.Setenv("APP_ENVIRONMENT", "INVALID_APP_ENVIRONMENT")
	soteErr = NewDatabase(run)
	AssertEqual(t, soteErr.FmtErrMsg, "209110: environment value (INVALID_APP_ENVIRONMENT) is invalid")
	os.Setenv("APP_ENVIRONMENT", run.Env.AppEnvironment) //reset
}

func TestDatabaseInvalidConn(t *testing.T) {
	run := newDbRun()
	soteErr := createDatabaseHelper(run, nil)
	AssertEqual(t, soteErr.FmtErrMsg, "209299: No database connection has been established")
}

func TestDatabaseSelectAll(t *testing.T) {
	query := Query{
		Table: "TABLE",
	}.Select()
	AssertEqual(t, query.sql.String(), "SELECT *")
}

func TestDatabaseSelectColumns(t *testing.T) {
	query := Query{
		Table:   "TABLE",
		Columns: []string{"COL1", "COL2", "COL3"},
	}.Select()
	AssertEqual(t, query.sql.String(), "SELECT COL1, COL2, COL3")
}

func TestDatabaseExec(t *testing.T) {
	run := newDbRun()
	soteErr := createDatabaseHelper(run, &Result{err: errors.New("CUSTOM DB ERROR")})
	AssertEqual(t, soteErr.ErrCode, nil)
	tRows, err := Query{
		Table: "TABLE",
	}.Select().Exec(run)
	AssertEqual(t, err.FmtErrMsg, "200999: SQL error - see Details ERROR DETAILS: >>Key: SQL ERROR Value: CUSTOM DB ERROR")
	AssertEqual(t, tRows, nil)
}

func TestDatabaseExecFullQuery(t *testing.T) {
	run := newDbRun()
	createDatabaseHelper(run, &Result{})
	query := Query{
		Table:   "TABLE1",
		Columns: []string{"COL1", "COL2", "COL3"},
		Join:    "INNER JOIN TABLE2 ON TABLE1.ID=TABLE2.CID",
		Where:   "ID IS NOT NULL",
		GroupBy: "NAME",
		OrderBy: "CREATE_DATE",
	}.Select()
	query.Exec(run)
	AssertEqual(t, query.sql.String(), "SELECT COL1, COL2, COL3 FROM sote.TABLE1 INNER JOIN TABLE2 ON TABLE1.ID=TABLE2.CID WHERE ID IS NOT NULL GROUP BY NAME ORDER BY CREATE_DATE")
}

func TestDatabaseQueryPanic(t *testing.T) {
	run := newDbRun()
	run.GetConnection = func(dbName, user, password, host, sslMode string, port, timeout int) (dbConnInfo sDatabase.ConnInfo, soteErr sError.SoteError) {
		dbConnInfo = sDatabase.ConnInfo{}
		return
	}
	NewDatabase(run)
	defer func() {
		r := recover()
		AssertEqual(t, fmt.Sprint(r), "runtime error: invalid memory address or nil pointer dereference")
	}()
	Query{
		Table: "TABLE",
	}.Select().Exec(run)
}
