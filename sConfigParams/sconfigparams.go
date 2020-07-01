package sConfigParams

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

const (
	// Environments
	staging     = "staging"
	development = "development"
	demo        = "demo"
	production  = "production"
	// Variables
	dbPasswordKey = "DATABASE_PASSWORD"
	dbHostKey     = ""
	dbUserKey     = ""
	dbPortKey     = ""
	dbNameKey     = ""
	dbSSLModeKey  = ""
	noKey         = ""
)

var (
	defaultPath                = "/sote/api/"
	setToTrue           bool   = true
	pTrue               *bool  = &setToTrue
	maxResult           int64  = 10
	pMaxResult          *int64 = &maxResult
	sSSMParameterValues *ssm.GetParametersByPathOutput
	awsService          *ssm.SSM
)

func init() {
	slogger.DebugMethod()

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatalln(err)
	}
	awsService = ssm.New(sess)
}

func initParameters(tAWS, key string) {
	slogger.DebugMethod()

	var (
		ssmPath           ssm.GetParametersByPathInput
		err               error
		queryPath         string  = defaultPath + tAWS
		pQueryPath        *string = &queryPath
		queryPathWithKey  string  = defaultPath + tAWS + key
		pQueryPathWithKey *string = &queryPathWithKey
	)
	if len(key) == 0 {
		ssmPath.Path = pQueryPath
	} else {
		ssmPath.Path = pQueryPathWithKey
	}
	ssmPath.Recursive = pTrue
	ssmPath.WithDecryption = pTrue
	ssmPath.MaxResults = pMaxResult
	sSSMParameterValues, err = awsService.GetParametersByPath(&ssmPath)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetParameters(tAWS string) (parameters []*ssm.Parameter, found bool) {
	slogger.DebugMethod()

	initParameters(tAWS, noKey)
	if len(sSSMParameterValues.Parameters) > 0 {
		parameters = sSSMParameterValues.Parameters
	}

	return
}

func GetDBPassword(tAWS string) (dbPassword string, found bool) {
	slogger.DebugMethod()

	initParameters(tAWS, dbPasswordKey)
	dbPassword, found = search(tAWS, dbPasswordKey)

	return
}

func GetDBHost(tAWS string) (dbHost string, found bool) {
	slogger.DebugMethod()

	initParameters(tAWS, dbHostKey)
	dbHost, found = search(tAWS, dbHostKey)

	return
}

func GetDBUser(tAWS string) (dbUser string, found bool) {
	slogger.DebugMethod()

	initParameters(tAWS, dbUserKey)
	dbUser, found = search(tAWS, dbUserKey)

	return
}

func GetDBPort(tAWS string) (dbPort string, found bool) {
	slogger.DebugMethod()

	initParameters(tAWS, dbPortKey)
	dbPort, found = search(tAWS, dbPortKey)

	return
}

func GetDBName(tAWS string) (dbName string, found bool) {
	slogger.DebugMethod()

	initParameters(tAWS, dbNameKey)
	dbName, found = search(tAWS, dbNameKey)

	return
}

func GetDBSSLMode(tAWS string) (dbSSLMode string, found bool) {
	slogger.DebugMethod()

	initParameters(tAWS, dbSSLModeKey)
	dbSSLMode, found = search(tAWS, dbSSLModeKey)

	return
}

func search(tAWS, key string) (value string, found bool) {
	slogger.DebugMethod()

	found = false
	for _, v := range sSSMParameterValues.Parameters {
		if strings.Contains(*v.Name, tAWS) {
			if strings.Contains(*v.Name, key) {
				found = true
				value = *v.Value
				break
			}
		}
	}

	return
}
