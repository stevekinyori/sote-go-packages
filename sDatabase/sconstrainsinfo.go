package sDatabase

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

func getColumnConstrainInfo(schemaName string, dbConnection *pgx.Conn) {
	sLogger.DebugMethod()

	qStmt1 := "SELECT tc.table_schema, tc.table_name, COUNT (tc.table_name) FROM information_schema.table_constraints tc " +
		"INNER JOIN information_schema.constraint_column_usage ccu ON tc.table_schema = ccu.table_schema and tc.table_name = ccu.table_name " +
		"WHERE tc.constraint_type = 'PRIMARY KEY' AND tc.table_schema = $1 " +
		"GROUP BY tc.table_schema, tc.table_name HAVING COUNT (tc.table_name) = 1;"

	tbRows, err := dbConnection.Query(context.Background(), qStmt1, myHarvest.SchemaName)
	if err != nil {
		log.Fatalln(err)
	}

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

	rows, err := dbConnection.Query(context.Background(), qStmt2.String())
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

	fmt.Println(pkList)
}

func getKeyColumnInfo(schemaName string, dbConnection *pgxpool.Conn) {
	sLogger.DebugMethod()

	qStmt1 := "SELECT tc.table_schema, tc.table_name, COUNT (tc.table_name) FROM information_schema.table_constraints tc " +
		"INNER JOIN information_schema.constraint_column_usage ccu ON tc.table_schema = ccu.table_schema and tc.table_name = ccu.table_name " +
		"WHERE tc.constraint_type = 'PRIMARY KEY' AND tc.table_schema = $1 " +
		"GROUP BY tc.table_schema, tc.table_name HAVING COUNT (tc.table_name) = 1;"

	tbRows, err := dbConnection.Query(context.Background(), qStmt1, myHarvest.SchemaName)
	if err != nil {
		log.Fatalln(err)
	}


}