package packages

import (
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sHTTPClient"
	"gitlab.com/soteapps/packages/v2020/sLogger"
	//"fmt"
	"testing"
)

func init() {
	sLogger.SetLogMessagePrefix("sHTTPConsumer/shttpclient_test.go")
}

const (
	TOKEN = "eyJraWQiOiJjbmhHanQwVXNqZFNHNzFvUWQ1cThTRjNTb29mOHBPNU1qTThMaDdNWDlrPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJkZDg5NmVhNi03NmI2LTQ1OGYtYWYyNC0zMTAyN2JiOGQzODMiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xXzAxOGQxNTUxLTg5NWQtNDZjYy1iNmQzLWFhMTllNTkyYWRjMSIsImV2ZW50X2lkIjoiNzRmMzVjNTYtN2I2YS00ZDBjLWIxMDYtYjIxOGFlY2MwNWFlIiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTYxMTgxMzYwOCwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfUVZQd3dDZzJjIiwiZXhwIjoxNjExODE3MjA4LCJpYXQiOjE2MTE4MTM2MDgsImp0aSI6ImUxOTgxMzEwLWZiMmYtNDUzYS1iY2U1LTI3NTIxMmJlZTgxMCIsImNsaWVudF9pZCI6InR0c21yNjUzcXRma3VjZ2dvZnBtcWJic3AiLCJ1c2VybmFtZSI6ImRkODk2ZWE2LTc2YjYtNDU4Zi1hZjI0LTMxMDI3YmI4ZDM4MyJ9.AvCI1fKfAE6XUDxHC4bi4nRjJvRk2tJY1JR3ivM2Hx9uHgFKx5lLpI0qTaTp9Ra_oC0ZqrA4b05hHSn0Hxqa7Yt0T0j4tG89-t65zlDeMF7HPDhJvzS6DUfpAneJ1m15plXH7ui0iRpGej2Z6Kk3pJEVroi40hQ29iHdVeVAHang0Xy_Vp0o7YUssXzU54H_-ds2RWdE-nRKiUIwTFQCCPDvnF5b-9HkrRkpU2WrYy3L9PGYxNwe8sFsJrzGfD1-yrhfkvYRKbmrQVyq-FGQsHkcNjxB0bfv1A1VOfJSCQ4-JPYS3BmxqldcqUwbkQ_V-FoV5LZCSe1P5b-VnQptvA"
	SURL  = "http://httpbin.org"
)

func TestHTTPClientNew(t *testing.T) {
	if _, soteErr := sHTTPClient.New(SURL, TOKEN); soteErr.ErrCode != nil {
		t.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}

	if _, soteErr := sHTTPClient.New("", TOKEN); soteErr.ErrCode != 609990 {
		t.Errorf("New Failed: Expected error code 609990 but got : %v", soteErr.ErrCode)
	}
}

func TestGet(t *testing.T) {
	var (
		soteErr sError.SoteError
		httpm   *sHTTPClient.HTTPManager
	)

	if httpm, soteErr = sHTTPClient.New(SURL, TOKEN); soteErr.ErrCode == nil {
		reqParams := map[string]interface{}{
			"q": "news",
		}

		if soteErr = httpm.Get("/get", reqParams); soteErr.ErrCode != 400200 {
			t.Errorf("Get failed: Expected error code 400200 but got %v", soteErr)
		}

		if soteErr = httpm.Get("/delete", reqParams); soteErr.ErrCode != 201000 {
			t.Errorf("Get failed: Expected error code 201000 but got %v", soteErr.ErrCode)
		}
	} else {
		t.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}
}

func TestDelete(t *testing.T) {
	var (
		soteErr sError.SoteError
		httpm   *sHTTPClient.HTTPManager
	)

	if httpm, soteErr = sHTTPClient.New(SURL, TOKEN); soteErr.ErrCode == nil {
		reqParams := make(map[string]interface{})

		if soteErr = httpm.Delete("/delete", reqParams); soteErr.ErrCode != 400200 {
			t.Errorf("Post failed: Expected error code 400200 but got %v", soteErr.ErrCode)
		}

		if soteErr = httpm.Delete("/get", reqParams); soteErr.ErrCode != 201000 {
			t.Errorf("Post failed: Expected error code 201000 but got %v", soteErr.ErrCode)
		}
	} else {
		t.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}
}

func TestPost(t *testing.T) {
	var (
		soteErr sError.SoteError
		httpm   *sHTTPClient.HTTPManager
	)

	if httpm, soteErr = sHTTPClient.New(SURL, TOKEN); soteErr.ErrCode == nil {
		reqParams := map[string]interface{}{
			"name": "value",
		}

		if soteErr = httpm.Post("/post", reqParams); soteErr.ErrCode != 400200 {
			t.Errorf("Post failed: Expected error code 400200 but got %v", soteErr.ErrCode)
		}

		if soteErr = httpm.Post("/get", reqParams); soteErr.ErrCode != 201000 {
			t.Errorf("Post failed: Expected error code 201000 but got %v", soteErr.ErrCode)
		}
	} else {
		t.Errorf("New failed: Expected error code to be nil but got %v", soteErr.ErrCode)
	}
}
