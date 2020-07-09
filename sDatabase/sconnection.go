// This is a wrapper for github.com/jackc/pgx/v4 connection.  We are wrapping this
// so that all Sote Go developers connect to Sote databases the same way.
package sDatabase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

const (
	// SSL Modes
	SSLMODEDISABLE  = "disable"
	SSLMODEALLOW    = "allow"
	SSLMODEPREFER   = "prefer"
	SSLMODEREQUIRED = "require"
	DSCONNFORMAT    = "dbname=%v user=%v password=%v host=%v port=%v connect_timeout=%v sslmode=%v"
)

type ConnInfo struct {
	dbPoolPtr    *pgxpool.Pool
	dsConnValues ConnValues
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

// This will create a connection to a database and populate ConnInfo
//
//   dbName   Name of the Postgres database
//   user     User that connection will use to authenticate
//   password Users password for authentication
//   host     Internet DNS or IP address of the server running the instance of Postgres
//   sslMode  Type of encryption used for the connection (https://www.postgresql.org/docs/12/libpq-ssl.html for version 12)
//   port     Interface the connection communicates with Postgres
//   timeout  Number of seconds a request must complete (3 seconds is normal setting)
func GetConnection(dbName, user, password, host, sslMode string, port, timeout int) (dbConnInfo ConnInfo, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if dbConnInfo.dsConnValues, soteErr = setConnectionValues(dbName, user, password, host, sslMode, port, timeout); soteErr.ErrCode != nil {
		panic("Invalid connection parameters for database: " + soteErr.FmtErrMsg)
	} else {
		var err error
		var dsConnString = fmt.Sprintf(DSCONNFORMAT, dbConnInfo.dsConnValues.DBName, dbConnInfo.dsConnValues.User, dbConnInfo.dsConnValues.Password, dbConnInfo.dsConnValues.Host,
			dbConnInfo.dsConnValues.Port, dbConnInfo.dsConnValues.Timeout, dbConnInfo.dsConnValues.SSLMode)
		dbConnInfo.dbPoolPtr, err = pgxpool.Connect(context.Background(), dsConnString)
		if err != nil {
			errDetails, soteErr := sError.ConvertErr(err)
			if soteErr.ErrCode != nil {
				sLogger.Info(soteErr.FmtErrMsg)
				panic("sError.ConvertErr Failed")
			}
			sLogger.Info(sError.GetSError(800100, nil, errDetails).FmtErrMsg)
			panic("MakeConnection Failed")
		}
		defer dbConnInfo.dbPoolPtr.Close()
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
		soteErr = sError.GetSError(602020, buildParams([]string{sslMode}), nil)
		sLogger.Info(soteErr.FmtErrMsg)
	}

	tConnInfo = ConnValues{dbName, user, password, host, port, timeout, sslMode}

	return
}

// This will convert the connection values used to connect to the Sote database into
// a JSON string.
func ToJSONString(dsConnValues ConnValues) (jsonString string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	jsonConnValues, err := json.Marshal(dsConnValues)
	if err != nil {
		sLogger.Info(sError.GetSError(400100, buildParams([]string{"dsConnValues", "struct"}), nil).FmtErrMsg)
	}

	jsonString = string(jsonConnValues)

	return
}

func buildParams(values []string) (s []interface{}) {
	sLogger.DebugMethod()

	s = make([]interface{}, len(values))
	for i, v := range values {
		s[i] = v
	}

	return
}
