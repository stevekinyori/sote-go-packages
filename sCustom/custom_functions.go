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
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/nats-io/nats.go"

	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
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
			soteErr = sError.GetSError(sError.ErrInvalidJSON, sError.BuildParams([]string{fmt.Sprint(v)}), sError.EmptyMap)
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
		soteErr = sError.GetSError(sError.ErrInvalidJSON, sError.BuildParams([]string{fmt.Sprint(v)}), sError.EmptyMap)
	}

	return
}

// JSONUnmarshal parses the JSON-encoded data and stores the result in the value pointed to by response
func JSONUnmarshal(ctx context.Context, data []byte, response interface{}) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(data) > 0 { // only unmarshall if message is not empty
		if err := json.Unmarshal(data, &response); err != nil {
			sLogger.Info(err.Error())
			soteErr = sError.GetSError(sError.ErrInvalidJSON, sError.BuildParams([]string{string(data)}), sError.EmptyMap)
		}
	}

	return
}

// CallUserFunc calls Exportable Methods using their method Name. This does not work on functions
func CallUserFunc(funcName string, receiver any, args ...any) (response []reflect.Value, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if method := reflect.ValueOf(receiver).MethodByName(funcName); method.IsValid() {
		var (
			expectedParams = method.Type().NumIn()
			providedParams = len(args)
			methodType     = method.Type()
		)
		if expectedParams != providedParams && !methodType.IsVariadic() {
			soteErr = sError.GetSError(sError.ErrInvalidParameterCount,
				sError.BuildParams([]string{fmt.Sprint(providedParams), fmt.Sprint(expectedParams)}), nil)
			return
		}

		inputs := make([]reflect.Value, providedParams)
		for i := range args {
			var (
				inType   reflect.Type
				argValue reflect.Value
			)
			if methodType.IsVariadic() && i >= expectedParams-1 {
				inType = methodType.In(expectedParams - 1).Elem()
			} else {
				inType = methodType.In(i)
			}

			if argValue = reflect.ValueOf(args[i]); !argValue.IsValid() {
				soteErr = sError.GetSError(sError.ErrInvalidParameterType,
					sError.BuildParams([]string{argValue.String(), fmt.Sprint(inType)}), nil)
				return
			}

			argValueType := argValue.Type()
			switch true {
			case inType.String() == "context.Context" && argValueType.String() == "*context.emptyCtx":
				// skip this check
			case argValue.Type() != inType:
				// (inType.String() != "context.Context")
				soteErr = sError.GetSError(sError.ErrInvalidParameterType,
					sError.BuildParams([]string{argValue.String(), fmt.Sprint(inType)}), nil)
				return
			}

			inputs[i] = argValue
		}

		response = method.Call(inputs)
		return
	}

	var errMsg string
	if receiver == nil {
		errMsg = fmt.Sprintf("Function %v", funcName)
	} else {
		errMsg = fmt.Sprintf("Method (%v)%v", reflect.ValueOf(receiver).Type(), funcName)
	}

	soteErr = sError.GetSError(sError.ErrItemNotFound, sError.BuildParams([]string{errMsg}), nil)

	return
}

func UserFuncExists(funcName string, receiver any) (soteErr sError.SoteError) {
	if method := reflect.ValueOf(receiver).MethodByName(funcName); method.IsValid() {
		return
	}

	var errMsg string
	if receiver == nil {
		errMsg = fmt.Sprintf("Function %v", funcName)
	} else {
		errMsg = fmt.Sprintf("Method (%v)%v", reflect.ValueOf(receiver).Type(), funcName)
	}

	soteErr = sError.GetSError(sError.ErrItemNotFound, sError.BuildParams([]string{errMsg}), nil)

	return
}

// CopyDir copies all files from a specific director to another
// empty ext means all files
func CopyDir(source, destination string, ext string) (soteErr sError.SoteError) {
	if err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		var relPath = strings.Replace(path, source, "", 1)
		if relPath == "" {
			fmt.Println(source, destination)

			return nil
		}
		if info.IsDir() {
			return os.MkdirAll(filepath.Join(destination, relPath), os.ModePerm)
		} else if ext == "" || filepath.Ext(info.Name()) == ext {
			var data, err1 = ioutil.ReadFile(filepath.Join(source, relPath))
			if err1 != nil {
				return err1
			}

			if err2 := ioutil.WriteFile(filepath.Join(destination, relPath), data, os.ModePerm); err2 != nil {
				return err2
			}

			for start := time.Now(); ; {
				if _, err3 := os.Stat(filepath.Join(destination, relPath)); err3 == nil || time.Since(start) >= time.Second {
					return nil
				}
			}

		} else {
			return nil
		}
	}); err != nil {
		soteErr = sError.GetSError(sError.ErrGenericError, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
	}

	return
}

// CopyFile copies a files from a specific director to another
func CopyFile(source, destination string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var err error

	if runtime.GOOS == "windows" {
		var (
			data []byte
		)

		if data, err = ioutil.ReadFile(source); err == nil {
			err = ioutil.WriteFile(destination, data, os.ModePerm)
		}
	} else {
		err = exec.Command("cp", source, destination).Run()
	}

	if err == nil {
		for start := time.Now(); ; {
			if _, err = os.Stat(destination); err == nil || time.Since(start) >= time.Second {
				return
			}

		}

	} else {
		soteErr = sError.GetSError(sError.ErrGenericError, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
	}

	return
}
