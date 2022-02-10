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
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgconn"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

type SoteError struct {
	ErrCode          interface{}
	ErrType          string
	ParamCount       int
	ParamDescription string
	FmtErrMsg        string
	ErrorDetails     map[string]string
	Loc              string
	Err              error
}

// Error categories
const (
	USERERROR          = "User_Error"
	PROCESSERROR       = "Process_Error"
	NATSERROR          = "NATS_Error"
	CONTENTERROR       = "Content_Error"
	PERMISSIONERROR    = "Permission_Error"
	CONFIGURATIONISSUE = "Configuration_Issue"
	APICONTRACTERROR   = "API_Contract_Error"
	GENERALERROR       = "General_Error"
	MARKDOWNTITLEBAR   = "| Error Code | Category | Parameter Description | Formatted Error Text |\n|--------|--------|--------|--------|\n"
	FUNCCOMMENTSHEADER = "\tError Code with requiring parameters:\n"
	SQLSTATE           = "SQLSTATE"
	PREFIX             = ""
	INDENT             = "  "
)

var (
	EmptyMap = make(map[string]string)
	/*
		Error code Ranges are not limited to a single error category.
		The required format for the FmtErrMsg value is ": <sprintf string>" where the ":" followed by the space must start the value.

		ErrType is for general grouping and doesn't affect the use of the error. If the error fits your needs, disregard the ErrType.
	*/
	soteErrors = map[int]SoteError{
		// Errors where the Front End can take action
		100000: {100000, USERERROR, 1, "Item Name", ": %v already exists", EmptyMap, "", nil},
		100100: {100100, USERERROR, 2, "List of users roles, Requested action", ": Your roles %v are not authorized to %v", EmptyMap, "", nil},
		100200: {100200, PROCESSERROR, 0, "None", ": Row has been updated since reading it, re-read the row", EmptyMap, "", nil},
		100500: {100500, PROCESSERROR, 1, "Thing being changed", ": You are making changes to a canceled or completed %v", EmptyMap, "", nil},
		100600: {100600, PROCESSERROR, 1, "Item is not active", ": You are making changes to an inactive %v", EmptyMap, "", nil},
		101010: {101010, PROCESSERROR, 1, "Service Name", ": %v timed out", EmptyMap, "", nil},
		109999: {109999, USERERROR, 1, "Item name", ": %v was/were not found", EmptyMap, "", nil},
		199999: {199999, GENERALERROR, 1, "Error Details", ": An error has occurred that is not expected. See Log! %v", EmptyMap, "", nil},
		// ======================================================================
		// Errors where the Back End can take action or the system needs to panic
		200100: {200100, PROCESSERROR, 0, "None", ": Table doesn't exist", EmptyMap, "", nil},
		200200: {200200, PROCESSERROR, 2, "Parameter name, Data type of parameter", ": %v must be of type %v", EmptyMap, "", nil},
		200250: {200250, PROCESSERROR, 3, "Parameter name, Parameter value, List of values allowed", ": %v (%v) must contain one of these values: %v",
			EmptyMap, "", nil},
		200260: {200260, PROCESSERROR, 3, "Other parameter name, Parameter name, Parameter value", ": %v must be provided when %v is set to (%v)",
			EmptyMap, "", nil},
		200510: {200510, PROCESSERROR, 3, "Parameter name, Field name, Field value", ": %v can't be updated because %v is set to %v", EmptyMap, "", nil},
		200511: {200511, PROCESSERROR, 2, "Parameter name, Another parameter name", ": %v and %v must both be populated or null", EmptyMap, "", nil},
		200512: {200512, PROCESSERROR, 2, "Parameter name, Another parameter name", ": %v and %v must both be populated", EmptyMap, "", nil},
		200513: {200513, PROCESSERROR, 1, "Parameter name", ": %v must be populated", EmptyMap, "", nil},
		200514: {200514, PROCESSERROR, 3, "Parameter name, Another parameter name, Another parameter name", ": %v, %v and %v must all be populated",
			EmptyMap, "", nil},
		200515: {200515, PROCESSERROR, 2, "Parameter name, Another parameter name", ": %v must be empty when %v is populated", EmptyMap, "", nil},
		200600: {200600, PROCESSERROR, 1, "Info returned from HTTP/HTTPS Request", ": Bad HTTP/HTTPS Request - %v", EmptyMap, "", nil},
		200700: {200700, PROCESSERROR, 1, "Environment Name", ": The API you are calling is not available in this environment (%v)", EmptyMap, "", nil},
		200800: {200800, PROCESSERROR, 0, "None", ": QuickSight error - see Details", EmptyMap, "", nil},
		200900: {200900, PROCESSERROR, 0, "None", ": Database constraint error - see Details", EmptyMap, "", nil},
		200999: {200999, PROCESSERROR, 0, "None", ": SQL error - see Details", EmptyMap, "", nil},
		201999: {201999, PROCESSERROR, 0, "None", ": Cognito error - see Details", EmptyMap, "", nil},
		203000: {203000, PROCESSERROR, 0, "None", ": The number of parameters provided for the error message does not match the required number",
			EmptyMap, "", nil},
		203050: {203050, PROCESSERROR, 2, "Name, Application/Package name", ": Number of parameters defined in the %v is not support by %v", EmptyMap,
			"", nil},
		203060: {203060, PROCESSERROR, 2, "Provided parameter count, Expected parameter count",
			": Number of parameters provided (%v) doesn't match the number expected (%v)", EmptyMap, "", nil},
		205000: {205000, PROCESSERROR, 0, "None", ": AWS SES error - see details in retPack", EmptyMap, "", nil},
		205005: {205005, PROCESSERROR, 0, "None", ": AWS STS error - see details in retPack", EmptyMap, "", nil},
		206000: {206000, NATSERROR, 0, "None", ": Jetstream is not enabled", EmptyMap, "", nil},
		206050: {206050, NATSERROR, 2, "Subscription Name, Subject", ": (%v) is an invalid subscription. Subject: %v", EmptyMap, "", nil},
		206100: {206100, NATSERROR, 1, "Key name", ": Upper or lower case %v key is missing", EmptyMap, "", nil},
		206105: {206105, NATSERROR, 1, "Key name", ": Upper or lower case %v keys value is missing", EmptyMap, "", nil},
		206200: {206200, NATSERROR, 1, "List of required parameters",
			": Message doesn't match signature. Sender must provide the following parameter names: %v", EmptyMap, "", nil},
		206300: {206300, NATSERROR, 0, "None", ": Stream pointer is nil. Must be a validate pointer to a stream.", EmptyMap, "", nil},
		206400: {206400, NATSERROR, 1, "Stream Name", ": Stream creation encountered an error that is not expected. Stream Name: %v", EmptyMap, "", nil},
		206600: {206600, NATSERROR, 2, "Stream Name, Consumer Name", ": Consumer creation encountered an error that is not expected. " +
			"Stream Name: %v Consumer Name: %v", EmptyMap, "", nil},
		206700: {206700, NATSERROR, 2, "Stream Name, Consumer Subject Filter",
			": The consumer subject filter must be a subset of the stream subject. " +
				"Stream Name: %v Consumer Subject Filter: %v", EmptyMap, "", nil},
		207000: {207000, CONTENTERROR, 2, "Field name, Field value", ": %v (%v) is not numeric", EmptyMap, "", nil},
		207005: {207005, CONTENTERROR, 2, "Field name, Minimal length", ": %v must have a value greater than %v", EmptyMap, "", nil},
		207010: {207010, CONTENTERROR, 2, "Field name, Field value", ": %v (%v) is not a string", EmptyMap, "", nil},
		207020: {207020, CONTENTERROR, 2, "Field name, Field value", ": %v (%v) is not a float", EmptyMap, "", nil},
		207030: {207030, CONTENTERROR, 2, "Field name, Field value", ": %v (%v) is not a array", EmptyMap, "", nil},
		207040: {207040, CONTENTERROR, 2, "Field name, Field value", ": %v (%v) is not a json string", EmptyMap, "", nil},
		207050: {207050, CONTENTERROR, 2, "Field name, Field value", ": %v (%v) is not a valid email address", EmptyMap, "", nil},
		207060: {207060, CONTENTERROR, 2, "Field name, Field value", ": %v (%v) contains special characters which are not allowed", EmptyMap, "", nil},
		207065: {207065, CONTENTERROR, 2, "Field name, Field value", ": %v (%v) contains special characters other than underscore", EmptyMap, "", nil},
		207070: {207070, CONTENTERROR, 2, "Field name, Field value", ": %v (%v) is not a valid date", EmptyMap, "", nil},
		207080: {207080, CONTENTERROR, 2, "Field name, Field value", ": %v (%v) is not a valid timestamp. Format's are UTC, GMT or Zulu", EmptyMap,
			"", nil},
		207090: {207090, CONTENTERROR, 6, "Field name, Field value, 'small' or 'large', 'Min' or 'Max', expected size, actual size",
			": %v (%v) is too %v. %v size: %v Actual size: %v", EmptyMap, "", nil},
		207095: {207095, CONTENTERROR, 4, "Field name, Field value, greater than value, less than value",
			": %v (%v) must be greater than %v and less than %v", EmptyMap, "", nil},
		207100: {207100, CONTENTERROR, 2, "Parameter name, Data Structure Type", ": %v couldn't be converted to an %v - JSON conversion error",
			EmptyMap, "", nil},
		207105: {207105, CONTENTERROR, 2, "Data Structure Name, Data Structure Type",
			": %v (%v) couldn't be converted to JSON - JSON conversion error", EmptyMap, "", nil},
		207110: {207110, CONTENTERROR, 1, "Parameter name", ": %v couldn't be parsed - Invalid JSON error", EmptyMap, "", nil},
		207111: {207111, CONTENTERROR, 2, "Parameter name, Application/Package name", ": %v couldn't be converted to a map/keyed array - %v",
			EmptyMap, "", nil},
		207200: {207200, CONTENTERROR, 2, "Parameter name, Data Structure Type", ": %v couldn't be converted to an %v", EmptyMap, "", nil},
		208000: {208000, CONTENTERROR, 0, "None", ": Column must have a non-null value. Details: ", EmptyMap, "", nil},
		208010: {208010, CONTENTERROR, 0, "None", ": Column data type is not support or invalid. Details: ", EmptyMap, "", nil},
		208110: {208110, CONTENTERROR, 2, "Thing being changed, System Id for the thing",
			": No update is needed. No fields where changed for %v with id %v", EmptyMap, "", nil},
		208120: {208120, CONTENTERROR, 3, "JSON array name, Thing being changed, System Id for the thing", ": The %v was empty for %v with id %v",
			EmptyMap, "", nil},
		208200: {208200, CONTENTERROR, 1, "Error message number", ": %v error message is missing from sError package", EmptyMap, "", nil},
		208300: {208300, PERMISSIONERROR, 0, "None", ": iss (Issuer) is not valid", EmptyMap, "", nil},
		208310: {208310, PERMISSIONERROR, 1, "Subject", ": sub (Subject: %v) was not present", EmptyMap, "", nil},
		208320: {208320, PERMISSIONERROR, 0, "None", ": token_use is not valid", EmptyMap, "", nil},
		208330: {208330, PERMISSIONERROR, 0, "None", ": client id is not valid", EmptyMap, "", nil},
		208340: {208340, PERMISSIONERROR, 0, "None", ": client id is not valid for this application", EmptyMap, "", nil},
		208350: {208350, PERMISSIONERROR, 0, "None", ": Token is expired", EmptyMap, "", nil},
		208355: {208355, PERMISSIONERROR, 0, "None", ": Token is invalid", EmptyMap, "", nil},
		208356: {208356, PERMISSIONERROR, 0, "None", ": Token contains an invalid number of segments", EmptyMap, "", nil},
		208360: {208360, PERMISSIONERROR, 1, "Claim names", ": These claims are invalid: %v", EmptyMap, "", nil},
		208370: {208370, PERMISSIONERROR, 0, "None", ": Required claim(s) is/are missing", EmptyMap, "", nil},
		209000: {209000, CONFIGURATIONISSUE, 0, "None", ": .env files are missing", EmptyMap, "", nil},
		209010: {209010, CONFIGURATIONISSUE, 2, "File name, Message returned from Open", ": %v file was not found. Message return: %v", EmptyMap, "", nil},
		209100: {209100, CONFIGURATIONISSUE, 1, "Environment name", ": environment variable is missing (%v)", EmptyMap, "", nil},
		209110: {209110, CONFIGURATIONISSUE, 1, "Environment name", ": environment value (%v) is invalid", EmptyMap, "", nil},
		209200: {209200, CONFIGURATIONISSUE, 3, "Database name, Database driver name, Port value",
			": Unable to connect to database %v using driver %v on port %v", EmptyMap, "", nil},
		209210: {209210, CONFIGURATIONISSUE, 0, "None", ": Unable to pass database authentication", EmptyMap, "", nil},
		209220: {209220, CONFIGURATIONISSUE, 1, "SSL Mode", ": Only disable, allow, prefer and required are supported.", EmptyMap, "", nil},
		209230: {209230, CONFIGURATIONISSUE, 1, "Connection Type", ": Only single or pool are supported.", EmptyMap, "", nil},
		209299: {209299, CONFIGURATIONISSUE, 0, "None", ": No database connection has been established", EmptyMap, "", nil},
		209398: {209398, CONFIGURATIONISSUE, 0, "None", ": no nkey seed found", EmptyMap, "", nil},
		209499: {209499, CONFIGURATIONISSUE, 0, "None", ": No nats connection has been established", EmptyMap, "", nil},
		209500: {209500, CONFIGURATIONISSUE, 0, "None", ": Unexpected signing method", EmptyMap, "", nil},
		209510: {209510, CONFIGURATIONISSUE, 0, "None", ": Kid header not found", EmptyMap, "", nil},
		209520: {209520, CONFIGURATIONISSUE, 1, "Kid", ": key (%v) was not found in token", EmptyMap, "", nil},
		209521: {209521, CONFIGURATIONISSUE, 1, "Kid", ": Kid (%v) was not found in public key set", EmptyMap, "", nil},
		210030: {210030, CONFIGURATIONISSUE, 2, "Region, Environment", ": Failed to fetch remote JWK (status = 404) for %v region %v environment",
			EmptyMap, "", nil},
		210090: {210090, CONFIGURATIONISSUE, 1, "Parameter name", ": URL is missing (%v)", EmptyMap, "", nil},
		210098: {210098, CONFIGURATIONISSUE, 1, "Parameter name", ": Start up parameter is out of value range (%v)", EmptyMap, "", nil},
		210099: {210099, CONFIGURATIONISSUE, 1, "Parameter name", ": Start up parameter is missing (%v)", EmptyMap, "", nil},
		210100: {210100, APICONTRACTERROR, 1, "List of required parameters",
			": Call doesn't match API signature. Caller must provide the following parameter names: %v", EmptyMap, "", nil},
		210200: {210200, GENERALERROR, 0, "None", ": Postgres error has occurred that is not expected.", EmptyMap, "", nil},
		210299: {210299, GENERALERROR, 0, "None", ": Postgres is not responding over TCP. Container may not be running.", EmptyMap, "", nil},
		210399: {210399, GENERALERROR, 0, "None", ": AWS session error has occurred that is not expected", EmptyMap, "", nil},
		210499: {210499, GENERALERROR, 0, "None", ": Synadia error has occurred that is not expected.", EmptyMap, "", nil},
		210599: {210599, GENERALERROR, 0, "None", ": Business Service error has occurred that is not expected.", EmptyMap, "", nil},
	}
)

// Error returns the string representation of the error message.
func (e *SoteError) Error() string {

	var buf bytes.Buffer

	// If wrapping an error, print its Error() message.
	// Otherwise print the error code & message.
	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.ErrCode != 0 {
			fmt.Fprintf(&buf, "<%s:%d> ", e.ErrType, e.ErrCode)
		}
		buf.WriteString(e.FmtErrMsg)
	}
	return buf.String()
}

func (e *SoteError) Unwrap() error { return e.Err }

/*
	This will return the formatted message using the supplied code and parameters

	Error Code with requiring parameters:
		100000	Item Name > : %v already exists
		100100	List of users roles, Requested action > : Your roles %v are not authorized to %v
		100200	None > : Row has been updated since reading it, re-read the row
		100500	Thing being changed > : You are making changes to a canceled or completed %v
		100600	Item is not active > : You are making changes to an inactive %v
		101010	Service Name > : %v timed out
		109999	Item name > : %v was/were not found
		199999	Error Details > : An error has occurred that is not expected. See Log! %v
		200100	None > : Table doesn't exist
		200200	Parameter name, Data type of parameter > : %v must be of type %v
		200250	Parameter name, Parameter value, List of values allowed > : %v (%v) must contain one of these values: %v
		200260	Other parameter name, Parameter name, Parameter value > : %v must be provided when %v is set to (%v)
		200510	Parameter name, Field name, Field value > : %v can't be updated because %v is set to %v
		200511	Parameter name, Another parameter name > : %v and %v must both be populated or null
		200512	Parameter name, Another parameter name > : %v and %v must both be populated
		200513	Parameter name > : %v must be populated
		200514	Parameter name, Another parameter name, Another parameter name > : %v, %v and %v must all be populated
		200515	Parameter name, Another parameter name > : %v must be empty when %v is populated
		200600	Info returned from HTTP/HTTPS Request > : Bad HTTP/HTTPS Request - %v
		200700	Environment Name > : The API you are calling is not available in this environment (%v)
		200800	None > : QuickSight error - see Details
		200900	None > : Database constraint error - see Details
		200999	None > : SQL error - see Details
		201999	None > : Cognito error - see Details
		203000	None > : The number of parameters provided for the error message does not match the required number
		203050	Name, Application/Package name > : Number of parameters defined in the %v is not support by %v
		203060	Provided parameter count, Expected parameter count > : Number of parameters provided (%v) doesn't match the number expected (%v)
		205000	None > : AWS SES error - see details in retPack
		205005	None > : AWS STS error - see details in retPack
		206000	None > : Jetstream is not enabled
		206050	Subscription Name, Subject > : (%v) is an invalid subscription. Subject: %v
		206100	Key name > : Upper or lower case %v key is missing
		206105	Key name > : Upper or lower case %v keys value is missing
		206200	List of required parameters > : Message doesn't match signature. Sender must provide the following parameter names: %v
		206300	None > : Stream pointer is nil. Must be a validate pointer to a stream.
		206400	Stream Name > : Stream creation encountered an error that is not expected. Stream Name: %v
		206600	Stream Name, Consumer Name > : Consumer creation encountered an error that is not expected. Stream Name: %v Consumer Name: %v
		206700	Stream Name, Consumer Subject Filter > : The consumer subject filter must be a subset of the stream subject. Stream Name: %v Consumer Subject Filter: %v
		207000	Field name, Field value > : %v (%v) is not numeric
		207005	Field name, Minimal length > : %v must have a value greater than %v
		207010	Field name, Field value > : %v (%v) is not a string
		207020	Field name, Field value > : %v (%v) is not a float
		207030	Field name, Field value > : %v (%v) is not a array
		207040	Field name, Field value > : %v (%v) is not a json string
		207050	Field name, Field value > : %v (%v) is not a valid email address
		207060	Field name, Field value > : %v (%v) contains special characters which are not allowed
		207065	Field name, Field value > : %v (%v) contains special characters other than underscore
		207070	Field name, Field value > : %v (%v) is not a valid date
		207080	Field name, Field value > : %v (%v) is not a valid timestamp. Format's are UTC, GMT or Zulu
		207090	Field name, Field value, 'small' or 'large', 'Min' or 'Max', expected size, actual size > : %v (%v) is too %v. %v size: %v Actual size: %v
		207095	Field name, Field value, greater than value, less than value > : %v (%v) must be greater than %v and less than %v
		207100	Parameter name, Data Structure Type > : %v couldn't be converted to an %v - JSON conversion error
		207105	Data Structure Name, Data Structure Type > : %v (%v) couldn't be converted to JSON - JSON conversion error
		207110	Parameter name > : %v couldn't be parsed - Invalid JSON error
		207111	Parameter name, Application/Package name > : %v couldn't be converted to a map/keyed array - %v
		207200	Parameter name, Data Structure Type > : %v couldn't be converted to an %v
		208000	None > : Column must have a non-null value. Details:
		208010	None > : Column data type is not support or invalid. Details:
		208110	Thing being changed, System Id for the thing > : No update is needed. No fields where changed for %v with id %v
		208120	JSON array name, Thing being changed, System Id for the thing > : The %v was empty for %v with id %v
		208200	Error message number > : %v error message is missing from sError package
		208300	None > : iss (Issuer) is not valid
		208310	Subject > : sub (Subject: %v) was not present
		208320	None > : token_use is not valid
		208330	None > : client id is not valid
		208340	None > : client id is not valid for this application
		208350	None > : Token is expired
		208355	None > : Token is invalid
		208356	None > : Token contains an invalid number of segments
		208360	Claim names > : These claims are invalid: %v
		208370	None > : Required claim(s) is/are missing
		209000	None > : .env files are missing
		209010	File name, Message returned from Open > : %v file was not found. Message return: %v
		209100	Environment name > : environment variable is missing (%v)
		209110	Environment name > : environment value (%v) is invalid
		209200	Database name, Database driver name, Port value > : Unable to connect to database %v using driver %v on port %v
		209210	None > : Unable to pass database authentication
		209220	SSL Mode > : Only disable, allow, prefer and required are supported.
		209230	Connection Type > : Only single or pool are supported.
		209299	None > : No database connection has been established
		209398	None > : no nkey seed found
		209499	None > : No nats connection has been established
		209500	None > : Unexpected signing method
		209510	None > : Kid header not found
		209520	Kid > : key (%v) was not found in token
		209521	Kid > : Kid (%v) was not found in public key set
		210030	Region, Environment > : Failed to fetch remote JWK (status = 404) for %v region %v environment
		210090	Parameter name > : URL is missing (%v)
		210098	Parameter name > : Start up parameter is out of value range (%v)
		210099	Parameter name > : Start up parameter is missing (%v)
		210100	List of required parameters > : Call doesn't match API signature. Caller must provide the following parameter names: %v
		210200	None > : Postgres error has occurred that is not expected.
		210299	None > : Postgres is not responding over TCP. Container may not be running.
		210399	None > : AWS session error has occurred that is not expected
		210499	None > : Synadia error has occurred that is not expected.
		210599	None > : Business Service error has occurred that is not expected.
*/
func GetSError(code int, params []interface{}, errorDetails map[string]string) (soteErr SoteError) {
	sLogger.DebugMethod()

	soteErr = soteErrors[code]
	if soteErr.ErrCode != code {
		s := make([]interface{}, 1)
		s[0] = code
		soteErr = GetSError(208200, s, errorDetails)
	} else if soteErr.ParamCount != len(params) {
		s := make([]interface{}, 2)
		s[0] = soteErr.ParamCount
		s[1] = len(params)
		soteErr = GetSError(203060, s, errorDetails)
	} else {
		if soteErr.ParamCount == 0 {
			soteErr.ErrorDetails = errorDetails
			soteErr.FmtErrMsg = strconv.Itoa(code) + fmt.Sprintf(soteErr.FmtErrMsg) + fmt.Sprint(formatErrorDetails(soteErr.ErrorDetails))
		} else {
			soteErr.ErrorDetails = errorDetails
			soteErr.FmtErrMsg = strconv.Itoa(code) + fmt.Sprintf(soteErr.FmtErrMsg, params...) + fmt.Sprint(formatErrorDetails(soteErr.ErrorDetails))
		}
	}
	_, file, no, ok := runtime.Caller(1)
	if ok {
		for i := 2; ok && strings.HasSuffix(file, "/packages/sHelper/error.go"); i++ {
			_, file, no, ok = runtime.Caller(i)
		}
		sLogger.Info(fmt.Sprintf("called from %s#%d", file, no))
	}
	return
}

// This will convert a postgresql error into error details for inclusion in SoteError
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
		soteErr = GetSError(207111, s, EmptyMap)
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
	this code the master source of Sote Error messages. Results are output to the console.
*/
func GenerateDocumentation() (markDown, funcComments string) {
	sLogger.DebugMethod()

	// Sort the Keys from SError map
	var errorKeys []int
	for _, i2 := range soteErrors {
		errorKeys = append(errorKeys, i2.ErrCode.(int))
	}
	sort.Ints(errorKeys)
	// Generate Documentation
	markDown = MARKDOWNTITLEBAR
	funcComments = FUNCCOMMENTSHEADER
	for _, i2 := range errorKeys {
		x := soteErrors[i2]
		markDown += fmt.Sprintf("| %v | %v | %v | %v |\n", x.ErrCode, x.ErrType, x.ParamDescription, x.FmtErrMsg)
		funcComments += fmt.Sprintf("\t\t%v\t%v > %v\n", x.ErrCode, x.ParamDescription, x.FmtErrMsg)
	}
	return
}

func OutputErrorJSON(inSoteErr SoteError) (outSoteErr []byte) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if outSoteErr, err = json.MarshalIndent(inSoteErr, PREFIX, INDENT+INDENT); err != nil {
		sLogger.Info(err.Error())
		soteErr := GetSError(207110, BuildParams([]string{"Sote Error"}), EmptyMap)
		PanicService(soteErr)
	}

	return
}

func PanicService(soteErr SoteError) {
	sLogger.DebugMethod()

	sLogger.Info(soteErr.FmtErrMsg)
	panic(soteErr.FmtErrMsg)
}

func formatErrorDetails(tErrDetails map[string]string) (errDetails string) {
	sLogger.DebugMethod()

	if len(tErrDetails) > 0 {
		errDetails = " ERROR DETAILS:"
		for key, value := range tErrDetails {
			errDetails = errDetails + " >>Key: " + key + " Value: " + value
		}
	}

	return
}
