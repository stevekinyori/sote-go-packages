package sError

import "fmt"

const (
	GENERAL_ERROR = "general-error"
)

func ErrUnexpected(logData interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode: ITEM_NOT_FOUND,
		ErrType: USER_ERROR,
		FmtErrMsg: fmt.Sprintf(
			"Error Details: An error has occurred that is not expected. See Log! %v", logData),
		Loc: operation,
		Err: Err,
	}
}
