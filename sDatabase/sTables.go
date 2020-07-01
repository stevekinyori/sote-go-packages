package sDatabase

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)



func getTables(schemaName string) {
	sLogger.DebugMethod()

	qStmt := "SELECT table_schema, table_name FROM information_schema.tables WHERE table_schema = $1;"

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

	var tableData []interface{}
	for tbRows.Next() {
		tableData, err = tbRows.Values()
		if err != nil {
			log.Fatalln(err)
		}
	}
	defer tbRows.Close()

	fmt.Println(tableData)
}


