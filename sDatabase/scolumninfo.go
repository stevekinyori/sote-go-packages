// All functions use the dbConnPtr or DBPoolPtr which are established using sconnection.
package sDatabase

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

type SColumnInfo struct {
	// Col was added to the name to prevent conflicts
	colName     string
	colDefault  string
	colNullable string
	colDataType string
}

// This function gets column information for the supplied schema and table.
func GetColumnInfo(schemaName, tableName string, tConnInfo ConnInfo) (tableColumnInfo []SColumnInfo, soteErr sError.SoteError) {
	sLogger.DebugMethod()

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
			tRowInfo.colName = columnRow[0].(string)
			if columnRow[1] == nil {
				tRowInfo.colDefault = ""
			} else {
				tRowInfo.colDefault = columnRow[1].(string)
			}
			tRowInfo.colNullable = columnRow[2].(string)
			tRowInfo.colDataType = columnRow[3].(string)
			tableColumnInfo = append(tableColumnInfo, tRowInfo)
		}
		defer colRows.Close()
	}

	return
}
