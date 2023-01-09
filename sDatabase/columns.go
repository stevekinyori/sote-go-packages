package sDatabase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v4"
	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
	"golang.org/x/exp/slices"
)

type ColumnInfo struct {
	Name            string
	Default         any
	DataType        string
	Length          int
	Description     any
	IsAutoIncrement bool
	IsNullable      bool
}

type ColumnsAdd struct {
	commaPrefix       bool
	tInitialColParams int
	columns           []ColumnInfo
	tableName         string
}

const (
	STRING             = "varchar"
	TEXT               = "text"
	BIGINTEGER         = "bigint"
	INTEGER            = "integer"
	BOOLEAN            = "boolean"
	DOUBLE             = "double"
	FLOAT              = "float precision"
	DATE               = "date"
	TIMESTAMP          = "timestamp"
	TIMESTAMPTZ        = "timestamp with time zone"
	CHARACTERVARYING   = "character varying"
	TSVECTOR           = "tsvector"
	CURRENTDATE        = "CURRENT_DATE"
	CURRENTTIMESTAMPTZ = "NOW()" // with timezone
	CURRENTTIMESTAMP   = "CURRENT_TIMESTAMP"
	CURRENTTIME        = "CURRENT_TIME"
)

// GetColumnInfoJSONFormat This function gets column information for the supplied schema and table and returns the data as JSON.
func GetColumnInfoJSONFormat(schemaName, tableName string, tConnInfo ConnInfo) (tableColumnInfoJSON []byte, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		tableColumnInfo []ColumnInfo
		err             error
	)

	if tableColumnInfo, soteErr = GetColumnInfo(schemaName, tableName, tConnInfo); soteErr.ErrCode == nil {
		if tableColumnInfoJSON, err = json.MarshalIndent(tableColumnInfo, sError.PREFIX, sError.INDENT); err != nil {
			sLogger.Info(err.Error())
			soteErr = sError.GetSError(sError.ErrInvalidJSON, sError.BuildParams([]string{"Sote Error"}), sError.EmptyMap)
		}
	}

	return
}

// GetColumnInfo This function gets column information for the supplied schema and table.
func GetColumnInfo(schemaName, tableName string, tConnInfo ConnInfo) (tableColumnInfo []ColumnInfo, soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var (
		tRows       pgx.Rows
		err         error
		tRowCounter int
	)

	if len(schemaName) == 0 {
		soteErr = sError.GetSError(sError.ErrMissingParameters, sError.BuildParams([]string{"Schema name: " + schemaName}), nil)
	}

	if len(tableName) == 0 {
		soteErr = sError.GetSError(sError.ErrMissingParameters, sError.BuildParams([]string{"Schema name: " + tableName}), nil)
	}

	if soteErr.ErrCode == nil {
		if soteErr = VerifyConnection(tConnInfo); soteErr.ErrCode == nil {
			qStmt := "SELECT column_name, column_default, LOWER(is_nullable) AS is_nullable, LOWER(data_type) as data_type," +
				"COALESCE(character_maximum_length,0) AS character_maximum_length FROM information_schema.columns WHERE table_schema = $1 AND table_name = $2;"
			tRows, soteErr = tConnInfo.QueryDBStmt(context.Background(), qStmt, "column info", schemaName, tableName)
			defer tRows.Close()
			if soteErr.ErrCode == nil {
				for tRows.Next() {
					columnRow := new(ColumnInfo)
					nullable := ""
					if err = tRows.Scan(
						&columnRow.Name,
						&columnRow.Default,
						&nullable,
						&columnRow.DataType,
						&columnRow.Length,
					); err != nil {
						log.Fatalln(err)
					}

					if nullable == "yes" {
						columnRow.IsNullable = true
					}

					if v, ok := columnRow.Default.(string); ok && strings.HasPrefix(v, "nextval('") {
						columnRow.IsAutoIncrement = true
					}

					tableColumnInfo = append(tableColumnInfo, *columnRow)
					tRowCounter++
				}

				if tRowCounter == 0 {
					soteErr = sError.GetSError(sError.ErrItemNotFound, sError.BuildParams([]string{schemaName + "." + tableName}), nil)
				}
			}
		}
	}

	return
}

// AddColumns add columns to existing table
func (dbConnInfo *ConnInfo) AddColumns(ctx context.Context, tableName string, columns []ColumnInfo, errorKey string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		queries      = make(map[string][]Query, 0)
		columnQuery  = fmt.Sprintf("ALTER TABLE %v ", tableName)
		columnParams = make([]any, 0)
	)

	if tableName == "" {
		soteErr = sError.GetSError(sError.ErrItemNotFound, sError.BuildParams([]string{"Table Name"}), sError.EmptyMap)
		return
	}

	if queries, _, soteErr = setAddColumnsQuery(true, ColumnsAdd{
		commaPrefix:       false,
		tInitialColParams: 0,
		columns:           columns,
		tableName:         tableName,
	}, errorKey); soteErr.ErrCode != nil {
		return
	}

	if len(queries[ColumnKey]) == 0 {
		soteErr = sError.GetSError(sError.ErrItemNotFound, sError.BuildParams([]string{"Columns"}), sError.EmptyMap)
		return
	}

	for _, v := range queries[ColumnKey] {
		columnQuery += v.Statement
		columnParams = append(columnParams, v.Params...)
	}

	queries[ColumnKey] = []Query{
		{
			Statement: columnQuery,
			Params:    columnParams,
			ErrorKey:  errorKey,
		},
	}

	soteErr = dbConnInfo.runTableQueries(ctx, queries)

	return
}

// DropColumns deletes columns from a table if they exist
func (dbConnInfo *ConnInfo) DropColumns(ctx context.Context, tableName string, columnNames []string, errorKey string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	for _, columnName := range columnNames {
		if soteErr = dbConnInfo.ExecDBStmt(ctx, fmt.Sprintf("ALTER TABLE %v DROP COLUMN IF EXISTS %v", tableName, columnName),
			errorKey); soteErr.ErrCode != nil {
			return
		}
	}

	return
}

// RenameColumns renames columns in a table
func (dbConnInfo *ConnInfo) RenameColumns(ctx context.Context, tableName string, columns []ObjRename, errorKey string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	for _, column := range columns {
		if soteErr = dbConnInfo.ExecDBStmt(ctx, fmt.Sprintf("ALTER TABLE %v RENAME COLUMN %v TO %v;", tableName, column.OldName, column.NewName),
			errorKey); soteErr.ErrCode != nil {
			return
		}
	}

	return
}

// generates query details for adding column(s)
func setAddColumnsQuery(isExistingTable bool, columnInfo ColumnsAdd, errorKey string) (queries map[string][]Query,
	initialColParams int, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		prefix              string
		variableLengthTypes = []string{STRING, CHARACTERVARYING}
		query               Query
		columnsLen          = len(columnInfo.columns)
		literalDefaults     = []string{CURRENTDATE, CURRENTTIMESTAMPTZ, CURRENTTIMESTAMP, CURRENTTIME}
	)

	queries = make(map[string][]Query, 0)
	if columnsLen > 0 {
		if isExistingTable {
			prefix = "ADD COLUMN "
		}

		if columnInfo.commaPrefix {
			query.Statement = ",\n"
		}

		queries[CommentKey] = make([]Query, 0)
		queries[ColumnKey] = make([]Query, 0)
		for k, column := range columnInfo.columns {
			if column.Name == "" {
				soteErr = sError.GetSError(sError.ErrItemNotFound, sError.BuildParams([]string{"Column Name"}), sError.EmptyMap)
				return
			}

			if soteErr = validateDataType(column.DataType); soteErr.ErrCode != nil {
				return
			}

			lengthStr := ""
			if ok := slices.Contains(variableLengthTypes, column.DataType); ok && column.Length > 0 {
				lengthStr = fmt.Sprintf("(%v)", column.Length)
			}

			query.Statement += fmt.Sprintf("%v%v %v%v", prefix, column.Name, column.DataType, lengthStr)
			if !column.IsNullable {
				query.Statement += " NOT NULL"
			}

			if column.Default != nil {
				d, _ := column.Default.(string)
				if column.IsAutoIncrement || slices.Contains(literalDefaults, d) {
					query.Statement += fmt.Sprintf(" DEFAULT %v", column.Default)
				} else {
					query.Statement += fmt.Sprintf(" DEFAULT '%v'", column.Default)
				}
			}

			if k < columnsLen-1 {
				query.Statement += ",\n"
			}

			if column.Description != nil {
				queries[CommentKey] = append(queries[CommentKey], setCommentQuery("COLUMN", columnInfo.tableName+"."+column.Name,
					column.Description, errorKey))
			}
		}

		queries[ColumnKey] = append(queries[ColumnKey], query)
	}

	initialColParams = columnInfo.tInitialColParams

	return
}

// validateDataType validates the SQL data type
func validateDataType(dataType string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var types = []string{STRING, TEXT, BIGINTEGER, INTEGER, BOOLEAN, DOUBLE, FLOAT, DATE, TIMESTAMP, TIMESTAMPTZ, CHARACTERVARYING}

	if !slices.Contains(types, strings.ToLower(dataType)) {
		soteErr = sError.GetSError(sError.ErrInvalidSQLDataType, nil, nil)
		sLogger.Info(soteErr.FmtErrMsg)
	}

	return
}
