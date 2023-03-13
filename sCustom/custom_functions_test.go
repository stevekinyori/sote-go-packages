/*
General description of the purpose of the go file.

RESTRICTIONS:
    AWS functions:
    * Program must have access to a .aws/credentials file in the default location.
    * This will only access system parameters that start with '/sote' (ROOTPATH).
    * {Enter other restrictions here for AWS

    {Other catagories of restrictions}
    * {List of restrictions for the catagory

NOTES:
    {Enter any additional notes that you believe will help the next developer.}
*/
package sCustom

import (
	// Add imports here

	"runtime"
	"testing"

	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
)

// List type's here

var (
// Add Variables here for the file (Remember, they are global)
)

func init() {
	sLogger.SetLogMessagePrefix(LOGMESSAGEPREFIX)
}

func TestMarshal(tPtr *testing.T) {
	var (
		v = map[string]string{"note": "Example & test note"}
	)

	if _, soteErr := JSONMarshalIndent(v, PREFIX, INDENT); soteErr.ErrCode != nil {
		tPtr.Errorf("TestMarshal Failed: Expected error code to be %v but got %v", nil, soteErr.FmtErrMsg)
	}

	if _, soteErr := JSONMarshal(v); soteErr.ErrCode != nil {
		tPtr.Errorf("TestMarshal Failed: Expected error code to be %v but got %v", nil, soteErr.FmtErrMsg)
	}
}
func TestParseEmail(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	tPtr.Run("valid email with special characters", func(tPtr *testing.T) {
		if _, soteErr = ParseEmail("\r\n\ttuser.test@example.com\n\t"); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("valid email with extra characters", func(tPtr *testing.T) {
		if _, soteErr = ParseEmail("User <user.test@example.com>"); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("invalid email", func(tPtr *testing.T) {
		if _, soteErr = ParseEmail("example"); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})
}
