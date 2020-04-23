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

	"gitlab.com/soteapps/packages/slogger"
)

type SoteError struct {
	ErrCode          interface{}
	ErrType          string
	ParamCount       int
	ParamDescription string
	FmtErrMsg        string
	Loc              string
}

const UserError string = "User_Error"
const ProcessError string = "Process_Error"
const NatsError string = "NATS_Error"
const ContentError string = "Content_Error"
const LogicIssue string = "Logic_Issue"
const ConfigurationIssue string = "Configuration_Issue"
const ApiContractError string = "API_Contract_Error"
const MakeDownTitleBar string = "| Error Code | Category | Parameter Description | Formatted Error Text |\n|--------|--------|--------|--------|\n"

var SErrors = map[int]SoteError{
	100000: {100000, UserError, 0, "None", "Item already exists", ""},
	100100: {100100, UserError, 2, "List Of users roles, requested action", "Your roles %v are not authorized to %v", ""},
	109999: {109999, UserError, 1, "Item name", "No %v was/were found", ""},
	//
	200000: {200000, ProcessError, 0, "None", "Row has been updated since reading it, re-read the row", ""},
	200100: {200100, ProcessError, 0, "None", "Table doesn't exist", ""},
	200200: {200200, ProcessError, 2, "Parameter name, data type of parameter", "%v must be of type %v", ""},
	200250: {200250, ProcessError, 3, "Parameter name, parameter value, list of values allowed", "%v (%v) must contain one of these values: %v", ""},
	200260: {200260, ProcessError, 3, "Other parameter name, parameter name, parameter value", "%v must be provided when %v is set to (%v)", ""},
	200500: {200500, ProcessError, 1, "Thing being changed", "You are making changes to a canceled or completed %v", ""},
	200510: {200510, ProcessError, 3, "Parameter name, field name, field value", "%v can't be updated because %v is set to %v", ""},
	200511: {200510, ProcessError, 2, "Parameter name, another parameter name", "%v and %v must both be populated or null", ""},
	201000: {201000, ProcessError, 1, "Info returned from HTTP/HTTPS Request", "Bad HTTP/HTTPS Request - %v", ""},
	201005: {201005, ProcessError, 0, "None", "Invalid Claim", ""},
	202000: {202000, ProcessError, 1, "Environment Name", "The API you are calling is not available in this environment (%v)", ""},
	209999: {209999, ProcessError, 0, "None", "SQL error - see details in retPack", ""},
	219999: {219999, ProcessError, 0, "None", "Cognito error - see details in retPack", ""},
	230000: {230000, ProcessError, 0, "None", "The number of parameters provided for the error message does not match the required number", ""},
	230050: {230050, ProcessError, 0, "None", "Number of parameters defined in the error message is not support by serror.GetSError", ""},
	250000: {250000, ProcessError, 0, "None", "AWS SES error - see details in retPack", ""},
	//
	300000: {300000, NatsError, 0, "None", "TBD", ""},
	310000: {310000, NatsError, 1, "Key name", "Upper or lower case %v key is missing", ""},
	310005: {310005, NatsError, 1, "Key name", "Upper or lower case %v keys value is missing", ""},
	320000: {320000, NatsError, 1, "List of required parameters", "Message doesn't match signature. Sender must provide the following parameter names: %v", ""},
	//
	400000: {400000, ContentError, 2, "Field name, field value", "%v (%v) is not numeric", ""},
	400010: {400010, ContentError, 2, "Field name, field value", "%v (%v) is not a string", ""},
	400020: {400020, ContentError, 2, "Field name, field value", "%v (%v) is not a float", ""},
	400030: {400030, ContentError, 2, "Field name, field value", "%v (%v) is not a array", ""},
	400040: {400040, ContentError, 2, "Field name, field value", "%v (%v) is not a json string", ""},
	400050: {400050, ContentError, 2, "Field name, field value", "%v (%v) is not a valid email address", ""},
	400060: {400060, ContentError, 2, "Field name, field value", "%v (%v) contains special characters which are not allowed", ""},
	400065: {400065, ContentError, 2, "Field name, field value", "%v (%v) contains special characters other than alpha and underscore", ""},
	400070: {400070, ContentError, 2, "Field name, field value", "%v (%v) is not a valid date", ""},
	400080: {400080, ContentError, 2, "Field name, field value", "%v (%v) is not a valid timestamp. Format's are UTC, GMT or Zulu", ""},
	400090: {400090, ContentError, 6, "Field name, field value, 'small' or 'large', 'Min' or 'Max', expected size, actual size", "%v (%v) is too %v. %v size: %v Actual size: %v", ""},
	400100: {400100, ContentError, 1, "Parameter Name", "%v could't be converted to an array - JSON conversion error", ""},
	400110: {400110, ContentError, 1, "Parameter Name", "%v could't be parsed - Invalid JSON error", ""},
	405110: {405110, ContentError, 2, "Thing being changed. System Id for the thing", "No update is needed. No fields where changed for %v with id %v", ""},
	405120: {405120, ContentError, 3, "JSON array name, thing being changed, System Id for the thing", "The %v was empty for %v with id %v", ""},
	410000: {410000, ContentError, 1, "Error Message Number", "%v error message is missing from serror package", ""},
	//
	500000: {500000, LogicIssue, 0, "None", "Code is exiting in unexpected way.  Investigate logs right away!", ""},
	//
	600000: {600000, ConfigurationIssue, 0, "None", ".env files are missing", ""},
	600010: {600010, ConfigurationIssue, 2, "File name, Message Return from Open", "%v file was not found. Message return: %v", ""},
	601000: {601000, ConfigurationIssue, 1, "Environment name", "environment variable is missing (%v)", ""},
	602000: {602000, ConfigurationIssue, 3, "Database name, Database driver name, port value", "Unable to connect to database %v using driver %v on port %v", ""},
	602010: {602010, ConfigurationIssue, 0, "None", "Unable to pass database authentication", ""},
	609999: {609999, ConfigurationIssue, 1, "Variable name", "Start up variable is missing (%v)", ""},
	//
	700000: {700000, ApiContractError, 1, "List of required parameters", "Call doesn't match API signature. Caller must provide the following parameter names: %v", ""},
}

/*
	This will return the formatted message using the supplied code and parameters

	Error Code requiring parameters:
		100100	List of Roles that use has been assigned, Action the user is trying to execute
		200200	Parameter Name that was passed, Data Type the parameter must have
		200250	Parameter Name that was passed, Parameter value that was passed, List of values that are valid
		200260	Parameter Name that is missing, Parameter name that was passed, Parameter value that was passed
		200500	Object Name that is being changed
		200510	Parameter Name that was passed, Field name that is set, Field value that was passed.
		200511	Parameter Name that was passed, Parameter name that passed
		201000	Detailed Message returned from request
		202000	Environment where the error occurred
		310000	Key name that is missing
		310005	Key name that is missing a value
		320000	The list of parameters that are needed
		400000	Field name being checked, Field value being passed
		400010	Field name being checked, Field value being passed
		400020	Field name being checked, Field value being passed
		400030	Field name being checked, Field value being passed
		400040	Field name being checked, Field value being passed
		400050	Field name being checked, Field value being passed
		400060	Field name being checked, Field value being passed
		400065	Field name being checked, Field value being passed
		400070	Field name being checked, Field value being passed
		400080	Field name being checked, Field value being passed
		400090	Field name being checked, Field value being passed, Value is too small or large, Min or Max, Size excepted, Size provided
		400100	Field value being passed
		405110	Object type that is being changed, System id of the instance of the Object Type
		405120	JSON array name that is empty, Object type being changed, System id of the instance of the Object Type
		601000	Missing environment variable name
		602000	Database Name being connected, Database driver name being used, Port number traffic is routed via
		609999	Start up variable that is missing
		700000  List of parameters needed for RESTFUL API being called
*/
func GetSError(code int, params []string) SoteError {
	slogger.DebugMethod()

	var fmttdError SoteError = SErrors[code]
	if fmttdError.ErrCode != code {
		fmttdError = GetSError(109999, params)
	} else if fmttdError.ParamCount != len(params) {
		fmttdError = GetSError(230000, params)
	} else {
		switch fmttdError.ParamCount {
		case 0:
			fmttdError.FmtErrMsg = fmt.Sprintf(fmttdError.FmtErrMsg)
		case 1:
			fmttdError.FmtErrMsg = fmt.Sprintf(fmttdError.FmtErrMsg, params[0])
		case 2:
			fmttdError.FmtErrMsg = fmt.Sprintf(fmttdError.FmtErrMsg, params[0], params[1])
		case 3:
			fmttdError.FmtErrMsg = fmt.Sprintf(fmttdError.FmtErrMsg, params[0], params[1], params[2])
		case 6:
			fmttdError.FmtErrMsg = fmt.Sprintf(fmttdError.FmtErrMsg, params[0], params[1], params[2], params[3], params[4], params[5])
		default:
			fmttdError = SErrors[230050]
		}
	}
	return fmttdError
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
	var markDown string = MakeDownTitleBar
	for _, i2 := range errorKeys {
		x := SErrors[i2]
		markDown += fmt.Sprintf("| %v | %v | %v | %v |\n", x.ErrCode, x.ErrType, x.ParamDescription, x.FmtErrMsg)
	}
	return markDown
}
