// All functions use the dbConnPtr or dbPoolPtr which are established using sconnection.
package sDatabase

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

// This function gets a list of tables for the supplied schema.
func getTables(schemaName string) (tableList []string, soteErr sError.SoteError){
	sLogger.DebugMethod()

	if soteErr = ConnectionEstablished(); soteErr.ErrCode == nil {

		qStmt := "SELECT table_name FROM information_schema.tables WHERE table_schema = $1;"

		var tbRows pgx.Rows
		var err error
		if dsConnValues.ConnType == SINGLECONN {
			tbRows, err = dbConnPtr.Query(context.Background(), qStmt, schemaName)
		} else {
			tbRows, err = dbPoolPtr.Query(context.Background(), qStmt, schemaName)
		}
		if err != nil {
			log.Fatalln(err)
		}

		var tableRow []interface{}
		for tbRows.Next() {
			tableRow, err = tbRows.Values()
			tableList = append(tableList, tableRow[0].(string))
			if err != nil {
				log.Fatalln(err)
			}

		}
		defer tbRows.Close()
	}

	return
}
