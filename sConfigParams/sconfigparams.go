/*
This will retrieve any configuration parameter that is used by Sote.  Areas are environment
variables and System Manager Parameters. For System Manager Parameter they must be stored in
the path and (optional) key provided.

RESTRICTIONS:
    AWS functions:
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
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

const (
	// Environment variables
	APPENV    = "APP_ENVIRONMENT"
	AWSREGION = "AWS_REGION"
	// Environments
	DEVELOPMENT = "development" // Also referred to as local
	STAGING     = "staging"
	DEMO        = "demo"
	PRODUCTION  = "production"
	// System Manager Parameter Keys
	AWSREGIONIKEY = "AWS_REGION"
	CLIENTIDKEY   = "COGNITO_CLIENT_ID"
	CREDENTIALS   = "credentials"
	DBHOSTKEY     = "DB_HOST"
	DBNAMEKEY     = "DB_NAME"
	DBPASSWORDKEY = "DATABASE_PASSWORD"
	DBPORTKEY     = "DB_PORT"
	DBSSLMODEKEY  = "DB_SSL_MODE"
	DBUSERKEY     = "DB_USERNAME"
	URL           = "url"
	USERPOOLIDKEY = "COGNITO_USER_POOL_ID"
	// Root Path
	ROOTPATH = "/sote"
)

var (
	awsService     *ssm.SSM
	setToTrue            = true       // This can not be a constant because we need a pointer.
	pTrue                = &setToTrue // pointer to the setToTrue variable
	maxResult      int64 = 10
	pMaxResult           = &maxResult
	myConfigParams configParams
)

type configParams struct {
	Application string
	Environment string
}

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
doesn't support pulling more than the first 10 parameters based on the path.

*/
func GetParameters(tApplication, tEnvironment string) (parameters map[string]interface{}, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var pSSMPathOutput *ssm.GetParametersByPathOutput

	if soteErr = setConfigParamsApplication(tApplication); soteErr.ErrCode == nil {
		if soteErr = setConfigParamsEnvironment(tEnvironment); soteErr.ErrCode == nil {
			parameters = make(map[string]interface{})
			if pSSMPathOutput, soteErr = listParameters(tApplication, strings.ToLower(tEnvironment)); soteErr.ErrCode == nil {
				for _, pParameter := range pSSMPathOutput.Parameters {
					parameters[*pParameter.Name] = *pParameter.Value
				}
			}
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

	var tDBPassword interface{}

	if soteErr = setConfigParamsApplication(tApplication); soteErr.ErrCode == nil {
		if soteErr = setConfigParamsEnvironment(tEnvironment); soteErr.ErrCode == nil {
			tDBPassword, soteErr = getParameter(tApplication, strings.ToLower(tEnvironment), DBPASSWORDKEY)
			if tDBPassword != nil {
				dbPassword = tDBPassword.(string)
			}
		}
	}

	return
}

/*
This will retrieve the database host parameter that is in AWS System Manager service for the ROOTPATH and
application.  Application and environment are required.
*/
func GetDBHost(tApplication, tEnvironment string) (dbHost string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tDBHost interface{}

	if soteErr = setConfigParamsApplication(tApplication); soteErr.ErrCode == nil {
		if soteErr = setConfigParamsEnvironment(tEnvironment); soteErr.ErrCode == nil {
			tDBHost, soteErr = getParameter(tApplication, strings.ToLower(tEnvironment), DBHOSTKEY)
			if tDBHost != nil {
				dbHost = tDBHost.(string)
			}
		}
	}

	return
}

/*
This will retrieve the database user parameter that is in AWS System Manager service for the ROOTPATH and
application.  Application and environment are required.
*/
func GetDBUser(tApplication, tEnvironment string) (dbUser string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tDBUser interface{}

	if soteErr = setConfigParamsApplication(tApplication); soteErr.ErrCode == nil {
		if soteErr = setConfigParamsEnvironment(tEnvironment); soteErr.ErrCode == nil {
			tDBUser, soteErr = getParameter(tApplication, strings.ToLower(tEnvironment), DBUSERKEY)
			if tDBUser != nil {
				dbUser = tDBUser.(string)
			}
		}
	}

	return
}

/*
This will retrieve the database port parameter that is in AWS System Manager service for the ROOTPATH and
application.  Application and environment are required.
*/
func GetDBPort(tApplication, tEnvironment string) (dbPort int, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tDBPort interface{}

	if soteErr = setConfigParamsApplication(tApplication); soteErr.ErrCode == nil {
		if soteErr = setConfigParamsEnvironment(tEnvironment); soteErr.ErrCode == nil {
			tDBPort, soteErr = getParameter(tApplication, strings.ToLower(tEnvironment), DBPORTKEY)
			if tDBPort != nil {
				dbPort, _ = strconv.Atoi(tDBPort.(string))
			} else {
				soteErr = sError.GetSError(109999, sError.BuildParams([]string{DBPORTKEY}), sError.EmptyMap)
			}
		}
	}

	return
}

/*
This will retrieve the database name parameter that is in AWS System Manager service for the ROOTPATH and
application.  Application and environment are required.
*/
func GetDBName(tApplication, tEnvironment string) (dbName string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tDBName interface{}

	if soteErr = setConfigParamsApplication(tApplication); soteErr.ErrCode == nil {
		if soteErr = setConfigParamsEnvironment(tEnvironment); soteErr.ErrCode == nil {
			tDBName, soteErr = getParameter(tApplication, strings.ToLower(tEnvironment), DBNAMEKEY)
			if tDBName != nil {
				dbName = tDBName.(string)
			}
		}
	}

	return
}

/*
This will retrieve the database SSL mode parameter that is in AWS System Manager service for the ROOTPATH and
application.  Application and environment are required.
*/
func GetDBSSLMode(tApplication, tEnvironment string) (dbSSLMode string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tDBSSLMode interface{}

	if soteErr = setConfigParamsApplication(tApplication); soteErr.ErrCode == nil {
		if soteErr = setConfigParamsEnvironment(tEnvironment); soteErr.ErrCode == nil {
			tDBSSLMode, soteErr = getParameter(tApplication, strings.ToLower(tEnvironment), DBSSLMODEKEY)
			if tDBSSLMode != nil {
				dbSSLMode = tDBSSLMode.(string)
			}
		}
	}

	return
}

/*
This will retrieve the AWS Region parameter that is in AWS System Manager service for the ROOTPATH
*/
func GetRegion() (region string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tRegion interface{}

	tRegion, soteErr = getParameter("", "", AWSREGIONIKEY)
	if tRegion != nil {
		region = tRegion.(string)
	}

	return
}

/*
This will retrieve the cognito user pool id parameter that is in AWS System Manager service for the ROOTPATH and
environment.  Environment are required.
*/
func GetUserPoolId(tEnvironment string) (userPoolId string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tUserPoolId interface{}

	if soteErr = setConfigParamsEnvironment(tEnvironment); soteErr.ErrCode == nil {
		tUserPoolId, soteErr = getParameter("", strings.ToLower(tEnvironment), USERPOOLIDKEY)
		if tUserPoolId != nil {
			userPoolId = tUserPoolId.(string)
		}
	}

	return
}

/*
This will retrieve the cognito client id for the allocation that is in AWS System Manager service for the ROOTPATH and
environment.  Application and environment are required.
*/
func GetClientId(tApplication, tEnvironment string) (clientId string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tClientId interface{}

	if soteErr = setConfigParamsApplication(tApplication); soteErr.ErrCode == nil {
		if soteErr = setConfigParamsEnvironment(tEnvironment); soteErr.ErrCode == nil {
			tClientId, soteErr = getParameter(tApplication, strings.ToLower(tEnvironment), CLIENTIDKEY)
			if tClientId != nil {
				clientId = tClientId.(string)
			}
		}
	}

	return
}

/*
This will retrieve the messaging credentials needed to authenticate that is in AWS System Manager service for the ROOTPATH and
environment.
*/
func GetNATSCredentials() (natsCredentials func(string, string) (interface{}, sError.SoteError)) {
	sLogger.DebugMethod()

	natsCredentials = getCreds()

	return
}

func getCreds() func(string, string) (interface{}, sError.SoteError) {
	return func(tApplication, tEnvironment string) (natsCredentials interface{}, soteErr sError.SoteError) {
		if soteErr = setConfigParamsApplication(tApplication); soteErr.ErrCode == nil {
			if soteErr = setConfigParamsEnvironment(tEnvironment); soteErr.ErrCode == nil {
				natsCredentials, soteErr = getParameter(tApplication, strings.ToLower(tEnvironment), CREDENTIALS)
			}
		}
		return
	}
}

/*
This will retrieve the messaging server URL needed to connect that is in AWS System Manager service for the ROOTPATH and
environment.  Application and environment are required.
*/
func GetNATSURL(tApplication, tEnvironment string) (natsURL string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tNatsURL interface{}

	if soteErr = setConfigParamsApplication(tApplication); soteErr.ErrCode == nil {
		if soteErr = setConfigParamsEnvironment(tEnvironment); soteErr.ErrCode == nil {
			tNatsURL, soteErr = getParameter(tApplication, strings.ToLower(tEnvironment), URL)
			if tNatsURL != nil {
				natsURL = tNatsURL.(string)
			}
		}
	}

	return
}

/*
The Environment is validated against 'development', 'staging', 'demo' and 'production'. The value supplied
will be forced to lower case.  Environment are required.
*/
func ValidateEnvironment(tEnvironment string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	switch strings.ToLower(tEnvironment) {
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
This will get the AWS Region that is set in the environment variables. If the environment variable is not found or the value is empty,
the function will return an error code for not found.
*/
func GetEnvironmentAWSRegion() (envValue string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	envValue, soteErr = GetEnvironmentVariable(AWSREGIONIKEY)

	return
}

/*
This will get the AWS Region that is set in the environment variables. If the environment variable is not found or the value is empty,
the function will return an error code for not found.
*/
func GetEnvironmentAppEnvironment() (envValue string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	envValue, soteErr = GetEnvironmentVariable(APPENV)

	return
}

/*
Get the requested environment variable. If the environment variable is not found or the value is empty,
the function will return an error code for not found.
*/
func GetEnvironmentVariable(key string) (envValue string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	envValue = os.Getenv(key)
	if envValue = os.Getenv(key); len(envValue) == 0 {
		soteErr = sError.GetSError(109999, sError.BuildParams([]string{key}), sError.EmptyMap)
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

	if soteErr = ValidateEnvironment(tEnvironment); soteErr.ErrCode == nil {
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
func getParameter(tApplication, tEnvironment, key string) (returnValue interface{}, soteErr sError.SoteError) {
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

func setConfigParamsApplication(tApplication string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if tApplication == "" {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{tApplication}), sError.EmptyMap)
	} else {
		myConfigParams.Application = tApplication
	}

	return
}

func GetConfigParamsApplication() string {
	sLogger.DebugMethod()

	return myConfigParams.Application
}

func setConfigParamsEnvironment(tEnvironment string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = ValidateEnvironment(tEnvironment); soteErr.ErrCode == nil {
		myConfigParams.Environment = tEnvironment
	}

	return
}

func GetConfigParamsEnvironment() string {
	sLogger.DebugMethod()

	return myConfigParams.Environment
}
