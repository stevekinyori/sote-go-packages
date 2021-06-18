package sAuthentication

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"math"
	"path/filepath"
	"strconv"
	"time"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

const DEVICE_FILE = "../coverage.out"
const DEVICE_TIMEOUT = 30

type RequestHeaderSchema struct {
	JsonWebToken   string   `json:"json-web-token"`
	MessageId      string   `json:"message-id"`
	AwsUserName    string   `json:"aws-user-name"`
	OrganizationId int      `json:"organizations-id"`
	RoleList       []string `json:"role-list"` //optional
	DeviceId       int64    `json:"device-id"` //optional
}

type RequestHeader struct {
	RequestHeaderSchema //version 0.1
	//version 1.0+
	Header RequestHeaderSchema `json:"request-header"`
}

func ValidateBody(data []byte, tApplication, tEnvironment string, isTestMode bool) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	rh := RequestHeader{}
	json.Unmarshal(data, &rh) //flush stream
	if rh.AwsUserName == "" && rh.Header.AwsUserName == "" {
		soteErr = sError.GetSError(206200, []interface{}{"#/properties/aws-user-name"}, nil)
	} else if rh.OrganizationId == 0 && rh.Header.OrganizationId == 0 {
		soteErr = sError.GetSError(206200, []interface{}{"#/properties/organizations-id"}, nil)
	}

	if soteErr.ErrCode != nil { //cannot send this error back to the client
		if isTestMode {
			panic(soteErr.FmtErrMsg)
		}
	}

	if isTestMode {
		isUnitest := flag.Lookup("test.count")
		if isUnitest != nil { // go test ./...
			val, _ := isUnitest.Value.(flag.Getter)
			if val.Get().(uint) != 0 { //skip validation for unittests except srequest_test.go
				return
			}
		}
		if rh.DeviceId != 0 || rh.Header.DeviceId != 0 { //access from test scripts
			path, _ := filepath.Abs(DEVICE_FILE)
			fileData, err := ioutil.ReadFile(path)
			if err != nil {
				sLogger.Debug(err.Error())
				return sError.GetSError(208355, nil, nil)
			} else {
				t, err := strconv.ParseInt(string(fileData), 10, 0)
				if err != nil {
					sLogger.Debug(err.Error())
					return sError.GetSError(208355, nil, nil)
				}
				deviceId := rh.DeviceId
				if deviceId == 0 {
					deviceId = rh.Header.DeviceId
				}
				if t != deviceId || math.Abs(float64(time.Now().Unix()-t)) >= DEVICE_TIMEOUT { //maxium duration of a device token
					return sError.GetSError(208350, nil, nil)
				} else {
					return
				}
			}
		}
	}

	if rh.JsonWebToken == "" && rh.Header.JsonWebToken == "" {
		soteErr = sError.GetSError(208355, nil, nil)
	} else {
		//https://auth0.com/docs/tokens?_ga=2.253547273.1898510496.1593591557-1741611737.1593591372
		accessToken := rh.Header.JsonWebToken
		if accessToken == "" {
			accessToken = rh.JsonWebToken
		}
		soteErr = ValidToken(tApplication, tEnvironment, accessToken)
	}
	return
}
