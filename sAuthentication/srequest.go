package sAuthentication

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

const DEVICE_FILE = ".git/device.info"
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

func ValidateHeader(h nats.Header, tEnvironment string, isTestMode bool) (RequestHeaderSchema, sError.SoteError) {
	rh := RequestHeader{}
	rh.Header.JsonWebToken = h.Get("json-web-token")
	rh.Header.MessageId = h.Get("message-id")
	rh.Header.AwsUserName = h.Get("aws-user-name")
	rh.Header.RoleList = strings.Split(regexp.MustCompile(`\[|\]`).ReplaceAllString(h.Get("role-list"), ""), ",")
	fmt.Sscan(h.Get("organizations-id"), &rh.Header.OrganizationId)
	fmt.Sscan(h.Get("device-id"), &rh.Header.DeviceId)
	return Validate(rh, tEnvironment, isTestMode)
}

func ValidateBody(data []byte, tEnvironment string, isTestMode bool) (RequestHeaderSchema, sError.SoteError) {
	sLogger.DebugMethod()
	rh := RequestHeader{}
	json.Unmarshal(data, &rh) //flush stream

	if rh.Header.AwsUserName == "" {
		json.Unmarshal(data, &rh.Header) //suports schema version 0.1
		return Validate(rh, tEnvironment, isTestMode)
	} else {
		return Validate(rh, tEnvironment, isTestMode)
	}
}

func Validate(rh RequestHeader, tEnvironment string, isTestMode bool) (RequestHeaderSchema, sError.SoteError) {
	soteErr := sError.SoteError{}
	if rh.Header.AwsUserName == "" {
		soteErr = sError.GetSError(206200, []interface{}{"#/properties/aws-user-name"}, nil)
	} else if rh.Header.OrganizationId == 0 {
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
				return rh.Header, soteErr
			}
		}
		if rh.Header.DeviceId != 0 { //access from test scripts
			git, _ := filepath.Abs(".git")
			os.MkdirAll(git, os.ModePerm)
			path, _ := filepath.Abs(DEVICE_FILE)
			fileData, err := ioutil.ReadFile(path)
			if err != nil {
				return rh.Header, sError.GetSError(209010, []interface{}{DEVICE_FILE, err.Error()}, nil)
			} else {
				t, err := strconv.ParseInt(strings.TrimSpace(string(fileData)), 10, 0)
				if err != nil {
					sLogger.Debug(err.Error())
					return rh.Header, sError.GetSError(208355, nil, nil)
				}
				if t != rh.Header.DeviceId || math.Abs(float64(time.Now().Unix()-t)) >= DEVICE_TIMEOUT { //maxium duration of a device token
					return rh.Header, sError.GetSError(208350, nil, nil)
				} else {
					return rh.Header, soteErr
				}
			}
		}
	}

	if rh.Header.JsonWebToken == "" {
		soteErr = sError.GetSError(208355, nil, nil)
	} else {
		//https://auth0.com/docs/tokens?_ga=2.253547273.1898510496.1593591557-1741611737.1593591372
		soteErr = ValidToken(tEnvironment, rh.Header.JsonWebToken)
	}
	return rh.Header, soteErr
}
