package sAuthentication

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"gitlab.com/soteapps/packages/v2022/sConfigParams"
	"gitlab.com/soteapps/packages/v2022/sError"
	"gitlab.com/soteapps/packages/v2022/sLogger"
)

func InitiateAuth(ctx context.Context, environment string) (authResp *cognitoidentityprovider.InitiateAuthOutput, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		cfg              aws.Config
		err              error
		username         string
		password         string
		clientId         string
		params           = make(map[string]string)
		cognitoClientPtr *cognitoidentityprovider.Client
	)

	if cfg, err = config.LoadDefaultConfig(ctx); err == nil {
		cognitoClientPtr = cognitoidentityprovider.NewFromConfig(cfg)
		if clientId, soteErr = sConfigParams.GetClientId(ctx, "sdcc", environment); soteErr.ErrCode == nil {
			if username, soteErr = sConfigParams.GetDataLoadUser(ctx, environment); soteErr.ErrCode == nil {
				params["USERNAME"] = username
				if password, soteErr = sConfigParams.GetDataLoadPassword(ctx, environment); soteErr.ErrCode == nil {
					params["PASSWORD"] = password

					authTry := &cognitoidentityprovider.InitiateAuthInput{
						AuthFlow:       "USER_PASSWORD_AUTH",
						ClientId:       aws.String(clientId),
						AuthParameters: params,
					}

					if authResp, err = cognitoClientPtr.InitiateAuth(ctx, authTry); err != nil {
						sLogger.Info(err.Error())
						soteErr = sError.GetSError(201999, sError.BuildParams([]string{}), sError.EmptyMap)
					}
				}
			}
		}
	} else {
		sLogger.Info(err.Error())
		soteErr = sError.GetSError(201999, sError.BuildParams([]string{}), sError.EmptyMap)
	}

	return
}
