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
	"bytes"
	"encoding/json"
	"fmt"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

const (
	LOGMESSAGEPREFIX = "sCustom"
	// Add Constants here
)

// List type's here

var (
// Add Variables here for the file (Remember, they are global)
)

func init() {
	sLogger.SetLogMessagePrefix(LOGMESSAGEPREFIX)
}

// JSONMarshalIndent is like JSONMarshal but applies Indent to format the output.
// Each JSON element in the output will begin on a new line beginning with prefix
// followed by one or more copies of indent according to the indentation nesting.
func JSONMarshalIndent(v interface{}, prefix, indent string) (buf []byte, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		tBuffer []byte
		buffer  bytes.Buffer
		err     error
	)

	if tBuffer, soteErr = JSONMarshal(v); soteErr.ErrCode == nil {
		if err = json.Indent(&buffer, tBuffer, prefix, indent); err == nil {
			buf = buffer.Bytes()
		} else {
			sLogger.Info(err.Error())
			soteErr = sError.GetSError(207110, sError.BuildParams([]string{fmt.Sprint(v)}), sError.EmptyMap)
		}
	}

	return
}

/*JSONMarshal converts interface into a JSON encoding*/
func JSONMarshal(v interface{}) (buffer []byte, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		tBuffer *bytes.Buffer
		err     error
	)

	tBuffer = &bytes.Buffer{}
	encoder := json.NewEncoder(tBuffer)
	encoder.SetEscapeHTML(false) // Disable escaping of HTML characters such as &, <, and >
	if err = encoder.Encode(v); err == nil {
		buffer = tBuffer.Bytes()
	} else {
		sLogger.Info(err.Error())
		soteErr = sError.GetSError(207110, sError.BuildParams([]string{fmt.Sprint(v)}), sError.EmptyMap)
	}

	return
}
