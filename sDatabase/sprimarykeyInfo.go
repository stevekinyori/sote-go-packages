package sDatabase

import (
	"context"
	"fmt"
	"log"
	"strings"

	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

var pkList = make(map[string]pkInfo)

type pkInfo struct {
	schemaName string
	tableName  string
}

func pkPrimer(schemaName string, dbConnInfo ConnInfo) {
	sLogger.DebugMethod()

	// This will return a list of tables for the given schema that have primary keys with only one column.
	// At this time, Harvester doesn't support multiple column Primary Key lookup
	qStmt1 := "SELECT tc.table_schema, tc.table_name, COUNT (tc.table_name) FROM information_schema.table_constraints tc " +
		"INNER JOIN information_schema.constraint_column_usage ccu ON tc.table_schema = ccu.table_schema and tc.table_name = ccu.table_name " +
		"WHERE tc.constraint_type = 'PRIMARY KEY' AND tc.table_schema = $1 " +
		"GROUP BY tc.table_schema, tc.table_name HAVING COUNT (tc.table_name) = 1;"

	tbRows, err := dbConnInfo.DBPoolPtr.Query(context.Background(), qStmt1, schemaName)
	if err != nil {
		log.Fatalln(err)
	}

	// This will build a query that will pull the column being used in the primary key
	var wSchema bool = true
	var qStmt2 strings.Builder
	qStmt2.WriteString("SELECT table_schema, table_name, column_name FROM information_schema.constraint_column_usage WHERE table_schema ")

	for tbRows.Next() {
		schTbColumns, err := tbRows.Values()
		if err != nil {
			log.Fatalln(err)
		}

		if wSchema {
			qStmt2.WriteString(fmt.Sprintf("= '%v' AND table_name IN (", schTbColumns[0]))
			qStmt2.WriteString(fmt.Sprintf("'%v'", schTbColumns[1]))
			wSchema = false
		} else {
			qStmt2.WriteString(fmt.Sprintf(", '%v'", schTbColumns[1]))
		}

	}
	defer tbRows.Close()
	qStmt2.WriteString(");")
	fmt.Println(qStmt2)

	rows, err := dbConnInfo.DBPoolPtr.Query(context.Background(), qStmt2.String())
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		columnValues, err := rows.Values()
		if err != nil {
			log.Fatalln(err)
		}
		pkList[columnValues[2].(string)] = pkInfo{schemaName: columnValues[0].(string), tableName: columnValues[1].(string)}
	}
	rows.Close()
}

func pkLookup(tSchemaName, sTableName, sColumnName string, dbConnInfo ConnInfo, test bool) (schemaName, tableName string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(pkList) == 0 {
		pkPrimer(tSchemaName, dbConnInfo)
	} else {
		// run lookup
	}

	// if pkList[sColumnName].tableName == sTableName {
	// 	schemaName = "self"
	// 	tableName = "self"
	// } else {
	// 	schemaName = pkList[sColumnName].schemaName
	// 	tableName = pkList[sColumnName].tableName
	// }

	return
}