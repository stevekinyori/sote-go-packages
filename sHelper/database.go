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
	Sql           *bytes.Buffer
	Filter        *FilterHeaderSchema
	Result        QueryResult
	Table         string
	Schema        string
	Columns       []string
	Values        []interface{}
	Join          string
	Where         string
	OrderBy       string
	GroupBy       string
	Limit         *int64
	Offset        *int64
	action        string
	returnColumns []string
}

type DatabaseHelper struct {
	dbConnInfo sDatabase.ConnInfo
	run        *Run
	query      func(sql string, args ...interface{}) (sDatabase.SRows, error)
}

type QueryResult struct {
	Items      []interface{} `json:"items"`
	Pagination *Pagination   `json:"pagination"`
}

type Pagination struct {
	Total  int64 `json:"total"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
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

func (q Query) Pagination() Query {
	q.Result.Pagination = &Pagination{}
	return q
}

func (q Query) Exec(r *Run) (sDatabase.SRows, sError.SoteError) {
	sLogger.DebugMethod()
	if q.action == "SELECT" {
		q.Sql.WriteString(" FROM " + getTable(&q))
	} else if q.action == "INSERT" || q.action == "UPDATE" {
		if len(q.Columns) > 0 && len(q.Columns) != len(q.Values) {
			return nil, NewError().SqlError("the number of columns in the query does not match the number of values")
		}
	}
	if q.Join != "" {
		q.Sql.WriteString(" " + q.Join)
	}
	if q.Where != "" {
		q.Sql.WriteString(" WHERE " + q.Where)
	}
	if q.GroupBy != "" {
		q.Sql.WriteString(" GROUP BY " + q.GroupBy)
	}
	if q.OrderBy != "" {
		q.Sql.WriteString(" ORDER BY " + q.OrderBy)
	}
	if q.Limit != nil {
		q.Sql.WriteString(fmt.Sprintf(" LIMIT %v", *q.Limit))
	}
	if q.Offset != nil {
		q.Sql.WriteString(fmt.Sprintf(" OFFSET %v", *q.Offset))
	}
	if len(q.returnColumns) > 0 {
		q.Sql.WriteString(" RETURNING " + strings.Join(q.returnColumns, ", "))
	}
	sql := q.Sql.String()
	sLogger.Info("Database::Exec - " + sql)
	tRows, err := r.dbHelper.query(sql, q.Values...)
	return tRows, q.GetError(err)
}

func getTable(q *Query) string {
	if q.Schema == "" {
		return "sote." + q.Table
	} else {
		return fmt.Sprintf("%v.%v", q.Schema, q.Table)
	}
}

func values(q *Query) []string {
	values := make([]string, len(q.Values))
	for i := 0; i < len(q.Values); i++ {
		values[i] = fmt.Sprintf("$%v", i+1)
	}
	return values
}

func (q *Query) where(op string, obj map[string]interface{}) {
	for name, val := range obj {
		if q.Where != "" {
			q.Where += " AND "
		}
		switch val.(type) {
		case string:
			q.Where += fmt.Sprintf("%v %v '%v'", name, op, val)
		default:
			q.Where += fmt.Sprintf("%v %v %v", name, op, val)
		}
	}
}

func (q Query) Select() Query {
	sLogger.DebugMethod()
	q.action = "SELECT"
	q.Sql = bytes.NewBufferString("SELECT ")
	if q.Filter != nil {
		if q.Result.Pagination != nil {
			q.Sql.WriteString("count(*) OVER(), ")
		}
		q.Sql.WriteString(strings.Join(q.Filter.Items, ", "))
		q.GroupBy = strings.Join(q.Filter.GroupBy, ", ")
		if len(q.Filter.SortAsc) > 0 {
			q.OrderBy = strings.Join(q.Filter.SortAsc, ", ") + " ASC"
		}
		if len(q.Filter.SortDesc) > 0 {
			if q.OrderBy != "" {
				q.OrderBy += ", "
			}
			q.OrderBy += strings.Join(q.Filter.SortDesc, ", ") + " DESC"
		}
		if q.Filter.Limit != nil {
			q.Limit = q.Filter.Limit
			if q.Result.Pagination != nil {
				q.Result.Pagination.Limit = *q.Limit
			}
		}
		if q.Filter.Offset != nil {
			q.Offset = q.Filter.Offset
			if q.Result.Pagination != nil {
				q.Result.Pagination.Offset = *q.Offset
			}
		}
		q.where("=", q.Filter.Equal)
		q.where("<", q.Filter.Less)
		q.where(">", q.Filter.Greater)
	} else if len(q.Columns) == 0 {
		q.Sql.WriteString("*")
	} else {
		q.Sql.WriteString(strings.Join(q.Columns, ", "))
	}
	return q
}

func (q Query) Update(returnColumns ...string) Query {
	sLogger.DebugMethod()
	q.action = "UPDATE"
	q.returnColumns = returnColumns
	q.Sql = bytes.NewBufferString("UPDATE " + getTable(&q) + " SET ")
	total := len(q.Columns)
	if total == len(q.Values) {
		values := values(&q)
		for i, name := range q.Columns {
			q.Sql.WriteString(fmt.Sprintf("%v = %v", name, values[i]))
			if i+1 < total {
				q.Sql.WriteString(", ")
			}
		}
	}
	return q
}

func (q Query) Insert(returnColumns ...string) Query {
	sLogger.DebugMethod()
	q.action = "INSERT"
	q.returnColumns = returnColumns
	q.Sql = bytes.NewBufferString("INSERT INTO " + getTable(&q))
	if len(q.Columns) > 0 {
		q.Sql.WriteString(fmt.Sprintf(" (%v)", strings.Join(q.Columns, ", ")))
	}
	q.Sql.WriteString(fmt.Sprintf(" VALUES(%v)", strings.Join(values(&q), ", ")))
	return q
}

func (q Query) Delete(returnColumns ...string) Query {
	sLogger.DebugMethod()
	q.action = "DELETE"
	q.returnColumns = returnColumns
	q.Sql = bytes.NewBufferString("DELETE FROM " + getTable(&q))
	return q
}

func (q Query) GetError(err error) (soteErr sError.SoteError) {
	if err != nil {
		soteErr = NewError().SqlError(fmt.Sprint(err))
	}
	return
}

func (q *Query) Scan(tRows sDatabase.SRows) (tCols []interface{}, soteErr sError.SoteError) {
	var (
		err error
	)
	tCols, err = tRows.Values()
	if err == nil {
		offset := 0
		if q.Result.Pagination != nil {
			offset = 1
			q.Result.Pagination.Total = tCols[0].(int64)
		}
		row := make(map[string]interface{})
		names := q.Columns
		if q.Filter != nil {
			names = q.Filter.Items
		}
		for i := offset; i < len(tCols) && i < len(names)+offset; i++ { //0 - total
			name := names[i-offset]
			row[name] = tCols[i]
		}
		q.Result.Items = append(q.Result.Items, row)
	} else {
		soteErr = NewError().SqlError(err.Error())
	}
	return
}

func (q Query) Close(tRows sDatabase.SRows, soteErr *sError.SoteError) {
	tRows.Close()
	if soteErr == nil || soteErr.ErrCode == nil {
		err := tRows.Err()
		if err != nil {
			*soteErr = NewError().SqlError(err.Error())
		}
	}
}
