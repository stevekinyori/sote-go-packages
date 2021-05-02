package sHelper

import (
	"testing"

	"gitlab.com/soteapps/packages/v2021/sError"
)

func TestErrorItem(t *testing.T) {
	// Zero arguments
	error := NewError().InternalError()
	verifyError(t, error, 210599, sError.GENERALERROR, "210599: Business Service error has occurred that is not expected.")
}

func TestErrorMultipleParameters(t *testing.T) {
	// Multiple Arguments
	error := NewError().AllowValues("a", "b", "c")
	verifyError(t, error, 200250, sError.PROCESSERROR, "200250: a (b) must contain one of these values: c")
}

func TestErrorItemDetails(t *testing.T) {
	// The best way is to define errorDetails key(s) as function argument(s). See - sErrorHelper :: SqlError
	error := NewError(map[string]string{"MyKey": "MyValue"}).ItemNotFound("Item")
	verifyError(t, error, 109999, sError.USERERROR, "109999: Item was/were not found ERROR DETAILS: >>Key: MyKey Value: MyValue")

	// Custom errorDetails
	error = NewError().SqlError("Connection Timeout")
	verifyError(t, error, 200999, sError.PROCESSERROR, "200999: SQL error - see Details ERROR DETAILS: >>Key: SQL ERROR Value: Connection Timeout")
}

func TestErrorCodes(t *testing.T) {
	// Any new error functions for code coverage should be defined here.
	verifyError(t, NewError().AlreadyExists("Item"), 100000, sError.USERERROR, "100000: Item already exists")
	verifyError(t, NewError().ItemNotFound("Item"), 109999, sError.USERERROR, "109999: Item was/were not found")
	verifyError(t, NewError().SqlError("Connection failed"), 200999, sError.PROCESSERROR,
		"200999: SQL error - see Details ERROR DETAILS: >>Key: SQL ERROR Value: Connection failed")
	verifyError(t, NewError().AllowValues("a", "b", []int{1, 2, 3}), 200250, sError.PROCESSERROR, "200250: a (b) must contain one of these values: [1 2 3]")
	verifyError(t, NewError().MustBePopulated("a, b, c"), 200513, sError.PROCESSERROR, "200513: a, b, c must be populated")
	verifyError(t, NewError().MustBeType("a, b, c", []int{1, 2, 3}), 200200, sError.PROCESSERROR, "200200: a, b, c must be of type [1 2 3]")
	verifyError(t, NewError().InvalidJson("./schema.json"), 207110, sError.CONTENTERROR, "207110: ./schema.json couldn't be parsed - Invalid JSON error")
	verifyError(t, NewError().InvalidParameters("Item"), 206200, sError.NATSERROR,
		"206200: Message doesn't match signature. Sender must provide the following parameter names: Item")
	verifyError(t, NewError().NoDbConnection(), 209299, sError.CONFIGURATIONISSUE, "209299: No database connection has been established")
	verifyError(t, NewError().FileNotFound("foo.json", "/foo.json"), 209010, sError.CONFIGURATIONISSUE, "209010: foo.json file was not found. Message return: /foo.json")
}

func verifyError(t *testing.T, error sError.SoteError, code int, categoty string, message string) {
	t.Helper()
	AssertEqual(t, error.ErrCode, code)
	AssertEqual(t, error.ErrType, categoty)
	AssertEqual(t, error.FmtErrMsg, message)
}
