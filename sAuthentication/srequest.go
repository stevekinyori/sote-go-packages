package sAuthentication

import (
	"encoding/json"

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

func ValidateBody(data []byte, tEnvironment string) (*RequestHeaderSchema, error) {
	sLogger.DebugMethod()
	rh := RequestHeader{}
	err := json.Unmarshal(data, &rh)
	if err != nil {
		return nil, err
	}
	if rh.Header.AwsUserName == "" {
		err := json.Unmarshal(data, &rh.Header) //suports schema version 0.1
		if err != nil {
			return nil, err
		}
		return Validate(rh, tEnvironment)
	} else {
		return Validate(rh, tEnvironment)
	}
}

func Validate(rh RequestHeader, tEnvironment string) (*RequestHeaderSchema, error) {

	if rh.Header.AwsUserName == "" {
		err := sError.GetSError(206200, []interface{}{"#/properties/aws-user-name"}, nil)
		return nil, err

	}

	if rh.Header.OrganizationId == 0 {
		err := sError.GetSError(206200, []interface{}{"#/properties/organizations-id"}, nil)
		return nil, err
	}

	if rh.Header.JsonWebToken == "" {
		err := sError.GetSError(208355, nil, nil)
		return nil, err
	} else {
		//https://auth0.com/docs/tokens?_ga=2.253547273.1898510496.1593591557-1741611737.1593591372
		err := ValidToken(tEnvironment, rh.Header.JsonWebToken)
		if err != nil {
			return nil, err
		}
	}
	return &rh.Header, nil
}
