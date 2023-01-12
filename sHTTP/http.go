package sHTTP

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/soteapps/packages/v2023/sAuthentication"
	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
)

const (
	BEARERSCHEMA = "Bearer "
)

// RequestHeaderParams is a Request JSON Header payload, holding authorization and authentication metadata.
type RequestHeaderParams struct {
	OrganizationsId   int      `json:"organizations-id" validate:"required"`
	AWSUsername       string   `json:"aws-user-name" validate:"required"`
	MessageId         string   `json:"message-id" validate:"required"`
	JSONWebToken      string   `json:"json-web-token" validate:"required"`
	RoleList          []string `json:"role-list" validate:"required"`
	ContentType       []string `json:"Content-Type"`
	Origin            string   `json:"origin"`
	ReplySubjectNode  string
	BusinessServiceId string
}

type Request struct {
	URL         string
	Headers     map[string][]string
	QueryParams map[string]string
	BodyParams  []byte
}

type ConnInfo struct {
	Pool *sync.Pool
}

// AuthenticationMiddleware runs JWT authentication check
func AuthenticationMiddleware(authConfig *sAuthentication.Config, testMode bool) gin.HandlerFunc {
	sLogger.DebugMethod()

	return func(ctx *gin.Context) {
		if !testMode {
			var (
				soteErr   sError.SoteError
				bearerLen = len(BEARERSCHEMA)
			)

			if token := ctx.GetHeader("Authorization"); len(token) > bearerLen {
				token = token[bearerLen:]
				ctx.Request.Header.Set("json-web-token", token)
				soteErr = sAuthentication.ValidToken(ctx, token, authConfig)
			} else {
				sLogger.Info("Key: 'RequestHeaderParams.JSONWebToken' Error:Field validation for 'JSONWebToken' failed on the 'required' tag")
				soteErr = sError.ConvertError(sError.GetSError(sError.ErrInvalidMsgSignature,
					sError.BuildParams([]string{"See Configuration worksheet"}), sError.EmptyMap), testMode)
			}

			if soteErr.ErrCode != nil {
				ctx.AbortWithStatusJSON(ConvertError(soteErr), Response{
					MessageId: ctx.Request.Header.Get("message-id"),
					Error:     soteErr,
				})
			}
		}
	}
}

// CORSMiddleware allow cross-origin reference from soteapps.com or localhost
func CORSMiddleware(targetEnvironment string) gin.HandlerFunc {
	sLogger.DebugMethod()

	return func(ctx *gin.Context) {
		origin := GetAllowedOrigins(ctx, targetEnvironment)
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		// ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET, PATCH, DELETE")
		ctx.Writer.Header().Set("Access-Control-Max-Age", fmt.Sprint(12*time.Hour))
		ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		ctx.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

// GetAllowedOrigins get all allowed domain origin
func GetAllowedOrigins(ctx *gin.Context, targetEnvironment string) (origin string) {
	sLogger.DebugMethod()

	var pattern string
	origin = strings.ToLower(ctx.GetHeader("Origin"))
	sLogger.Info("Checking CORS: " + origin)
	switch targetEnvironment {
	case "development", "staging":
		pattern = `(?i)((^|^[^:]+:\/\/|[^\.]+\.)localhost((:[0-9]{1,4})?(\/[-\w]*\/?)?)$)|(^https:\/\/([a-z0-9]+([a-z0-9-]{1,61}[a-z0-9])?\.)*?staging\.(soteapps|sote)\.com$)`
	case "production":
		pattern = `(?i)^https:\/\/([a-z0-9]+([a-z0-9-]{1,61}[a-z0-9])?\.){0,1}soteapps|sote\.com$`
	default:
		pattern = fmt.Sprintf(`(?i)^https:\/\/([a-z0-9]+([a-z0-9-]{1,61}[a-z0-9])?\.)*?%v\.soteapps|sote\.com$`, targetEnvironment)
	}

	if len(regexp.MustCompile(pattern).FindStringIndex(origin)) == 0 {
		sLogger.Info("CORS failed: " + origin)
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	return
}

// Delete makes a DELETE request and return the response body
func (ConnInfo *ConnInfo) Delete(request *Request, testMode bool) (body []byte, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		req *http.Request
		err error
	)

	if request != nil {
		if req, err = http.NewRequest("DELETE", request.URL, nil); err != nil { // create a new DELETE request
			sLogger.Info(err.Error())
			soteErr = sError.GetSError(sError.ErrGenericError, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
			return
		}

		SetReqHeaders(req, request.Headers)   // set headers
		SetReqQuery(req, request.QueryParams) // set query parameters

		// execute the request and read the entire response body
		body, soteErr = ConnInfo.DoRequest(req, testMode)
	}

	return
}

// Get makes a GET request and return the response body
func (ConnInfo *ConnInfo) Get(request *Request, testMode bool) (body []byte, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		req *http.Request
		err error
	)

	if request != nil {
		if req, err = http.NewRequest("GET", request.URL, nil); err != nil { // create a new GET request
			sLogger.Info(err.Error())
			soteErr = sError.GetSError(sError.ErrGenericError, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)

			return
		}

		SetReqHeaders(req, request.Headers)   // set headers
		SetReqQuery(req, request.QueryParams) // set query parameters

		body, soteErr = ConnInfo.DoRequest(req, testMode) // execute the request and read the entire response body
	}

	return
}

// Patch makes a PATCH request and return the response body
func (ConnInfo *ConnInfo) Patch(request *Request, testMode bool) (body []byte, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		req       *http.Request
		err       error
		tBodyByte []byte
		tBody     io.Reader
	)

	if request != nil {
		if len(request.BodyParams) > 0 {
			tBody = bytes.NewReader(tBodyByte)
		}

		if req, err = http.NewRequest("PATCH", request.URL, tBody); err != nil { // create a new PATCH request
			sLogger.Info(err.Error())
			soteErr = sError.GetSError(sError.ErrGenericError, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)

			return
		}

		SetReqHeaders(req, request.Headers)               // set headers
		body, soteErr = ConnInfo.DoRequest(req, testMode) // execute the request and read the entire response body
	}

	return
}

// Post makes a POST request and return the response body
func (ConnInfo *ConnInfo) Post(request *Request, testMode bool) (body []byte, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		req       *http.Request
		err       error
		tBodyByte []byte
		tBody     io.Reader
	)

	if request != nil {
		if len(request.BodyParams) > 0 {
			tBody = bytes.NewReader(tBodyByte)
		}

		if req, err = http.NewRequest("POST", request.URL, tBody); err != nil { // create a new POST request
			sLogger.Info(err.Error())
			soteErr = sError.GetSError(sError.ErrGenericError, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)

			return
		}

		SetReqHeaders(req, request.Headers)               // set headers
		body, soteErr = ConnInfo.DoRequest(req, testMode) // execute the request and read the entire response body
	}

	return
}

// DoRequest makes an HTTP request, read the entire response body, and return the response
func (ConnInfo *ConnInfo) DoRequest(req *http.Request, testMode bool) (body []byte, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		client *http.Client
		resp   *http.Response
		err    error
	)

	client = ConnInfo.Pool.Get().(*http.Client) // get a client from the pool
	if resp, err = client.Do(req); err == nil { // execute the request and get the response
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body) // read the entire response body

		body, err = ioutil.ReadAll(resp.Body)
		if resp.StatusCode >= http.StatusBadRequest && len(body) == 0 {
			sLogger.Info(fmt.Sprintf("External API %v Call Error", req.URL.RawQuery))
			switch resp.StatusCode {
			case http.StatusUnauthorized:
				soteErr = sError.GetSError(sError.ErrStatusUnauthorized, sError.BuildParams([]string{}), sError.EmptyMap)
			default:
				soteErr = sError.GetSError(sError.ErrGenericError, sError.BuildParams([]string{fmt.Sprintf("Error %v,Message:%v", resp.StatusCode,
					string(body))}),
					sError.EmptyMap)
			}

			return
		}
	}

	// put the client back into the pool
	ConnInfo.Pool.Put(client)
	if err != nil {
		soteErr = sError.GetSError(sError.ErrGenericError, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
	}

	return
}

// SetReqQuery sets request query parameters
func SetReqQuery(req *http.Request, params map[string]string) {
	sLogger.DebugMethod()
	if params != nil {
		var (
			q = req.URL.Query()
		)

		for k, v := range params {
			q.Add(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}
}

// SetReqHeaders sets request headers
func SetReqHeaders(req *http.Request, headers map[string][]string) {
	// set headers
	for k, v := range headers {
		req.Header.Del(k)
		for i := range v {
			req.Header.Add(k, v[i])
		}
	}
}

func ConvertError(soteErr sError.SoteError) (statusCode int) {
	sLogger.DebugMethod()

	sLogger.Info(soteErr.FmtErrMsg)
	// convert Sote Error to HTTP Errors
	switch soteErr.ErrCode {
	case 206200:
		statusCode = http.StatusBadRequest
	case 207110:
		statusCode = http.StatusUnprocessableEntity
	case sError.ErrItemNotFound:
		statusCode = http.StatusNotFound
	default:
		statusCode = http.StatusInternalServerError
	}

	return
}
