package sHelper

import (
	"bytes"
	"fmt"
	"strings"

	"gitlab.com/soteapps/packages/v2021/sDatabase"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

type Query struct {
	Table   string
	Schema  string
	Columns []string
	Values  []interface{}
	Join    string
	Where   string
	OrderBy string
	GroupBy string
	action  string
	sql     *bytes.Buffer
}

type DatabaseHelper struct {
	dbConnInfo sDatabase.ConnInfo
	run        *Run
	query      func(sql string, args ...interface{}) (sDatabase.SRows, error)
}

func NewDatabase(run *Run) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var (
		dbConnInfo sDatabase.ConnInfo
	)
	if soteErr = sDatabase.GetAWSParams(); soteErr.ErrCode == nil {
		if dbConnInfo, soteErr = run.GetConnection(sDatabase.DBName, sDatabase.DBUser, sDatabase.DBPassword, sDatabase.DBHost,
			sDatabase.DBSSLMode, sDatabase.DBPort, 3); soteErr.ErrCode == nil {
			run.dbHelper = &DatabaseHelper{
				run:        run,
				dbConnInfo: dbConnInfo,
				query: func(sql string, args ...interface{}) (sDatabase.SRows, error) {
					return dbConnInfo.DBPoolPtr.Query(run.dbHelper.dbConnInfo.DBContext, sql, args...)
				},
			}
		}
	}
	return
}

func (q Query) Exec(r *Run) (sDatabase.SRows, sError.SoteError) {
	sLogger.DebugMethod()
	if q.action == "SELECT" {
		q.sql.WriteString(" FROM " + getTable(q))
	} else if q.action == "INSERT" || q.action == "UPDATE" {
		if len(q.Columns) > 0 && len(q.Columns) != len(q.Values) {
			return nil, NewError().SqlError("the number of columns in the query does not match the number of values")
		}
	}
	if q.Join != "" {
		q.sql.WriteString(" " + q.Join)
	}
	if q.Where != "" {
		q.sql.WriteString(" WHERE " + q.Where)
	}
	if q.GroupBy != "" {
		q.sql.WriteString(" GROUP BY " + q.GroupBy)
	}
	if q.OrderBy != "" {
		q.sql.WriteString(" ORDER BY " + q.OrderBy)
	}
	sql := q.sql.String()
	sLogger.Info("Database::Exec - " + sql)
	tRows, err := r.dbHelper.query(sql, q.Values...)
	return tRows, q.GetError(err)
}

func getTable(q Query) string {
	if q.Schema == "" {
		return "sote." + q.Table
	} else {
		return fmt.Sprintf("%v.%v", q.Schema, q.Table)
	}
}

func values(q Query) []string {
	values := make([]string, len(q.Values))
	for i := 0; i < len(q.Values); i++ {
		values[i] = fmt.Sprintf("$%v", i+1)
	}
	return values
}

func (q Query) Select() Query {
	sLogger.DebugMethod()
	q.action = "SELECT"
	q.sql = bytes.NewBufferString("SELECT ")
	if len(q.Columns) == 0 {
		q.sql.WriteString("*")
	} else {
		q.sql.WriteString(strings.Join(q.Columns, ", "))
	}
	return q
}

func (q Query) Update() Query {
	sLogger.DebugMethod()
	q.action = "UPDATE"
	q.sql = bytes.NewBufferString("UPDATE " + getTable(q) + " SET ")
	total := len(q.Columns)
	if total == len(q.Values) {
		values := values(q)
		for i, name := range q.Columns {
			q.sql.WriteString(fmt.Sprintf("%v = %v", name, values[i]))
			if i+1 < total {
				q.sql.WriteString(", ")
			}
		}
	}
	return q
}

func (q Query) Insert(returnColumns ...string) Query {
	sLogger.DebugMethod()
	q.action = "INSERT"
	q.sql = bytes.NewBufferString("INSERT INTO " + getTable(q))
	if len(q.Columns) > 0 {
		q.sql.WriteString(fmt.Sprintf(" (%v)", strings.Join(q.Columns, ", ")))
	}
	q.sql.WriteString(fmt.Sprintf(" VALUES(%v)", strings.Join(values(q), ", ")))
	if len(returnColumns) > 0 {
		q.sql.WriteString(" RETURNING " + strings.Join(returnColumns, ", "))
	}
	return q
}

func (q Query) Delete() Query {
	sLogger.DebugMethod()
	q.action = "DELETE"
	q.sql = bytes.NewBufferString("DELETE FROM " + getTable(q))
	return q
}

func (q Query) GetError(err error) (soteErr sError.SoteError) {
	if err != nil {
		soteErr = NewError().SqlError(fmt.Sprint(err))
	}
	return
}
