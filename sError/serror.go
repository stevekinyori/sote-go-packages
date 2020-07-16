/*
	This is a wrapper for errors messages used by Sote GO software developers.

	This package provides a number of functions that help development and generate documentation.

	The fields that make up the Sote Error structure are the following:
		ErrCode          This is the assigned error number or nil.  Nil means that there was no error
		ErrType          The category of the error message
		ParamCount       The number of parameters have are needed for the message
		ParamDescription Description of the parameters that need to be supplied
		FmtErrMsg        This is the raw formatted message before the parameters are applied
		ErrorDetails     This can be used for anything including look up value errors.
		Loc              The location in the code where the error occurred
*/
package sError

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgconn"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

type SoteError struct {
	ErrCode          interface{}
	ErrType          string
	ParamCount       int
	ParamDescription string
	FmtErrMsg        string
	ErrorDetails     map[string]string
	Loc              string
}

// Error categories
const USERERROR string = "User_Error"
const PROCESSERROR string = "Process_Error"
const NATSERROR string = "NATS_Error"
const CONTENTERROR string = "Content_Error"
const PERMISSIONERROR string = "Permission_Error"
const CONFIGURATIONISSUE string = "Configuration_Issue"
const APICONTRACTERROR string = "API_Contract_Error"
const GENERALERROR string = "General_Error"
const MARKDOWNTITLEBAR string = "| Error Code | Category | Parameter Description | Formatted Error Text |\n|--------|--------|--------|--------|\n"
const FUNCCOMMENTSHEADER string = "\tError Code with requiring parameters:\n"
const SQLSTATE string = "SQLSTATE"

var (
	EmptyMap = make(map[string]string)
)

// Ranges are not limited to a single error category.
var soteErrors = map[int]SoteError{
	100000: {100000, USERERROR, 0, "None", "Item already exists", EmptyMap, ""},
	100100: {100100, USERERROR, 2, "List Of users roles, Requested action", "Your roles %v are not authorized to %v", EmptyMap, ""},
	109999: {109999, USERERROR, 1, "Item name", "No %v was/were found", EmptyMap, ""},
	//
	200000: {200000, PROCESSERROR, 0, "None", "Row has been updated since reading it, re-read the row", EmptyMap, ""},
	200100: {200100, PROCESSERROR, 0, "None", "Table doesn't exist", EmptyMap, ""},
	200200: {200200, PROCESSERROR, 2, "Parameter name, Data type of parameter", "%v must be of type %v", EmptyMap, ""},
	200250: {200250, PROCESSERROR, 3, "Parameter name, Parameter value, List of values allowed", "%v (%v) must contain one of these values: %v", EmptyMap, ""},
	200260: {200260, PROCESSERROR, 3, "Other parameter name, Parameter name, Parameter value", "%v must be provided when %v is set to (%v)", EmptyMap, ""},
	200500: {200500, PROCESSERROR, 1, "Thing being changed", "You are making changes to a canceled or completed %v", EmptyMap, ""},
	200510: {200510, PROCESSERROR, 3, "Parameter name, Field name, Field value", "%v can't be updated because %v is set to %v", EmptyMap, ""},
	200511: {200511, PROCESSERROR, 2, "Parameter name, Another parameter name", "%v and %v must both be populated or null", EmptyMap, ""},
	200512: {200512, PROCESSERROR, 2, "Parameter name, Another parameter name", "%v and %v must both be populated", EmptyMap, ""},
	200513: {200513, PROCESSERROR, 1, "Parameter name", "%v must be populated", EmptyMap, ""},
	201000: {201000, PROCESSERROR, 1, "Info returned from HTTP/HTTPS Request", "Bad HTTP/HTTPS Request - %v", EmptyMap, ""},
	201005: {201005, PROCESSERROR, 0, "None", "Invalid Claim", EmptyMap, ""},
	202000: {202000, PROCESSERROR, 1, "Environment Name", "The API you are calling is not available in this environment (%v)", EmptyMap, ""},
	209500: {209500, PROCESSERROR, 0, "None", "QuickSight error - see Details", EmptyMap, ""},
	209998: {209998, PROCESSERROR, 0, "None", "Database constraint error - see Details", EmptyMap, ""},
	209999: {209999, PROCESSERROR, 0, "None", "SQL error - see Details", EmptyMap, ""},
	219999: {219999, PROCESSERROR, 0, "None", "Cognito error - see Details", EmptyMap, ""},
	230000: {230000, PROCESSERROR, 0, "None", "The number of parameters provided for the error message does not match the required number", EmptyMap, ""},
	230050: {230050, PROCESSERROR, 2, "Name, Application/Package name", "Number of parameters defined in the %v is not support by %v", EmptyMap, ""},
	230060: {230060, PROCESSERROR, 2, "Provided parameter count, Expected parameter count", "Number of parameters provided (%v) doesn't match the number expected (%v)", EmptyMap, ""},
	250000: {250000, PROCESSERROR, 0, "None", "AWS SES error - see details in retPack", EmptyMap, ""},
	250005: {250005, PROCESSERROR, 0, "None", "AWS STS error - see details in retPack", EmptyMap, ""},
	//
	300000: {300000, NATSERROR, 0, "None", "TBD", EmptyMap, ""},
	310000: {310000, NATSERROR, 1, "Key name", "Upper or lower case %v key is missing", EmptyMap, ""},
	310005: {310005, NATSERROR, 1, "Key name", "Upper or lower case %v keys value is missing", EmptyMap, ""},
	320000: {320000, NATSERROR, 1, "List of required parameters", "Message doesn't match signature. Sender must provide the following parameter names: %v", EmptyMap, ""},
	//
	400000: {400000, CONTENTERROR, 2, "Field name, Field value", "%v (%v) is not numeric", EmptyMap, ""},
	400005: {400005, CONTENTERROR, 2, "Field name, Minimal length", "%v must a value greater than %v", EmptyMap, ""},
	400010: {400010, CONTENTERROR, 2, "Field name, Field value", "%v (%v) is not a string", EmptyMap, ""},
	400020: {400020, CONTENTERROR, 2, "Field name, Field value", "%v (%v) is not a float", EmptyMap, ""},
	400030: {400030, CONTENTERROR, 2, "Field name, Field value", "%v (%v) is not a array", EmptyMap, ""},
	400040: {400040, CONTENTERROR, 2, "Field name, Field value", "%v (%v) is not a json string", EmptyMap, ""},
	400050: {400050, CONTENTERROR, 2, "Field name, Field value", "%v (%v) is not a valid email address", EmptyMap, ""},
	400060: {400060, CONTENTERROR, 2, "Field name, Field value", "%v (%v) contains special characters which are not allowed", EmptyMap, ""},
	400065: {400065, CONTENTERROR, 2, "Field name, Field value", "%v (%v) contains special characters other than underscore", EmptyMap, ""},
	400070: {400070, CONTENTERROR, 2, "Field name, Field value", "%v (%v) is not a valid date", EmptyMap, ""},
	400080: {400080, CONTENTERROR, 2, "Field name, Field value", "%v (%v) is not a valid timestamp. Format's are UTC, GMT or Zulu", EmptyMap, ""},
	400090: {400090, CONTENTERROR, 6, "Field name, Field value, 'small' or 'large', 'Min' or 'Max', expected size, actual size", "%v (%v) is too %v. %v size: %v Actual size: %v", EmptyMap, ""},
	400100: {400100, CONTENTERROR, 2, "Parameter name, Data Structure Type", "%v could't be converted to an %v - JSON conversion error", EmptyMap, ""},
	400105: {400105, CONTENTERROR, 2, "Data Structure Name, Data Structure Type", "%v (%v) could't be converted to JSON - JSON conversion error", EmptyMap, ""},
	400110: {400110, CONTENTERROR, 1, "Parameter name", "%v could't be parsed - Invalid JSON error", EmptyMap, ""},
	400111: {400111, CONTENTERROR, 2, "Parameter name, Application/Package name", "%v could't be converted to a map/keyed array - %v", EmptyMap, ""},
	401000: {401000, CONTENTERROR, 0, "None", "Column must have a non-null value. Details: ", EmptyMap, ""},
	405110: {405110, CONTENTERROR, 2, "Thing being changed, System Id for the thing", "No update is needed. No fields where changed for %v with id %v", EmptyMap, ""},
	405120: {405120, CONTENTERROR, 3, "JSON array name, Thing being changed, System Id for the thing", "The %v was empty for %v with id %v", EmptyMap, ""},
	410000: {410000, CONTENTERROR, 1, "Error message number", "%v error message is missing from sError package", EmptyMap, ""},
	//
	500000: {500000, PERMISSIONERROR, 0, "None", "iss (Issuer) is not valid", EmptyMap, ""},
	500010: {500010, PERMISSIONERROR, 0, "None", "sub (Subject) was not present", EmptyMap, ""},
	500020: {500020, PERMISSIONERROR, 0, "None", "token_use is not valid", EmptyMap, ""},
	500030: {500030, PERMISSIONERROR, 0, "None", "client is not valid", EmptyMap, ""},
	500040: {500040, PERMISSIONERROR, 0, "None", "client is not valid for this application", EmptyMap, ""},
	500050: {500050, PERMISSIONERROR, 0, "None", "Token is expired", EmptyMap, ""},
	//
	600000: {600000, CONFIGURATIONISSUE, 0, "None", ".env files are missing", EmptyMap, ""},
	600010: {600010, CONFIGURATIONISSUE, 2, "File name, Message returned from Open", "%v file was not found. Message return: %v", EmptyMap, ""},
	601000: {601000, CONFIGURATIONISSUE, 1, "Environment name", "environment variable is missing (%v)", EmptyMap, ""},
	601010: {601010, CONFIGURATIONISSUE, 1, "Environment name", "environment value (%v) is invalid", EmptyMap, ""},
	602000: {602000, CONFIGURATIONISSUE, 3, "Database name, Database driver name, Port value", "Unable to connect to database %v using driver %v on port %v", EmptyMap, ""},
	602010: {602010, CONFIGURATIONISSUE, 0, "None", "Unable to pass database authentication", EmptyMap, ""},
	602020: {602020, CONFIGURATIONISSUE, 1, "SSL Mode", "Only disable, allow, prefer and required are supported.", EmptyMap, ""},
	602030: {602030, CONFIGURATIONISSUE, 1, "Connection Type", "Only single or pool are supported.", EmptyMap, ""},
	602999: {602999, CONFIGURATIONISSUE, 0, "None", "No database connection has been established", EmptyMap, ""},
	605000: {605000, CONFIGURATIONISSUE, 0, "None", "Unexpected signing method", EmptyMap, ""},
	605010: {605010, CONFIGURATIONISSUE, 0, "None", "kid header not found", EmptyMap, ""},
	605020: {605020, CONFIGURATIONISSUE, 1, "Kid", "key (%v) was not found in token", EmptyMap, ""},
	605021: {605021, CONFIGURATIONISSUE, 1, "Kid", "Kid (%v) was not found in public key set", EmptyMap, ""},
	605030: {605030, CONFIGURATIONISSUE, 1, "Environment name", "Failed to fetch remote JWK (status = 404) for %v", EmptyMap, ""},
	609999: {609999, CONFIGURATIONISSUE, 1, "Parameter name", "Start up parameter is missing (%v)", EmptyMap, ""},
	//
	700000: {700000, APICONTRACTERROR, 1, "List of required parameters", "Call doesn't match API signature. Caller must provide the following parameter names: %v", EmptyMap, ""},
	//
	800000: {800000, GENERALERROR, 0, "None", "An error has occurred that is not expected.", EmptyMap, ""},
	800100: {800100, GENERALERROR, 0, "None", "Postgres error has occurred that is not expected.", EmptyMap, ""},
	800199: {800199, GENERALERROR, 0, "None", "Postgres is not responding over TCP. Container may not be running.", EmptyMap, ""},
}

/*
	This will return the formatted message using the supplied code and parameters

	Error Code with requiring parameters:
		100100	List Of users roles, Requested action
		109999	Item name
		200200	Parameter name, Data type of parameter
		200250	Parameter name, Parameter value, List of values allowed
		200260	Other parameter name, Parameter name, Parameter value
		200500	Thing being changed
		200510	Parameter name, Field name, Field value
		200511	Parameter name, Another parameter name
		200512	Parameter name, Another parameter name
		200513	Parameter name
		201000	Info returned from HTTP/HTTPS Request
		202000	Environment Name
		230050	Name, Application/Package name
		230060	Provided parameter count, Expected parameter count
		310000	Key name
		310005	Key name
		320000	List of required parameters
		400000	Field name, Field value
		400005	Field name, Minimal length
		400010	Field name, Field value
		400020	Field name, Field value
		400030	Field name, Field value
		400040	Field name, Field value
		400050	Field name, Field value
		400060	Field name, Field value
		400065	Field name, Field value
		400070	Field name, Field value
		400080	Field name, Field value
		400090	Field name, Field value, 'small' or 'large', 'Min' or 'Max', expected size, actual size
		400100	Parameter name, Data Structure Type
		400105	Data Structure Name, Data Structure Type
		400110	Parameter name
		400111	Parameter name, Application/Package name
		405110	Thing being changed, System Id for the thing
		405120	JSON array name, Thing being changed, System Id for the thing
		410000	Error message number
		600010	File name, Message returned from Open
		601000	Environment name
		601010	Environment name
		602000	Database name, Database driver name, Port value
		602020	SSL Mode
		602030	Connection Type
		605020	Kid
		605021	Kid
		605030	Environment name
		609999	Parameter name
		700000	List of required parameters
*/
func GetSError(code int, params []interface{}, errorDetails map[string]string) (soteErr SoteError) {
	sLogger.DebugMethod()

	soteErr = soteErrors[code]
	if soteErr.ErrCode != code {
		s := make([]interface{}, 1)
		s[0] = code
		soteErr = GetSError(410000, s, errorDetails)
	} else if soteErr.ParamCount != len(params) {
		s := make([]interface{}, 2)
		s[0] = soteErr.ParamCount
		s[1] = len(params)
		soteErr = GetSError(230060, s, errorDetails)
	} else {
		if soteErr.ParamCount == 0 {
			soteErr.FmtErrMsg = fmt.Sprintf(soteErr.FmtErrMsg)
			soteErr.ErrorDetails = errorDetails
		} else {
			soteErr.FmtErrMsg = fmt.Sprintf(soteErr.FmtErrMsg, params)
			soteErr.ErrorDetails = errorDetails
		}
	}
	return
}

// This will convert a postgresql error into error details for include in SoteError
func ConvertErr(err error) (errorDetails map[string]string, soteErr SoteError) {
	sLogger.DebugMethod()

	if strings.Contains(err.Error(), SQLSTATE) {
		pgErr := err.(*pgconn.PgError)

		errorDetails = map[string]string{
			"Code":             pgErr.Code,
			"ColumnName":       pgErr.ColumnName,
			"ConstraintName":   pgErr.ConstraintName,
			"DataTypeName":     pgErr.DataTypeName,
			"Error":            pgErr.Error(),
			"File":             pgErr.File,
			"Hint":             pgErr.Hint,
			"InternalPosition": strconv.Itoa(int(pgErr.InternalPosition)),
			"InternalQuery":    pgErr.InternalQuery,
			"Line":             strconv.Itoa(int(pgErr.Line)),
			"Message":          pgErr.Message,
			"Position":         strconv.Itoa(int(pgErr.Position)),
			"Routine":          pgErr.Routine,
			"SchemaName":       pgErr.SchemaName,
			"Severity":         pgErr.Severity,
			"SQLSTATE":         pgErr.SQLState(),
			"TableName":        pgErr.TableName,
			"Where":            pgErr.Where,
		}
	} else {
		s := make([]interface{}, 2)
		s[0] = "err"
		s[1] = "sError"
		soteErr = GetSError(400111, s, EmptyMap)
	}
	return
}

/*
This will convert an array of strings to a param list for sError.GetSError
*/
func BuildParams(values []string) (s []interface{}) {
	sLogger.DebugMethod()

	s = make([]interface{}, len(values))
	for i, v := range values {
		s[i] = v
	}

	return
}

/*
	This will generate the markdown syntax that can be published on a Wiki page.  This makes
	this code the master source of Sote Error messages
*/
func GenMarkDown() (markDown string) {
	sLogger.DebugMethod()

	// Sort the Keys from SError map
	var errorKeys []int
	for _, i2 := range soteErrors {
		errorKeys = append(errorKeys, i2.ErrCode.(int))
	}
	sort.Ints(errorKeys)
	// Generate the markdown syntax
	markDown = MARKDOWNTITLEBAR
	for _, i2 := range errorKeys {
		x := soteErrors[i2]
		markDown += fmt.Sprintf("| %v | %v | %v | %v |\n", x.ErrCode, x.ErrType, x.ParamDescription, x.FmtErrMsg)
	}
	return
}

/*
	This will generate plain text comments about error code that require parameters.  This can be used
	to update the GetSError function comments
*/
func GenErrorLisRequiredParams() (funcComments string) {
	sLogger.DebugMethod()

	// Sort the Keys from SError map
	var errorKeys []int
	for _, i2 := range soteErrors {
		errorKeys = append(errorKeys, i2.ErrCode.(int))
	}
	sort.Ints(errorKeys)
	// Generate the plain text
	funcComments = FUNCCOMMENTSHEADER
	for _, i2 := range errorKeys {
		if x := soteErrors[i2]; x.ParamCount > 0 {
			funcComments += fmt.Sprintf("\t\t%v\t%v\n", x.ErrCode, x.ParamDescription)
		}
	}
	return
}
