package sError

import "fmt"

const (
	// Error Category
	USER_ERROR = "user-error"

	// Error Codes
	ITEM_ALREADY_EXISTS = 100000
	NOT_AUTHORIZED      = 100100
	ITEM_NOT_FOUND      = 109999
)

func ErrItemAlreadyExists(item interface{}, operation string) error {
	return &SoteError{
		ErrCode:   ITEM_ALREADY_EXISTS,
		ErrType:   USER_ERROR,
		FmtErrMsg: fmt.Sprintf("Item Name: %v already exists", item),
		Loc:       operation,
	}
}

func ErrNotAuthorized(userRoles, permissions interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode: NOT_AUTHORIZED,
		ErrType: USER_ERROR,
		FmtErrMsg: fmt.Sprintf(
			"List of users roles, Requested action: Your roles %v are not authorized to %v",
			userRoles, permissions),
		Loc: operation,
		Err: Err,
	}
}

func ErrItemNotFound(item interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode:   ITEM_NOT_FOUND,
		ErrType:   USER_ERROR,
		FmtErrMsg: fmt.Sprintf("Item name: %v was/were not found", item),
		Loc:       operation,
		Err:       Err,
	}
}
