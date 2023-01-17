package sHTTP

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gitlab.com/soteapps/packages/v2023/sAuthentication"
	"gitlab.com/soteapps/packages/v2023/sCustom"
	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
)

// LeafReqParams Holds leaf-request parameters
type LeafReqParams struct {
	ReqServiceURL    string
	ReqMessage       interface{}
	ReqHeaderMessage *RequestHeaderParams
	CurrentServiceId string
}

// RequestParams Holds request message parameters
type RequestParams struct {
	RequestMsg   []byte
	Headers      map[string][]string
	JSONWebToken string
	TestMode     bool
	AuthConfig   *sAuthentication.Config
}

// PrepareLeafRequest Prepares request to be sent to leaf service
func PrepareLeafRequest(params *LeafReqParams, testMode bool) (leafManagerPtr *Request, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		newSubject = uuid.New().String()
		buffer     []byte
	)

	if buffer, soteErr = sCustom.JSONMarshal(params.ReqMessage); soteErr.ErrCode != nil {
		soteErr = sError.ConvertError(soteErr, testMode)
		return
	}

	leafManagerPtr = &Request{
		URL: params.ReqServiceURL,
		Headers: map[string][]string{
			"message-id":       {params.CurrentServiceId + "." + newSubject},
			"reply-subject":    {params.CurrentServiceId, newSubject},
			"organizations-id": {fmt.Sprint(params.ReqHeaderMessage.OrganizationsId)},
			"Authorization":    {BEARERSCHEMA + params.ReqHeaderMessage.JSONWebToken},
			"aws-user-name":    {params.ReqHeaderMessage.AWSUsername},
			"role-list":        params.ReqHeaderMessage.RoleList,
			"Content-Type":     params.ReqHeaderMessage.ContentType,
			"origin":           {params.ReqHeaderMessage.Origin},
		},
		QueryParams: map[string]string{
			"params": base64.StdEncoding.EncodeToString(buffer),
		},
		BodyParams: buffer,
	}

	return
}

// PrepareReqMessage validates HTTP request and prepares the request headers
func PrepareReqMessage(ctx context.Context, req *RequestParams, requestMap interface{}) (reqHeaderMessage *RequestHeaderParams,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	reqHeaderMessage = &RequestHeaderParams{}
	if soteErr = reqHeaderMessage.PrepareReqHeader(req.Headers, req.TestMode); soteErr.ErrCode == nil {
		if soteErr = ValidateRequestMessage(ctx, reqHeaderMessage, req.TestMode); soteErr.ErrCode == nil {
			soteErr = PrepareMessage(ctx, req, requestMap)
		}
	}
	return
}

// PrepareReqHeader formats Request Headers to the header struct
func (rh *RequestHeaderParams) PrepareReqHeader(header map[string][]string, testMode bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
forLoop:
	for k, v := range header {
		switch k = strings.ToLower(k); k {
		case "organizations-id":
			if i, err := strconv.Atoi(v[0]); err == nil {
				rh.OrganizationsId = i
			} else {
				soteErr = sError.ConvertError(sError.GetSError(sError.ErrNotNumeric, sError.BuildParams([]string{k, v[0]}),
					sError.EmptyMap), testMode)
				break forLoop
			}

		case "role-list":
			rh.RoleList = v
		case "aws-user-name":
			rh.AWSUsername = v[0]
		case "message-id":
			rh.MessageId = v[0]
			sLogger.Info(fmt.Sprintf("starting processing message-id %v", v[0]))
		case "json-web-token":
			rh.JSONWebToken = v[0]
		case "reply-subject": // for business to business service calls
			rh.BusinessServiceId = v[0]
			rh.ReplySubjectNode = v[1]
		case "content-type":
			rh.ContentType = v
		case "origin":
			rh.Origin = v[0]
		}
	}

	if len(rh.ContentType) == 0 {
		rh.ContentType = []string{"application/json"}
	}

	return
}

// PrepareMessage unmarshal the incoming request message and validate it
func PrepareMessage(ctx context.Context, req *RequestParams, requestMap interface{}) (
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if requestMap != nil {
		if soteErr = sCustom.JSONUnmarshal(ctx, req.RequestMsg, &requestMap); soteErr.ErrCode != nil {
			soteErr = sError.ConvertError(soteErr, req.TestMode)
			return
		}

		if soteErr.ErrCode == nil {
			soteErr = ValidateRequestMessage(ctx, requestMap,
				req.TestMode) // validate the request against the request struct even when message is empty
		}
	}

	if soteErr.ErrCode == nil && req.JSONWebToken != "" && !req.TestMode {
		soteErr = sAuthentication.ValidToken(ctx, req.JSONWebToken, req.AuthConfig)
	}

	return
}

// ReadRequest reads data from HTTP request
func ReadRequest(ctx *gin.Context) (reqMsg []byte, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	switch ctx.Request.Method {
	case http.MethodGet, http.MethodDelete:
		qParams := ctx.Request.URL.Query().Get("params")
		if qParams != "" {
			reqMsg, err = base64.StdEncoding.DecodeString(qParams)
		}
	default:
		reqMsg, err = ioutil.ReadAll(ctx.Request.Body)
	}

	if err != nil {
		sLogger.Info(err.Error())
		soteErr = sError.GetSError(sError.ErrGenericError, sError.BuildParams([]string{""}), sError.EmptyMap)
	}

	ctx.Header("origin", ctx.Request.Host)

	return
}

// ValidateRequestMessage checks if the request message meets the necessary validation criteria/rules. (see validate struct tags)
func ValidateRequestMessage(ctx context.Context, requestMessage interface{}, testMode bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var (
		validate = validator.New()
		tMold    = modifiers.New()
		err      error
	)

	if err = tMold.Struct(ctx, requestMessage); err == nil {
		if err = validate.Struct(requestMessage); err != nil {
			sLogger.Info(err.Error())
			soteErr = sError.GetSError(sError.ErrInvalidMsgSignature, sError.BuildParams([]string{"See Configuration worksheet"}), sError.EmptyMap)
		}
	} else {
		sLogger.Info(err.Error())
		soteErr = sError.GetSError(sError.ErrGenericError, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
	}

	if soteErr.ErrCode != nil {
		soteErr = sError.ConvertError(soteErr, testMode)
	}

	return
}
