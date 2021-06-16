package sAuthentication

import (
	"encoding/json"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sHelper"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

type RequestHeader struct {
	sHelper.RequestHeaderSchema //version 0.1
	//version 1.0+
	Header sHelper.RequestHeaderSchema `json:"request-header"`
}

func ValidateBody(data []byte, tApplication, tEnvironment string, isTestMode bool) (soteErr sError.SoteError) {
	rh := RequestHeader{}
	json.Unmarshal(data, &rh) //flush stream
	if rh.AwsUserName == "" && rh.Header.AwsUserName == "" {
		soteErr = sHelper.NewError().InvalidParameters("#/properties/aws-user-name")
	} else if rh.OrganizationId == 0 && rh.Header.OrganizationId == 0 {
		soteErr = sHelper.NewError().InvalidParameters("#/properties/organizations-id")
	}

	if soteErr.ErrCode != nil { //cannot send this error back to the client
		if isTestMode {
			sLogger.DebugMethod()
			panic(soteErr.FmtErrMsg)
		}
	} else if rh.JsonWebToken == "" && rh.Header.JsonWebToken == "" {
		soteErr = sHelper.NewError().InvalidToken()
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
