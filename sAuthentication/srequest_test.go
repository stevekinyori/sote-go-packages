package sAuthentication

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"gitlab.com/soteapps/packages/v2023/sConfigParams"
	"gitlab.com/soteapps/packages/v2023/sError"
)

var devExpToken = "eyJraWQiOiJjbmhHanQwVXNqZFNHNzFvUWQ1cThTRjNTb29mOHBPNU1qTThMaDdNWDlrPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJkZDg5NmVhNi03NmI2LTQ1OGYtYWYyNC0zMTAyN2JiOGQzODMiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xX2ZhMzUzYzMxLWM5ZGYtNGY3MC04YWUyLTMwZThjMTQ4ZTczNyIsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC5ldS13ZXN0LTEuYW1hem9uYXdzLmNvbVwvZXUtd2VzdC0xX1FWUHd3Q2cyYyIsImNsaWVudF9pZCI6InR0c21yNjUzcXRma3VjZ2dvZnBtcWJic3AiLCJldmVudF9pZCI6ImQyNzYzZTYzLWEzZjAtNDA4Zi1hOGE1LTIyYjAxZjdhOGE1MiIsInRva2VuX3VzZSI6ImFjY2VzcyIsInNjb3BlIjoiYXdzLmNvZ25pdG8uc2lnbmluLnVzZXIuYWRtaW4iLCJhdXRoX3RpbWUiOjE2NTgzMjk3NzUsImV4cCI6MTY1ODMzMzM3NSwiaWF0IjoxNjU4MzI5Nzc1LCJqdGkiOiJjZGU0ODljMy1kYTA5LTRlYWEtYWQ5NC04ZjQ0Nzk1YTcwMTYiLCJ1c2VybmFtZSI6ImRkODk2ZWE2LTc2YjYtNDU4Zi1hZjI0LTMxMDI3YmI4ZDM4MyJ9.Z_Pt10Mi3QyI8TbBaV068gHXDQYobuHmwMyBp0lkaU0mxvI66fN0Cyebrsnp6sXhidygxUEgWX7BmHADCwPjsi7hwmSujnXEi5T2VauCdWS05l1jNR1g0sse7uM8SUZpaK9jlMKPXzrqyiEOPG34E8gwJylLAfPYlFfgfkymU_6trygn2s5ZDjzSO8FNCrgkl1sLWXdBqKy7Tn9NmAJdJ8YccP4Ax5CS2iMo9BEc4ZD89oYbf0CABU9tlpB_Rnam2GN5xqLdIZUBCUQrWAYqct8oTsPf-xxRRa17nS_80SXh-WJA-8FdOWdd9rTyrmUurknzoadTV1TTghGRFWpMBA"

func AssertEqual(t *testing.T, actual, expected interface{}) {
	t.Helper() // get caller function a line number
	if expected != actual {
		t.Fatal(fmt.Sprintf("Not equal:\nexpected: %v\nactual:   %v", expected, actual))
	}
}

func validateBodyTest(environment string, data []byte) sError.SoteError {
	_, soteError := ValidateBody(parentCtx, data, environment, true)
	return soteError
}

func validateBodyMock(data []byte, tEnvironment string) sError.SoteError {
	_, soteError := ValidateBody(parentCtx, data, tEnvironment, true)
	return soteError
}

func TestInit(t *testing.T) {
	soteErr := validateBodyTest(sConfigParams.DEVELOPMENT, []byte(`{
		"aws-user-name": "soteuser",
		"organizations-id": 10003
	}`))
	AssertEqual(t, soteErr.ErrCode, nil)                   // ignored by default
	flag.Lookup("test.count").Value.(flag.Getter).Set("0") // enable unittest validations
}

func TestScriptAccessMissingFile(t *testing.T) {
	os.Remove(DEVICE_FILE)
	soteErr := validateBodyTest(sConfigParams.DEVELOPMENT, []byte(`{
		"aws-user-name": "soteuser",
		"organizations-id": 10003,
		"device-id": 123456789
	}`))
	AssertEqual(t, strings.Split(soteErr.FmtErrMsg, " Message return:")[0],
		fmt.Sprintf("%v: .git/device.info file was not found.", sError.ErrMissingFile))
}

func TestScriptAccessInvalidEntry(t *testing.T) {
	os.Remove(DEVICE_FILE)
	ioutil.WriteFile(DEVICE_FILE, []byte("HELLO WORLD"), 0644)
	soteErr := validateBodyTest(sConfigParams.DEVELOPMENT, []byte(`{
		"aws-user-name": "soteuser",
		"organizations-id": 10003,
		"device-id": 123456789
	}`))
	AssertEqual(t, soteErr.FmtErrMsg, fmt.Sprintf("%v: Token is invalid", sError.ErrInvalidToken))
}

func TestScriptAccessTimeoutToken(t *testing.T) {
	now := fmt.Sprint(time.Now().Unix() - DEVICE_TIMEOUT)
	os.Remove(DEVICE_FILE)
	ioutil.WriteFile(DEVICE_FILE, []byte(now), 0644)
	soteErr := validateBodyTest(sConfigParams.DEVELOPMENT, []byte(`{
		"request-header": {
			"aws-user-name": "soteuser",
			"organizations-id": 10003,
			"device-id": `+now+`
		}
	}`))
	AssertEqual(t, soteErr.FmtErrMsg, fmt.Sprintf("%v: Token is expired", sError.ErrExpiredToken))
}

func TestScriptAccess(t *testing.T) {
	now := fmt.Sprint(time.Now().Unix())
	os.Remove(DEVICE_FILE)
	ioutil.WriteFile(DEVICE_FILE, []byte(now), 0644)
	soteErr := validateBodyTest(sConfigParams.DEVELOPMENT, []byte(`{
		"aws-user-name": "soteuser",
		"organizations-id": 10003,
		"device-id": `+now+`
	}`))
	AssertEqual(t, soteErr.FmtErrMsg, "")
}

func TestRequestMissingAwsUserName(t *testing.T) {
	soteErr := validateBodyTest(sConfigParams.DEVELOPMENT, []byte(`{}`))
	AssertEqual(t, soteErr.FmtErrMsg, fmt.Sprintf("%v: Message doesn't match signature. "+
		"Sender must provide the following parameter names: #/properties/aws-user-name", sError.ErrInvalidMsgSignature))

}

func TestRequestMissingOrganizationId(t *testing.T) {
	soteErr := validateBodyTest(sConfigParams.DEVELOPMENT, []byte(`{
		"aws-user-name": "soteuser"
	}`))
	AssertEqual(t, soteErr.FmtErrMsg, fmt.Sprintf("%v: Message doesn't match signature. "+
		"Sender must provide the following parameter names: #/properties/organizations-id",
		sError.ErrInvalidMsgSignature))

}

func TestRequestMissingJsonWebToken(t *testing.T) {
	soteErr := validateBodyTest(sConfigParams.DEVELOPMENT, []byte(`{
		"aws-user-name": "soteuser",
		"organizations-id": 10003
	}`))
	AssertEqual(t, soteErr.FmtErrMsg, fmt.Sprintf("%v: Token is invalid", sError.ErrInvalidToken))
}

func TestRequestInvalidJsonWebToken(t *testing.T) {
	soteErr := validateBodyTest(sConfigParams.DEVELOPMENT, []byte(`{
		"json-web-token": "eyJraWQiOvxxx",
		"aws-user-name": "soteuser",
		"organizations-id": 10003
	}`))
	AssertEqual(t, soteErr.FmtErrMsg, fmt.Sprintf("%v: Token contains an invalid number of segments",
		sError.ErrMissingTokenSegments))
}

func TestRequestExpiredToken(t *testing.T) {
	soteErr := validateBodyTest(sConfigParams.DEVELOPMENT, []byte(`{
		"json-web-token": "`+devExpToken+`",
		"aws-user-name": "soteuser",
		"organizations-id": 10003
	}`))
	AssertEqual(t, soteErr.FmtErrMsg, fmt.Sprintf("%v: Token is expired", sError.ErrExpiredToken))
}

func TestRequestRequesrHeaderReleaseOne(t *testing.T) {
	soteErr := validateBodyMock([]byte(`{
		"json-web-token": "`+devExpToken+`",
		"aws-user-name": "soteuser",
		"organizations-id": 10003
	}`), sConfigParams.DEVELOPMENT)
	fmt.Println(soteErr)
	AssertEqual(t, soteErr.ErrCode, sError.ErrExpiredToken)
}

func TestRequestRequesrHeaderFutureReleases(t *testing.T) {
	soteErr := validateBodyMock([]byte(`{
		"request-header": {
			"json-web-token": "`+devExpToken+`",
			"aws-user-name": "soteuser",
			"organizations-id": 10003
		}
	}`), sConfigParams.DEVELOPMENT)
	AssertEqual(t, soteErr.ErrCode, sError.ErrExpiredToken)
}
