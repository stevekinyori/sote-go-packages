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
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"gitlab.com/soteapps/packages/v2022/sError"
	"gitlab.com/soteapps/packages/v2022/sLogger"
)

const (
	// Environment variables
	APPENV = "APP_ENVIRONMENT"
	// Environments
	DEVELOPMENT = "development" // Also referred to as local
	STAGING     = "staging"
	DEMO        = "demo"
	PRODUCTION  = "production"
	// System Manager Parameter Keys
	AWSACCOUNTIDKEY              = "AWS_ACCOUNT_ID"
	AWSREGIONIKEY                = "AWS_REGION"
	AWSS3BUCKETKEY               = "AWS_S3_BUCKET"
	CLIENTIDKEY                  = "COGNITO_CLIENT_ID"
	COGNITOUSER                  = "USER"
	COGNITOCLIENTID              = "CLIENT_ID"
	COGNITOPASSWORD              = "DATA_LOAD_PASSWORD"
	CREDENTIALS                  = "credentials"
	DBHOSTKEY                    = "DB_HOST"
	DBNAMEKEY                    = "DB_NAME"
	DBPASSWORDKEY                = "DATABASE_PASSWORD"
	DBPORTKEY                    = "DB_PORT"
	DBSSLMODEKEY                 = "DB_SSL_MODE"
	DBUSERKEY                    = "DB_USERNAME"
	URL                          = "url"
	TLSURLMASK                   = "tls-urlmask"
	UNPROCESSEDDOCUMENTSKEY      = "inbound/name"
	PROCESSEDDOCUMENTSKEY        = "processed/name"
	USERPOOLIDKEY                = "COGNITO_USER_POOL_ID"
	SMTPUSERNAME                 = "USERNAME"
	SMTPPASSWORD                 = "PASSWORD"
	SMTPPORT                     = "PORT"
	SMTPHOST                     = "HOST"
	QUICKBOOKSCLIENTID           = "CLIENT_ID"
	QUICKBOOKSCLIENTSECRET       = "CLIENT_SECRET"
	QUICKBOOKSWEBHOOKTOKEN       = "WEBHOOK_TOKEN"
	QUICKBOOKSHOST               = "HOST"
	QUICKBOOKSCONFIGURL          = "CONFIG_URL"
	QUICKBOOKSREFRESHTOKEN       = "REFRESH_TOKEN"
	QUICKBOOKSREFRESHTOKENEXPIRY = "REFRESH_TOKEN_EXPIRY"

	// Application values
	API        string = "api"
	SDCC       string = "sdcc"
	SYNADIA    string = "synadia"
	DOCUMENTS  string = "documents"
	QUICKBOOKS        = "quickbooks"
	SMTP              = "smtp"
	COGNITO           = "cognito"
	// Root Path
	ROOTPATH = "/sote"
)

var (
	awsService *ssm.Client
	pTrue            = true // pointer to the setToTrue variable
	pMaxResult int32 = 10
)

type SMTPConfig struct {
	Host     string
	UserName string
	Password string
	Port     string
}

type QuickbooksConfig struct {
	Host               string
	ClientId           string
	ClientSecret       string
	WebhookToken       string
	ConfigURL          string
	RefreshToken       string
	RefreshTokenExpiry string
}

type QuickBooksRefreshToken struct {
	Token      string
	ExpiryDate time.Time
}

type CognitoConfig struct {
	ClientId string
	UserName string
	Password string
}

type SSMParameter struct {
	Key        string
	Content    string
	TargetType types.ParameterType
}

type Database struct {
	Name     string
	User     string
	Password string
	Host     string
	SSLMode  string
	Port     int
}

/*
This will establish a session using the default .aws location
*/
func init() {
	sLogger.DebugMethod()

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalln(err)
	}
	awsService = ssm.NewFromConfig(cfg)
}

/*
GetParameters will retrieve the parameters that are in the AWS System Manager service for the ROOTPATH and the supplied
application and environment.  AWS limits the maximum number of parameters to 10 in a single query.  sconfigparams
doesn't support pulling more than the first 10 parameters based on the path.
*/
func GetParameters(ctx context.Context, application, environment string) (parameters map[string]interface{}, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var pSSMPathOutput *ssm.GetParametersByPathOutput

	if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
		if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
			parameters = make(map[string]interface{})
			if pSSMPathOutput, soteErr = listParameters(ctx, application, strings.ToLower(environment)); soteErr.ErrCode == nil {
				for _, pParameter := range pSSMPathOutput.Parameters {
					parameters[*pParameter.Name] = *pParameter.Value
				}
			}
		}
	}

	return
}

// GetSMTPConfig retrieves all SMTP configurations  from SSM
func GetSMTPConfig(ctx context.Context, application, environment string) (parameters *SMTPConfig, soteErr sError.SoteError) {
	var (
		smtpUserNameKey  = setPath(application, environment) + "/" + SMTPUSERNAME
		smtpPasswordKey  = setPath(application, environment) + "/" + SMTPPASSWORD
		smtpHostKey      = setPath(application, "") + "/" + SMTPHOST
		smtpPortKey      = setPath(application, "") + "/" + SMTPPORT
		pSSMParamsOutput = &ssm.GetParametersOutput{}
		err              error
	)

	parameters = &SMTPConfig{}
	if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
		if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
			environment = strings.ToLower(environment)
			pSSMParamsInput := &ssm.GetParametersInput{
				Names:          []string{smtpUserNameKey, smtpPasswordKey, smtpHostKey, smtpPortKey},
				WithDecryption: pTrue,
			}
			if pSSMParamsOutput, err = awsService.GetParameters(ctx, pSSMParamsInput); err == nil {
				if len(pSSMParamsOutput.Parameters) < len(pSSMParamsInput.Names) {
					soteErr = sError.GetSError(109999, sError.BuildParams([]string{"smtp configuration"}), sError.EmptyMap)
				} else {
					for _, pParameter := range pSSMParamsOutput.Parameters {
						switch *pParameter.Name {
						case smtpUserNameKey:
							parameters.UserName = *pParameter.Value
						case smtpPasswordKey:
							parameters.Password = *pParameter.Value
						case smtpHostKey:
							parameters.Host = *pParameter.Value
						case smtpPortKey:
							parameters.Port = *pParameter.Value
						}
					}
				}
			} else {
				soteErr = sError.GetSError(199999, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
			}
		}
	}

	return
}

// GetQuickbooksConfig retrieves all Quickbooks configurations  from SSM
func GetQuickbooksConfig(ctx context.Context, application, environment string) (parameters *QuickbooksConfig, soteErr sError.SoteError) {
	var (
		clientIdKey        = setPath(application, environment) + "/" + QUICKBOOKSCLIENTID
		clientSecretKey    = setPath(application, environment) + "/" + QUICKBOOKSCLIENTSECRET
		hostKey            = setPath(application, environment) + "/" + QUICKBOOKSHOST
		configURLKey       = setPath(application, environment) + "/" + QUICKBOOKSCONFIGURL
		webhookToken       = setPath(application, environment) + "/" + QUICKBOOKSWEBHOOKTOKEN
		refreshToken       = setPath(application, environment) + "/" + QUICKBOOKSREFRESHTOKEN
		refreshTokenExpiry = setPath(application, environment) + "/" + QUICKBOOKSREFRESHTOKENEXPIRY
		pSSMParamsOutput   = &ssm.GetParametersOutput{}
		err                error
	)

	parameters = &QuickbooksConfig{}
	if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
		if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
			environment = strings.ToLower(environment)
			pSSMParamsInput := &ssm.GetParametersInput{
				Names:          []string{clientIdKey, clientSecretKey, hostKey, configURLKey, webhookToken, refreshToken, refreshTokenExpiry},
				WithDecryption: pTrue,
			}
			if pSSMParamsOutput, err = awsService.GetParameters(ctx, pSSMParamsInput); err == nil {
				if len(pSSMParamsOutput.Parameters) < len(pSSMParamsInput.Names) {
					soteErr = sError.GetSError(109999, sError.BuildParams([]string{"quickbooks configuration"}), sError.EmptyMap)
				} else {
					for _, pParameter := range pSSMParamsOutput.Parameters {
						switch *pParameter.Name {
						case clientIdKey:
							parameters.ClientId = *pParameter.Value
						case clientSecretKey:
							parameters.ClientSecret = *pParameter.Value
						case hostKey:
							parameters.Host = *pParameter.Value
						case configURLKey:
							parameters.ConfigURL = *pParameter.Value
						case webhookToken:
							parameters.WebhookToken = *pParameter.Value
						case refreshToken:
							parameters.RefreshToken = *pParameter.Value
						case refreshTokenExpiry:
							parameters.RefreshTokenExpiry = *pParameter.Value
						}
					}
				}
			} else {
				soteErr = sError.GetSError(199999, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
			}
		}
	}

	return
}

// GetCognitoConfig retrieves all Cognito configurations  from SSM
func GetCognitoConfig(ctx context.Context, application, environment string) (parameters *CognitoConfig, soteErr sError.SoteError) {
	var (
		clientIdKey      = setPath(application, environment) + "/" + COGNITOCLIENTID
		userKey          = setPath(application, environment) + "/" + COGNITOUSER
		passwordKey      = setPath(application, environment) + "/" + COGNITOPASSWORD
		pSSMParamsOutput = &ssm.GetParametersOutput{}
		err              error
	)

	parameters = &CognitoConfig{}
	if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
		if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
			environment = strings.ToLower(environment)
			pSSMParamsInput := &ssm.GetParametersInput{
				Names:          []string{clientIdKey, userKey, passwordKey},
				WithDecryption: pTrue,
			}
			if pSSMParamsOutput, err = awsService.GetParameters(ctx, pSSMParamsInput); err == nil {
				if len(pSSMParamsOutput.Parameters) < len(pSSMParamsInput.Names) {
					soteErr = sError.GetSError(109999, sError.BuildParams([]string{"cognito configuration"}), sError.EmptyMap)
				} else {
					for _, pParameter := range pSSMParamsOutput.Parameters {
						switch *pParameter.Name {
						case clientIdKey:
							parameters.ClientId = *pParameter.Value
						case userKey:
							parameters.UserName = *pParameter.Value
						case passwordKey:
							parameters.Password = *pParameter.Value
						}
					}
				}
			} else {
				soteErr = sError.GetSError(199999, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
			}
		}
	}

	return
}

/*
GetSmtpUsername will retrieve the SMTP username parameter that is in AWS System Manager service for the ROOTPATH,
application and environment.  Application and environment are required.
*/
func GetSmtpUsername(ctx context.Context, application, environment string) (smtpUsername string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		tSmtpUsername interface{}
	)

	if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
		if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
			tSmtpUsername, soteErr = getParameter(ctx, application, strings.ToLower(environment), SMTPUSERNAME)
			if tSmtpUsername != nil {
				smtpUsername = tSmtpUsername.(string)
			}
		}
	}

	return
}

/*
GetSmtpPassword will retrieve the SMTP password parameter that is in AWS System Manager service for the ROOTPATH,
application and environment.  Application and environment are required.
*/
func GetSmtpPassword(ctx context.Context, application, environment string) (smtpPassword string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		tSmtpPassword interface{}
	)

	if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
		if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
			tSmtpPassword, soteErr = getParameter(ctx, application, strings.ToLower(environment), SMTPPASSWORD)
			if tSmtpPassword != nil {
				smtpPassword = tSmtpPassword.(string)
			}
		}
	}

	return
}

func GetAWSParams(ctx context.Context, application, environment string) (parameters *Database, soteErr sError.SoteError) {
	var (
		nameKey          = setPath(application, environment) + "/" + DBNAMEKEY
		userKey          = setPath(application, environment) + "/" + DBUSERKEY
		passwordKey      = setPath(application, environment) + "/" + DBPASSWORDKEY
		hostKey          = setPath(application, environment) + "/" + DBHOSTKEY
		sslModeKey       = setPath(application, environment) + "/" + DBSSLMODEKEY
		portKey          = setPath(application, environment) + "/" + DBPORTKEY
		pSSMParamsOutput = &ssm.GetParametersOutput{}
		err              error
	)

	parameters = &Database{}
	if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
		if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
			environment = strings.ToLower(environment)
			pSSMParamsInput := &ssm.GetParametersInput{
				Names:          []string{nameKey, userKey, passwordKey, hostKey, sslModeKey, portKey},
				WithDecryption: pTrue,
			}
			if pSSMParamsOutput, err = awsService.GetParameters(ctx, pSSMParamsInput); err == nil {
				if len(pSSMParamsOutput.Parameters) < len(pSSMParamsInput.Names) {
					soteErr = sError.GetSError(109999, sError.BuildParams([]string{"cognito configuration"}), sError.EmptyMap)
				} else {
					for _, pParameter := range pSSMParamsOutput.Parameters {
						switch *pParameter.Name {
						case nameKey:
							parameters.Name = *pParameter.Value
						case userKey:
							parameters.User = *pParameter.Value
						case passwordKey:
							parameters.Password = *pParameter.Value
						case hostKey:
							parameters.Host = *pParameter.Value
						case sslModeKey:
							parameters.SSLMode = *pParameter.Value
						case portKey:
							if parameters.Port, err = strconv.Atoi(*pParameter.Value); err != nil {
								break
							}
						}
					}
				}
			}

			if err != nil {
				soteErr = sError.GetSError(199999, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
			}
		}
	}

	return
}

/*
GetDBPassword will retrieve the database password parameter that is in AWS System Manager service for the ROOTPATH,
application and environment.  Application and environment are required.
*/
func GetDBPassword(ctx context.Context, application, environment string) (dbPassword string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tDBPassword interface{}

	if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
		if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
			tDBPassword, soteErr = getParameter(ctx, application, strings.ToLower(environment), DBPASSWORDKEY)
			if tDBPassword != nil {
				dbPassword = tDBPassword.(string)
			}
		}
	}

	return
}

/*
GetDBHost will retrieve the database host parameter that is in AWS System Manager service for the ROOTPATH and
application.  Application and environment are required.
*/
func GetDBHost(ctx context.Context, application, environment string) (dbHost string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tDBHost interface{}

	if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
		if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
			tDBHost, soteErr = getParameter(ctx, application, strings.ToLower(environment), DBHOSTKEY)
			if tDBHost != nil {
				dbHost = tDBHost.(string)
			}
		}
	}

	return
}

/*
GetDBUser will retrieve the database user parameter that is in AWS System Manager service for the ROOTPATH and
application.  Application and environment are required.
*/
func GetDBUser(ctx context.Context, application, environment string) (dbUser string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tDBUser interface{}

	if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
		if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
			tDBUser, soteErr = getParameter(ctx, application, strings.ToLower(environment), DBUSERKEY)
			if tDBUser != nil {
				dbUser = tDBUser.(string)
			}
		}
	}

	return
}

/*
GetDBPort will retrieve the database port parameter that is in AWS System Manager service for the ROOTPATH and
application.  Application and environment are required.
*/
func GetDBPort(ctx context.Context, application, environment string) (dbPort int, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tDBPort interface{}

	if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
		if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
			tDBPort, soteErr = getParameter(ctx, application, strings.ToLower(environment), DBPORTKEY)
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
GetDBName will retrieve the database name parameter that is in AWS System Manager service for the ROOTPATH and
application.  Application and environment are required.
*/
func GetDBName(ctx context.Context, application, environment string) (dbName string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tDBName interface{}

	if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
		if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
			tDBName, soteErr = getParameter(ctx, application, strings.ToLower(environment), DBNAMEKEY)
			if tDBName != nil {
				dbName = tDBName.(string)
			}
		}
	}

	return
}

/*
GetDBSSLMode will retrieve the database SSL mode parameter that is in AWS System Manager service for the ROOTPATH and
application.  Application and environment are required.
*/
func GetDBSSLMode(ctx context.Context, application, environment string) (dbSSLMode string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tDBSSLMode interface{}

	if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
		if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
			tDBSSLMode, soteErr = getParameter(ctx, application, strings.ToLower(environment), DBSSLMODEKEY)
			if tDBSSLMode != nil {
				dbSSLMode = tDBSSLMode.(string)
			}
		}
	}

	return
}

/*
GetRegion will retrieve the AWS Region parameter that is in AWS System Manager service for the ROOTPATH
*/
func GetRegion(ctx context.Context) (region string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tRegion interface{}

	tRegion, soteErr = getParameter(ctx, "", "", AWSREGIONIKEY)
	if tRegion != nil {
		region = tRegion.(string)
	}

	return
}

/*
GetUserPoolId will retrieve the cognito user pool id parameter that is in AWS System Manager service for the ROOTPATH and
environment.  Environment are required.
*/
func GetUserPoolId(ctx context.Context, environment string) (userPoolId string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tUserPoolId interface{}

	if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
		tUserPoolId, soteErr = getParameter(ctx, "", strings.ToLower(environment), USERPOOLIDKEY)
		if tUserPoolId != nil {
			userPoolId = tUserPoolId.(string)
		}
	}

	return
}

/*
GetClientId will retrieve the cognito client id for the allocation that is in AWS System Manager service for the ROOTPATH and
environment.  Application and environment are required.
*/
func GetClientId(ctx context.Context, clientName, environment string) (clientId string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tClientId interface{}

	if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
		tClientId, soteErr = getParameter(ctx, clientName, strings.ToLower(environment), CLIENTIDKEY)
		if tClientId != nil {
			clientId = tClientId.(string)
		}
	}

	return
}

/*
GetNATSCredentials will retrieve the messaging credentials needed to authenticate that is in AWS System Manager service for the ROOTPATH and
environment.
*/
func GetNATSCredentials(ctx context.Context) (natsCredentials func(string, string) (interface{}, sError.SoteError)) {
	sLogger.DebugMethod()

	natsCredentials = getCreds(ctx)

	return
}

func getCreds(ctx context.Context) func(string, string) (interface{}, sError.SoteError) {
	return func(application, environment string) (natsCredentials interface{}, soteErr sError.SoteError) {
		if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
			if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
				natsCredentials, soteErr = getParameter(ctx, application, strings.ToLower(environment), CREDENTIALS)
			}
		}
		return
	}
}

/*
GetNATSURL will retrieve the messaging server URL needed to connect that is in AWS System Manager service for the ROOTPATH and
environment.  Application and environment are required.
*/
func GetNATSURL(ctx context.Context, application, environment string) (natsURL string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tNatsURL interface{}

	if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
		if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
			tNatsURL, soteErr = getParameter(ctx, application, strings.ToLower(environment), URL)
			if tNatsURL != nil {
				natsURL = tNatsURL.(string)
			}
		}
	}

	return
}

/*
GetNATSTLSURLMask will retrieve the messaging server TLS URL mask needed. Application is required.
*/
func GetNATSTLSURLMask(ctx context.Context, application string) (natsTLSURLMask string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tNATSTLSURLMask interface{}

	if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
		tNATSTLSURLMask, soteErr = getParameter(ctx, application, "", TLSURLMASK)
		if tNATSTLSURLMask != nil {
			natsTLSURLMask = tNATSTLSURLMask.(string)
		}
	}

	return
}

/*
GetAWSS3Bucket will retrieve the AWS S3 Bucket parameter that is in AWS System Manager service for the ROOTPATH
*/
func GetAWSS3Bucket(ctx context.Context, application string) (AWSS3Bucket string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tAWSS3Bucket interface{}

	tAWSS3Bucket, soteErr = getParameter(ctx, application, "", AWSS3BUCKETKEY)
	if tAWSS3Bucket != nil {
		AWSS3Bucket = tAWSS3Bucket.(string)
	}

	return
}

/*
GetAWSAccountId will retrieve the AWS Client ID parameter that is in AWS System Manager service for the ROOTPATH
*/
func GetAWSAccountId(ctx context.Context) (AWSAccountId string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tAWSAccountId interface{}

	tAWSAccountId, soteErr = getParameter(ctx, "", "", AWSACCOUNTIDKEY)
	if tAWSAccountId != nil {
		AWSAccountId = tAWSAccountId.(string)
	}

	return
}

/*
	SGetS3BucketURL will retrieve the AWS S3 server URL found in AWS System Manager service for the ROOTPATH and
	environment. The URL is needed to access Sote's unprocessed/ processed documents.  Application,
	environment and key are required.
*/
func SGetS3BucketURL(ctx context.Context, application, environment, key string) (sS3BucketURL string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var tS3BucketURL interface{}
	if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
		if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
			if tS3BucketURL, soteErr = getParameter(ctx, application, strings.ToLower(environment), key); tS3BucketURL != nil {
				sS3BucketURL = tS3BucketURL.(string)
			}
		}
	}

	return
}

/*
ValidateApplication is validated against empty string.  Application is required.
*/
func ValidateApplication(application string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if application == "" {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{application}), sError.EmptyMap)
	}

	return
}

/*
ValidateEnvironment is validated against 'development', 'staging', 'demo' and 'production'. The value supplied
will be forced to lower case.  Environment are required.
*/
func ValidateEnvironment(environment string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	switch strings.ToLower(environment) {
	case DEVELOPMENT, STAGING, DEMO, PRODUCTION:
		return
	default:
		soteErr = sError.GetSError(209110, sError.BuildParams([]string{environment}), sError.EmptyMap)
	}

	return
}

/*
GetEnvironmentAppEnvironment will get the AWS Region that is set in the environment variables. If the environment variable is not found or the value is empty,
the function will return an error code for not found.
*/
func GetEnvironmentAppEnvironment() (envValue string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	envValue, soteErr = GetEnvironmentVariable(APPENV)

	return
}

/*
GetEnvironmentVariable the requested environment variable. If the environment variable is not found or the value is empty,
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

// UpdateQuickbooksRefreshToken updates the SSM value of quickbook refresh token and expiry date
func UpdateQuickbooksRefreshToken(ctx context.Context, application, environment string,
	parameters QuickBooksRefreshToken) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = ValidateApplication(application); soteErr.ErrCode == nil {
		if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
			var (
				ssmParam = []SSMParameter{
					{
						Key:        QUICKBOOKSREFRESHTOKEN,
						Content:    parameters.Token,
						TargetType: types.ParameterTypeSecureString,
					},
					{
						Key:        QUICKBOOKSREFRESHTOKENEXPIRY,
						Content:    parameters.ExpiryDate.String(),
						TargetType: types.ParameterTypeString,
					},
				}
			)

			soteErr = updateParameter(ctx, application, environment, ssmParam)
		}
	}

	return
}

/*
setPath will build the query path based on the ROOTPATH, Application and Environment.
*/
func setPath(application, environment string) (path string) {
	sLogger.DebugMethod()

	if application == "" && environment == "" {
		path = ROOTPATH
	} else {
		if application == "" && environment != "" {
			path = ROOTPATH + "/" + environment
		} else {
			if application != "" && environment == "" {
				path = ROOTPATH + "/" + application
			} else {
				path = ROOTPATH + "/" + application + "/" + environment
			}
		}
	}

	return
}

/*
listParameters will query up to the first 10 parameters for the ROOTPATH with some combination of application
and environment variable values.  Application and environment can be empty.
*/
func listParameters(ctx context.Context, application, environment string) (pSSMPathOutput *ssm.GetParametersByPathOutput, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var err error

	if soteErr = ValidateEnvironment(environment); soteErr.ErrCode == nil {
		var (
			path         = setPath(application, environment)
			ssmPathInput ssm.GetParametersByPathInput
		)

		ssmPathInput.Path = &path
		ssmPathInput.Recursive = pTrue
		ssmPathInput.WithDecryption = pTrue
		ssmPathInput.MaxResults = pMaxResult
		// If there are any parameters that matches the path, a result set will be return by the GetParametersByPath call.
		if pSSMPathOutput, err = awsService.GetParametersByPath(ctx, &ssmPathInput); len(pSSMPathOutput.Parameters) == 0 {
			soteErr = sError.GetSError(109999, sError.BuildParams([]string{path}), sError.EmptyMap)
		}
	}
	if err != nil {
		log.Fatalln(err)
	}

	return
}

/*
getParameter will query the first 10 parameters for the ROOTPATH with some combination of application
and environment variable values.  Application and environment can be empty.
*/
func getParameter(ctx context.Context, application, environment, key string) (returnValue interface{}, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var ssmParamInput ssm.GetParameterInput

	ssmParamInput.WithDecryption = pTrue
	name := setPath(application, environment) + "/" + key
	ssmParamInput.Name = &name

	// If there are any parameters that match the path, a result set will be return by the GetParametersByPath call.
	if pSSMParamOutput, err := awsService.GetParameter(ctx, &ssmParamInput); err != nil {
		soteErr = sError.GetSError(109999, sError.BuildParams([]string{name}), sError.EmptyMap)
	} else {
		returnValue = *pSSMParamOutput.Parameter.Value
	}

	return
}

// Updates one or more values for an SSM document.
func updateParameter(ctx context.Context, application, environment string, parameters []SSMParameter) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	for _, parameter := range parameters {
		if _, soteErr = getParameter(ctx, application, environment, parameter.Key); soteErr.ErrCode != nil {
			break
		}

		name := setPath(application, environment) + "/" + parameter.Key
		pSSMParamInput := &ssm.PutParameterInput{
			Name:      &name,
			Value:     &parameter.Content,
			Type:      parameter.TargetType,
			Overwrite: true,
		}

		if _, err := awsService.PutParameter(ctx, pSSMParamInput); err != nil {
			soteErr = sError.GetSError(199999, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
			break
		}
	}

	return
}
