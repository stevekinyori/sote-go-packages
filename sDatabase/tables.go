package sDatabase

import (
	"context"
	"fmt"
	"strings"

	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
)

type TableInfo struct {
	Name        string
	PrimaryKey  *PrimaryKeyInfo
	Description any
}

type PrimaryKeyInfo struct {
	Columns []string
	*AutoIncrementInfo
	Description any
	tableName   string
}

type AutoIncrementInfo struct {
	IsAutoIncrement       bool
	AutoIncrementStartBy  int
	AutoIncrementInterval int
}

type TableData struct {
	Data map[string]any // map key is the column name and value is the column data value
}

// CreateTable creates new table in the database
func (dbConnInfo *ConnInfo) CreateTable(ctx context.Context, tableInfo TableInfo, columns []ColumnInfo, errorKey string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	soteErr = dbConnInfo.createTable(ctx, tableInfo, columns, false, errorKey)

	return
}

// CreateTableIfNotExists creates new table if it doesn't exist in the database
func (dbConnInfo *ConnInfo) CreateTableIfNotExists(ctx context.Context, tableInfo TableInfo, columns []ColumnInfo,
	errorKey string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	soteErr = dbConnInfo.createTable(ctx, tableInfo, columns, true, errorKey)

	return
}

// DropTable deletes a table from the database if it exists
func (dbConnInfo *ConnInfo) DropTable(ctx context.Context, tableName string, errorKey string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	soteErr = dbConnInfo.ExecDBStmt(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %v", tableName), errorKey)

	return
}

// RenameTable renames a table if it exists
func (dbConnInfo *ConnInfo) RenameTable(ctx context.Context, table ObjRename, errorKey string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	soteErr = dbConnInfo.ExecDBStmt(ctx, fmt.Sprintf("ALTER TABLE IF EXISTS %v RENAME TO %v;", table.OldName, table.NewName), errorKey)

	return
}

// HasTable checks if a table exists
func (dbConnInfo *ConnInfo) HasTable(ctx context.Context, tableName string, errorKey string) (exists bool, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if _, soteErr = dbConnInfo.QueryOneColumnStmt(ctx, "select 1 from information_schema.tables where table_name=$1", errorKey,
		tableName); soteErr.ErrCode == nil {
		exists = true
	} else if soteErr.ErrCode == sError.ErrItemNotFound {
		soteErr = sError.SoteError{}
	}

	return
}

// InsertTableData inserts data to a database table
//  var data the map key is the column name and value is the column data value
func (dbConnInfo *ConnInfo) InsertTableData(ctx context.Context, tableName string, data []map[string]any,
	errorKey string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	for _, d := range data {
		if len(d) > 0 {
			var (
				params = make([]any, 0)
				kQstmt string
				vQstmt string
				i      int
			)

			for k, v := range d {
				i++
				kQstmt += fmt.Sprintf("%v,", k)
				vQstmt += fmt.Sprintf("$%v,", i)
				params = append(params, v)
			}

			qStmt := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", tableName, strings.TrimSuffix(kQstmt, ","), strings.TrimSuffix(vQstmt, ","))
			if soteErr = dbConnInfo.ExecDBStmt(ctx, qStmt, errorKey, params...); soteErr.ErrCode != nil {
				return
			}
		}
	}

	return
}

// internal function for creating a table
func (dbConnInfo *ConnInfo) createTable(ctx context.Context, tableInfo TableInfo, columns []ColumnInfo,
	ignoreTableIfExists bool, errorKey string) (soteErr sError.SoteError) {
	var (
		queries          = make(map[string][]Query, 0)
		pkQueries        = make(map[string][]Query, 0)
		colQueries       = make(map[string][]Query, 0)
		initialColParams = 0
		commaPrefix      bool
	)

	if tableInfo.Name == "" {
		soteErr = sError.GetSError(sError.ErrItemNotFound, sError.BuildParams([]string{"Table Name"}), sError.EmptyMap)
		return
	}

	if tableInfo.PrimaryKey != nil {
		tableInfo.PrimaryKey.tableName = tableInfo.Name
		if pkQueries, initialColParams, soteErr = setPrimaryKeyQuery(false, commaPrefix, initialColParams,
			*tableInfo.PrimaryKey, errorKey); soteErr.ErrCode != nil {
			return
		}

		if len(tableInfo.PrimaryKey.Columns) > 0 {
			commaPrefix = true
		}

	}

	if colQueries, initialColParams, soteErr = setAddColumnsQuery(false, ColumnsAdd{
		commaPrefix:       commaPrefix,
		tInitialColParams: initialColParams,
		columns:           columns,
		tableName:         tableInfo.Name,
	}, errorKey); soteErr.ErrCode != nil {
		return
	}

	if len(pkQueries[ColumnKey]) == 0 && len(colQueries[ColumnKey]) == 0 {
		soteErr = sError.GetSError(sError.ErrItemNotFound, sError.BuildParams([]string{"Columns"}), sError.EmptyMap)
		return
	}
	// append primary key queries
	for k, v := range pkQueries {
		queries[k] = append(queries[k], v...)
	}

	// append column queries
	for k, v := range colQueries {
		queries[k] = append(queries[k], v...)
	}

	queries[ColumnKey] = []Query{setCreateTableQuery(tableInfo.Name, queries[ColumnKey], ignoreTableIfExists, errorKey)}
	if tableInfo.Description != nil {
		if _, ok := queries[CommentKey]; !ok {
			queries[CommentKey] = make([]Query, 0)
		}

		queries[CommentKey] = append([]Query{setCommentQuery("TABLE", tableInfo.Name, tableInfo.Description, errorKey)}, queries[CommentKey]...)
	}

	soteErr = dbConnInfo.runTableQueries(ctx, queries)

	return
}

// sets query format for adding a primary key to a table
func setPrimaryKeyQuery(isExistingTable bool, commaPrefix bool, tInitialColParams int, primaryKey PrimaryKeyInfo,
	errorKey string) (queries map[string][]Query, initialColParams int, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		prefix   string
		tQueries = make(map[string][]Query, 0)
	)
	if len(primaryKey.Columns) > 0 {
		queries = make(map[string][]Query, 0)
		if isExistingTable {
			prefix = "ADD "
		}

		queries[ColumnKey] = make([]Query, 0)
		if primaryKey.AutoIncrementInfo != nil && primaryKey.IsAutoIncrement {
			queries[SequenceKey] = make([]Query, 0)
			var (
				seq       = fmt.Sprintf("%v_%v_seq", primaryKey.tableName, primaryKey.Columns[0])
				start     = 1
				increment = 1
			)

			if primaryKey.AutoIncrementStartBy > 0 {
				start = primaryKey.AutoIncrementStartBy
			}

			if primaryKey.AutoIncrementInterval > 0 {
				increment = primaryKey.AutoIncrementInterval
			}

			queries[SequenceKey] = append(queries[SequenceKey], Query{
				Statement: fmt.Sprintf("CREATE SEQUENCE IF NOT EXISTS %v START %v INCREMENT %v;\n", seq, start, increment),
				ErrorKey:  errorKey,
			})

			if tQueries, tInitialColParams, soteErr = setAddColumnsQuery(isExistingTable, ColumnsAdd{
				commaPrefix:       commaPrefix,
				tInitialColParams: tInitialColParams,
				columns: []ColumnInfo{
					{
						Name:            primaryKey.Columns[0],
						DataType:        BIGINTEGER,
						IsNullable:      false,
						Default:         fmt.Sprintf("NEXTVAL('%v'::REGCLASS)", seq),
						Description:     primaryKey.Description,
						IsAutoIncrement: true,
					},
				},
				tableName: primaryKey.tableName,
			}, errorKey); soteErr.ErrCode != nil {
				return
			}

			for k, v := range tQueries {
				queries[k] = append(queries[k], v...)
			}

			primaryKey.Columns = primaryKey.Columns[:1]
		} else if commaPrefix {
			prefix = ",\n" + prefix
		}

		if tInitialColParams > 0 {
			prefix = ",\n" + prefix
		} else if v, ok := queries[ColumnKey]; ok && len(v) > 0 {
			prefix = ",\n" + prefix
		}

		queries[ColumnKey] = append(queries[ColumnKey], Query{
			Statement: fmt.Sprintf("%vPRIMARY KEY (%v)", prefix, strings.Join(primaryKey.Columns, ",")),
			Params:    nil,
			ErrorKey:  errorKey,
		})
	}

	initialColParams = tInitialColParams

	return
}

// generates query details for creating a table
func setCreateTableQuery(tableName string, tQueries []Query, ignoreTableIfExists bool, errorKey string) (query Query) {
	sLogger.DebugMethod()

	var (
		columnQuery  string
		columnParams = make([]any, 0)
	)

	if ignoreTableIfExists {
		columnQuery = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (", tableName)
	} else {
		columnQuery = fmt.Sprintf("CREATE TABLE %v (", tableName)
	}

	for _, v := range tQueries {
		columnQuery += v.Statement
		columnParams = append(columnParams, v.Params...)
	}

	if columnQuery != "" {
		query = Query{
			Statement: columnQuery + ");",
			Params:    columnParams,
			ErrorKey:  errorKey,
		}
	}

	return
}

// runs the prepared queries
func (dbConnInfo *ConnInfo) runTableQueries(ctx context.Context, queries map[string][]Query) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = dbConnInfo.ExecDBStmts(ctx, queries[SequenceKey]); soteErr.ErrCode != nil {
		return
	}

	if soteErr = dbConnInfo.ExecDBStmts(ctx, queries[ColumnKey]); soteErr.ErrCode != nil {
		return
	}

	delete(queries, SequenceKey)
	delete(queries, ColumnKey)

	for _, q := range queries {
		if soteErr = dbConnInfo.ExecDBStmts(ctx, q); soteErr.ErrCode != nil {
			return
		}
	}

	return
}
