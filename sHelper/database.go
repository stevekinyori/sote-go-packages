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
	Columns []string
	Join    string
	Where   string
	OrderBy string
	GroupBy string
	sql     *bytes.Buffer
}

type DatabaseHelper struct {
	dbConnInfo sDatabase.ConnInfo
	run        *Run
	query      func(sql string) (sDatabase.SRows, error)
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
				query: func(sql string) (sDatabase.SRows, error) {
					return dbConnInfo.DBPoolPtr.Query(run.dbHelper.dbConnInfo.DBContext, sql)
				},
			}
		}
	}
	return
}

func (q Query) Exec(r *Run) (sDatabase.SRows, sError.SoteError) {
	sLogger.DebugMethod()
	q.sql.WriteString(" FROM sote." + q.Table)
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
	tRows, err := r.dbHelper.query(sql)
	return tRows, q.GetError(err)
}

func (q Query) Select() Query {
	sLogger.DebugMethod()
	q.sql = bytes.NewBufferString("SELECT ")
	if len(q.Columns) == 0 {
		q.sql.WriteString("*")
	} else {
		q.sql.WriteString(strings.Join(q.Columns, ", "))
	}
	return q
}

func (q Query) GetError(err error) (soteErr sError.SoteError) {
	if err != nil {
		soteErr = NewError().SqlError(fmt.Sprint(err))
	}
	return
}
