package sConfigParams

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"gitlab.com/soteapps/packages/slogger"
)

const (
	// Environments
	staging     = "staging"
	development = "development"
	demo        = "demo"
	production  = "production"
	// Variables
	dbNameParam = "DATABASE_PASSWORD"
)

var (
	defaultPath                 = "/sote/api"
	dPath               *string = &defaultPath
	setToTrue           bool    = true
	sTrue               *bool   = &setToTrue
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

func init() {
	slogger.DebugMethod()

	var (
		ssmPath ssm.GetParametersByPathInput
		err     error
	)
	ssmPath.Path = dPath
	ssmPath.Recursive = sTrue
	ssmPath.WithDecryption = sTrue
	sSSMParameterValues, err = awsService.GetParametersByPath(&ssmPath)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetDBName(tAWS string) (dbName string) {
	slogger.DebugMethod()

	found := false
	for _, v := range sSSMParameterValues.Parameters {
		if strings.Contains(*v.Name, tAWS) {
			if strings.Contains(*v.Name, dbNameParam) {
				found = true
				fmt.Println(v)
			}
		}
	}

	if !found {

	}
	return
}
