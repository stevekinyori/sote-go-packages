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
func GetTableList(schemaName string, tConnInfo ConnInfo) (tableList []string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = VerifyConnection(tConnInfo); soteErr.ErrCode == nil {
		qStmt := "SELECT table_name FROM information_schema.tables WHERE table_schema = $1;"

		var tbRows pgx.Rows
		var err error
		tbRows, err = tConnInfo.DBPoolPtr.Query(context.Background(), qStmt, schemaName)
		if err != nil {
			log.Fatalln(err)
		}

		var tableRow []interface{}
		for tbRows.Next() {
			tableRow, err = tbRows.Values()
			if err != nil {
				log.Fatalln(err)
			}
			tableList = append(tableList, tableRow[0].(string))

		}
		defer tbRows.Close()
	}

	return
}
