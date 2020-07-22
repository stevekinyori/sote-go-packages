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
	// System Manager Parameter Keys
	AWSREGIONIKEY = "AWS_REGION"
	USERPOOLIDKEY = "COGNITO_USER_POOL_ID"
	CLIENTIDKEY   = "COGNITO_CLIENT_ID"
	DBPASSWORDKEY = "DATABASE_PASSWORD"
	DBHOSTKEY     = "DB_HOST"
	DBUSERKEY     = "DB_USERNAME"
	DBPORTKEY     = "DB_PORT"
	DBNAMEKEY     = "DB_NAME"
	DBSSLMODEKEY  = "DB_SSL_MODE"
	// Root Path
	ROOTPATH = "/sote"
)

var (
	awsService *ssm.SSM
	setToTrue        = true       // This can not be a constant because we need a pointer.
	pTrue            = &setToTrue // pointer to the setToTrue variable
	maxResult  int64 = 10
	pMaxResult       = &maxResult
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
This will retrieve the parameters that are in the AWS System Manager service for the ROOTPATH and the supplied
application and environment.  AWS limits the maximum number of parameters to 10 in a single query.  sconfigparams
doesn't support pulling more than the first 10 parameters based on the path
*/
func GetParameters(tApplication, tEnvironment string) (parameters map[string]interface{}, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var pSSMPathOutput *ssm.GetParametersByPathOutput

	parameters = make(map[string]interface{})
	if pSSMPathOutput, soteErr = listParameters(tApplication, strings.ToLower(tEnvironment)); soteErr.ErrCode == nil {
		for _, pParameter := range pSSMPathOutput.Parameters {
			parameters[*pParameter.Name] = *pParameter.Value
		}
	}

	return
}

/*
This will retrieve the database password parameter that is in AWS System Manager service for the ROOTPATH,
application and environment.  Application and environment are required.
*/
func GetDBPassword(tApplication, tEnvironment string) (dbPassword string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if tApplication == "" || tEnvironment == "" {
		soteErr = sError.GetSError(200512, sError.BuildParams([]string{tApplication, tEnvironment}), sError.EmptyMap)
	} else {
		dbPassword, soteErr = getParameter(tApplication, strings.ToLower(tEnvironment), DBPASSWORDKEY)
	}

	return
}

/*
This will retrieve the database host parameter that is in AWS System Manager service for the ROOTPATH and
application.
*/
func GetDBHost(tApplication, tEnvironment string) (dbHost string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if tApplication == "" || tEnvironment == "" {
		soteErr = sError.GetSError(200512, sError.BuildParams([]string{tApplication, tEnvironment}), sError.EmptyMap)
	} else {
		dbHost, soteErr = getParameter(tApplication, strings.ToLower(tEnvironment), DBHOSTKEY)
	}

	return
}

/*
This will retrieve the database user parameter that is in AWS System Manager service for the ROOTPATH and
application.
*/
func GetDBUser(tApplication, tEnvironment string) (dbUser string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if tApplication == "" || tEnvironment == "" {
		soteErr = sError.GetSError(200512, sError.BuildParams([]string{tApplication, tEnvironment}), sError.EmptyMap)
	} else {
		dbUser, soteErr = getParameter(tApplication, strings.ToLower(tEnvironment), DBUSERKEY)
	}

	return
}

/*
This will retrieve the database port parameter that is in AWS System Manager service for the ROOTPATH and
application.
*/
func GetDBPort(tApplication, tEnvironment string) (dbPort string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if tApplication == "" || tEnvironment == "" {
		soteErr = sError.GetSError(200512, sError.BuildParams([]string{tApplication, tEnvironment}), sError.EmptyMap)
	} else {
		dbPort, soteErr = getParameter(tApplication, strings.ToLower(tEnvironment), DBPORTKEY)
	}

	return
}

/*
This will retrieve the database name parameter that is in AWS System Manager service for the ROOTPATH and
application.
*/
func GetDBName(tApplication, tEnvironment string) (dbName string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if tApplication == "" || tEnvironment == "" {
		soteErr = sError.GetSError(200512, sError.BuildParams([]string{tApplication, tEnvironment}), sError.EmptyMap)
	} else {
		dbName, soteErr = getParameter(tApplication, strings.ToLower(tEnvironment), DBNAMEKEY)
	}

	return
}

/*
This will retrieve the database SSL mode parameter that is in AWS System Manager service for the ROOTPATH and
application.
*/
func GetDBSSLMode(tApplication, tEnvironment string) (dbSSLMode string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if tApplication == "" || tEnvironment == "" {
		soteErr = sError.GetSError(200512, sError.BuildParams([]string{tApplication, tEnvironment}), sError.EmptyMap)
	} else {
		dbSSLMode, soteErr = getParameter(tApplication, strings.ToLower(tEnvironment), DBSSLMODEKEY)
	}

	return
}

/*
This will retrieve the AWS Region parameter that is in AWS System Manager service for the ROOTPATH
*/
func GetRegion() (region string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	region, soteErr = getParameter("", "", AWSREGIONIKEY)

	return
}

/*
This will retrieve the cognito user pool id parameter that is in AWS System Manager service for the ROOTPATH and
environment.
*/
func GetUserPoolId(tEnvironment string) (userPoolId string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if tEnvironment == "" {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{tEnvironment}), sError.EmptyMap)
	} else {
		userPoolId, soteErr = getParameter("", strings.ToLower(tEnvironment), USERPOOLIDKEY)
	}

	return
}

/*
This will retrieve the cognito client id for the allocation that is in AWS System Manager service for the ROOTPATH and
environment.
*/
func GetClientId(tApplication, tEnvironment string) (clientId string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if tApplication == "" || tEnvironment == "" {
		soteErr = sError.GetSError(200512, sError.BuildParams([]string{tApplication, tEnvironment}), sError.EmptyMap)
	} else {
		clientId, soteErr = getParameter(tApplication, strings.ToLower(tEnvironment), CLIENTIDKEY)
	}

	return
}

/*
The Environment is validated against 'development', 'staging', 'demo' and 'production' (case sensitive values)
*/
func validateEnvironment(tEnvironment string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	switch tEnvironment {
	case DEVELOPMENT:
	case STAGING:
	case DEMO:
	case PRODUCTION:
	default:
		soteErr = sError.GetSError(601010, sError.BuildParams([]string{tEnvironment}), sError.EmptyMap)
	}

	return
}

/*
This will build the query path based on the ROOTPATH, Application and Environment.
*/
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

/*
This will query up to the first 10 parameters for the ROOTPATH with some combination of application
and environment variable values.  Application and environment can be empty.
*/
func listParameters(tApplication, tEnvironment string) (pSSMPathOutput *ssm.GetParametersByPathOutput, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var err error

	if soteErr = validateEnvironment(tEnvironment); soteErr.ErrCode == nil {
		var (
			path         = setPath(tApplication, tEnvironment)
			ssmPathInput ssm.GetParametersByPathInput
		)

		ssmPathInput.SetPath(path)
		ssmPathInput.Recursive = pTrue
		ssmPathInput.WithDecryption = pTrue
		ssmPathInput.MaxResults = pMaxResult
		// If there are any parameters that matches the path, a result set will be return by the GetParametersByPath call.
		if pSSMPathOutput, err = awsService.GetParametersByPath(&ssmPathInput); len(pSSMPathOutput.Parameters) == 0 {
			soteErr = sError.GetSError(109999, sError.BuildParams([]string{path}), sError.EmptyMap)
		}
	}
	if err != nil {
		log.Fatalln(err)
	}

	return
}

/*
This will query the first 10 parameters for the ROOTPATH with some combination of application
and environment variable values.  Application and environment can be empty.
*/
func getParameter(tApplication, tEnvironment, key string) (returnValue string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var ssmParamInput ssm.GetParameterInput
	ssmParamInput.WithDecryption = pTrue
	name := setPath(tApplication, tEnvironment) + "/" + key
	ssmParamInput.Name = &name

	// If there are any parameters that match the path, a result set will be return by the GetParametersByPath call.
	if pSSMParamOutput, err := awsService.GetParameter(&ssmParamInput); err != nil {
		soteErr = sError.GetSError(109999, sError.BuildParams([]string{name}), sError.EmptyMap)
	} else {
		returnValue = *pSSMParamOutput.Parameter.Value
	}

	return
}
