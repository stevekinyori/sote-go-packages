/*
This is a wrapper for github.com/jackc/pgx/v4 connection.  We are wrapping this
so that all Sote Go developers connect to Sote databases the same way.
*/
package sDatabase

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.com/soteapps/packages/v2023/sConfigParams"
	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
)

const (
	// SSL Modes
	SSLMODEDISABLE  = "disable"
	SSLMODEALLOW    = "allow"
	SSLMODEPREFER   = "prefer"
	SSLMODEREQUIRED = "require"
	DSCONNFORMAT    = "dbname=%v search_path=%v user=%v password=%v host=%v port=%v connect_timeout=%v sslmode=%v"
)

// ConnInfo SRow and SRows are so pgx package doesn't need to be imported in everywhere there are queries to the database.
type ConnInfo struct {
	DBPoolPtr    *pgxpool.Pool
	DSConnValues ConnValues
	DBContext    context.Context
}

type ConnValues struct {
	DBName   string `json:"dbName"`
	Schema   string `json:"schema"`
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

// New creates a new database connection based on environment
func New(ctx context.Context, environment string) (dbConnInfo ConnInfo, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		dbConfig *sConfigParams.Database
	)

	if dbConfig, soteErr = sConfigParams.GetAWSParams(ctx, sConfigParams.API, environment); soteErr.ErrCode != nil {
		return
	}

	if dbConnInfo, soteErr = GetConnection(dbConfig.Name, dbConfig.Schema, dbConfig.User, dbConfig.Password, dbConfig.Host,
		dbConfig.SSLMode, dbConfig.Port, 3); soteErr.ErrCode != nil {
		return
	}

	if soteErr = VerifyConnection(dbConnInfo); soteErr.ErrCode != nil {
		return
	}

	return
}

// GetConnection This will create a connection to a database and populate ConnInfo
//
//  dbName   - Name of the Postgres database
//  dbSchema   - The default Schema of the Postgres database
//  user     - User that connection will use to authenticate
//  password - Users password for authentication
//  host     - Internet DNS or IP address of the server running the instance of Postgres
//  sslMode  - Type of encryption used for the connection (https://www.postgresql.org/docs/12/libpq-ssl.html for version 12)
//  port     - Interface the connection communicates with Postgres
//  timeout  - Number of seconds a request must complete (3 seconds is normal setting)
//
//  DBContext is also set to context.Background() an empty context./*
func GetConnection(dbName, dbSchema, user, password, host, sslMode string, port, timeout int) (dbConnInfo ConnInfo, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if dbConnInfo.DSConnValues, soteErr = setConnectionValues(dbName, dbSchema, user, password, host, sslMode, port,
		timeout); soteErr.ErrCode != nil {
		sLogger.Info("Invalid connection parameters for database: " + soteErr.FmtErrMsg)

		return
	} else {
		var (
			err    error
			DBPort int
		)

		if tDbPort, ok := os.LookupEnv("DB_PORT"); ok {
			if DBPort, err = strconv.Atoi(tDbPort); err == nil {
				dbConnInfo.DSConnValues.Port = DBPort
			}
		}

		var dsConnString = fmt.Sprintf(DSCONNFORMAT, dbConnInfo.DSConnValues.DBName, dbConnInfo.DSConnValues.Schema, dbConnInfo.DSConnValues.User,
			dbConnInfo.DSConnValues.Password, dbConnInfo.DSConnValues.Host, dbConnInfo.DSConnValues.Port, dbConnInfo.DSConnValues.Timeout,
			dbConnInfo.DSConnValues.SSLMode)
		if dbConnInfo.DBPoolPtr, err = pgxpool.Connect(context.Background(), dsConnString); err != nil {
			if strings.Contains(err.Error(), "dial") {
				soteErr = sError.GetSError(209299, nil, sError.EmptyMap)
				sLogger.Info(soteErr.FmtErrMsg)
			} else {
				var errDetails = make(map[string]string)
				errDetails, soteErr = sError.ConvertErr(err)
				if soteErr.ErrCode != nil {
					sLogger.Info(soteErr.FmtErrMsg)
					sLogger.Info("sError.ConvertErr Failed")
					return
				}

				sLogger.Info(sError.GetSError(210200, nil, errDetails).FmtErrMsg)
				sLogger.Info("sDatabase.sconnection.GetConnection Failed")
				return
			}
		}

		dbConnInfo.DBContext = context.Background()
	}

	return
}

// This will set the connection values so GetConnection can be executed.
func setConnectionValues(dbName, dbSchema, user, password, host, sslMode string, port, timeout int) (tConnInfo ConnValues, soteErr sError.SoteError) {
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

	tConnInfo = ConnValues{dbName, dbSchema, user, password, host, port, timeout, sslMode}
	return
}

// ToJSONString This will convert the connection values used to connect to the Sote database into
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

// VerifyConnection Verify that the pointer to the database connection is active.
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
