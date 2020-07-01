package sDatabase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

const (
	// Connection types
	SINGLECONN = "SINGLE"
	POOLCONN   = "POOL"
	// SSL Modes
	SSLMODEDISABLE  = "disable"
	SSLMODEALLOW    = "allow"
	SSLMODEPREFER   = "prefer"
	SSLMODEREQUIRED = "require"
	DSCONNFORMAT    = "dbname=%v user=%v password=%v host=%v port=%v connect_timeout=%v sslmode=%v"
)

var (
	dsConnValues ConnValues
	dbPoolPtr    *pgxpool.Pool
	dbConnPtr    *pgx.Conn
)

type ConnValues struct {
	ConnType string `json:"connType"`
	DBName   string `json:"dbName"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Timeout  int    `json:"timeout"`
	SSLMode  string `json:"sslMode"`
}

// This will return a pointer to the database connection based on the type requested.
//
//   connType (Forced to upper case)
//     single for a connection that only allows one action at a time
//     pool   for a connection that only multiple concurrent uses
//   dbName   Name of the Postgres database
//   user     User that connection will use to authenticate
//   password Users password for authentication
//   host     Internet DNS or IP address of the server running the instance of Postgres
//   sslMode  Type of encryption used for the connection (https://www.postgresql.org/docs/12/libpq-ssl.html for version 12)
//   port     Interface the connection communicates with Postgres
//   timeout  Number of seconds a request must complete (3 seconds is normal setting)
func GetConnection(connType, dbName, user, password, host, sslMode string, port, timeout int) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr := setConnectionValues(connType, dbName, user, password, host, sslMode, port, timeout); soteErr.ErrCode != nil {
		panic("Invalid connection parameters for database: " + soteErr.FmtErrMsg)
	} else {
		var err error
		var dsConnString = fmt.Sprintf(DSCONNFORMAT, dsConnValues.DBName, dsConnValues.User, dsConnValues.Password, dsConnValues.Host, dsConnValues.Port, dsConnValues.Timeout, dsConnValues.SSLMode)
		switch dsConnValues.ConnType {
		case SINGLECONN:
			dbConnPtr, err = pgx.Connect(context.Background(), dsConnString)
		case POOLCONN:
			dbPoolPtr, err = pgxpool.Connect(context.Background(), dsConnString)
		default:
			params := make([]interface{}, 1)
			params[0] = connType
			soteErr = sError.GetSError(602100, params, sError.EmptyMap)
		}
		if err != nil {
			errDetails, soteErr := sError.ConvertErr(err)
			if soteErr.ErrCode != nil {
				sLogger.Info(soteErr.FmtErrMsg)
				panic("sError.ConvertErr Failed")
			}
			sLogger.Info(sError.GetSError(800100, nil, errDetails).FmtErrMsg)
			panic("GetConnection Failed")
		}
	}

	return
}

// This will set the connection values so GetConnection can be called.
func setConnectionValues(connType, dbName, user, password, host, sslMode string, port, timeout int) (soteErr sError.SoteError) {
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

	dsConnValues = ConnValues{strings.ToUpper(connType), dbName, user, password, host, port, timeout, sslMode}

	return
}

// This will return the connection values in JSON format
func GetConnectionValuesJSON() (jsonString string) {
	sLogger.DebugMethod()

	json, err := json.Marshal(dsConnValues)
	if err != nil {
		sLogger.Info(sError.GetSError(400100, buildParams([]string{"dsConnValues", "struct"}), nil).FmtErrMsg)
	}

	jsonString = string(json)

	return
}

func buildParams(values []string) (s []interface{}) {
	s = make([]interface{}, len(values))
	for i, v := range values {
		s[i] = v
	}

	return
}
