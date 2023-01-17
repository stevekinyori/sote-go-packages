package sHTTP

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.com/soteapps/packages/v2023/sAuthentication"
	"gitlab.com/soteapps/packages/v2023/sConfigParams"
	"gitlab.com/soteapps/packages/v2023/sCustom"
	"gitlab.com/soteapps/packages/v2023/sError"
)

const (
	TESTMODE  = true
	TESTTOKEN = "TOKEN"
	LOCALHOST = "http://localhost:1234"
)

var (
	parentCtx     = context.Background()
	awsRegion, _  = sConfigParams.GetRegion(parentCtx)
	userPoolId, _ = sConfigParams.GetUserPoolId(parentCtx, sConfigParams.DEVELOPMENT)
)

type httpReqTest struct {
	method string
	path   string
	body   string
	params map[string]string
	header map[string][]string
}

type errReader int

// MockDo makes an HTTP Do request using mocked settings
type ClientMock struct {
	MockDo func(*http.Request) (*http.Response, error)
}

func (clientPtr *ClientMock) Do(req *http.Request) (*http.Response, error) {
	return clientPtr.MockDo(req)
}

func (errReader) Read([]byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestAuthenticationMiddleware(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	tPtr.Run("testMode authentication", func(tPtr *testing.T) {
		httpReqPtr := &httpReqTest{
			method: http.MethodDelete,
			path:   "/",
		}

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, AuthenticationMiddleware(&sAuthentication.Config{
			AppEnvironment: sConfigParams.DEVELOPMENT,
			AwsRegion:      awsRegion,
			UserPoolId:     userPoolId,
		}, TESTMODE))
		httpResp := httptest.NewRecorder()
		routerPtr.ServeHTTP(httpResp, reqPtr)
		if httpResp.Code != http.StatusOK {
			tPtr.Errorf("%v Failed: Expected response error to be %v got  %v", testName, http.StatusOK, httpResp.Code)
		}
	})

	tPtr.Run("missing token authentication", func(tPtr *testing.T) {
		httpReqPtr := &httpReqTest{
			method: http.MethodDelete,
			path:   "/",
		}

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, AuthenticationMiddleware(&sAuthentication.Config{
			AppEnvironment: sConfigParams.DEVELOPMENT,
			AwsRegion:      awsRegion,
			UserPoolId:     userPoolId,
		}, false))
		httpResp := httptest.NewRecorder()
		routerPtr.ServeHTTP(httpResp, reqPtr)
		if httpResp.Code != http.StatusInternalServerError {
			tPtr.Errorf("%v Failed: Expected response error to be %v got  %v", testName, http.StatusInternalServerError, httpResp.Code)
		}

		res := &Response{}
		sCustom.JSONUnmarshal(parentCtx, httpResp.Body.Bytes(), res)
		if soteError := res.Error.(map[string]interface{}); soteError["ErrCode"].(float64) != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteError["FmtErrMsg"])
		}
	})

	tPtr.Run("invalid token authentication", func(tPtr *testing.T) {
		httpReqPtr := &httpReqTest{
			method: http.MethodDelete,
			path:   "/",
		}

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, AuthenticationMiddleware(&sAuthentication.Config{
			AppEnvironment: sConfigParams.DEVELOPMENT,
			AwsRegion:      awsRegion,
			UserPoolId:     userPoolId,
		}, false))
		reqPtr.Header = make(map[string][]string, 0)
		reqPtr.Header.Add("Authorization", fmt.Sprintf("%v %v", BEARERSCHEMA, TESTTOKEN))
		httpResp := httptest.NewRecorder()
		routerPtr.ServeHTTP(httpResp, reqPtr)
		if httpResp.Code != http.StatusInternalServerError {
			tPtr.Errorf("%v Failed: Expected response error to be %v got  %v", testName, http.StatusInternalServerError, httpResp.Code)
		}

		res := &Response{}
		sCustom.JSONUnmarshal(parentCtx, httpResp.Body.Bytes(), res)
		if soteError := res.Error.(map[string]interface{}); soteError["ErrCode"].(float64) != sError.ErrMissingTokenSegments {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrMissingTokenSegments, soteError["FmtErrMsg"])
		}
	})
}

func TestCORSMiddleware(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	tPtr.Run("valid CORS "+sConfigParams.DEVELOPMENT, func(tPtr *testing.T) {
		httpReqPtr := &httpReqTest{
			method: http.MethodDelete,
			path:   "/",
		}

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, CORSMiddleware(sConfigParams.DEVELOPMENT))
		reqPtr.Header = make(map[string][]string, 0)
		reqPtr.Header.Add("Origin", "http://localhost")
		httpResp := httptest.NewRecorder()
		routerPtr.ServeHTTP(httpResp, reqPtr)
		if httpResp.Code != http.StatusOK {
			tPtr.Errorf("%v Failed: Expected response error to be %v got  %v", testName, http.StatusOK, httpResp.Code)
		}
	})

	tPtr.Run("valid CORS "+sConfigParams.STAGING+" or "+sConfigParams.DEMO, func(tPtr *testing.T) {
		env := sConfigParams.DEMO
		httpReqPtr := &httpReqTest{
			method: http.MethodDelete,
			path:   "/",
		}

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, CORSMiddleware(env))
		reqPtr.Header = make(map[string][]string, 0)
		reqPtr.Header.Add("Origin", fmt.Sprintf("https://test.%v.soteapps.com", env))
		httpResp := httptest.NewRecorder()
		routerPtr.ServeHTTP(httpResp, reqPtr)
		if httpResp.Code != http.StatusOK {
			tPtr.Errorf("%v Failed: Expected response error to be %v got  %v", testName, http.StatusOK, httpResp.Code)
		}
	})

	tPtr.Run("valid CORS "+sConfigParams.PRODUCTION, func(tPtr *testing.T) {
		httpReqPtr := &httpReqTest{
			method: http.MethodDelete,
			path:   "/",
		}

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, CORSMiddleware(sConfigParams.PRODUCTION))
		reqPtr.Header = make(map[string][]string, 0)
		reqPtr.Header.Add("Origin", "https://test.soteapps.com")
		httpResp := httptest.NewRecorder()
		routerPtr.ServeHTTP(httpResp, reqPtr)
		if httpResp.Code != http.StatusOK {
			tPtr.Errorf("%v Failed: Expected response error to be %v got  %v", testName, http.StatusOK, httpResp.Code)
		}
	})

	tPtr.Run("valid CORS any other environment", func(tPtr *testing.T) {
		env := "env"
		httpReqPtr := &httpReqTest{
			method: http.MethodDelete,
			path:   "/",
		}

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, CORSMiddleware(env))
		reqPtr.Header = make(map[string][]string, 0)
		reqPtr.Header.Add("Origin", fmt.Sprintf("https://test.%v.soteapps.com", env))
		httpResp := httptest.NewRecorder()
		routerPtr.ServeHTTP(httpResp, reqPtr)
		if httpResp.Code != http.StatusOK {
			tPtr.Errorf("%v Failed: Expected response error to be %v got  %v", testName, http.StatusOK, httpResp.Code)
		}
	})

	tPtr.Run("valid CORS with http.MethodOptions method", func(tPtr *testing.T) {
		httpReqPtr := &httpReqTest{
			method: http.MethodOptions,
			path:   "/",
		}

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, CORSMiddleware(sConfigParams.DEVELOPMENT))
		reqPtr.Header = make(map[string][]string, 0)
		reqPtr.Header.Add("Origin", "http://localhost")
		httpResp := httptest.NewRecorder()
		routerPtr.ServeHTTP(httpResp, reqPtr)
		if httpResp.Code != http.StatusNoContent {
			tPtr.Errorf("%v Failed: Expected response error to be %v got  %v", testName, http.StatusNoContent, httpResp.Code)
		}
	})

	tPtr.Run("invalid CORS", func(tPtr *testing.T) {
		httpReqPtr := &httpReqTest{
			method: http.MethodDelete,
			path:   "/",
		}

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, CORSMiddleware(sConfigParams.DEVELOPMENT))
		reqPtr.Header = make(map[string][]string, 0)
		httpResp := httptest.NewRecorder()
		routerPtr.ServeHTTP(httpResp, reqPtr)
		if httpResp.Code != http.StatusUnauthorized {
			tPtr.Errorf("%v Failed: Expected response error to be %v got  %v", testName, http.StatusUnauthorized, httpResp.Code)
		}
	})
}

func TestConnInfo_Delete(tPtr *testing.T) {
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

	tPtr.Run("non-mock", func(tPtr *testing.T) {
		var (
			pool = &sync.Pool{
				New: func() interface{} {
					return &http.Client{}
				},
			}

			clientPtr = &ClientPool{Pool: pool}
			poolPtr   = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Delete(&Request{}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("valid mock", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
			reqMessage = testRequestMessage{
				Id:   1,
				Name: TESTAWSNAME,
				Age:  1,
			}
			req = &Request{}
		)

		if req, soteErr = PrepareLeafRequest(&LeafReqParams{
			ReqServiceURL:    LOCALHOST,
			ReqMessage:       &reqMessage,
			ReqHeaderMessage: &header,
			CurrentServiceId: sConfigParams.SHIPMENTSBSID,
		}, TESTMODE); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}

		if _, soteErr = poolPtr.Delete(req, TESTMODE); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}

	})

	tPtr.Run("invalid request", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Delete(&Request{URL: fmt.Sprintf("{%v}", LOCALHOST)}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}

	})

	tPtr.Run("invalid response", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, errors.New("mock error")
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Delete(&Request{},
			TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("authentication failure", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusUnauthorized,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Delete(&Request{}, TESTMODE); soteErr.ErrCode != sError.ErrStatusUnauthorized {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrStatusUnauthorized, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("other failures", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Delete(&Request{}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})
}

func TestConnInfo_Get(tPtr *testing.T) {
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

	tPtr.Run("non-mock", func(tPtr *testing.T) {
		var (
			pool = &sync.Pool{
				New: func() interface{} {
					return &http.Client{}
				},
			}

			clientPtr = &ClientPool{Pool: pool}
			poolPtr   = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Get(&Request{}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("valid mock", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
			reqMessage = testRequestMessage{
				Id:   1,
				Name: TESTAWSNAME,
				Age:  1,
			}
			req = &Request{}
		)

		if req, soteErr = PrepareLeafRequest(&LeafReqParams{
			ReqServiceURL:    LOCALHOST,
			ReqMessage:       &reqMessage,
			ReqHeaderMessage: &header,
			CurrentServiceId: sConfigParams.SHIPMENTSBSID,
		}, TESTMODE); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}

		if _, soteErr = poolPtr.Get(req, TESTMODE); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}

	})

	tPtr.Run("invalid request", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Get(&Request{URL: fmt.Sprintf("{%v}", LOCALHOST)}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}

	})

	tPtr.Run("invalid response", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, errors.New("mock error")
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Get(&Request{},
			TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("authentication failure", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusUnauthorized,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Get(&Request{}, TESTMODE); soteErr.ErrCode != sError.ErrStatusUnauthorized {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrStatusUnauthorized, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("other failures", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Get(&Request{}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})
}

func TestConnInfo_Patch(tPtr *testing.T) {
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

	tPtr.Run("non-mock", func(tPtr *testing.T) {
		var (
			pool = &sync.Pool{
				New: func() interface{} {
					return &http.Client{}
				},
			}

			clientPtr = &ClientPool{Pool: pool}
			poolPtr   = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Patch(&Request{}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("valid mock", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
			reqMessage = testRequestMessage{
				Id:   1,
				Name: TESTAWSNAME,
				Age:  1,
			}
			req = &Request{}
		)

		if req, soteErr = PrepareLeafRequest(&LeafReqParams{
			ReqServiceURL:    LOCALHOST,
			ReqMessage:       &reqMessage,
			ReqHeaderMessage: &header,
			CurrentServiceId: sConfigParams.SHIPMENTSBSID,
		}, TESTMODE); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}

		if _, soteErr = poolPtr.Patch(req, TESTMODE); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}

	})

	tPtr.Run("invalid request", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Patch(&Request{URL: fmt.Sprintf("{%v}", LOCALHOST)}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}

	})

	tPtr.Run("invalid response", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, errors.New("mock error")
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Patch(&Request{},
			TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("authentication failure", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusUnauthorized,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Patch(&Request{}, TESTMODE); soteErr.ErrCode != sError.ErrStatusUnauthorized {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrStatusUnauthorized, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("other failures", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Patch(&Request{}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})
}

func TestConnInfo_Post(tPtr *testing.T) {
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

	tPtr.Run("non-mock", func(tPtr *testing.T) {
		var (
			pool = &sync.Pool{
				New: func() interface{} {
					return &http.Client{}
				},
			}

			clientPtr = &ClientPool{Pool: pool}
			poolPtr   = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Post(&Request{}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("valid mock", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
			reqMessage = testRequestMessage{
				Id:   1,
				Name: TESTAWSNAME,
				Age:  1,
			}
			req = &Request{}
		)

		if req, soteErr = PrepareLeafRequest(&LeafReqParams{
			ReqServiceURL:    LOCALHOST,
			ReqMessage:       &reqMessage,
			ReqHeaderMessage: &header,
			CurrentServiceId: sConfigParams.SHIPMENTSBSID,
		}, TESTMODE); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}

		if _, soteErr = poolPtr.Post(req, TESTMODE); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}

	})

	tPtr.Run("invalid request", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Post(&Request{URL: fmt.Sprintf("{%v}", LOCALHOST)}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}

	})

	tPtr.Run("invalid response", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, errors.New("mock error")
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Post(&Request{},
			TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("authentication failure", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusUnauthorized,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Post(&Request{}, TESTMODE); soteErr.ErrCode != sError.ErrStatusUnauthorized {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrStatusUnauthorized, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("other failures", func(tPtr *testing.T) {
		var (
			clientPtr Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Post(&Request{}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})
}

func TestConvertError(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	tPtr.Run("sError.ErrInvalidMsgSignature", func(tPtr *testing.T) {
		if statusCode := ConvertError(sError.GetSError(sError.ErrInvalidMsgSignature, sError.BuildParams([]string{""}),
			sError.EmptyMap)); statusCode != http.StatusBadRequest {
			tPtr.Errorf("%v Failed: Expected HTTP status code to be %v got %v", testName, http.StatusBadRequest, statusCode)
		}
	})

	tPtr.Run("sError.ErrInvalidJSON", func(tPtr *testing.T) {
		if statusCode := ConvertError(sError.GetSError(sError.ErrInvalidJSON, sError.BuildParams([]string{""}),
			sError.EmptyMap)); statusCode != http.StatusUnprocessableEntity {
			tPtr.Errorf("%v Failed: Expected HTTP status code to be %v got %v", testName, http.StatusUnprocessableEntity, statusCode)
		}
	})

	tPtr.Run("sError.ErrItemNotFound", func(tPtr *testing.T) {
		if statusCode := ConvertError(sError.GetSError(sError.ErrItemNotFound, sError.BuildParams([]string{""}),
			sError.EmptyMap)); statusCode != http.StatusNotFound {
			tPtr.Errorf("%v Failed: Expected HTTP status code to be %v got %v", testName, http.StatusNotFound, statusCode)
		}
	})
}

// setRequestTestHandler start a test server and return the server request options
func setRequestTestHandler(tPtr *testing.T, httpReqPtr *httpReqTest, middleware ...gin.HandlerFunc) (routerPtr *gin.Engine, reqPtr *http.Request) {
	tPtr.Helper()

	var (
		handler gin.HandlerFunc
	)

	routerPtr = gin.Default()
	routerPtr.Use(middleware...)
	handler = func(ctx *gin.Context) {}
	routerPtr.Handle(httpReqPtr.method, httpReqPtr.path, handler)
	reqPtr = httptest.NewRequest(httpReqPtr.method, httpReqPtr.path, io.NopCloser(strings.NewReader(httpReqPtr.body)))
	reqPtr.Header = httpReqPtr.header
	if httpReqPtr.params != nil {
		q := reqPtr.URL.Query()
		for k, v := range httpReqPtr.params {
			q.Add(k, v)
		}

		reqPtr.URL.RawQuery = q.Encode()
	}

	return
}
