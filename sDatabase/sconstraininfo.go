// All functions use the dbConnPtr or DBPoolPtr which are established using sconnection.
package sDatabase

type SConstraint struct {
	tableName  string
	columnName string
}

// This will only return column name for constrains that have one column as the primary key for the table
// func GetSingleColumnConstraintInfo(schemaName string, tConnInfo ConnInfo) (SConstraintInfo []SConstraint, soteErr sError.SoteError) {
// 	sLogger.DebugMethod()
//
// 	if len(schemaName) == 0 {
// 		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"Schema name: " + schemaName}), nil)
// 	} else {
// 		if soteErr = VerifyConnection(tConnInfo); soteErr.ErrCode == nil {
// 			qStmt1 := "SELECT tc.table_schema, tc.table_name, COUNT (tc.table_name) FROM information_schema.table_constraints tc " +
// 				"INNER JOIN information_schema.constraint_column_usage ccu ON tc.table_schema = ccu.table_schema and tc.table_name = ccu.table_name " +
// 				"WHERE tc.constraint_type = 'PRIMARY KEY' AND tc.table_schema = $1 " +
// 				"GROUP BY tc.table_schema, tc.table_name HAVING COUNT (tc.table_name) = 1;"
//
// 			var tbRows pgx.Rows
// 			var err error
// 			tbRows, err = tConnInfo.DBPoolPtr.Query(context.Background(), qStmt1, schemaName)
// 			if err != nil {
// 				log.Fatalln(err)
// 			}
//
// 			var firstList bool = true
// 			var qStmt2 strings.Builder
// 			qStmt2.WriteString("SELECT table_name, column_name FROM information_schema.constraint_column_usage WHERE table_schema ")
//
// 			for tbRows.Next() {
// 				schTbColumns, err := tbRows.Values()
// 				if err != nil {
// 					log.Fatalln(err)
// 				}
//
// 				if firstList {
// 					qStmt2.WriteString(fmt.Sprintf("= '%v' AND table_name IN (", schTbColumns[0]))
// 					qStmt2.WriteString(fmt.Sprintf("'%v'", schTbColumns[1]))
// 					firstList = false
// 				} else {
// 					qStmt2.WriteString(fmt.Sprintf(", '%v'", schTbColumns[1]))
// 				}
//
// 			}
// 			defer tbRows.Close()
// 			qStmt2.WriteString(");")
//
// 			if len(tbRows.RawValues()) > 0 {
// 				tbRows, err = tConnInfo.DBPoolPtr.Query(context.Background(), qStmt2.String())
// 				if err != nil {
// 					log.Fatalln(err)
// 				}
//
// 				var constraintRow []interface{}
// 				var tRowInfo SConstraint
// 				for tbRows.Next() {
// 					constraintRow, err = tbRows.Values()
// 					if err != nil {
// 						log.Fatalln(err)
// 					}
// 					tRowInfo.tableName = constraintRow[0].(string)
// 					tRowInfo.columnName = constraintRow[1].(string)
// 					SConstraintInfo = append(SConstraintInfo, tRowInfo)
// 				}
// 				defer tbRows.Close()
// 			}
// 		}
// 	}
//
// 	return
// }

// func getKeyColumnInfo(schemaName string, dbConnection *pgxpool.Conn) {
// 	sLogger.DebugMethod()

	// qStmt1 := "SELECT tc.table_schema, tc.table_name, COUNT (tc.table_name) FROM information_schema.table_constraints tc " +
	// 	"INNER JOIN information_schema.constraint_column_usage ccu ON tc.table_schema = ccu.table_schema and tc.table_name = ccu.table_name " +
	// 	"WHERE tc.constraint_type = 'PRIMARY KEY' AND tc.table_schema = $1 " +
	// 	"GROUP BY tc.table_schema, tc.table_name HAVING COUNT (tc.table_name) = 1;"

	// tbRows, err := dbConnection.Query(context.Background(), qStmt1, schemaName)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

// }
