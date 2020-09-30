package sDatabase

var pkList = make(map[string]pkInfo)

type pkInfo struct {
	schemaName string
	tableName  string
}

// func pkPrimer(schemaName string, dbConnInfo ConnInfo) {
// 	sLogger.DebugMethod()
//
// 	// This will return a list of tables for the given schema that have primary keys with only one column.
// 	// At this time, Harvester doesn't support multiple column Primary Key lookup
// 	qStmt1 := "SELECT tc.table_schema, tc.table_name, COUNT (tc.table_name) FROM information_schema.table_constraints tc " +
// 		"INNER JOIN information_schema.constraint_column_usage ccu ON tc.table_schema = ccu.table_schema and tc.table_name = ccu.table_name " +
// 		"WHERE tc.constraint_type = 'PRIMARY KEY' AND tc.table_schema = $1 " +
// 		"GROUP BY tc.table_schema, tc.table_name HAVING COUNT (tc.table_name) = 1;"
//
// 	tbRows, err := dbConnInfo.DBPoolPtr.Query(context.Background(), qStmt1, schemaName)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
//
// 	// This will build a query that will pull the column being used in the primary key
// 	var wSchema bool = true
// 	var qStmt2 strings.Builder
// 	qStmt2.WriteString("SELECT table_schema, table_name, column_name FROM information_schema.constraint_column_usage WHERE table_schema ")
//
// 	for tbRows.Next() {
// 		schTbColumns, err := tbRows.Values()
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
//
// 		if wSchema {
// 			qStmt2.WriteString(fmt.Sprintf("= '%v' AND table_name IN (", schTbColumns[0]))
// 			qStmt2.WriteString(fmt.Sprintf("'%v'", schTbColumns[1]))
// 			wSchema = false
// 		} else {
// 			qStmt2.WriteString(fmt.Sprintf(", '%v'", schTbColumns[1]))
// 		}
//
// 	}
// 	defer tbRows.Close()
// 	qStmt2.WriteString(");")
//
// 	rows, err := dbConnInfo.DBPoolPtr.Query(context.Background(), qStmt2.String())
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
//
// 	for rows.Next() {
// 		columnValues, err := rows.Values()
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 		pkList[columnValues[2].(string)] = pkInfo{schemaName: columnValues[0].(string), tableName: columnValues[1].(string)}
// 	}
// 	rows.Close()
// }
//
// /*
// Using this function will return the table where the column is a primary key
// */
// func PKLookup(tSchemaName, sTableName, sColumnName string, dbConnInfo ConnInfo) (tableName string, soteErr sError.SoteError) {
// 	sLogger.DebugMethod()
//
// 	if len(tSchemaName) == 0 {
// 		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"Schema name: " + tSchemaName}), nil)
// 	}
// 	if soteErr.ErrCode == nil && len(sTableName) == 0 {
// 		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"Table name: " + sTableName}), nil)
// 	}
// 	if soteErr.ErrCode == nil && len(sColumnName) == 0 {
// 		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"Column Name name: " + sColumnName}), nil)
// 	}
//
// 	if soteErr.ErrCode == nil {
// 		if soteErr = VerifyConnection(dbConnInfo); soteErr.ErrCode == nil {
// 			if len(pkList) == 0 {
// 				pkPrimer(tSchemaName, dbConnInfo)
// 				if len(pkList) == 0 {
// 					soteErr = sError.GetSError(109999, sError.BuildParams([]string{"PK INFO for" + tSchemaName}), sError.EmptyMap)
// 				}
// 			}
//
// 			if soteErr.ErrCode == nil || tSchemaName == pkList[sColumnName].schemaName {
// 				if pkList[sColumnName].tableName == sTableName {
// 					tableName = "self"
// 				} else {
// 					tableName = pkList[sColumnName].tableName
// 				}
// 			}
// 		}
// 	}
//
// 	return
// }
