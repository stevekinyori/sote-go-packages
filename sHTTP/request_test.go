package sHTTP

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.com/soteapps/packages/v2023/sAuthentication"
	"gitlab.com/soteapps/packages/v2023/sConfigParams"
	"gitlab.com/soteapps/packages/v2023/sError"
	"golang.org/x/exp/maps"
)

const (
	TESTAWSNAME        = "DELETE_ME_TEST"
	TESTORGANIZATIONID = 1
)

type testRequestMessage struct {
	Id   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required"`
}

func TestPrepareLeafRequest(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		header            = RequestHeaderParams{
			OrganizationsId: TESTORGANIZATIONID,
			AWSUsername:     TESTAWSNAME,
			MessageId:       "test-message",
			JSONWebToken:    fmt.Sprintf("%v %v", BEARERSCHEMA, TESTTOKEN),
			RoleList:        []string{"ROLE-1", "ROLE-2"},
		}
	)

	tPtr.Run("valid", func(tPtr *testing.T) {
		var (
			reqMessage = testRequestMessage{
				Id:   1,
				Name: TESTAWSNAME,
				Age:  1,
			}
		)
		if _, soteErr = PrepareLeafRequest(&LeafReqParams{
			ReqServiceURL:    LOCALHOST,
			ReqMessage:       &reqMessage,
			ReqHeaderMessage: &header,
			CurrentServiceId: sConfigParams.SHIPMENTSBSID,
		}, TESTMODE); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("invalid JSON", func(tPtr *testing.T) {
		var (
			reqMessage = make(chan int)
		)

		if _, soteErr = PrepareLeafRequest(&LeafReqParams{
			ReqServiceURL:    LOCALHOST,
			ReqMessage:       &reqMessage,
			ReqHeaderMessage: &header,
			CurrentServiceId: sConfigParams.SHIPMENTSBSID,
		}, TESTMODE); soteErr.ErrCode != sError.ErrInvalidJSON {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrInvalidJSON, soteErr.FmtErrMsg)
		}
	})
}

func TestPrepareReqMessage(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		header            = map[string][]string{
			"organizations-id": {fmt.Sprint(TESTORGANIZATIONID)},
			"aws-user-name":    {TESTAWSNAME},
			"json-web-token":   {fmt.Sprintf("%v %v", BEARERSCHEMA, TESTTOKEN)},
			"role-list":        {"ROLE-1", "ROLE-2"},
			"message-id":       {"test-message"},
			"reply-subject":    {sConfigParams.SHIPMENTSBSID, "test-message"},
		}
	)

	tPtr.Run("valid", func(tPtr *testing.T) {
		var (
			testMessage = testRequestMessage{}
			reqMessage  = fmt.Sprintf("{\"id\": 1,\"name\": %q,\"age\": 1}", TESTAWSNAME)
		)
		if _, soteErr = PrepareReqMessage(parentCtx, &RequestParams{
			RequestMsg: []byte(reqMessage),
			Headers:    header,
			TestMode:   TESTMODE,
		}, &testMessage); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("invalid header", func(tPtr *testing.T) {
		var (
			testMessage = testRequestMessage{}
			tHeader     = make(map[string][]string, 0)
		)

		maps.Copy(tHeader, header)
		tHeader["organizations-id"] = []string{TESTAWSNAME}
		if _, soteErr = PrepareReqMessage(parentCtx, &RequestParams{
			Headers:  tHeader,
			TestMode: TESTMODE,
		}, &testMessage); soteErr.ErrCode != sError.ErrNotNumeric {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrNotNumeric, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("missing header", func(tPtr *testing.T) {
		var (
			testMessage = testRequestMessage{}
			tHeader     = make(map[string][]string, 0)
		)

		maps.Copy(tHeader, header)
		delete(tHeader, "organizations-id")
		if _, soteErr = PrepareReqMessage(parentCtx, &RequestParams{
			Headers:  tHeader,
			TestMode: TESTMODE,
		}, &testMessage); soteErr.ErrCode != sError.ErrInvalidMsgSignature {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrInvalidMsgSignature, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("invalid JSON", func(tPtr *testing.T) {
		var (
			testMessage = testRequestMessage{}
			reqMessage  = fmt.Sprintf("{\"id\": %q,\"name\": %q,\"age\": 1}", "age", TESTAWSNAME)
		)

		if _, soteErr = PrepareReqMessage(parentCtx, &RequestParams{
			RequestMsg: []byte(reqMessage),
			Headers:    header,
			TestMode:   TESTMODE,
		}, &testMessage); soteErr.ErrCode != sError.ErrInvalidJSON {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrInvalidJSON, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("missing body parameter", func(tPtr *testing.T) {
		var (
			testMessage = testRequestMessage{}
			reqMessage  = fmt.Sprintf("{\"id\": 1,\"name\": %q}", TESTAWSNAME)
		)

		if _, soteErr = PrepareReqMessage(parentCtx, &RequestParams{
			RequestMsg: []byte(reqMessage),
			Headers:    header,
			TestMode:   TESTMODE,
		}, &testMessage); soteErr.ErrCode != sError.ErrInvalidMsgSignature {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrInvalidMsgSignature, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("failed Token validation", func(tPtr *testing.T) {
		var (
			testMessage = testRequestMessage{}
			reqMessage  = fmt.Sprintf("{\"id\": 1,\"name\": %q,\"age\": 1}", TESTAWSNAME)
		)

		if _, soteErr = PrepareReqMessage(parentCtx, &RequestParams{
			RequestMsg:   []byte(reqMessage),
			Headers:      header,
			TestMode:     false,
			JSONWebToken: TESTTOKEN,
			AuthConfig: &sAuthentication.Config{
				AppEnvironment: sConfigParams.DEVELOPMENT,
				AwsRegion:      awsRegion,
				UserPoolId:     userPoolId,
			},
		}, &testMessage); soteErr.ErrCode != sError.ErrMissingTokenSegments {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrMissingTokenSegments, soteErr.FmtErrMsg)
		}
	})
}

func TestRequestHeaderParams_PrepareReqHeader(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		header            = map[string][]string{
			"organizations-id": {fmt.Sprint(TESTORGANIZATIONID)},
			"aws-user-name":    {TESTAWSNAME},
			"json-web-token":   {fmt.Sprintf("%v %v", BEARERSCHEMA, TESTTOKEN)},
			"role-list":        {"ROLE-1", "ROLE-2"},
			"message-id":       {"test-message"},
			"reply-subject":    {sConfigParams.SHIPMENTSBSID, "test-message"},
			"content-type":     {"application/json"},
			"origin":           {LOCALHOST},
		}
	)

	tPtr.Run("valid", func(tPtr *testing.T) {
		var (
			reqHeader = &RequestHeaderParams{}
		)

		if soteErr = reqHeader.PrepareReqHeader(header, TESTMODE); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("default origin", func(tPtr *testing.T) {
		var (
			reqHeader = &RequestHeaderParams{}
			tHeader   = make(map[string][]string, 0)
		)

		maps.Copy(tHeader, header)
		delete(tHeader, "origin")
		if soteErr = reqHeader.PrepareReqHeader(tHeader, TESTMODE); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("invalid header", func(tPtr *testing.T) {
		var (
			reqHeader = &RequestHeaderParams{}
			tHeader   = make(map[string][]string, 0)
		)

		maps.Copy(tHeader, header)
		tHeader["organizations-id"] = []string{TESTAWSNAME}
		if soteErr = reqHeader.PrepareReqHeader(tHeader, TESTMODE); soteErr.ErrCode != sError.ErrNotNumeric {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrNotNumeric, soteErr.FmtErrMsg)
		}
	})
}

func TestPrepareMessage(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		header            = map[string][]string{
			"organizations-id": {fmt.Sprint(TESTORGANIZATIONID)},
			"aws-user-name":    {TESTAWSNAME},
			"json-web-token":   {fmt.Sprintf("%v %v", BEARERSCHEMA, TESTTOKEN)},
			"role-list":        {"ROLE-1", "ROLE-2"},
			"message-id":       {"test-message"},
			"reply-subject":    {sConfigParams.SHIPMENTSBSID, "test-message"},
		}
	)

	tPtr.Run("valid", func(tPtr *testing.T) {
		var (
			testMessage = testRequestMessage{}
			reqMessage  = fmt.Sprintf("{\"id\": 1,\"name\": %q,\"age\": 1}", TESTAWSNAME)
		)
		if soteErr = PrepareMessage(parentCtx, &RequestParams{
			RequestMsg: []byte(reqMessage),
			Headers:    header,
			TestMode:   TESTMODE,
		}, &testMessage); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("invalid JSON", func(tPtr *testing.T) {
		var (
			testMessage = testRequestMessage{}
			reqMessage  = fmt.Sprintf("{\"id\": %q,\"name\": %q,\"age\": 1}", "age", TESTAWSNAME)
		)

		if soteErr = PrepareMessage(parentCtx, &RequestParams{
			RequestMsg: []byte(reqMessage),
			Headers:    header,
			TestMode:   TESTMODE,
		}, &testMessage); soteErr.ErrCode != sError.ErrInvalidJSON {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrInvalidJSON, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("missing body parameter", func(tPtr *testing.T) {
		var (
			testMessage = testRequestMessage{}
			reqMessage  = fmt.Sprintf("{\"id\": 1,\"name\": %q}", TESTAWSNAME)
		)

		if soteErr = PrepareMessage(parentCtx, &RequestParams{
			RequestMsg: []byte(reqMessage),
			Headers:    header,
			TestMode:   TESTMODE,
		}, &testMessage); soteErr.ErrCode != sError.ErrInvalidMsgSignature {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrInvalidMsgSignature, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("failed Token validation", func(tPtr *testing.T) {
		var (
			testMessage = testRequestMessage{}
			reqMessage  = fmt.Sprintf("{\"id\": 1,\"name\": %q,\"age\": 1}", TESTAWSNAME)
		)

		if soteErr = PrepareMessage(parentCtx, &RequestParams{
			RequestMsg:   []byte(reqMessage),
			Headers:      header,
			TestMode:     false,
			JSONWebToken: TESTTOKEN,
			AuthConfig: &sAuthentication.Config{
				AppEnvironment: sConfigParams.DEVELOPMENT,
				AwsRegion:      awsRegion,
				UserPoolId:     userPoolId,
			},
		}, &testMessage); soteErr.ErrCode != sError.ErrMissingTokenSegments {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrMissingTokenSegments, soteErr.FmtErrMsg)
		}
	})
}

func TestReadRequest(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	tPtr.Run("valid", func(tPtr *testing.T) {
		var (
			reqMessage = fmt.Sprintf("{\"id\": 1,\"name\": %q,\"age\": 1}", TESTAWSNAME)
			httpReqPtr = &httpReqTest{
				method: http.MethodPost,
				path:   "/",
				body:   reqMessage,
			}
		)

		_, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, CORSMiddleware(sConfigParams.DEVELOPMENT))
		httpResp := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(httpResp)
		ctx.Request = reqPtr
		if _, soteErr = ReadRequest(ctx); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("invalid", func(tPtr *testing.T) {
		var (
			reqMessage = "XXXXXaGVsbG8="
			httpReqPtr = &httpReqTest{
				method: http.MethodDelete,
				path:   "/",
				params: map[string]string{"params": reqMessage},
			}
		)

		_, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, CORSMiddleware(sConfigParams.DEVELOPMENT))

		httpResp := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(httpResp)
		ctx.Request = reqPtr
		if _, soteErr = ReadRequest(ctx); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})
}

func TestValidateRequestMessage(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	tPtr.Run("valid", func(tPtr *testing.T) {
		var (
			reqMessage = testRequestMessage{
				Id:   1,
				Name: TESTAWSNAME,
				Age:  1,
			}
		)

		if soteErr = ValidateRequestMessage(parentCtx, &reqMessage, TESTMODE); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("invalid parameters", func(tPtr *testing.T) {
		if soteErr = ValidateRequestMessage(parentCtx, "test", TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("missing parameters", func(tPtr *testing.T) {
		var (
			reqMessage = testRequestMessage{
				Id:   1,
				Name: TESTAWSNAME,
			}
		)

		if soteErr = ValidateRequestMessage(parentCtx, &reqMessage, TESTMODE); soteErr.ErrCode != sError.ErrInvalidMsgSignature {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrInvalidMsgSignature, soteErr.FmtErrMsg)
		}
	})
}

func TestBindMultipart(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	type testStruct struct {
		Id   int    `form:"id"`
		Name string `form:"name"`
		Age  int    `form:"age"`
	}
	tPtr.Run("valid", func(tPtr *testing.T) {
		var (
			reqMessage = fmt.Sprintf("{\"id\": 1,\"name\": %q,\"age\": 1}", "name")
			resp       = &testStruct{}
		)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormField("message")
		_, _ = part.Write([]byte(reqMessage))
		_ = writer.Close()
		reqPtr := httptest.NewRequest(http.MethodPost, "/test", body)
		reqPtr.Header.Add("Content-Type", writer.FormDataContentType())
		httpResp := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(httpResp)
		ctx.Request = reqPtr
		if soteErr = BindMultipart(ctx, resp); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("invalid", func(tPtr *testing.T) {
		var (
			reqMessage = fmt.Sprintf("{\"id\": 1,\"name\":1,\"age\":  %q}", TESTAWSNAME)
			httpReqPtr = &httpReqTest{
				method: http.MethodDelete,
				path:   "/",
				params: map[string]string{"params": reqMessage},
			}
			resp = &testStruct{}
		)

		_, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, CORSMiddleware(sConfigParams.DEVELOPMENT))
		reqPtr.Header = make(map[string][]string, 0)
		reqPtr.Header.Add("Content-Type", "multipart/form-data")
		httpResp := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(httpResp)
		ctx.Request = reqPtr
		if soteErr = BindMultipart(ctx, resp); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

}
