package sError

import (
	"encoding/json"
	"fmt"
	"sort"
)

type Exception struct {
	errCode          int
	paramCount       int
	paramDescription string
	fmtErrMsg        string
	errorDetails     string
}

func (ex Exception) Code() int {
	return ex.errCode
}

func (ex Exception) Params(params ...interface{}) Exception {
	if ex.paramCount != len(params) {
		return InvalidParameterCount.Params(ex.paramCount, len(params))
	}
	newEx := ex
	newEx.fmtErrMsg = fmt.Sprint(ex.errCode) + fmt.Sprintf(ex.fmtErrMsg, params...)
	return newEx
}

func (ex Exception) Details(errorDetails string) Exception {
	newEx := ex
	newEx.errorDetails = errorDetails
	return newEx
}

func (ex Exception) Json() string {
	out, _ := json.MarshalIndent(map[string]interface{}{
		"ErrCode":          ex.errCode,
		"ParamCount":       ex.paramCount,
		"ParamDescription": ex.paramDescription,
		"FmtErrMsg":        ex.fmtErrMsg,
		"ErrorDetails":     ex.errorDetails,
	}, PREFIX, INDENT+INDENT)
	return string(out)
}

func (ex Exception) String() string {
	if ex.errorDetails != "" {
		return ex.fmtErrMsg + "\n" + ex.errorDetails
	}
	return ex.fmtErrMsg
}

func GetException(errCode int) Exception {
	return exceptionCodes[errCode]
}

func ExceptionDocumentation() (markDown, funcComments string) {
	// Sort the Keys from SError map
	var errorKeys []int
	for _, ex := range exceptionCodes {
		errorKeys = append(errorKeys, ex.errCode)
	}
	sort.Ints(errorKeys)
	// Generate Documentation
	markDown = MARKDOWNTITLEBAR
	funcComments = FUNCCOMMENTSHEADER
	for _, i := range errorKeys {
		x := exceptionCodes[i]
		markDown += fmt.Sprintf("| %v | %v | %v |\n", x.errCode, x.paramDescription, x.fmtErrMsg)
		funcComments += fmt.Sprintf("\t\t%v\t%v > %v\n", x.errCode, x.paramDescription, x.fmtErrMsg)
	}
	return
}

func NewException(errCode int, paramCount int, paramDescription, fmtErrMsg string) Exception {
	ex := Exception{
		errCode:          errCode,
		fmtErrMsg:        fmtErrMsg,
		paramDescription: paramDescription,
		paramCount:       paramCount,
	}
	exceptionCodes[errCode] = ex
	return ex
}

var (
	exceptionCodes = make(map[int]Exception)
	// User_Error
	ItemAlreadyExists = NewException(100000, 1, "Item Name", ": %v already exists")
	NotAuthorized     = NewException(100100, 2, "List of users roles, Requested action",
		": Your roles %v are not authorized to %v")
	// Process_Error
	DirtyRead        = NewException(100200, 0, "None", ": Row has been updated since reading it, re-read the row")
	CanceledComplete = NewException(100500, 1, "Thing being changed",
		": You are making changes to a canceled or completed %v")
	TimeOut              = NewException(101010, 1, "Service Name", ": %v timed out")
	ItemNotFound         = NewException(109999, 1, "Item name", ": %v was/were not found")
	UnexpectedError      = NewException(199999, 1, "Error Details", ": An error has occurred that is not expected. See Log! %v")
	TableMissing         = NewException(200100, 0, "None", ": Table doesn't exist")
	InvalidDataType      = NewException(200200, 2, "Parameter name, Data type of parameter", ": %v must be of type %v")
	RequiredValueMissing = NewException(200250, 3, "Parameter name, Parameter value, List of values allowed",
		": %v (%v) must contain one of these values: %v")
	LinkedParameterValueMissing = NewException(200260, 3, "Other parameter name, Parameter name, Parameter value",
		": %v must be provided when %v is set to (%v)")
	ParameterLockOtherParameterSet = NewException(200510, 3, "Parameter name, Field name, Field value",
		": %v can't be updated because %v is set to %v")
	ParameterMustBeSetOrNull = NewException(200511, 2, "Parameter name, Another parameter name",
		": %v and %v must both be populated or null")
	ParametersMustBeProvided = NewException(200512, 2, "Parameter name, Another parameter name",
		": %v and %v must both be populated")
	ParameterMustBeSet       = NewException(200513, 1, "Parameter name", ": %v must be populated")
	ThreeParametersMustBeSet = NewException(200514, 3, "Parameter name, Another parameter name, "+
		"Another parameter name",
		": %v, %v and %v must all be populated")
	ParameterMustBeEmptyWhenParameterSet = NewException(200515, 2, "Parameter name, Another parameter name",
		": %v must be empty when %v is populated")
	BadHTTPRequest           = NewException(200600, 1, "Info returned from HTTP/HTTPS Request", ": Bad HTTP/HTTPS Request - %v")
	InvalidEnvironmentForAPI = NewException(200700, 1, "Environment Name",
		": The API you are calling is not available in this environment (%v)")
	QuickSightError       = NewException(200800, 0, "None", ": QuickSight error - see Details")
	DatabaseError         = NewException(200900, 0, "None", ": Database constraint error - see Details")
	SqlError              = NewException(200999, 0, "None", ": SQL error - see Details")
	CognitoError          = NewException(201999, 0, "None", ": Cognito error - see Details")
	InvalidParameterCount = NewException(203060, 2, "Provided parameter count, Expected parameter count",
		": Number of parameters provided (%v) doesn't match the number expected (%v)")
	AwsSESError = NewException(205000, 0, "None", ": AWS SES error - see details in retPack")
	AwsSTSError = NewException(205005, 0, "None", ": AWS STS error - see details in retPack")
	// NATS_Error
	JetStreamError           = NewException(206000, 0, "None", ": Jetstream is not enabled")
	NatsSubscriptionError    = NewException(206050, 2, "Subscription Name, Subject", ": (%v) is an invalid subscription. Subject: %v")
	NatsStreamPointerMissing = NewException(206300, 0, "None", ": Stream pointer is nil. Must be a validate pointer to a stream.")
	NatsStreamCreateError    = NewException(206400, 1, "Stream Name",
		": Stream creation encountered an error that is not expected. Stream Name: %v")
	NatsConsumerCreateError = NewException(206600, 2, "Stream Name, Consumer Name",
		": Consumer creation encountered an error that is not expected. Stream Name: %v Consumer Name: %v")
	NatsInvalidConsumerSubjectFilter = NewException(206700, 2, "Stream Name, Consumer Subject Filter",
		": The consumer subject filter must be a subset of the stream subject. Stream Name: %v Consumer Subject Filter: %v")
	// Content_Error
	ParameterNotNumeric    = NewException(207000, 2, "Field name, Field value", ": %v (%v) is not numeric")
	ParameterToSmall       = NewException(207005, 2, "Field name, Minimal length", ": %v must have a value greater than %v")
	ParameterNotString     = NewException(207010, 2, "Field name, Field value", ": %v (%v) is not a string")
	ParameterNotFloat      = NewException(207020, 2, "Field name, Field value", ": %v (%v) is not a float")
	ParameterNotArray      = NewException(207030, 2, "Field name, Field value", ": %v (%v) is not a array")
	ParameterNotJsonString = NewException(207040, 2, "Field name, Field value", ": %v (%v) is not a json string")
	InvalidEmailFormat     = NewException(207050, 2, "Field name, Field value", ": %v (%v) is not a valid email address")
	ParameterNotDate       = NewException(207070, 2, "Field name, Field value", ": %v (%v) is not a valid date")
	ParameterNotTimestamp  = NewException(207080, 2, "Field name, Field value",
		": %v (%v) is not a valid timestamp. Format's are UTC, GMT or Zulu")
	ParameterInvalidSize = NewException(207090, 6,
		"Field name, Field value, 'small' | 'large' | 'Min' | 'Max' | 'low' | 'high', expected size, actual size",
		": %v (%v) is too %v. %v size: %v Actual size: %v")
	JsonConversionError = NewException(207105, 2, "Data Structure Name, Data Structure Type",
		": %v (%v) couldn't be converted to JSON - JSON conversion error")
	ParameterNotMap = NewException(207111, 2, "Parameter name, Application/Package name",
		": %v couldn't be converted to a map/keyed array - %v")
	MissingErrorNumber = NewException(208200, 1, "Error message number", ": %v error message is missing from sError package")
	// Permission_Error
	InvalidISS           = NewException(208300, 0, "None", ": iss (Issuer) is not valid")
	InvalidSubject       = NewException(208310, 1, "Subject", ": sub (Subject: %v) was not present")
	InvalidToken         = NewException(208320, 0, "None", ": token_use is not valid")
	InvalidAppClientId   = NewException(208340, 0, "None", ": client id is not valid for this application")
	TokenExpired         = NewException(208350, 0, "None", ": Token is expired")
	TokenInvalid         = NewException(208355, 0, "None", ": Token is invalid")
	SegmentsCountInvalid = NewException(208356, 0, "None", ": Token contains an invalid number of segments")
	InvalidClaim         = NewException(208360, 1, "Claim names", ": These claims are invalid: %v")
	MissingClaim         = NewException(208370, 0, "None", ": Required claim(s) is/are missing")
	// Configuration_Issue
	EnvFileMissing = NewException(209000, 0, "None", ": .env files are missing")
	FileNotFound   = NewException(209010, 2, "File name, Message returned from Open",
		": %v file was not found. Message return: %v")
	EnvironmentMissing  = NewException(209100, 1, "Environment name", ": environment variable is missing (%v)")
	EnvironmentInvalid  = NewException(209110, 1, "Environment name", ": environment value (%v) is invalid")
	InvalidDBConnection = NewException(209200, 3, "Database name, Database driver name, Port value",
		": Unable to connect to database %v using driver %v on port %v")
	InvalidDBAuthentication     = NewException(209210, 0, "None", ": Unable to pass database authentication")
	InvalidDBSSLMode            = NewException(209220, 1, "SSL Mode", ": Only disable, allow, prefer and required are supported.")
	InvalidConnectionType       = NewException(209230, 1, "Connection Type", ": Only single or pool are supported.")
	NoDBConnection              = NewException(209299, 0, "None", ": No database connection has been established")
	NatsNkeyMissing             = NewException(209398, 0, "None", ": no nkey seed found")
	NoNATSConnection            = NewException(209499, 0, "None", ": No nats connection has been established")
	UnexpectedSign              = NewException(209500, 0, "None", ": Unexpected signing method")
	KidNotFound                 = NewException(209510, 0, "None", ": Kid header not found")
	KidMissingFromToken         = NewException(209520, 1, "Kid", ": key (%v) was not found in token")
	KidDoesNotMatchPublicKeySet = NewException(209521, 1, "Kid", ": Kid (%v) was not found in public key set")
	InvalidRegion               = NewException(210030, 2, "Region, Environment",
		": Failed to fetch remote JWK (status = 404) for %v region %v environment")
	InvalidURL      = NewException(210090, 1, "Parameter name", ": URL is missing (%v)")
	OutOfValidRange = NewException(210098, 1, "Parameter name", ": Start up parameter is out of value range (%v)")
)
