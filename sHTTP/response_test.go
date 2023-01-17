package sHTTP

import (
	"runtime"
	"testing"

	"gitlab.com/soteapps/packages/v2023/sCustom"
	"gitlab.com/soteapps/packages/v2023/sError"
)

func TestProcessLeafResponse(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		testMessage       = testRequestMessage{}
		soteErr           sError.SoteError
		response          []byte
	)

	tPtr.Run("valid", func(tPtr *testing.T) {
		reqMessage := Response{
			MessageId: "",
			Message: testRequestMessage{
				Id:   1,
				Name: "name",
				Age:  1,
			},
		}
		if response, soteErr = sCustom.JSONMarshal(reqMessage); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}

		if soteErr = ProcessLeafResponse(parentCtx, response, &testMessage, TESTMODE); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("invalid response message", func(tPtr *testing.T) {
		if soteErr = ProcessLeafResponse(parentCtx, []byte("test-message"), &testMessage, TESTMODE); soteErr.ErrCode != sError.ErrInvalidJSON {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrInvalidJSON, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("invalid sote error", func(tPtr *testing.T) {
		reqMessage := Response{
			MessageId: "",
			Error:     TESTAWSNAME,
		}
		if response, soteErr = sCustom.JSONMarshal(reqMessage); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}

		if soteErr = ProcessLeafResponse(parentCtx, response, &testMessage, TESTMODE); soteErr.ErrCode != sError.ErrInvalidJSON {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrInvalidJSON, soteErr.FmtErrMsg)
		}
	})
}
