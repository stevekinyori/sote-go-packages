/*
	This is a wrapper for github.com/ddliu/go-httpclient.
	We are wrapping this so Sote Go developers can make Get,Post and Delete http calls the same way
*/

package sHTTPClient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"sync"

	"github.com/ddliu/go-httpclient"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

func init() {
	sLogger.SetLogMessagePrefix("sHTTPClient/shttpclient.go")
}

type HTTPManager struct {
	httpclient *httpclient.HttpClient
	sURL       string
	reqParams  map[string]string
	PayloadManager
	sync.Mutex
}

type PayloadManager struct {
	sHTTPResponse    *httpclient.Response
	sHTTPBytePayload []byte
	sHTTPMapPayload  map[string]interface{}
	RetPack          interface{}
	sync.Mutex
}

/*
	New will create a Sote HTTP Manager. The domain and token are required
*/
func New(sURL string, token string) (pHTTPManager *HTTPManager, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	//http client is instantiated with the line below
	pHTTPManager = &HTTPManager{httpclient: httpclient.Defaults(httpclient.Map{"Authorization": "Bearer " + token})}
	soteErr = pHTTPManager.setURL(sURL)
	return
}

/*
	This will parse the domain url
*/
func (httpm *HTTPManager) setURL(sURL string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if _, err := url.Parse(sURL); err != nil || sURL == "" {
		soteErr = sError.GetSError(210090, sError.BuildParams([]string{sURL}), nil)
	} else {
		httpm.sURL = sURL
	}

	return
}

/*
	This will make an HTTP DELETE call to the set route with the supplied request parameters
*/
func (httpm *HTTPManager) Delete(route string, reqParams map[string]interface{}, parseResult bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = httpm.paramFormatting(reqParams); soteErr.ErrCode == nil {
		if soteErr = httpm.sHTTPCall("DELETE", route); soteErr.ErrCode == nil {
			if soteErr = httpm.readHTTPResponse(); soteErr.ErrCode == nil {
				soteErr = httpm.parseJSONResult(parseResult)
			}
		}
	}

	return
}

/*
	This will make an HTTP GET call to the set route with the supplied request parameters
	Using the parseResult parameter, you can load the raw result or the JSON result
	into the struct payloadManager.RetPack using the parseResult parameter (false = load row, true = parse JSON)
*/
func (httpm *HTTPManager) Get(route string, reqParams map[string]interface{}, parseResult bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = httpm.paramFormatting(reqParams); soteErr.ErrCode == nil {
		if soteErr = httpm.sHTTPCall("GET", route); soteErr.ErrCode == nil {
			if soteErr = httpm.readHTTPResponse(); soteErr.ErrCode == nil {
				soteErr = httpm.parseJSONResult(parseResult)
			}
		}
	}

	return
}

/*
	This will make an HTTP POST call to the set route with the supplied request parameters
*/
func (httpm *HTTPManager) Post(route string, reqParams map[string]interface{}, parseResult bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if soteErr = httpm.paramFormatting(reqParams); soteErr.ErrCode == nil {
		if soteErr = httpm.sHTTPCall("POST", route); soteErr.ErrCode == nil {
			if soteErr = httpm.readHTTPResponse(); soteErr.ErrCode == nil {
				soteErr = httpm.parseJSONResult(parseResult)
			}
		}
	}

	return
}

/*
	This is a dynamic function for making HTTP calls using te httpclient
*/
func (httpm *HTTPManager) sHTTPCall(method string, route string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var err error

	// Success is indicated with 2xx status codes
	switch method {
	case "DELETE":
		if httpm.sHTTPResponse, err = httpm.httpclient.Delete(httpm.sURL+route, httpm.reqParams); err != nil || httpm.sHTTPResponse.StatusCode < 200 || httpm.sHTTPResponse.StatusCode >= 300 {
			soteErr = sError.GetSError(200600, sError.BuildParams([]string{httpm.sHTTPResponse.Status}), sError.EmptyMap)
			sLogger.Debug(soteErr.FmtErrMsg)
		}
	case "GET":
		if httpm.sHTTPResponse, err = httpm.httpclient.Get(httpm.sURL+route, httpm.reqParams); err != nil || httpm.sHTTPResponse.StatusCode < 200 || httpm.sHTTPResponse.StatusCode >= 300 {
			soteErr = sError.GetSError(200600, sError.BuildParams([]string{httpm.sHTTPResponse.Status}), sError.EmptyMap)
			sLogger.Debug(soteErr.FmtErrMsg)
		}
	case "POST":
		if httpm.sHTTPResponse, err = httpm.httpclient.PostJson(httpm.sURL+route, httpm.reqParams); err != nil || httpm.sHTTPResponse.StatusCode < 200 || httpm.sHTTPResponse.StatusCode >= 300 {
			soteErr = sError.GetSError(200600, sError.BuildParams([]string{httpm.sHTTPResponse.Status}), sError.EmptyMap)
			sLogger.Debug(soteErr.FmtErrMsg)
		}
	default:
		soteErr = sError.GetSError(200600, sError.BuildParams([]string{"HTTP method Not supported"}), sError.EmptyMap)
	}

	return
}

/*
	This converts the request parameters from reqParams map[string]interface{} to reqParams map[string]string
*/
func (httpm *HTTPManager) paramFormatting(reqParams map[string]interface{}) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	httpm.reqParams = make(map[string]string)

	// This is not a validation. It's for conversion if parameters exist
	if len(reqParams) > 0 {
		var (
			x  string
			ok bool
		)
	forLoop:
		for k, v := range reqParams {
			switch v.(type) {
			case nil:
				break
			case int:
				x = strconv.Itoa(v.(int))
			default:
				if x, ok = v.(string); !ok {
					soteErr = sError.GetSError(207200, sError.BuildParams([]string{"v", "string"}), sError.EmptyMap)
					sLogger.Debug(soteErr.FmtErrMsg)
					break forLoop
				}
			}
			httpm.reqParams[k] = x
		}
	}

	return
}

/*
	This reads HTTP response body to []byte
*/
func (httpm *HTTPManager) readHTTPResponse() (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var err error

	if httpm.sHTTPBytePayload, err = ioutil.ReadAll(httpm.sHTTPResponse.Body); err != nil {
		soteErr = sError.GetSError(200600, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
		sLogger.Debug(soteErr.FmtErrMsg)
	}

	err = httpm.sHTTPResponse.Body.Close()
	if err != nil {
		panic(err.Error())
	}

	return
}

/*
	This formats the HTTP []byte payload to Sote error if errCode is not 0 else puts the response in HTTPManager RetPack which is a map[string]interface{}
*/
func (httpm *HTTPManager) parseJSONResult(parseResult bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	if parseResult {
		if err = json.Unmarshal(httpm.sHTTPBytePayload, &httpm.sHTTPMapPayload); err != nil {
			soteErr = sError.GetSError(207105, sError.BuildParams([]string{"sPayloadBody", "[]byte"}), sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		} else if soteErr = httpm.convertErrors(); soteErr.ErrCode == nil {
			httpm.RetPack = httpm.sHTTPMapPayload["retPack"]
		}
	} else {
		httpm.RetPack = httpm.sHTTPBytePayload
	}

	return
}

/*
	This converts Sote PHP error codes to Sote Go error codes
*/
func (httpm *HTTPManager) convertErrors() (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if errCode, ok := httpm.sHTTPMapPayload["errCode"].(float64); ok {
		switch errCode {
		case 0: // 0 is not an error
			break
		case 201005:
			soteErr.ErrCode = 208360
			break
		case 2080, 2097:
			soteErr.ErrCode = 200600
			break
		case 500000:
			soteErr.ErrCode = 800000
			break
		case 800000:
			soteErr.ErrCode = 206200
			break
		default:
			soteErr.ErrCode = httpm.sHTTPMapPayload["errCode"]
			break
		}

		if soteErr.ErrCode != nil {
			soteErr.Loc = httpm.sHTTPMapPayload["codeLoc"].(string)
			soteErr.FmtErrMsg = httpm.sHTTPMapPayload["statusText"].(string)
			soteErr.ErrorDetails = map[string]string{"microservice_error": fmt.Sprintf("%v", httpm.sHTTPMapPayload["retPack"])}
		}
	} else {
		soteErr = sError.GetSError(207200, sError.BuildParams([]string{"errCode", "float64"}), sError.EmptyMap)
		sLogger.Debug(soteErr.FmtErrMsg)
	}

	return
}
