package sHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"

	"testing"

	"gitlab.com/soteapps/packages/v2021/sDatabase"
	"gitlab.com/soteapps/packages/v2021/sError"
)

type TestPagination struct {
	Filter FilterHeaderSchema `json:"filter-header"`
}

type NoErrorRow struct {
	sDatabase.Rows
}

func (r NoErrorRow) Err() error { return nil }

type ErrorRow struct {
	sDatabase.Rows
}

func (r ErrorRow) Err() error { return errors.New("duplicate key value violates unique constraint") }

type ScanValues struct {
	sDatabase.Rows
}

func (r ScanValues) Values() ([]interface{}, error) {
	var total int64 = 1
	return []interface{}{total, "Hello World", true, 100}, nil
}

type ScanValuesError struct {
	sDatabase.Rows
}

func (r ScanValuesError) Values() ([]interface{}, error) {
	return nil, errors.New("invalid column name")
}

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
		r.dbHelper.query = func(sql string, args ...interface{}) (sDatabase.SRows, error) {
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
	if soteErr.FmtErrMsg != "209299: No database connection has been established" &&
		soteErr.FmtErrMsg != "109999: /sote/api/staging/DB_NAME was/were not found" { //Jenkins
		AssertEqual(t, soteErr.FmtErrMsg, "209299: No database connection has been established")
	}
}

func TestDatabaseSelectAll(t *testing.T) {
	query := Query{
		Table: "TABLE",
	}.Select()
	AssertEqual(t, query.Sql.String(), "SELECT *")
}

func TestDatabaseSelectColumns(t *testing.T) {
	query := Query{
		Table:   "TABLE",
		Columns: []string{"COL1", "COL2", "COL3"},
	}.Select()
	AssertEqual(t, query.Sql.String(), "SELECT COL1, COL2, COL3")
}

func TestDatabaseExec(t *testing.T) {
	run := newDbRun()
	createDatabaseHelper(run, &Result{err: errors.New("CUSTOM DB ERROR")})
	tRows, err := Query{
		Table: "TABLE",
	}.Select().Exec(run)
	AssertEqual(t, err.FmtErrMsg, "200999: SQL error - see Details ERROR DETAILS: >>Key: SQL ERROR Value: CUSTOM DB ERROR")
	AssertEqual(t, tRows, nil)
}

func TestDatabasePaginationExec(t *testing.T) {
	run := newDbRun()
	createDatabaseHelper(run, &Result{})
	var (
		limit  int64 = 1
		offset int64 = 2
	)
	pagination := TestPagination{
		Filter: FilterHeaderSchema{
			Items:    []string{"COL1", "COL2", "COL3"},
			Limit:    &limit,
			Offset:   &offset,
			SortAsc:  []string{"COL1"},
			SortDesc: []string{"COL2"},
			GroupBy:  []string{"COL2", "COL3"},
			Equal:    map[string]interface{}{"COL1": "Hello World"},
			Greater:  map[string]interface{}{"COL2": 20},
			Less:     map[string]interface{}{"COL3": 30},
		},
	}
	query := Query{
		Table:  "TABLE1",
		Filter: &pagination.Filter,
	}.Pagination().Select()
	query.Exec(run)
	AssertEqual(t, query.Sql.String(), "SELECT count(*) OVER(), COL1, COL2, COL3 FROM sote.TABLE1 WHERE COL1 = 'Hello World' AND COL3 < 30 AND COL2 > 20 GROUP BY COL2, COL3 ORDER BY COL1 ASC, COL2 DESC LIMIT 1 OFFSET 2")
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
	AssertEqual(t, query.Sql.String(), "SELECT COL1, COL2, COL3 FROM sote.TABLE1 INNER JOIN TABLE2 ON TABLE1.ID=TABLE2.CID WHERE ID IS NOT NULL GROUP BY NAME ORDER BY CREATE_DATE")
}

func TestDatabaseListScan(t *testing.T) {
	query := Query{
		Filter: &FilterHeaderSchema{
			Items: []string{"_", "COL1", "COL2", "COL3"},
		},
	}
	query.Scan(ScanValues{})
	data, _ := json.MarshalIndent(query.Result.Items, "", "")
	re := regexp.MustCompile(`\r?\n`)
	AssertEqual(t, re.ReplaceAllString(string(data), ""), `[{"COL1": "Hello World","COL2": true,"COL3": 100,"_": 1}]`)
}

func TestDatabasePaginationScan(t *testing.T) {
	query := Query{
		Filter: &FilterHeaderSchema{
			Items: []string{"COL1", "COL2", "COL3"},
		},
		Result: QueryResult{
			Pagination: &Pagination{},
		},
	}
	query.Scan(ScanValues{})
	data, _ := json.MarshalIndent(query.Result, "", "")
	re := regexp.MustCompile(`\r?\n`)
	AssertEqual(t, re.ReplaceAllString(string(data), ""), `{"items": [{"COL1": "Hello World","COL2": true,"COL3": 100}],"pagination": {"total": 1,"limit": 0,"offset": 0}}`)
}

func TestDatabaseScanError(t *testing.T) {
	query := Query{
		Filter: &FilterHeaderSchema{},
	}
	_, soteErr := query.Scan(ScanValuesError{})
	AssertEqual(t, soteErr.FmtErrMsg, "200999: SQL error - see Details ERROR DETAILS: >>Key: SQL ERROR Value: invalid column name")
}

func TestDatabaseUpdate(t *testing.T) {
	run := newDbRun()
	createDatabaseHelper(run, &Result{})
	query := Query{
		Table:   "TABLE1",
		Columns: []string{"COL1", "COL2", "COL3"},
		Values:  []interface{}{"Hello", true, 123.45},
	}.Update("COL1")
	_, soteErr := query.Exec(run)
	AssertEqual(t, soteErr.FmtErrMsg, "")
	AssertEqual(t, query.Sql.String(), "UPDATE sote.TABLE1 SET COL1 = $1, COL2 = $2, COL3 = $3 RETURNING COL1")
}

func TestDatabaseDelete(t *testing.T) {
	run := newDbRun()
	createDatabaseHelper(run, &Result{})
	query := Query{
		Table: "TABLE1",
		Where: "COL1=123",
	}.Delete()
	_, soteErr := query.Exec(run)
	AssertEqual(t, soteErr.FmtErrMsg, "")
	AssertEqual(t, query.Sql.String(), "DELETE FROM sote.TABLE1 WHERE COL1=123")
}

func TestDatabaseInsert(t *testing.T) {
	run := newDbRun()
	createDatabaseHelper(run, &Result{})
	query := Query{
		Table:   "TABLE1",
		Schema:  "myschema",
		Columns: []string{"COL1", "COL2", "COL3"},
		Values:  []interface{}{"Hello", true, 123.45},
	}.Insert("COL1, COL2")
	_, soteErr := query.Exec(run)
	AssertEqual(t, soteErr.FmtErrMsg, "")
	AssertEqual(t, query.Sql.String(), "INSERT INTO myschema.TABLE1 (COL1, COL2, COL3) VALUES($1, $2, $3) RETURNING COL1, COL2")
}

func TestDatabaseInsertByValues(t *testing.T) {
	run := newDbRun()
	createDatabaseHelper(run, &Result{})
	query := Query{
		Table:  "TABLE1",
		Values: []interface{}{"Hello", true, 123.45},
	}.Insert()
	_, soteErr := query.Exec(run)
	AssertEqual(t, soteErr.FmtErrMsg, "")
	AssertEqual(t, query.Sql.String(), "INSERT INTO sote.TABLE1 VALUES($1, $2, $3)")
}

func TestDatabaseInsertError(t *testing.T) {
	run := newDbRun()
	createDatabaseHelper(run, &Result{})
	query := Query{
		Table:   "TABLE1",
		Columns: []string{"COL1", "COL2"},
		Values:  []interface{}{"Hello", true, 123.45},
	}.Insert()
	_, soteErr := query.Exec(run)
	AssertEqual(t, soteErr.FmtErrMsg, "200999: SQL error - see Details ERROR DETAILS: >>Key: SQL ERROR Value: the number of columns in the query does not match the number of values")
}

func TestDatabaseClose(t *testing.T) {
	soteErr := sError.SoteError{}
	query := Query{}
	query.Close(NoErrorRow{}, &soteErr)
	AssertEqual(t, soteErr.FmtErrMsg, "")

	query.Close(ErrorRow{}, &soteErr)
	AssertEqual(t, soteErr.FmtErrMsg, "200999: SQL error - see Details ERROR DETAILS: >>Key: SQL ERROR Value: duplicate key value violates unique constraint")
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
