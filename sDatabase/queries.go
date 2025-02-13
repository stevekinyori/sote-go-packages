package sDatabase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
	"golang.org/x/exp/slices"
)

const (
	ColumnKey     = "columns"
	SequenceKey   = "sequences"
	CommentKey    = "comments"
	MigrationType = "migration"
	SeedingType   = "seeding"
)

type Query struct {
	Statement string
	Params    []any
	ErrorKey  string
}

type ObjRename struct {
	OldName string
	NewName string
}

// QueryDBStmt query multiple rows
func (dbConnInfo *ConnInfo) QueryDBStmt(ctx context.Context, qStmt string, errorKey string, args ...interface{}) (tRows SRows,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if !slices.Contains([]string{MigrationType, SeedingType}, errorKey) {
		sLogger.Info(fmt.Sprintf("Executing: %s", qStmt))
	}

	tRows, err = dbConnInfo.DBPoolPtr.Query(ctx, qStmt, args...)
	soteErr = convertSQLError(ctx, err, errorKey)

	return
}

// QueryOneColumn query single column to destination(pointer)
func (dbConnInfo *ConnInfo) QueryOneColumn(ctx context.Context, qStmt string, dest interface{}, errorKey string,
	args ...interface{}) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if !slices.Contains([]string{MigrationType, SeedingType}, errorKey) {
		sLogger.Info(fmt.Sprintf("Executing: %s", qStmt))
	}

	err := dbConnInfo.DBPoolPtr.QueryRow(ctx, qStmt, args...).Scan(dest)
	soteErr = convertSQLError(ctx, err, errorKey)

	return
}

// QueryOneRow query single row of n columns (unknown number of columns)
func (dbConnInfo *ConnInfo) QueryOneRow(ctx context.Context, qStmt string, errorKey string, args ...interface{}) (row []interface{},
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		rows        SRows
		tRow        []interface{}
		columnCount int
		err         error
	)

	if rows, soteErr = dbConnInfo.QueryDBStmt(ctx, qStmt, errorKey, args...); soteErr.ErrCode != nil {
		return
	}

	defer rows.Close()
	if rows.Next() {
		columnCount = len(rows.FieldDescriptions())
		row = make([]interface{}, columnCount)
		tRow = make([]interface{}, columnCount)
		for i := range row {
			tRow[i] = &row[i]
		}

		err = rows.Scan(tRow...)
		soteErr = convertSQLError(ctx, err, errorKey)
	}

	return
}

// QueryOneRowWithDest query single row of known number of columns and scan them to the set destination(every slice value must be a pointer)
func (dbConnInfo *ConnInfo) QueryOneRowWithDest(ctx context.Context, qStmt string, dest []interface{}, errorKey string,
	args ...interface{}) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if !slices.Contains([]string{MigrationType, SeedingType}, errorKey) {
		sLogger.Info(fmt.Sprintf("Executing: %s", qStmt))
	}

	err := dbConnInfo.DBPoolPtr.QueryRow(ctx, qStmt, args...).Scan(dest...)
	soteErr = convertSQLError(ctx, err, errorKey)

	return
}

// ExecDBStmt query database without expecting a response other than error
func (dbConnInfo *ConnInfo) ExecDBStmt(ctx context.Context, qStmt string, errorKey string, args ...interface{}) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if !slices.Contains([]string{MigrationType, SeedingType}, errorKey) {
		sLogger.Info(fmt.Sprintf("Executing: %s", qStmt))
	}

	_, err := dbConnInfo.DBPoolPtr.Exec(ctx, qStmt, args...)
	soteErr = convertSQLError(ctx, err, errorKey)

	return
}

// ExecDBStmts executes multiple queries using ExecDBStmt
func (dbConnInfo *ConnInfo) ExecDBStmts(ctx context.Context, queries []Query) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	for _, q := range queries {
		if soteErr = dbConnInfo.ExecDBStmt(ctx, q.Statement, q.ErrorKey, q.Params...); soteErr.ErrCode != nil {
			return
		}
	}

	return
}

// convertSQLError converts SQL errors to Sote Error
func convertSQLError(ctx context.Context, err error, errorKey string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		dbError map[string]string
	)
	if err != nil {
		// if context cancellation signal was issued,log it and give it a special sote Error code
		if errors.Is(ctx.Err(), context.Canceled) {
			sLogger.Info(errorKey + " " + err.Error())
			soteErr = sError.GetSError(sError.ErrContextCancelled, sError.BuildParams([]string{errorKey}), sError.EmptyMap)
		} else if err == pgx.ErrNoRows {
			sLogger.Info(err.Error() + " " + errorKey)
			soteErr = sError.GetSError(sError.ErrItemNotFound, sError.BuildParams([]string{errorKey}), sError.EmptyMap)
		} else {
			sLogger.Info(err.Error())
			if dbError, soteErr = sError.ConvertSQLError(err); len(dbError) == 0 {
				dbError = make(map[string]string)
				dbError["Error"] = err.Error()
			}
			switch dbError["Code"] {
			case "23505":
				soteErr = sError.GetSError(sError.ErrDuplicateItems, sError.BuildParams([]string{"Value"}), dbError)
			default:
				soteErr = sError.GetSError(sError.ErrSQLError, nil, dbError)
			}
		}
	}

	return
}

// sets query format for adding a comment. Can be on table or column
func setCommentQuery(objName string, objValue string, description any, errorKey string) (query Query) {
	sLogger.DebugMethod()

	query = Query{
		Statement: fmt.Sprintf("COMMENT ON %v %v IS '%v';\n", strings.ToUpper(objName), objValue, description),
		ErrorKey:  errorKey,
	}

	return
}
