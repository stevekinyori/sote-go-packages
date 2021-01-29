// All functions use the dbConnPtr or DBPoolPtr which are established using sconnection.
package sDatabase

// This function gets a list of tables for the supplied schema.
// func GetTableList(schemaName string, tConnInfo ConnInfo) (tableList []string, soteErr sError.SoteError) {
// 	sLogger.DebugMethod()
//
// 	if len(schemaName) == 0 {
// 		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"Schema name: " + schemaName}), nil)
// 	}
//
// 	if soteErr.ErrCode == nil {
// 		if soteErr = VerifyConnection(tConnInfo); soteErr.ErrCode == nil {
// 			qStmt := "SELECT table_name FROM information_schema.tables WHERE table_schema = $1;"
//
// 			var tbRows pgx.Rows
// 			var err error
// 			tbRows, err = tConnInfo.DBPoolPtr.Query(context.Background(), qStmt, schemaName)
// 			if err != nil {
// 				log.Fatalln(err)
// 			}
//
// 			var tableRow []interface{}
// 			for tbRows.Next() {
// 				tableRow, err = tbRows.Values()
// 				if err != nil {
// 					log.Fatalln(err)
// 				}
// 				tableList = append(tableList, tableRow[0].(string))
//
// 			}
// 			defer tbRows.Close()
// 		}
// 	}
//
// 	return
// }
