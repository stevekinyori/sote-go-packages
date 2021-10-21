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

	"testing"

	"gitlab.com/soteapps/packages/v2021/sLogger"
)

const (
	// Add Constants here
	PREFIX string = ""
	INDENT string = "  "
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
