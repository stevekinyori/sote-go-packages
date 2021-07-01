package sAuthentication

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/dgrijalva/jwt-go"
	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
)

var stagingExpToken = "eyJraWQiOiJlOCt4TW4rOGYrZmlIXC9OZDNDZGNxOVRvU3FPKzdZYldcL1wvSUxCYVJyTElNPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJjYzJiNGYwYS0xYmE0LTQzNTEtYmRmMS0wYmM3NTRhN2NlNjMiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xX2E1MjFiZTA5LTQxMDQtNDc1MC1iZTQwLTQ2NTExYzczYzA2MCIsImNvZ25pdG86Z3JvdXBzIjpbIjEwMDM2Il0sImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC5ldS13ZXN0LTEuYW1hem9uYXdzLmNvbVwvZXUtd2VzdC0xX2ZwWkN5cWxRTiIsImNsaWVudF9pZCI6IjNlMzN0NGVjb2Vpam5ibTNscGVoZmNuaWcwIiwiZXZlbnRfaWQiOiIwOGVkNjhmZC1hZDk4LTRiMjEtYjNhZC0wN2U0MTk0YzFkOTkiLCJ0b2tlbl91c2UiOiJhY2Nlc3MiLCJzY29wZSI6ImF3cy5jb2duaXRvLnNpZ25pbi51c2VyLmFkbWluIiwiYXV0aF90aW1lIjoxNjIzODcwNTIzLCJleHAiOjE2MjM4NzA4MjMsImlhdCI6MTYyMzg3MDUyNCwianRpIjoiMjFmNWViZmItNDIxMC00MWU1LTllNTAtMWI2ZmNlOWExZDA0IiwidXNlcm5hbWUiOiJjYzJiNGYwYS0xYmE0LTQzNTEtYmRmMS0wYmM3NTRhN2NlNjMifQ.jngSSScEF3Nm1-XghY6-TGljngjKHk8LaXvBAofjLJB1NzYVh3CEk8UC4JdPswjhqc0_xdEclI6XRxHyM_44uxhFzFyzpHU39x1ly5ROfls_rSvucxgVbHGoaOkmtenOOlBwFarszgzcCdngkY_rujr-r_YhIrAYH-Y1JnpxZ_idPBzgfQ9W1O0UPmFOZFxyDPsC2kAGU0Zl_cBcPrqXMULywx6nc8Vxo3sqDTy1UFE80Ysyz8d_SmeCDihL-YXa2Lm4p_wGnk54rUJ5VG_eZJ6n5FEopdhXoPqY0lvkuWz01NbwtunPKENGTA-qi-o__UBHfpYcOA7Nx5HkI48j9Q"

func AssertEqual(t *testing.T, actual, expected interface{}) {
	t.Helper() // get caller function a line number
	if expected != actual {
		t.Fatal(fmt.Sprintf("Not equal:\nexpected: %v\nactual:   %v", expected, actual))
	}
}

func validateBodyTest(data []byte) sError.SoteError {
	_, soteError := ValidateBody(data, sConfigParams.STAGING, true)
	return soteError
}

func validateBodyMock(data []byte, tEnvironment string) sError.SoteError {
	var (
		validPatch  *monkey.PatchGuard
		verifyPatch *monkey.PatchGuard
		rsa         *jwt.SigningMethodRSA
	)
	validPatch = monkey.Patch(jwt.MapClaims.Valid, func(jwt.MapClaims) error {
		validPatch.Unpatch()
		return nil
	})
	verifyPatch = monkey.PatchInstanceMethod(reflect.TypeOf(rsa), "Verify", func(*jwt.SigningMethodRSA, string, string, interface{}) error {
		verifyPatch.Unpatch()
		return nil
	})
	_, soteError := ValidateBody(data, tEnvironment, true)
	return soteError
}

func TestInit(t *testing.T) {
	soteErr := validateBodyTest([]byte(`{
		"aws-user-name": "soteuser",
		"organizations-id": 10003
	}`))
	AssertEqual(t, soteErr.ErrCode, nil)                   //ignored by default
	flag.Lookup("test.count").Value.(flag.Getter).Set("0") //enable unittest validations
}

func TestScriptAccessMissingFile(t *testing.T) {
	os.Remove(DEVICE_FILE)
	soteErr := validateBodyTest([]byte(`{
		"aws-user-name": "soteuser",
		"organizations-id": 10003,
		"device-id": 123456789
	}`))
	AssertEqual(t, strings.Split(soteErr.FmtErrMsg, " Message return:")[0], "209010: .git/device.info file was not found.")
}

func TestScriptAccessInvalidEntry(t *testing.T) {
	os.Remove(DEVICE_FILE)
	ioutil.WriteFile(DEVICE_FILE, []byte("HELLO WORLD"), 0644)
	soteErr := validateBodyTest([]byte(`{
		"aws-user-name": "soteuser",
		"organizations-id": 10003,
		"device-id": 123456789
	}`))
	AssertEqual(t, soteErr.FmtErrMsg, "208355: Token is invalid")
}

func TestScriptAccessTimeoutToken(t *testing.T) {
	now := fmt.Sprint(time.Now().Unix() - DEVICE_TIMEOUT)
	os.Remove(DEVICE_FILE)
	ioutil.WriteFile(DEVICE_FILE, []byte(now), 0644)
	soteErr := validateBodyTest([]byte(`{
		"request-header": {
			"aws-user-name": "soteuser",
			"organizations-id": 10003,
			"device-id": ` + now + `
		}
	}`))
	AssertEqual(t, soteErr.FmtErrMsg, "208350: Token is expired")
}

func TestScriptAccess(t *testing.T) {
	now := fmt.Sprint(time.Now().Unix())
	os.Remove(DEVICE_FILE)
	ioutil.WriteFile(DEVICE_FILE, []byte(now), 0644)
	soteErr := validateBodyTest([]byte(`{
		"aws-user-name": "soteuser",
		"organizations-id": 10003,
		"device-id": ` + now + `
	}`))
	AssertEqual(t, soteErr.FmtErrMsg, "")
}

func TestRequestMissingAwsUserName(t *testing.T) {
	defer func() {
		r := recover()
		AssertEqual(t, r, "206200: Message doesn't match signature. Sender must provide the following parameter names: #/properties/aws-user-name")
	}()
	validateBodyTest([]byte(`{}`))
}

func TestRequestMissingOrganizationId(t *testing.T) {
	defer func() {
		r := recover()
		AssertEqual(t, r, "206200: Message doesn't match signature. Sender must provide the following parameter names: #/properties/organizations-id")
	}()
	validateBodyTest([]byte(`{
		"aws-user-name": "soteuser"
	}`))
}

func TestRequestMissingJsonWebToken(t *testing.T) {
	soteErr := validateBodyTest([]byte(`{
		"aws-user-name": "soteuser",
		"organizations-id": 10003
	}`))
	AssertEqual(t, soteErr.FmtErrMsg, "208355: Token is invalid")
}

func TestRequestInvalidJsonWebToken(t *testing.T) {
	soteErr := validateBodyTest([]byte(`{
		"json-web-token": "eyJraWQiOvxxx",
		"aws-user-name": "soteuser",
		"organizations-id": 10003
	}`))
	AssertEqual(t, soteErr.FmtErrMsg, "208356: Token contains an invalid number of segments")
}

func TestRequestExpiredToken(t *testing.T) {
	soteErr := validateBodyTest([]byte(`{
		"json-web-token": "` + stagingExpToken + `",
		"aws-user-name": "soteuser",
		"organizations-id": 10003
	}`))
	AssertEqual(t, soteErr.FmtErrMsg, "208350: Token is expired")
}

/*func TestRequestInvalidEnvironment(t *testing.T) {
	soteErr := ValidateBody([]byte(`{
		"json-web-token": "`+stagingExpToken+`",
		"aws-user-name": "soteuser",
		"organizations-id": 10003
	}`), "internal-clearance", sConfigParams.DEVELOPMENT, true)
	AssertEqual(t, soteErr.FmtErrMsg, "208355: Token is invalid")
}*/

func TestRequestInvalidKid(t *testing.T) {
	soteErr := validateBodyMock([]byte(`{
		"json-web-token": "`+stagingExpToken+`",
		"aws-user-name": "soteuser",
		"organizations-id": 10003
	}`), sConfigParams.DEVELOPMENT)
	AssertEqual(t, soteErr.FmtErrMsg, "209521: Kid (e8+xMn+8f+fiH/Nd3Cdcq9ToSqO+7YbW//ILBaRrLIM=) was not found in public key set")
}

/*func TestRequestInvalidAppName(t *testing.T) {
	soteErr := validateBodyMock([]byte(`{
		"json-web-token": "`+stagingExpToken+`",
		"aws-user-name": "soteuser",
		"organizations-id": 10003
	}`), sConfigParams.STAGING)
	AssertEqual(t, soteErr.FmtErrMsg, "208340: client id is not valid for this application")
}*/

func TestRequestRequesrHeaderReleaseOne(t *testing.T) {
	soteErr := validateBodyMock([]byte(`{
		"json-web-token": "`+stagingExpToken+`",
		"aws-user-name": "soteuser",
		"organizations-id": 10003
	}`), sConfigParams.STAGING)
	AssertEqual(t, soteErr.FmtErrMsg, "")
}

func TestRequestRequesrHeaderFutureReleases(t *testing.T) {
	soteErr := validateBodyMock([]byte(`{
		"request-header": {
			"json-web-token": "`+stagingExpToken+`",
			"aws-user-name": "soteuser",
			"organizations-id": 10003
		}
	}`), sConfigParams.STAGING)
	AssertEqual(t, soteErr.FmtErrMsg, "")
}
