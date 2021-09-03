package sHTTPClient

import (
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
	"testing"
)

func init() {
	sLogger.SetLogMessagePrefix("sHTTPConsumer/shttpclient_test.go")
}

const (
	TOKEN = "eyJraWQiOiJjbmhHanQwVXNqZFNHNzFvUWQ1cThTRjNTb29mOHBPNU1qTThMaDdNWDlrPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJkZDg5NmVhNi03NmI2LTQ1OGYtYWYyNC0zMTAyN2JiOGQzODMiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xXzAxOGQxNTUxLTg5NWQtNDZjYy1iNmQzLWFhMTllNTkyYWRjMSIsImV2ZW50X2lkIjoiNzRmMzVjNTYtN2I2YS00ZDBjLWIxMDYtYjIxOGFlY2MwNWFlIiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTYxMTgxMzYwOCwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfUVZQd3dDZzJjIiwiZXhwIjoxNjExODE3MjA4LCJpYXQiOjE2MTE4MTM2MDgsImp0aSI6ImUxOTgxMzEwLWZiMmYtNDUzYS1iY2U1LTI3NTIxMmJlZTgxMCIsImNsaWVudF9pZCI6InR0c21yNjUzcXRma3VjZ2dvZnBtcWJic3AiLCJ1c2VybmFtZSI6ImRkODk2ZWE2LTc2YjYtNDU4Zi1hZjI0LTMxMDI3YmI4ZDM4MyJ9.AvCI1fKfAE6XUDxHC4bi4nRjJvRk2tJY1JR3ivM2Hx9uHgFKx5lLpI0qTaTp9Ra_oC0ZqrA4b05hHSn0Hxqa7Yt0T0j4tG89-t65zlDeMF7HPDhJvzS6DUfpAneJ1m15plXH7ui0iRpGej2Z6Kk3pJEVroi40hQ29iHdVeVAHang0Xy_Vp0o7YUssXzU54H_-ds2RWdE-nRKiUIwTFQCCPDvnF5b-9HkrRkpU2WrYy3L9PGYxNwe8sFsJrzGfD1-yrhfkvYRKbmrQVyq-FGQsHkcNjxB0bfv1A1VOfJSCQ4-JPYS3BmxqldcqUwbkQ_V-FoV5LZCSe1P5b-VnQptvA"
	SURL  = "http://httpbin.org"
)

func TestNew(tPtr *testing.T) {
	if _, soteErr := New(SURL, TOKEN); soteErr.ErrCode != nil {
		tPtr.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}
}
func TestNewExpect210090(tPtr *testing.T) {
	if _, soteErr := New("", TOKEN); soteErr.ErrCode != 210090 {
		tPtr.Errorf("New Failed: Expected error code 210090 but got : %v", soteErr.ErrCode)
	}
}
func TestParamFormatting(tPtr *testing.T) {
	var (
		soteErr   sError.SoteError
		httpm     *HTTPManager
		reqParams map[string]interface{}
	)

	if httpm, soteErr = New(SURL, TOKEN); soteErr.ErrCode == nil {
		reqParams = map[string]interface{}{
			"param1": "value1",
			"param2": "value2",
		}

		if soteErr = httpm.paramFormatting(reqParams); soteErr.ErrCode != nil {
			tPtr.Errorf("paramFormatting failed: Expected error code to be nil but got %v", soteErr.ErrCode)
		}
	} else {
		tPtr.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}
}
func TestParamFormattingExpect207200(tPtr *testing.T) {
	var (
		soteErr   sError.SoteError
		httpm     *HTTPManager
		reqParams map[string]interface{}
	)

	if httpm, soteErr = New(SURL, TOKEN); soteErr.ErrCode == nil {
		reqParams = map[string]interface{}{
			"param1": "value1",
			"param2": struct {
				param1 string
				param2 float64
			}{"2value2", 1.75},
		}

		if soteErr = httpm.paramFormatting(reqParams); soteErr.ErrCode == nil {
			tPtr.Errorf("paramFormatting failed: Expected error code 207200 but got %v", soteErr.ErrCode)
		}
	} else {
		tPtr.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}
}
func TestConvertErrors(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		httpm   *HTTPManager
	)

	if httpm, soteErr = New(SURL, TOKEN); soteErr.ErrCode == nil {
		httpm.sHTTPMapPayload = map[string]interface{}{
			"errCode":    float64(0),
			"statusText": "Success",
			"codeLoc":    "sHTTPConsumer > shttpclient_test.go",
			"retPack":    make([]interface{}, 0),
		}

		if soteErr = httpm.convertErrors(); soteErr.ErrCode != nil {
			tPtr.Errorf("convertErrors failed: Expected error code to be nil but got %v", soteErr.ErrCode)
		}
	} else {
		tPtr.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}
}
func TestConvertErrorsExpect200600(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		httpm   *HTTPManager
	)

	if httpm, soteErr = New(SURL, TOKEN); soteErr.ErrCode == nil {
		httpm.sHTTPMapPayload = map[string]interface{}{
			"errCode":    float64(2097),
			"statusText": "HTTP Error",
			"codeLoc":    "sHTTPConsumer > shttpclient_test.go",
			"retPack":    "",
		}

		if soteErr = httpm.convertErrors(); soteErr.ErrCode != 200600 {
			tPtr.Errorf("convertErrors failed: Expected error code 200600 but got %v", soteErr.ErrCode)
		}
	} else {
		tPtr.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}
}
func TestConvertErrorsExpect207200(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		httpm   *HTTPManager
	)

	if httpm, soteErr = New(SURL, TOKEN); soteErr.ErrCode == nil {
		httpm.sHTTPMapPayload = map[string]interface{}{
			"errCode":    "",
			"statusText": "Error",
			"codeLoc":    "sHTTPConsumer > shttpclient_test.go",
			"retPack":    "",
		}

		if soteErr = httpm.convertErrors(); soteErr.ErrCode != 207200 {
			tPtr.Errorf("convertErrors failed: Expected error code 207200 but got %v", soteErr.ErrCode)
		}
	} else {
		tPtr.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}
}
func TestHTTPCall(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		httpm   *HTTPManager
	)

	if httpm, soteErr = New(SURL, TOKEN); soteErr.ErrCode == nil {
		httpm.reqParams = map[string]string{
			"q": "news",
		}

		if soteErr = httpm.sHTTPCall("GET", "/get"); soteErr.ErrCode != nil {
			tPtr.Errorf("sHTTPCall failed: Expected error code to be nil but got %v", soteErr.ErrCode)
		}
	} else {
		tPtr.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}
}
func TestHTTPCallExpect200600(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		httpm   *HTTPManager
	)

	if httpm, soteErr = New(SURL, TOKEN); soteErr.ErrCode == nil {
		httpm.reqParams = map[string]string{
			"q": "news",
		}

		if soteErr = httpm.sHTTPCall("DELETE", "/post"); soteErr.ErrCode != 200600 {
			tPtr.Errorf("sHTTPCall failed: Expected error code 200600 but got %v", soteErr.ErrCode)
		}
	} else {
		tPtr.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}
}
func TestGetExpect207200(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		httpm   *HTTPManager
	)

	if httpm, soteErr = New(SURL, TOKEN); soteErr.ErrCode == nil {
		reqParams := map[string]interface{}{
			"q": "news",
		}

		if soteErr = httpm.Get("/get", reqParams, true); soteErr.ErrCode != 207200 {
			tPtr.Errorf("Get failed: Expected error code 207200 but got %v", soteErr)
		}
	} else {
		tPtr.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}
}
func TestGetExpect200600(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		httpm   *HTTPManager
	)

	if httpm, soteErr = New(SURL, TOKEN); soteErr.ErrCode == nil {
		reqParams := map[string]interface{}{
			"q": "news",
		}

		if soteErr = httpm.Get("/delete", reqParams,true); soteErr.ErrCode != 200600 {
			tPtr.Errorf("Get failed: Expected error code 200600 but got %v", soteErr.ErrCode)
		}
	} else {
		tPtr.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}
}
func TestDeleteExpect207200(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		httpm   *HTTPManager
	)

	if httpm, soteErr = New(SURL, TOKEN); soteErr.ErrCode == nil {
		reqParams := make(map[string]interface{})

		if soteErr = httpm.Delete("/delete", reqParams, true); soteErr.ErrCode != 207200 {
			tPtr.Errorf("Post failed: Expected error code 207200 but got %v", soteErr.ErrCode)
		}
	} else {
		tPtr.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}
}
func TestDeleteExpect200600(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		httpm   *HTTPManager
	)

	if httpm, soteErr = New(SURL, TOKEN); soteErr.ErrCode == nil {
		reqParams := make(map[string]interface{})

		if soteErr = httpm.Delete("/get", reqParams, true); soteErr.ErrCode != 200600 {
			tPtr.Errorf("Post failed: Expected error code 200600 but got %v", soteErr.ErrCode)
		}
	} else {
		tPtr.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}
}
func TestPostExpect207200(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		httpm   *HTTPManager
	)

	if httpm, soteErr = New(SURL, TOKEN); soteErr.ErrCode == nil {
		reqParams := map[string]interface{}{
			"name": "value",
		}

		if soteErr = httpm.Post("/post", reqParams, true); soteErr.ErrCode != 207200 {
			tPtr.Errorf("Post failed: Expected error code 207200 but got %v", soteErr.ErrCode)
		}
	} else {
		tPtr.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}
}
func TestPostExpect200600(tPtr *testing.T) {
	var (
		soteErr sError.SoteError
		httpm   *HTTPManager
	)

	if httpm, soteErr = New(SURL, TOKEN); soteErr.ErrCode == nil {
		reqParams := map[string]interface{}{
			"name": "value",
		}

		if soteErr = httpm.Post("/get", reqParams, true); soteErr.ErrCode != 200600 {
			tPtr.Errorf("Post failed: Expected error code 200600 but got %v", soteErr.ErrCode)
		}
		if soteErr = httpm.Post("/get", reqParams, false); soteErr.ErrCode != 200600 {
			tPtr.Errorf("Post failed: Expected error code 200600 but got %v", soteErr.ErrCode)
		}
	} else {
		tPtr.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}
}
