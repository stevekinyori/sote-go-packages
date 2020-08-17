// All functions use the dbConnPtr or DBPoolPtr which are established using sconnection.
package sDatabase

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

// This function gets a list of tables for the supplied schema.
func GetColumnInfo(schemaName, tableName string, tConnInfo ConnInfo) (columnList []string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = VerifyConnection(tConnInfo); soteErr.ErrCode == nil {
		qStmt := "SELECT table_name FROM information_schema.columns WHERE table_schema = $1 AND table_name = $2;"

		var colRows pgx.Rows
		var err error
		colRows, err = tConnInfo.DBPoolPtr.Query(context.Background(), qStmt, schemaName, tableName)
		if err != nil {
			log.Fatalln(err)
		}

		var columnRow []interface{}
		for colRows.Next() {
			columnRow, err = colRows.Values()
			if err != nil {
				log.Fatalln(err)
			}
			columnList = append(columnList, columnRow[0].(string))

		}
		defer colRows.Close()
	}

	return
}
