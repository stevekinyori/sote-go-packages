package sHelper

import (
	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

const ENVDEFAULTAPPNAME = "synadia"
const ENVDEFAULTTARGET = "staging"

type Environment struct {
	ApplicationName   string
	TargetEnvironment string
	AppEnvironment    string
	TestMode          bool
}

func NewEnvironment(applicationName, targetEnvironment, appEnvironment string) (Environment, sError.SoteError) {
	sLogger.DebugMethod()
	var (
		env Environment
	)
	soteErr := sConfigParams.ValidateEnvironment(targetEnvironment)
	if soteErr.ErrCode == nil {
		env = Environment{
			ApplicationName:   applicationName,
			TargetEnvironment: targetEnvironment,
			AppEnvironment:    appEnvironment,
			TestMode:          targetEnvironment != "production",
		}
	}
	return env, soteErr
}
