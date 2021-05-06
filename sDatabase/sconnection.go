/*
This is a wrapper for github.com/jackc/pgx/v4 connection.  We are wrapping this
so that all Sote Go developers connect to Sote databases the same way.
*/
package sDatabase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

const (
	// SSL Modes
	SSLMODEDISABLE  = "disable"
	SSLMODEALLOW    = "allow"
	SSLMODEPREFER   = "prefer"
	SSLMODEREQUIRED = "require"
	DSCONNFORMAT    = "dbname=%v user=%v password=%v host=%v port=%v connect_timeout=%v sslmode=%v"
)

// SRow and SRows are so pgx package doesn't need to be imported in every where there are queries to the database.
type ConnInfo struct {
	DBPoolPtr    *pgxpool.Pool
	DSConnValues ConnValues
	DBContext    context.Context
}

type ConnValues struct {
	DBName   string `json:"dbName"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Timeout  int    `json:"timeout"`
	SSLMode  string `json:"sslMode"`
}

type STransaction pgx.Tx
type SRows pgx.Rows
type SRow pgx.Row

/*
This will create a connection to a database and populate ConnInfo

  dbName   Name of the Postgres database
  user     User that connection will use to authenticate
  password Users password for authentication
  host     Internet DNS or IP address of the server running the instance of Postgres
  sslMode  Type of encryption used for the connection (https://www.postgresql.org/docs/12/libpq-ssl.html for version 12)
  port     Interface the connection communicates with Postgres
  timeout  Number of seconds a request must complete (3 seconds is normal setting)

  DBContext is also set to context.Background() an empty context.
*/
func GetConnection(dbName, user, password, host, sslMode string, port, timeout int) (dbConnInfo ConnInfo, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if dbConnInfo.DSConnValues, soteErr = setConnectionValues(dbName, user, password, host, sslMode, port, timeout); soteErr.ErrCode != nil {
		panic("Invalid connection parameters for database: " + soteErr.FmtErrMsg)
	} else {
		var err error
		var dsConnString = fmt.Sprintf(DSCONNFORMAT, dbConnInfo.DSConnValues.DBName, dbConnInfo.DSConnValues.User, dbConnInfo.DSConnValues.Password,
			dbConnInfo.DSConnValues.Host,
			dbConnInfo.DSConnValues.Port, dbConnInfo.DSConnValues.Timeout, dbConnInfo.DSConnValues.SSLMode)
		if dbConnInfo.DBPoolPtr, err = pgxpool.Connect(context.Background(), dsConnString); err != nil {
			if strings.Contains(err.Error(), "dial") {
				soteErr = sError.GetSError(209299, nil, sError.EmptyMap)
				sLogger.Info(soteErr.FmtErrMsg)
			} else {
				var errDetails = make(map[string]string)
				errDetails, soteErr = sError.ConvertErr(err)
				if soteErr.ErrCode != nil {
					sLogger.Info(soteErr.FmtErrMsg)
					panic("sError.ConvertErr Failed")
				}
				sLogger.Info(sError.GetSError(210200, nil, errDetails).FmtErrMsg)
				panic("sDatabase.sconnection.GetConnection Failed")
			}
		}
		dbConnInfo.DBContext = context.Background()
	}
	return
}

// This will set the connection values so GetConnection can be execute.
func setConnectionValues(dbName, user, password, host, sslMode string, port, timeout int) (tConnInfo ConnValues, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	switch sslMode {
	case SSLMODEDISABLE:
	case SSLMODEALLOW:
	case SSLMODEPREFER:
	case SSLMODEREQUIRED:
	default:
		soteErr = sError.GetSError(209220, sError.BuildParams([]string{sslMode}), nil)
		sLogger.Info(soteErr.FmtErrMsg)
	}

	tConnInfo = ConnValues{dbName, user, password, host, port, timeout, sslMode}
	return
}

// This will convert the connection values used to connect to the Sote database into
// a JSON string.
func ToJSONString(DSConnValues ConnValues) (jsonString string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	jsonConnValues, err := json.Marshal(DSConnValues)
	if err != nil {
		soteErr = sError.GetSError(207010, sError.BuildParams([]string{"DSConnValues", "struct"}), nil)
		sLogger.Info(soteErr.FmtErrMsg)
	} else {
		jsonString = string(jsonConnValues)
	}

	return
}

// Verify that the pointer to the database connection is active.
func VerifyConnection(tConnInfo ConnInfo) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if tConnInfo.DBPoolPtr == nil {
		soteErr = sError.GetSError(209299, nil, nil)
		sLogger.Info(soteErr.FmtErrMsg)
	} else {
		qStmt := "SELECT * FROM pg_stat_activity WHERE datname = $1 and state = 'active';"

		var tbRows pgx.Rows
		var err error
		tbRows, err = tConnInfo.DBPoolPtr.Query(context.Background(), qStmt, tConnInfo.DSConnValues.DBName)
		if err != nil {
			soteErr = sError.GetSError(209299, nil, nil)
			sLogger.Info(soteErr.FmtErrMsg)
		}
		defer tbRows.Close()
	}

	return
}
