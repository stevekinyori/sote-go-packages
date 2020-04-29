/*
	This is a wrapper for errors messages used by Sote GO software developers.

	This package provides a number of functions that help development and generate documentation.

	The fields that make up the Sote Error structure are the following:
		ErrType          The category of the error message
		ParamCount       The number of parameters have are needed for the message
		ParamDescription Description of the parameters that need to be supplied
		FmtErrMsg        This is the raw formatted message before the parameters are applied
		Loc              The location in the code where the error occurred
*/
package serror

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgconn"
	"gitlab.com/soteapps/packages/slogger"
)

type SoteError struct {
	ErrCode          interface{}
	ErrType          string
	ParamCount       int
	ParamDescription string
	FmtErrMsg        string
	OrgErrorDetails  map[string]string
	Loc              string
}

const UserError string = "User_Error"
const ProcessError string = "Process_Error"
const NatsError string = "NATS_Error"
const ContentError string = "Content_Error"
const LogicIssue string = "Logic_Issue"
const ConfigurationIssue string = "Configuration_Issue"
const ApiContractError string = "API_Contract_Error"
const MarkDownTitleBar string = "| Error Code | Category | Parameter Description | Formatted Error Text |\n|--------|--------|--------|--------|\n"
const FuncCommentsHeader string = "\tError Code with requiring parameters:\n"
const SQLState string = "SQLSTATE"

var EmptyMap = make(map[string]string)

var SErrors = map[int]SoteError{
	100000: {100000, UserError, 0, "None", "Item already exists", emptyMap, ""},
	100100: {100100, UserError, 2, "List Of users roles, Requested action", "Your roles %v are not authorized to %v", emptyMap, ""},
	109999: {109999, UserError, 1, "Item name", "No %v was/were found", emptyMap, ""},
	//
	200000: {200000, ProcessError, 0, "None", "Row has been updated since reading it, re-read the row", emptyMap, ""},
	200100: {200100, ProcessError, 0, "None", "Table doesn't exist", emptyMap, ""},
	200200: {200200, ProcessError, 2, "Parameter name, Data type of parameter", "%v must be of type %v", emptyMap, ""},
	200250: {200250, ProcessError, 3, "Parameter name, Parameter value, List of values allowed", "%v (%v) must contain one of these values: %v", emptyMap, ""},
	200260: {200260, ProcessError, 3, "Other parameter name, Parameter name, Parameter value", "%v must be provided when %v is set to (%v)", emptyMap, ""},
	200500: {200500, ProcessError, 1, "Thing being changed", "You are making changes to a canceled or completed %v", emptyMap, ""},
	200510: {200510, ProcessError, 3, "Parameter name, Field name, Field value", "%v can't be updated because %v is set to %v", emptyMap, ""},
	200511: {200511, ProcessError, 2, "Parameter name, Another parameter name", "%v and %v must both be populated or null", emptyMap, ""},
	201000: {201000, ProcessError, 1, "Info returned from HTTP/HTTPS Request", "Bad HTTP/HTTPS Request - %v", emptyMap, ""},
	201005: {201005, ProcessError, 0, "None", "Invalid Claim", emptyMap, ""},
	202000: {202000, ProcessError, 1, "Environment Name", "The API you are calling is not available in this environment (%v)", emptyMap, ""},
	209998: {209998, ProcessError, 0, "None", "Database constraint error - see Details", emptyMap, ""},
	209999: {209999, ProcessError, 0, "None", "SQL error - see Details", emptyMap, ""},
	219999: {219999, ProcessError, 0, "None", "Cognito error - see Details", emptyMap, ""},
	230000: {230000, ProcessError, 0, "None", "The number of parameters provided for the error message does not match the required number", emptyMap, ""},
	230050: {230050, ProcessError, 2, "Name, Application/Package name", "Number of parameters defined in the %v is not support by %v", emptyMap, ""},
	230060: {230060, ProcessError, 2, "Provided parameter count, Expected parameter count", "Number of parameters provided (%v) doesn't match the number expected (%v)", emptyMap, ""},
	250000: {250000, ProcessError, 0, "None", "AWS SES error - see details in retPack", emptyMap, ""},
	//
	300000: {300000, NatsError, 0, "None", "TBD", emptyMap, ""},
	310000: {310000, NatsError, 1, "Key name", "Upper or lower case %v key is missing", emptyMap, ""},
	310005: {310005, NatsError, 1, "Key name", "Upper or lower case %v keys value is missing", emptyMap, ""},
	320000: {320000, NatsError, 1, "List of required parameters", "Message doesn't match signature. Sender must provide the following parameter names: %v", emptyMap, ""},
	//
	400000: {400000, ContentError, 2, "Field name, Field value", "%v (%v) is not numeric", emptyMap, ""},
	400010: {400010, ContentError, 2, "Field name, Field value", "%v (%v) is not a string", emptyMap, ""},
	400020: {400020, ContentError, 2, "Field name, Field value", "%v (%v) is not a float", emptyMap, ""},
	400030: {400030, ContentError, 2, "Field name, Field value", "%v (%v) is not a array", emptyMap, ""},
	400040: {400040, ContentError, 2, "Field name, Field value", "%v (%v) is not a json string", emptyMap, ""},
	400050: {400050, ContentError, 2, "Field name, Field value", "%v (%v) is not a valid email address", emptyMap, ""},
	400060: {400060, ContentError, 2, "Field name, Field value", "%v (%v) contains special characters which are not allowed", emptyMap, ""},
	400065: {400065, ContentError, 2, "Field name, Field value", "%v (%v) contains special characters other than alpha and underscore", emptyMap, ""},
	400070: {400070, ContentError, 2, "Field name, Field value", "%v (%v) is not a valid date", emptyMap, ""},
	400080: {400080, ContentError, 2, "Field name, Field value", "%v (%v) is not a valid timestamp. Format's are UTC, GMT or Zulu", emptyMap, ""},
	400090: {400090, ContentError, 6, "Field name, Field value, 'small' or 'large', 'Min' or 'Max', expected size, actual size", "%v (%v) is too %v. %v size: %v Actual size: %v", emptyMap, ""},
	400100: {400100, ContentError, 1, "Parameter name", "%v could't be converted to an array - JSON conversion error", emptyMap, ""},
	400110: {400110, ContentError, 1, "Parameter name", "%v could't be parsed - Invalid JSON error", emptyMap, ""},
	400111: {400111, ContentError, 2, "Parameter name, Application/Package name", "%v could't be converted to a map/keyed array - %v", emptyMap, ""},
	405110: {405110, ContentError, 2, "Thing being changed. System Id for the thing", "No update is needed. No fields where changed for %v with id %v", emptyMap, ""},
	405120: {405120, ContentError, 3, "JSON array name, Thing being changed, System Id for the thing", "The %v was empty for %v with id %v", emptyMap, ""},
	410000: {410000, ContentError, 1, "Error message number", "%v error message is missing from serror package", emptyMap, ""},
	//
	500000: {500000, LogicIssue, 0, "None", "Code is exiting in unexpected way.  Investigate logs right away!", emptyMap, ""},
	//
	600000: {600000, ConfigurationIssue, 0, "None", ".env files are missing", emptyMap, ""},
	600010: {600010, ConfigurationIssue, 2, "File name, Message returned from Open", "%v file was not found. Message return: %v", emptyMap, ""},
	601000: {601000, ConfigurationIssue, 1, "Environment name", "environment variable is missing (%v)", emptyMap, ""},
	602000: {602000, ConfigurationIssue, 3, "Database name, Database driver name, Port value", "Unable to connect to database %v using driver %v on port %v", emptyMap, ""},
	602010: {602010, ConfigurationIssue, 0, "None", "Unable to pass database authentication", emptyMap, ""},
	609999: {609999, ConfigurationIssue, 1, "Variable name", "Start up variable is missing (%v)", emptyMap, ""},
	//
	700000: {700000, ApiContractError, 1, "List of required parameters", "Call doesn't match API signature. Caller must provide the following parameter names: %v", emptyMap, ""},
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
		201000	Info returned from HTTP/HTTPS Request
		202000	Environment Name
		230050	Name, Application/Package name
		230060	Provided parameter count, Expected parameter count
		310000	Key name
		310005	Key name
		320000	List of required parameters
		400000	Field name, Field value
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
		400100	Parameter Name
		400110	Parameter Name
		405110	Thing being changed. System Id for the thing
		405120	JSON array name, Thing being changed, System Id for the thing
		410000	Error Message Number
		600010	File name, Message returned from Open
		601000	Environment name
		602000	Database name, Database driver name, Port value
		609999	Variable name
		700000	List of required parameters
*/
func GetSError(code int, params []string, orgErrorDetails map[string]string) SoteError {
	slogger.DebugMethod()

	var fmttdError SoteError = SErrors[code]
	if fmttdError.ErrCode != code {
		fmttdError = GetSError(410000, []string{strconv.Itoa(code)}, orgErrorDetails)
	} else if fmttdError.ParamCount != len(params) {
		fmttdError = GetSError(230060, []string{strconv.Itoa(fmttdError.ParamCount), strconv.Itoa(len(params))}, orgErrorDetails)
	} else {
		switch fmttdError.ParamCount {
		case 0:
			fmttdError.FmtErrMsg = fmt.Sprintf(fmttdError.FmtErrMsg)
			fmttdError.OrgErrorDetails = orgErrorDetails
		case 1:
			fmttdError.FmtErrMsg = fmt.Sprintf(fmttdError.FmtErrMsg, params[0])
			fmttdError.OrgErrorDetails = orgErrorDetails
		case 2:
			fmttdError.FmtErrMsg = fmt.Sprintf(fmttdError.FmtErrMsg, params[0], params[1])
			fmttdError.OrgErrorDetails = orgErrorDetails
		case 3:
			fmttdError.FmtErrMsg = fmt.Sprintf(fmttdError.FmtErrMsg, params[0], params[1], params[2])
			fmttdError.OrgErrorDetails = orgErrorDetails
		case 6:
			fmttdError.FmtErrMsg = fmt.Sprintf(fmttdError.FmtErrMsg, params[0], params[1], params[2], params[3], params[4], params[5])
			fmttdError.OrgErrorDetails = orgErrorDetails
		default:
			fmttdError = GetSError(230050, []string{"Error message", "serror.GetSError"}, orgErrorDetails)
		}
	}
	return fmttdError
}

func ConvertErr(err error) (orgErrorDetails []string, soteErr SoteError) {
	slogger.DebugMethod()

	if strings.Contains(err.Error(), SQLState) {
		pgErr := err.(*pgconn.PgError)
		orgErrorDetails = append(orgErrorDetails, pgErr.Code, pgErr.ColumnName, pgErr.ConstraintName, pgErr.DataTypeName, pgErr.Error(), pgErr.File, pgErr.Hint, strconv.Itoa(int(pgErr.InternalPosition)),
			pgErr.InternalQuery, strconv.Itoa(int(pgErr.Line)), pgErr.Message, strconv.Itoa(int(pgErr.Position)), pgErr.Routine, pgErr.SchemaName, pgErr.Severity, pgErr.SQLState(), pgErr.TableName, pgErr.Where)
	} else {
		soteErr = GetSError(400111, []string{"err", "serror"}, emptyMap)
	}
	return orgErrorDetails, soteErr
}

/*
	This will generate the markdown syntax that can be published on a Wiki page.  This makes
	this code the master source of Sote Error messages
*/
func GenMarkDown() string {
	slogger.DebugMethod()

	// Sort the Keys from SError map
	var errorKeys []int
	for _, i2 := range SErrors {
		errorKeys = append(errorKeys, i2.ErrCode.(int))
	}
	sort.Ints(errorKeys)
	// Generate the markdown syntax
	var markDown string = MarkDownTitleBar
	for _, i2 := range errorKeys {
		x := SErrors[i2]
		markDown += fmt.Sprintf("| %v | %v | %v | %v |\n", x.ErrCode, x.ErrType, x.ParamDescription, x.FmtErrMsg)
	}
	return markDown
}

/*
	This will generate plain text comments about error code that require parameters.  This can be used
	to update the GetSError function comments
*/
func GenErrorLisRequiredParams() string {
	slogger.DebugMethod()

	// Sort the Keys from SError map
	var errorKeys []int
	for _, i2 := range SErrors {
		errorKeys = append(errorKeys, i2.ErrCode.(int))
	}
	sort.Ints(errorKeys)
	// Generate the plain text
	var funcComments string = FuncCommentsHeader
	for _, i2 := range errorKeys {
		if x := SErrors[i2]; x.ParamCount > 0 {
			funcComments += fmt.Sprintf("\t\t%v\t%v\n", x.ErrCode, x.ParamDescription)
		}
	}
	return funcComments
}
