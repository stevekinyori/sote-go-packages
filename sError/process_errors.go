package sError

import (
	"fmt"
)

const (
	// Error Category
	PROCESS_ERROR = "config-error"

	DIRTY_READ                             = 100200
	CANCELLED_COMPLETE                     = 100500
	INACTIVE_ITEM                          = 100600
	TIME_OUT                               = 101010
	INVALID_DATA_TYPE                      = 200200
	REQUIRED_VALUE                         = 200250
	LINKED_PARAM_VALUE_MISSING             = 200260
	PARAM_LOCK_OTHER_PARAM_SET             = 200510
	LINKED_PARAMS_MUST_BOTH_BE_SET_OR_NULL = 200511
	PARAMS_MUST_BE_PROVIDED                = 200512
	PARAM_MUST_BE_SET                      = 200513
	LINKED_PARAM_EMPTY_WHEN_PARAM_SET      = 200515
	BAD_HTTP_REQUEST                       = 200600
	INVALID_API_ENVIRONMENT                = 200700
	QUICK_SIGHT_ERROR                      = 200800
	DATABASE_CONSTRAINT_ERROR              = 200900
	SQL_ERROR                              = 200999
	COGNITO_ERROR                          = 201999
	AWS_SES_ERROR                          = 205000
	AWS_STS_ERROR                          = 205005
)

func ErrDirtyRead(operation string, err error) error {
	return &SoteError{
		ErrCode:   DIRTY_READ,
		ErrType:   PROCESS_ERROR,
		FmtErrMsg: "Row has been updated since reading it, re-read the row",
		Err:       err,
		Loc:       operation,
	}
}

func ErrCancelledComplete(item interface{}, operation string, err error) error {
	return &SoteError{
		ErrCode: CANCELLED_COMPLETE,
		ErrType: PROCESS_ERROR,
		FmtErrMsg: fmt.Sprintf(
			"Item being changed: You are making changes to a canceled or completed %v",
			item,
		),
		Err: err,
		Loc: operation,
	}
}

func ErrInactiveItem(item interface{}, operation string, err error) error {
	return &SoteError{
		ErrCode: INACTIVE_ITEM,
		ErrType: PROCESS_ERROR,
		FmtErrMsg: fmt.Sprintf(
			"Item is not active: You are making changes to an inactive %v",
			item,
		),
		Err: err,
		Loc: operation,
	}
}

func ErrTimeOut(service interface{}, operation string, err error) error {
	return &SoteError{
		ErrCode: TIME_OUT,
		ErrType: PROCESS_ERROR,
		FmtErrMsg: fmt.Sprintf(
			"Service: [%v] timed out",
			service,
		),
		Err: err,
		Loc: operation,
	}
}

func ErrTableDoesNotExist(operation string, err error) error {
	return &SoteError{
		ErrCode:   TIME_OUT,
		ErrType:   PROCESS_ERROR,
		FmtErrMsg: "Table does not exist",
		Err:       err,
		Loc:       operation,
	}
}

func ErrInvalidDataType(expectedType, param interface{}, operation string, err error) error {
	return &SoteError{
		ErrCode: INVALID_DATA_TYPE,
		ErrType: PROCESS_ERROR,
		FmtErrMsg: fmt.Sprintf(
			"Parameter name, Data type of parameter: %v must be of type %v",
			param, expectedType),
		Loc: operation,
		Err: err,
	}
}

func ErrMissingRequiredValue(paramName, paramValue, listofValuesAllowed interface{}, err error) error {
	return &SoteError{
		ErrCode: REQUIRED_VALUE,
		ErrType: PROCESS_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": %v (%v) must contain one of these values: %v",
			paramName, paramValue, listofValuesAllowed),
		Err: err,
	}
}

func ErrLinkedParameterValueMissing(requiredParam, linkedParam, linkedParamValue interface{}, err error) error {
	return &SoteError{
		ErrCode: LINKED_PARAM_VALUE_MISSING,
		ErrType: PROCESS_ERROR,
		FmtErrMsg: fmt.Sprintf(": %v must be provided when %v is set to (%v)",
			requiredParam, linkedParam, linkedParamValue),
		Err: err,
	}
}

func ErrParameterLockOtherParameterSet(lockedParam, otherParam, otherParamValue interface{}, err error) error {
	return &SoteError{
		ErrCode: PARAM_LOCK_OTHER_PARAM_SET,
		ErrType: PROCESS_ERROR,
		FmtErrMsg: fmt.Sprintf(": %v can't be updated because %v is set to %v",
			lockedParam, otherParam, otherParamValue),
		Err: err,
	}
}

func ErrLinkedParamsMustBothBeSetOrNull(paramName, linkedParam interface{}, err error) error {
	return &SoteError{
		ErrCode: LINKED_PARAMS_MUST_BOTH_BE_SET_OR_NULL,
		ErrType: PROCESS_ERROR,
		FmtErrMsg: fmt.Sprintf(": %v and %v must both be populated or null",
			paramName, linkedParam),
		Err: err,
	}
}

func ErrParamsMustBeProvided(params interface{}, err error) error {
	return &SoteError{
		ErrCode: PARAMS_MUST_BE_PROVIDED,
		ErrType: PROCESS_ERROR,
		FmtErrMsg: fmt.Sprintf(": %v :- must be populated",
			params),
		Err: err,
	}
}

func ErrParamMustBeSet(param interface{}, err error) error {
	return &SoteError{
		ErrCode: PARAMS_MUST_BE_PROVIDED,
		ErrType: PROCESS_ERROR,
		FmtErrMsg: fmt.Sprintf(": %v must be populated",
			param),
		Err: err,
	}
}

func ErrParameterMustBeEmptyWhenParameterSet(param, linkedParam interface{}, err error) error {
	return &SoteError{
		ErrCode: LINKED_PARAM_EMPTY_WHEN_PARAM_SET,
		ErrType: PROCESS_ERROR,
		FmtErrMsg: fmt.Sprintf(": %v must be empty when %v is populated",
			param, linkedParam),
		Err: err,
	}
}

func ErrBadHTTPRequest(req interface{}, err error) error {
	return &SoteError{
		ErrCode: BAD_HTTP_REQUEST,
		ErrType: PROCESS_ERROR,
		FmtErrMsg: fmt.Sprintf(": Bad HTTP/HTTPS Request - %v",
			req),
		Err: err,
	}
}

func ErrInvalidAPIEnvironment(environmentName interface{}, err error) error {
	return &SoteError{
		ErrCode: INVALID_API_ENVIRONMENT,
		ErrType: PROCESS_ERROR,
		FmtErrMsg: fmt.Sprintf(": The API you are calling is not available in this environment (%v)",
			environmentName),
		Err: err,
	}
}

func ErrQuickSight(err error) error {
	return &SoteError{
		ErrCode:   QUICK_SIGHT_ERROR,
		ErrType:   PROCESS_ERROR,
		FmtErrMsg: ": QuickSight error - see Details",
		Err:       err,
	}
}

func ErrDatabaseError(err error) error {
	return &SoteError{
		ErrCode:   DATABASE_CONSTRAINT_ERROR,
		ErrType:   PROCESS_ERROR,
		FmtErrMsg: ": Database constraint error - see Details",
		Err:       err,
	}
}

func ErrSQLError(err error) error {
	return &SoteError{
		ErrCode:   SQL_ERROR,
		ErrType:   PROCESS_ERROR,
		FmtErrMsg: ": SQL error - see Details",
		Err:       err,
	}
}

func ErrCognitoError(err error) error {
	return &SoteError{
		ErrCode:   COGNITO_ERROR,
		ErrType:   PROCESS_ERROR,
		FmtErrMsg: ": Cognito error - see Details",
		Err:       err,
	}
}

func ErrAwsSESError(err error) error {
	return &SoteError{
		ErrCode:   AWS_SES_ERROR,
		ErrType:   PROCESS_ERROR,
		FmtErrMsg: ": AWS SES error - see details in retPack",
		Err:       err,
	}
}

func ErrAwsSTSError(err error) error {
	return &SoteError{
		ErrCode:   AWS_STS_ERROR,
		ErrType:   PROCESS_ERROR,
		FmtErrMsg: ": AWS STS error - see details in retPack",
		Err:       err,
	}
}
