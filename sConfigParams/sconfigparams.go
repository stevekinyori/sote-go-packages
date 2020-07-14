/*
This will query System Manager Parameter store for the path and (optional) key provided.

RESTRICTIONS:
    * Program must have access to a .aws/credentials file in the default locate.
    * This will only access system parameters that start with '/sote' (ROOTPATH).
    * You can only request a single key per call

NOTES:
    When the filter is not found in the result set from the GetParametersByPath call, the whole result
    set is returned.
*/
package sConfigParams

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

const (
	// Environments
	STAGING     = "staging"
	DEVELOPMENT = "development"
	DEMO        = "demo"
	PRODUCTION  = "production"
	//
	AWSREGIONIKEY = "AWS_REGION"
	USERPOOLIDKEY = "COGNITO_USER_POOL_ID"
	DBPASSWORDKEY = "DATABASE_PASSWORD"
	DBHOSTKEY     = "DB_HOST"
	DBUSERKEY     = "DB_USERNAME"
	DBPORTKEY     = "DB_PORT"
	DBNAMEKEY     = "DB_NAME"
	DBSSLMODEKEY  = "DB_SSL_MODE"
	NOKEY         = ""
	ROOTPATH      = "/sote"
)

var (
	awsService     *ssm.SSM
	setToTrue            = true       // This can not be a constant because we need a pointer.
	pTrue                = &setToTrue // pointer to the setToTrue variable
	maxResult      int64 = 10
	pMaxResult           = &maxResult
	pSSMPathOutput *ssm.GetParametersByPathOutput
)

/*
This will establish a session using the default .aws location
*/
func init() {
	sLogger.DebugMethod()

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatalln(err)
	}
	awsService = ssm.New(sess)
}

/*
This will build the query path based on the ROOTPATH, Application and Environment.  The
Environment is validated against 'development', 'staging', 'demo' and 'production' (case sensitive values)
*/
func initParameters(tApplication, tEnvironment, key string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var err error

	if tEnvironment != "" {
		switch tEnvironment {
		case DEVELOPMENT:
		case STAGING:
		case DEMO:
		case PRODUCTION:
		default:
			soteErr = sError.GetSError(601010, sError.BuildParams([]string{tEnvironment}), sError.EmptyMap)
		}
	}

	if soteErr.ErrCode == nil {
		var (
			path         string
			ssmPathInput ssm.GetParametersByPathInput
			filter       ssm.ParameterStringFilter
			pFilter      []*ssm.ParameterStringFilter
			pKeys        []*string
		)
		path = setPath(tApplication, tEnvironment)

		ssmPathInput.SetPath(path)
		ssmPathInput.Recursive = pTrue
		ssmPathInput.WithDecryption = pTrue
		ssmPathInput.MaxResults = pMaxResult
		if len(key) > 0 {
			pKeys = append(pKeys, &key)
			filter.SetKey("string")
			filter.SetOption("string")
			filter.SetValues(pKeys)
			ssmPathInput.SetParameterFilters(pFilter)
		}
		// If there are any parameters that match the path, a result set will be return by the GetParametersByPath call.
		if pSSMPathOutput, err = awsService.GetParametersByPath(&ssmPathInput); len(pSSMPathOutput.Parameters) == 0 {
			soteErr = sError.GetSError(109999, sError.BuildParams([]string{path}), sError.EmptyMap)
		} else {
			// This is used for testing to display the output
			fmt.Println(pSSMPathOutput)
		}
	}
	if err != nil {
		log.Fatalln(err)
	}

	return
}

/*
This will retrieve the parameters that are in the AWS System Manager service for the ROOTPATH and the supplied
application and environment.  AWS limits the maximum number of parameters to 10 in a single query.  sconfigparams
doesn't support pulling more than the first 10 parameters.
*/
func GetParameters(tApplication, tEnvironment string) (parameters map[string]interface{}, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	parameters = make(map[string]interface{})
	if soteErr = initParameters(tApplication, strings.ToLower(tEnvironment), NOKEY); soteErr.ErrCode == nil {
		for _, pParameter := range pSSMPathOutput.Parameters {
			parameters[*pParameter.Name] = *pParameter.Value
		}
	}

	return
}

/*
This will retrieve the database password parameter that is in AWS System Manager service for the ROOTPATH and
application.
*/
func GetDBPassword(tApplication, tEnvironment string) (dbPassword string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = initParameters(tApplication, strings.ToLower(tEnvironment), DBPASSWORDKEY); soteErr.ErrCode == nil {
		path := ROOTPATH + "/" + tApplication + "/" + tEnvironment + "/" + DBPASSWORDKEY
		if name := *pSSMPathOutput.Parameters[0].Name; name == path {
			dbPassword = *pSSMPathOutput.Parameters[0].Value
		} else {
			soteErr = sError.GetSError(109999, sError.BuildParams([]string{path}), sError.EmptyMap)
		}
	}

	return
}

/*
This will retrieve the database host parameter that is in AWS System Manager service for the ROOTPATH and
application.
*/
func GetDBHost(tApplication, tEnvironment string) (dbHost string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = initParameters(tApplication, strings.ToLower(tEnvironment), DBHOSTKEY); soteErr.ErrCode == nil {
		path := ROOTPATH + "/" + tApplication + "/" + tEnvironment + "/" + DBHOSTKEY
		if name := *pSSMPathOutput.Parameters[0].Name; name == path {
			dbHost = *pSSMPathOutput.Parameters[0].Value
		} else {
			soteErr = sError.GetSError(109999, sError.BuildParams([]string{path}), sError.EmptyMap)
		}
	}

	return
}

/*
This will retrieve the database user parameter that is in AWS System Manager service for the ROOTPATH and
application.
*/
func GetDBUser(tApplication, tEnvironment string) (dbUser string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = initParameters(tApplication, strings.ToLower(tEnvironment), DBUSERKEY); soteErr.ErrCode == nil {
		path := ROOTPATH + "/" + tApplication + "/" + tEnvironment + "/" + DBUSERKEY
		if name := *pSSMPathOutput.Parameters[0].Name; name == path {
			dbUser = *pSSMPathOutput.Parameters[0].Value
		} else {
			soteErr = sError.GetSError(109999, sError.BuildParams([]string{path}), sError.EmptyMap)
		}
	}

	return
}

/*
This will retrieve the database port parameter that is in AWS System Manager service for the ROOTPATH and
application.
*/
func GetDBPort(tApplication, tEnvironment string) (dbPort string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = initParameters(tApplication, strings.ToLower(tEnvironment), DBPORTKEY); soteErr.ErrCode == nil {
		path := ROOTPATH + "/" + tApplication + "/" + tEnvironment + "/" + DBPORTKEY
		if name := *pSSMPathOutput.Parameters[0].Name; name == path {
			dbPort = *pSSMPathOutput.Parameters[0].Value
		} else {
			soteErr = sError.GetSError(109999, sError.BuildParams([]string{path}), sError.EmptyMap)
		}
	}

	return
}

/*
This will retrieve the database name parameter that is in AWS System Manager service for the ROOTPATH and
application.
*/
func GetDBName(tApplication, tEnvironment string) (dbName string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = initParameters(tApplication, strings.ToLower(tEnvironment), DBNAMEKEY); soteErr.ErrCode == nil {
		path := ROOTPATH + "/" + tApplication + "/" + tEnvironment + "/" + DBNAMEKEY
		if name := *pSSMPathOutput.Parameters[0].Name; name == path {
			dbName = *pSSMPathOutput.Parameters[0].Value
		} else {
			soteErr = sError.GetSError(109999, sError.BuildParams([]string{path}), sError.EmptyMap)
		}
	}

	return
}

/*
This will retrieve the database SSL mode parameter that is in AWS System Manager service for the ROOTPATH and
application.
*/
func GetDBSSLMode(tApplication, tEnvironment string) (dbSSLMode string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = initParameters(tApplication, tEnvironment, DBSSLMODEKEY); soteErr.ErrCode == nil {
		path := ROOTPATH + "/" + tApplication + "/" + tEnvironment + "/" + DBSSLMODEKEY
		if name := *pSSMPathOutput.Parameters[0].Name; name == path {
			dbSSLMode = *pSSMPathOutput.Parameters[0].Value
		} else {
			soteErr = sError.GetSError(109999, sError.BuildParams([]string{path}), sError.EmptyMap)
		}
	}

	return
}

/*
This will retrieve the AWS Region parameter that is in AWS System Manager service for the ROOTPATH
*/
func GetRegion() (region string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = initParameters("", "", AWSREGIONIKEY); soteErr.ErrCode == nil {
		path := ROOTPATH + "/" + AWSREGIONIKEY
		if name := *pSSMPathOutput.Parameters[0].Name; name == path {
			region = *pSSMPathOutput.Parameters[0].Value
		} else {
			soteErr = sError.GetSError(109999, sError.BuildParams([]string{path}), sError.EmptyMap)
		}
	}

	return
}

/*
This will retrieve the cognito user pool id parameter that is in AWS System Manager service for the ROOTPATH and
environment.
*/
func GetUserPoolId(tEnvironment string) (userPoolId string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = initParameters("", tEnvironment, USERPOOLIDKEY); soteErr.ErrCode == nil {
		path := ROOTPATH + "/" + tEnvironment + "/" + USERPOOLIDKEY
		if name := *pSSMPathOutput.Parameters[0].Name; name == path {
			userPoolId = *pSSMPathOutput.Parameters[0].Value
		} else {
			soteErr = sError.GetSError(109999, sError.BuildParams([]string{path}), sError.EmptyMap)
		}
	}

	return
}

func setPath(tApplication, tEnvironment string) (path string) {
	sLogger.DebugMethod()

	if tApplication == "" && tEnvironment == "" {
		path = ROOTPATH
	} else {
		if tApplication == "" && tEnvironment != "" {
			path = ROOTPATH + "/" + tEnvironment
		} else {
			if tApplication != "" && tEnvironment == "" {
				path = ROOTPATH + "/" + tApplication
			} else {
				path = ROOTPATH + "/" + tApplication + "/" + tEnvironment
			}
		}
	}

	return
}
