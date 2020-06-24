package sDatabase

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.com/soteapps/packages/serror"
	"gitlab.com/soteapps/packages/slogger"
)

const (
	// Connection types
	SINGLECONN      = "SINGLE"
	POOLCONN        = "POOL"
	// SSL Modes
	SSLMODEDISABLE  = "disable"
	SSLMODEALLOW    = "allow"
	SSLMODEPREFER   = "prefer"
	SSLMODEREQUIRED = "require"
)

var (
	dsConnFormat = "dbname=%v user=%v password=%v host=%v port=%v connect_timeout=%v sslmode=%v"
	dsConnString string
	dbConn       *pgx.Conn
	dbPool       *pgxpool.Pool
)

// This will return a pointer to the database connection based on the type requested.
//
//   connType
//     SINGLE, single for a connection that only allows one action at a time
//     POOL, pool     for a connection that only multiple concurrent uses <Param/>
//
func GetConnection(connType string) (dbConnPtr *pgx.Conn, dbPoolPtr *pgxpool.Pool, soteErr serror.SoteError) {
	slogger.DebugMethod()

	switch strings.ToUpper(connType) {
	case SINGLECONN:
		soteErr = singleConnectPGSQL()
		dbConnPtr = dbConn
	case POOLCONN:
		soteErr = poolConnectPGSQL()
		dbPoolPtr = dbPool
	default:
		params := make([]interface{}, 0)
		// 	CREATE ERROR FOR UNSUPPORT CONNECTION TYPE
		soteErr = serror.GetSError(0, params, serror.EmptyMap)
	}

	return
}

func SetConnectionValues(dbName, user, password, host, sslMode string, port, timeout int) {
	slogger.DebugMethod()

	switch sslMode {
	case SSLMODEDISABLE:
	case SSLMODEALLOW:
	case SSLMODEPREFER:
	case SSLMODEREQUIRED:
	default:
		log.Println(serror.GetSError(10, nil, nil).FmtErrMsg)
	// 	CREATE ERROR FOR UNSUPPORT SSL MODE
	}
	dsConnString = fmt.Sprintf(dsConnFormat, dbName, user, password, host, port, timeout, sslMode)
}

func singleConnectPGSQL() (errDetails serror.SoteError) {
	slogger.DebugMethod()

	var err error
	dbConn, err = pgx.Connect(context.Background(), dsConnString)
	if err != nil {
		errDetails, soteErr := serror.ConvertErr(err)
		if soteErr.ErrCode != nil {

			log.Println(soteErr.FmtErrMsg)
		}
		log.Println(serror.GetSError(800100, nil, errDetails).FmtErrMsg)
	}

	return
}

func poolConnectPGSQL() (errDetails serror.SoteError) {
	slogger.DebugMethod()

	var err error
	dbPool, err = pgxpool.Connect(context.Background(), dsConnString)
	if err != nil {
		errDetails, soteErr := serror.ConvertErr(err)
		if soteErr.ErrCode != nil {

			log.Println(soteErr.FmtErrMsg)
		}
		log.Println(serror.GetSError(800100, nil, errDetails).FmtErrMsg)
	}

	return
}
