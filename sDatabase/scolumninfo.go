// All functions use the dbConnPtr or DBPoolPtr which are established using sconnection.
package sDatabase

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v4"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

type SColumnInfo struct {
	// Col was added to the name to prevent conflicts
	ColName     string
	ColDefault  string
	ColNullable string
	ColDataType string
}

const (
	// S was added to the name to prevent conflicts
	SBOOLEAN          = "boolean"
	SDATE             = "date"
	SBIGINT           = "bigint"
	STIMESTAMPZONE    = "timestamp with time zone"
	SCHARACTERVARYING = "character varying"
	SINTEGER          = "integer"
	STEXT             = "text"
)

// This function gets column information for the supplied schema and table and returns the data as JSON.
func GetColumnInfoJSONFormat(schemaName, tableName string, tConnInfo ConnInfo) (tableColumnInfoJSON []byte, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		tableColumnInfo []SColumnInfo
		err             error
	)

	if tableColumnInfo, soteErr = GetColumnInfo(schemaName, tableName, tConnInfo); soteErr.ErrCode == nil {
		if tableColumnInfoJSON, err = json.MarshalIndent(tableColumnInfo, sError.PREFIX, sError.INDENT); err != nil {
			sLogger.Info(err.Error())
			soteErr := sError.GetSError(207110, sError.BuildParams([]string{"Sote Error"}), sError.EmptyMap)
			sError.PanicService(soteErr)
		}
	}

	return
}

// This function gets column information for the supplied schema and table.
func GetColumnInfo(schemaName, tableName string, tConnInfo ConnInfo) (tableColumnInfo []SColumnInfo, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(schemaName) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"Schema name: " + schemaName}), nil)
	}

	if len(tableName) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"Schema name: " + tableName}), nil)
	}

	if soteErr.ErrCode == nil {
		if soteErr = VerifyConnection(tConnInfo); soteErr.ErrCode == nil {
			qStmt := "SELECT column_name, column_default, is_nullable, data_type FROM information_schema.columns WHERE table_schema = $1 AND table_name = $2;"

			var colRows pgx.Rows
			var err error
			colRows, err = tConnInfo.DBPoolPtr.Query(context.Background(), qStmt, schemaName, tableName)
			if err != nil {
				log.Fatalln(err)
			}

			var columnRow []interface{}
			var tRowInfo SColumnInfo
			for colRows.Next() {
				columnRow, err = colRows.Values()
				if err != nil {
					log.Fatalln(err)
				}
				tRowInfo.ColName = columnRow[0].(string)
				if columnRow[1] == nil {
					tRowInfo.ColDefault = ""
				} else {
					tRowInfo.ColDefault = columnRow[1].(string)
				}
				tRowInfo.ColNullable = columnRow[2].(string)
				tRowInfo.ColDataType = columnRow[3].(string)
				tableColumnInfo = append(tableColumnInfo, tRowInfo)
			}
			if len(colRows.RawValues()) == 0 {
				soteErr = sError.GetSError(109999, sError.BuildParams([]string{schemaName + "." + tableName}), nil)
			}
			defer colRows.Close()
		}
	}

	return
}

func ValidateDataType(dataType string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	switch dataType {
	case SBIGINT:
	case SBOOLEAN:
	case SCHARACTERVARYING:
	case SDATE:
	case SINTEGER:
	case STEXT:
	case STIMESTAMPZONE:
	default:
		soteErr = sError.GetSError(401010, nil, nil)
		sLogger.Info(soteErr.FmtErrMsg)
	}

	return
}
