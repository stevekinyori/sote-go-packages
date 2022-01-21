package sError

import "fmt"

const (
	// Error Category
	CONTENT_ERROR = "content-error"

	NOT_NUMERIC           = 207000
	TOO_SMALL             = 207005
	NOT_STRING            = 207010
	NOT_FLOAT             = 207020
	NOT_ARRAY             = 207030
	NOT_JSON_STRING       = 207040
	INVALID_EMAIL         = 207050
	NOT_DATE              = 207070
	NOT_TIMESTAMP         = 207080
	INVALID_SIZE          = 207090
	JSON_CONVERSION_ERROR = 207105
	NOT_MAP               = 207111
	MISSING_ERROR_NUMBER  = 208200
)

func ErrNotNumeric(fieldName, fieldValue interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode: NOT_NUMERIC,
		ErrType: CONTENT_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": %v (%v) is not numeric",
			fieldName, fieldValue),
		Loc: operation,
		Err: Err,
	}
}

func ErrTooSmall(fieldName, minSize interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode: TOO_SMALL,
		ErrType: CONTENT_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": %v must have a value greater than %v",
			fieldName, minSize),
		Loc: operation,
		Err: Err,
	}
}

func ErrNotString(fieldName, fieldValue interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode: NOT_STRING,
		ErrType: CONTENT_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": %v (%v) is not a string",
			fieldName, fieldValue),
		Loc: operation,
		Err: Err,
	}
}

func ErrNotFloat(fieldName, fieldValue interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode: NOT_FLOAT,
		ErrType: CONTENT_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": %v (%v) is not a float",
			fieldName, fieldValue),
		Loc: operation,
		Err: Err,
	}
}

func ErrNotArray(fieldName, fieldValue interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode: NOT_ARRAY,
		ErrType: CONTENT_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": %v (%v) is not a array",
			fieldName, fieldValue),
		Loc: operation,
		Err: Err,
	}
}

func ErrNotJSONString(fieldName, fieldValue interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode: NOT_JSON_STRING,
		ErrType: CONTENT_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": %v (%v) is not a json string",
			fieldName, fieldValue),
		Loc: operation,
		Err: Err,
	}
}

func ErrInvalidEmail(fieldName, fieldValue interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode: INVALID_EMAIL,
		ErrType: CONTENT_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": %v (%v) is not a valid email address",
			fieldName, fieldValue),
		Loc: operation,
		Err: Err,
	}
}

func ErrNotDate(fieldName, fieldValue interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode: NOT_DATE,
		ErrType: CONTENT_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": %v (%v) is not a valid date",
			fieldName, fieldValue),
		Loc: operation,
		Err: Err,
	}
}

func ErrNotTimestamp(fieldName, fieldValue interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode: NOT_TIMESTAMP,
		ErrType: CONTENT_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": %v (%v) is not a valid timestamp. Format's are UTC, GMT or Zulu",
			fieldName, fieldValue),
		Loc: operation,
		Err: Err,
	}
}

func ErrInvalidSize(
	fieldName, fieldValue, relativeSize, expectedSize,
	actualSize interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode: INVALID_SIZE,
		ErrType: CONTENT_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": %v (%v) is too %v.size: %v Actual size: %v",
			fieldName, fieldValue, relativeSize, expectedSize, actualSize),
		Loc: operation,
		Err: Err,
	}
}

func ErrJsonConversionError(structName, structType interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode: JSON_CONVERSION_ERROR,
		ErrType: CONTENT_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": %v (%v) couldn't be converted to JSON - JSON conversion error",
			structName, structType),
		Loc: operation,
		Err: Err,
	}
}

func ErrNotMap(paramName interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode: NOT_MAP,
		ErrType: CONTENT_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": %v couldn't be converted to a map/keyed array",
			paramName),
		Loc: operation,
		Err: Err,
	}
}

func ErrMissingErrorNumber(errorMsgNumber interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode: NOT_MAP,
		ErrType: CONTENT_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": %v error message is missing from sError package",
			errorMsgNumber),
		Loc: operation,
		Err: Err,
	}
}
