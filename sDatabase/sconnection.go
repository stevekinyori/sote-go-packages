package sDatabase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.com/soteapps/packages/sError"
	"gitlab.com/soteapps/packages/sLogger"
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
)

var (
	dsSingleConnFormat = "dbname=%v user=%v password=%v host=%v port=%v connect_timeout=%v sslmode=%v"
	// dsPoolConnFormat postgresql://[user[:password]@][host][:port][/dbname][?param1=value1&...]
	dsPoolConnFormat = "postgresql://%v:%v@%v:%v/%v?connect_timeout=v%&sslmode=%v"
	dbConn           *pgx.Conn
	dbPool           *pgxpool.Pool
)

type ConnValues struct {
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
//
//   If the connType provided is not support:
//     dbConnPtr will be nil
//     dbPoolPtr will be nil
//     sError.SoteError will be populated
//   dbName   Name of the Postgres database
//   user     User that connection will use to authenticate
//   password Users password for authentication
//   host     Internet DNS or IP address of the server running the instance of Postgres
//   sslMode  Type of encryption used for the connection (https://www.postgresql.org/docs/12/libpq-ssl.html for version 12)
//   port     Interface the connection communicates with Postgres
//   timeout  Number of seconds a request must complete
func GetConnection(connType, dbName, user, password, host, sslMode string, port, timeout int) (dbConnPtr *pgx.Conn, dbPoolPtr *pgxpool.Pool, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if dsConnValues, soteErr := setConnectionValues(dbName, user, password, host, sslMode, port, timeout); soteErr.ErrCode != nil {
		panic("Invalid connection parameters for database: " + soteErr.FmtErrMsg)
	} else {
		switch strings.ToUpper(connType) {
		case SINGLECONN:
			soteErr = singleConnectionPGSQL(dsConnValues)
			dbConnPtr = dbConn
		case POOLCONN:
			// 	db url: postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
			soteErr = poolConnectionPGSQL(dsConnValues)
			dbPoolPtr = dbPool
		default:
			params := make([]interface{}, 1)
			params[0] = connType
			soteErr = sError.GetSError(602100, params, sError.EmptyMap)
		}
	}

	return
}

// This will set the connection values so GetConnection can be called.
func setConnectionValues(dbName, user, password, host, sslMode string, port, timeout int) (dsConnValues ConnValues, soteErr sError.SoteError) {
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

	dsConnValues = ConnValues{dbName, user, password, host, port, timeout, sslMode}

	return
}

// This will return the connection values in JSON format
func GetConnectionValuesJSON(dsConnValues ConnValues) string {
	sLogger.DebugMethod()

	jsonString, err := json.Marshal(dsConnValues)
	if err != nil {
		sLogger.Info(sError.GetSError(400100, buildParams([]string{"dsConnValues", "struct"}), nil).FmtErrMsg)
	}

	return string(jsonString)
}

func singleConnectionPGSQL(dsConnValues ConnValues) (errDetails sError.SoteError) {
	sLogger.DebugMethod()

	var err error
	dbConn, err = pgx.Connect(context.Background(), fmt.Sprintf(dsSingleConnFormat, dsConnValues.DBName, dsConnValues.User, dsConnValues.Password, dsConnValues.Host, dsConnValues.Port,
		dsConnValues.Timeout, dsConnValues.SSLMode))
	if err != nil {
		errDetails, soteErr := sError.ConvertErr(err)
		if soteErr.ErrCode != nil {
			log.Println(soteErr.FmtErrMsg)
		}
		log.Println(sError.GetSError(800100, nil, errDetails).FmtErrMsg)
	}

	return
}

func poolConnectionPGSQL(dsConnValues ConnValues) (errDetails sError.SoteError) {
	sLogger.DebugMethod()

	var err error
	dbPool, err = pgxpool.Connect(context.Background(), fmt.Sprintf(dsSingleConnFormat, dsConnValues.DBName, dsConnValues.User, dsConnValues.Password, dsConnValues.Host, dsConnValues.Port,
		dsConnValues.Timeout, dsConnValues.SSLMode))
	if err != nil {
		errDetails, soteErr := sError.ConvertErr(err)
		if soteErr.ErrCode != nil {
			log.Println(soteErr.FmtErrMsg)
		}
		log.Println(sError.GetSError(800100, nil, errDetails).FmtErrMsg)
	}

	return
}

func buildParams(values []string) (s []interface{}) {
	s = make([]interface{}, len(values))
	for i, v := range values {
		s[i] = v
	}

	return
}
