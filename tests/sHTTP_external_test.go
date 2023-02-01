package tests_test

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
	"gitlab.com/soteapps/packages/v2023/sHTTP"
	"golang.org/x/exp/maps"
)

const (
	TESTMODE           = true
	TESTTOKEN          = "TOKEN"
	TESTAWSNAME        = "DELETE_ME_TEST"
	TESTORGANIZATIONID = 1
	LOCALHOST          = "http://localhost:1234"
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

type testRequestMessage struct {
	Id   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required"`
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

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, sHTTP.AuthenticationMiddleware(&sAuthentication.Config{
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

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, sHTTP.AuthenticationMiddleware(&sAuthentication.Config{
			AppEnvironment: sConfigParams.DEVELOPMENT,
			AwsRegion:      awsRegion,
			UserPoolId:     userPoolId,
		}, false))
		httpResp := httptest.NewRecorder()
		routerPtr.ServeHTTP(httpResp, reqPtr)
		if httpResp.Code != http.StatusInternalServerError {
			tPtr.Errorf("%v Failed: Expected response error to be %v got  %v", testName, http.StatusInternalServerError, httpResp.Code)
		}

		res := &sHTTP.Response{}
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

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, sHTTP.AuthenticationMiddleware(&sAuthentication.Config{
			AppEnvironment: sConfigParams.DEVELOPMENT,
			AwsRegion:      awsRegion,
			UserPoolId:     userPoolId,
		}, false))
		reqPtr.Header = make(map[string][]string, 0)
		reqPtr.Header.Add("Authorization", fmt.Sprintf("%v %v", sHTTP.BEARERSCHEMA, TESTTOKEN))
		httpResp := httptest.NewRecorder()
		routerPtr.ServeHTTP(httpResp, reqPtr)
		if httpResp.Code != http.StatusInternalServerError {
			tPtr.Errorf("%v Failed: Expected response error to be %v got  %v", testName, http.StatusInternalServerError, httpResp.Code)
		}

		res := &sHTTP.Response{}
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

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, sHTTP.CORSMiddleware(sConfigParams.DEVELOPMENT))
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

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, sHTTP.CORSMiddleware(env))
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

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, sHTTP.CORSMiddleware(sConfigParams.PRODUCTION))
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

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, sHTTP.CORSMiddleware(env))
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

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, sHTTP.CORSMiddleware(sConfigParams.DEVELOPMENT))
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

		routerPtr, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, sHTTP.CORSMiddleware(sConfigParams.DEVELOPMENT))
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
		header            = sHTTP.RequestHeaderParams{
			OrganizationsId: TESTORGANIZATIONID,
			AWSUsername:     TESTAWSNAME,
			MessageId:       "test-message",
			JSONWebToken:    fmt.Sprintf("%v %v", sHTTP.BEARERSCHEMA, TESTTOKEN),
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

			clientPtr = &sHTTP.ClientPool{Pool: pool}
			poolPtr   = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Delete(&sHTTP.Request{}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("valid mock", func(tPtr *testing.T) {
		var (
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
			reqMessage = testRequestMessage{
				Id:   1,
				Name: TESTAWSNAME,
				Age:  1,
			}
			req = &sHTTP.Request{}
		)

		if req, soteErr = sHTTP.PrepareLeafRequest(&sHTTP.LeafReqParams{
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
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Delete(&sHTTP.Request{URL: fmt.Sprintf("{%v}", LOCALHOST)}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}

	})

	tPtr.Run("invalid response", func(tPtr *testing.T) {
		var (
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, errors.New("mock error")
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Delete(&sHTTP.Request{},
			TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("authentication failure", func(tPtr *testing.T) {
		var (
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusUnauthorized,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Delete(&sHTTP.Request{}, TESTMODE); soteErr.ErrCode != sError.ErrStatusUnauthorized {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrStatusUnauthorized, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("other failures", func(tPtr *testing.T) {
		var (
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Delete(&sHTTP.Request{}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})
}

func TestConnInfo_Get(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		header            = sHTTP.RequestHeaderParams{
			OrganizationsId: TESTORGANIZATIONID,
			AWSUsername:     TESTAWSNAME,
			MessageId:       "test-message",
			JSONWebToken:    fmt.Sprintf("%v %v", sHTTP.BEARERSCHEMA, TESTTOKEN),
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

			clientPtr = &sHTTP.ClientPool{Pool: pool}
			poolPtr   = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Get(&sHTTP.Request{}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("valid mock", func(tPtr *testing.T) {
		var (
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
			reqMessage = testRequestMessage{
				Id:   1,
				Name: TESTAWSNAME,
				Age:  1,
			}
			req = &sHTTP.Request{}
		)

		if req, soteErr = sHTTP.PrepareLeafRequest(&sHTTP.LeafReqParams{
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
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Get(&sHTTP.Request{URL: fmt.Sprintf("{%v}", LOCALHOST)}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}

	})

	tPtr.Run("invalid response", func(tPtr *testing.T) {
		var (
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, errors.New("mock error")
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Get(&sHTTP.Request{},
			TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("authentication failure", func(tPtr *testing.T) {
		var (
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusUnauthorized,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Get(&sHTTP.Request{}, TESTMODE); soteErr.ErrCode != sError.ErrStatusUnauthorized {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrStatusUnauthorized, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("other failures", func(tPtr *testing.T) {
		var (
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Get(&sHTTP.Request{}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})
}

func TestConnInfo_Patch(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		header            = sHTTP.RequestHeaderParams{
			OrganizationsId: TESTORGANIZATIONID,
			AWSUsername:     TESTAWSNAME,
			MessageId:       "test-message",
			JSONWebToken:    fmt.Sprintf("%v %v", sHTTP.BEARERSCHEMA, TESTTOKEN),
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

			clientPtr = &sHTTP.ClientPool{Pool: pool}
			poolPtr   = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Patch(&sHTTP.Request{}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("valid mock", func(tPtr *testing.T) {
		var (
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
			reqMessage = testRequestMessage{
				Id:   1,
				Name: TESTAWSNAME,
				Age:  1,
			}
			req = &sHTTP.Request{}
		)

		if req, soteErr = sHTTP.PrepareLeafRequest(&sHTTP.LeafReqParams{
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
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Patch(&sHTTP.Request{URL: fmt.Sprintf("{%v}", LOCALHOST)}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}

	})

	tPtr.Run("invalid response", func(tPtr *testing.T) {
		var (
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, errors.New("mock error")
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Patch(&sHTTP.Request{},
			TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("authentication failure", func(tPtr *testing.T) {
		var (
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusUnauthorized,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Patch(&sHTTP.Request{}, TESTMODE); soteErr.ErrCode != sError.ErrStatusUnauthorized {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrStatusUnauthorized, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("other failures", func(tPtr *testing.T) {
		var (
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Patch(&sHTTP.Request{}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})
}

func TestConnInfo_Post(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		header            = sHTTP.RequestHeaderParams{
			OrganizationsId: TESTORGANIZATIONID,
			AWSUsername:     TESTAWSNAME,
			MessageId:       "test-message",
			JSONWebToken:    fmt.Sprintf("%v %v", sHTTP.BEARERSCHEMA, TESTTOKEN),
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

			clientPtr = &sHTTP.ClientPool{Pool: pool}
			poolPtr   = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Post(&sHTTP.Request{}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("valid mock", func(tPtr *testing.T) {
		var (
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
			reqMessage = testRequestMessage{
				Id:   1,
				Name: TESTAWSNAME,
				Age:  1,
			}
			req = &sHTTP.Request{}
		)

		if req, soteErr = sHTTP.PrepareLeafRequest(&sHTTP.LeafReqParams{
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
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Post(&sHTTP.Request{URL: fmt.Sprintf("{%v}", LOCALHOST)}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}

	})

	tPtr.Run("invalid response", func(tPtr *testing.T) {
		var (
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, errors.New("mock error")
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Post(&sHTTP.Request{},
			TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrGenericError, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("authentication failure", func(tPtr *testing.T) {
		var (
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusUnauthorized,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Post(&sHTTP.Request{}, TESTMODE); soteErr.ErrCode != sError.ErrStatusUnauthorized {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrStatusUnauthorized, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("other failures", func(tPtr *testing.T) {
		var (
			clientPtr sHTTP.Client = &ClientMock{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       ioutil.NopCloser(bytes.NewReader(nil)),
					}, nil
				},
			}
			poolPtr = &sHTTP.ClientPoolImpl{
				Client: clientPtr,
			}
		)

		if _, soteErr = poolPtr.Post(&sHTTP.Request{}, TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
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
		if statusCode := sHTTP.ConvertError(sError.GetSError(sError.ErrInvalidMsgSignature, sError.BuildParams([]string{""}),
			sError.EmptyMap)); statusCode != http.StatusBadRequest {
			tPtr.Errorf("%v Failed: Expected HTTP status code to be %v got %v", testName, http.StatusBadRequest, statusCode)
		}
	})

	tPtr.Run("sError.ErrInvalidJSON", func(tPtr *testing.T) {
		if statusCode := sHTTP.ConvertError(sError.GetSError(sError.ErrInvalidJSON, sError.BuildParams([]string{""}),
			sError.EmptyMap)); statusCode != http.StatusUnprocessableEntity {
			tPtr.Errorf("%v Failed: Expected HTTP status code to be %v got %v", testName, http.StatusUnprocessableEntity, statusCode)
		}
	})

	tPtr.Run("sError.ErrItemNotFound", func(tPtr *testing.T) {
		if statusCode := sHTTP.ConvertError(sError.GetSError(sError.ErrItemNotFound, sError.BuildParams([]string{""}),
			sError.EmptyMap)); statusCode != http.StatusNotFound {
			tPtr.Errorf("%v Failed: Expected HTTP status code to be %v got %v", testName, http.StatusNotFound, statusCode)
		}
	})
}

func TestPrepareLeafRequest(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
		header            = sHTTP.RequestHeaderParams{
			OrganizationsId: TESTORGANIZATIONID,
			AWSUsername:     TESTAWSNAME,
			MessageId:       "test-message",
			JSONWebToken:    fmt.Sprintf("%v %v", sHTTP.BEARERSCHEMA, TESTTOKEN),
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
		if _, soteErr = sHTTP.PrepareLeafRequest(&sHTTP.LeafReqParams{
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

		if _, soteErr = sHTTP.PrepareLeafRequest(&sHTTP.LeafReqParams{
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
			"json-web-token":   {fmt.Sprintf("%v %v", sHTTP.BEARERSCHEMA, TESTTOKEN)},
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
		if _, soteErr = sHTTP.PrepareReqMessage(parentCtx, &sHTTP.RequestParams{
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
		if _, soteErr = sHTTP.PrepareReqMessage(parentCtx, &sHTTP.RequestParams{
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
		if _, soteErr = sHTTP.PrepareReqMessage(parentCtx, &sHTTP.RequestParams{
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

		if _, soteErr = sHTTP.PrepareReqMessage(parentCtx, &sHTTP.RequestParams{
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

		if _, soteErr = sHTTP.PrepareReqMessage(parentCtx, &sHTTP.RequestParams{
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

		if _, soteErr = sHTTP.PrepareReqMessage(parentCtx, &sHTTP.RequestParams{
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
			"json-web-token":   {fmt.Sprintf("%v %v", sHTTP.BEARERSCHEMA, TESTTOKEN)},
			"role-list":        {"ROLE-1", "ROLE-2"},
			"message-id":       {"test-message"},
			"reply-subject":    {sConfigParams.SHIPMENTSBSID, "test-message"},
			"content-type":     {"application/json"},
			"origin":           {LOCALHOST},
		}
	)

	tPtr.Run("valid", func(tPtr *testing.T) {
		var (
			reqHeader = &sHTTP.RequestHeaderParams{}
		)

		if soteErr = reqHeader.PrepareReqHeader(header, TESTMODE); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("default origin", func(tPtr *testing.T) {
		var (
			reqHeader = &sHTTP.RequestHeaderParams{}
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
			reqHeader = &sHTTP.RequestHeaderParams{}
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
			"json-web-token":   {fmt.Sprintf("%v %v", sHTTP.BEARERSCHEMA, TESTTOKEN)},
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
		if soteErr = sHTTP.PrepareMessage(parentCtx, &sHTTP.RequestParams{
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

		if soteErr = sHTTP.PrepareMessage(parentCtx, &sHTTP.RequestParams{
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

		if soteErr = sHTTP.PrepareMessage(parentCtx, &sHTTP.RequestParams{
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

		if soteErr = sHTTP.PrepareMessage(parentCtx, &sHTTP.RequestParams{
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

		_, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, sHTTP.CORSMiddleware(sConfigParams.DEVELOPMENT))
		httpResp := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(httpResp)
		ctx.Request = reqPtr
		if _, soteErr = sHTTP.ReadRequest(ctx); soteErr.ErrCode != nil {
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

		_, reqPtr := setRequestTestHandler(tPtr, httpReqPtr, sHTTP.CORSMiddleware(sConfigParams.DEVELOPMENT))

		httpResp := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(httpResp)
		ctx.Request = reqPtr
		if _, soteErr = sHTTP.ReadRequest(ctx); soteErr.ErrCode != sError.ErrGenericError {
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

		if soteErr = sHTTP.ValidateRequestMessage(parentCtx, &reqMessage, TESTMODE); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("invalid parameters", func(tPtr *testing.T) {
		if soteErr = sHTTP.ValidateRequestMessage(parentCtx, "test", TESTMODE); soteErr.ErrCode != sError.ErrGenericError {
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

		if soteErr = sHTTP.ValidateRequestMessage(parentCtx, &reqMessage, TESTMODE); soteErr.ErrCode != sError.ErrInvalidMsgSignature {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrInvalidMsgSignature, soteErr.FmtErrMsg)
		}
	})
}

func TestProcessLeafResponse(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		testMessage       = testRequestMessage{}
		soteErr           sError.SoteError
		response          []byte
	)

	tPtr.Run("valid", func(tPtr *testing.T) {
		reqMessage := sHTTP.Response{
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

		if soteErr = sHTTP.ProcessLeafResponse(parentCtx, response, &testMessage, TESTMODE); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("invalid response message", func(tPtr *testing.T) {
		if soteErr = sHTTP.ProcessLeafResponse(parentCtx, []byte("test-message"), &testMessage, TESTMODE); soteErr.ErrCode != sError.ErrInvalidJSON {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrInvalidJSON, soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("invalid sote error", func(tPtr *testing.T) {
		reqMessage := sHTTP.Response{
			MessageId: "",
			Error:     TESTAWSNAME,
		}
		if response, soteErr = sCustom.JSONMarshal(reqMessage); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}

		if soteErr = sHTTP.ProcessLeafResponse(parentCtx, response, &testMessage, TESTMODE); soteErr.ErrCode != sError.ErrInvalidJSON {
			tPtr.Errorf("%v Failed: Expected return to be %v got %v", testName, sError.ErrInvalidJSON, soteErr.FmtErrMsg)
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
