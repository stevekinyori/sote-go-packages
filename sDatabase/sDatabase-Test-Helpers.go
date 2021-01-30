package sDatabase

import (
	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
)

const (
	API = "api"
)

var (
	DBName         string
	DBUser         string
	DBPassword     string
	DBHost         string
	DBSSLMode      string
	AppEnvironment string
	DBPort         int
)

func GetAWSParams() (soteErr sError.SoteError) {
	AppEnvironment, soteErr = sConfigParams.GetEnvironmentAppEnvironment()
	if soteErr.ErrCode == nil {
		DBName, soteErr = sConfigParams.GetDBName(API, AppEnvironment)
		if soteErr.ErrCode == nil {
			DBUser, soteErr = sConfigParams.GetDBUser(API, AppEnvironment)
			if soteErr.ErrCode == nil {
				DBPassword, soteErr = sConfigParams.GetDBPassword(API, AppEnvironment)
				if soteErr.ErrCode == nil {
					DBHost, soteErr = sConfigParams.GetDBHost(API, AppEnvironment)
					if soteErr.ErrCode == nil {
						DBSSLMode, soteErr = sConfigParams.GetDBSSLMode(API, AppEnvironment)
						if soteErr.ErrCode == nil {
							DBPort, soteErr = sConfigParams.GetDBPort(API, AppEnvironment)
						}
					}
				}
			}
		}
	} 

	return
}
