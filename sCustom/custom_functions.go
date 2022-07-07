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
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"

	"gitlab.com/soteapps/packages/v2022/sError"
	"gitlab.com/soteapps/packages/v2022/sLogger"
)

const (
	LOGMESSAGEPREFIX = "sCustom"
	// Add Constants here
)

// Options are the options expected by the panic service function
type Options struct {
	Testmode           bool
	AcknowledgeNatsMsg func(ctx context.Context) (*nats.Msg, bool)
	Server             string
	AppEnvironment     string
}

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

// PanicService panic when not in test mode/production/demo
func PanicService(ctx context.Context, inSoteErr sError.SoteError, opts Options) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	sLogger.Info(inSoteErr.FmtErrMsg)
	soteErr = inSoteErr
	if !opts.Testmode {
		if natsMsgPtr, ok := opts.AcknowledgeNatsMsg(ctx); ok { // at this point the message has been processed and if it's a NATS message, it should be acknowledged
			sLogger.Info(fmt.Sprintf("PANIC - FilterSubject: %v Message Body: %v", natsMsgPtr.Subject,
				string(natsMsgPtr.Data)))
		}

		if opts.Server == "nats" && (opts.AppEnvironment == "development" || opts.AppEnvironment == "staging") {
			defer func() {
				if r := recover(); r != nil {
					sLogger.Info(fmt.Sprintf("Recovered from panic %v", r))
				}
			}()

			panic(soteErr.FmtErrMsg)
		}

		if inSoteErr.ErrCode != 199999 {
			soteErr = sError.GetSError(199999, sError.BuildParams([]string{""}), sError.EmptyMap)
		}
	}

	return
}
