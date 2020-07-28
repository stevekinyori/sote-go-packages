package sDatabase

import (
	"gitlab.com/soteapps/packages/v2020/sConfigParams"
	"gitlab.com/soteapps/packages/v2020/sError"
)

var (
	AppEnvironment, DBName, DBUser, DBPassword, DBHost, DBSSLMode string
	DBPort                                                   int
)

func getAWSParams() (soteErr sError.SoteError) {
	AppEnvironment, soteErr = sConfigParams.GetEnvironmentAppEnvironment()
	if soteErr.ErrCode == nil {
		DBName, soteErr = sConfigParams.GetDBName("API", AppEnvironment)
		if soteErr.ErrCode == nil {
			DBUser, soteErr = sConfigParams.GetDBUser("API", AppEnvironment)
			if soteErr.ErrCode == nil {
				DBPassword, soteErr = sConfigParams.GetDBPassword("API", AppEnvironment)
				if soteErr.ErrCode == nil {
					DBHost, soteErr = sConfigParams.GetDBHost("API", AppEnvironment)
					if soteErr.ErrCode == nil {
						DBSSLMode, soteErr = sConfigParams.GetDBSSLMode("API", AppEnvironment)
						if soteErr.ErrCode == nil {
							DBPort, soteErr = sConfigParams.GetDBPort("API", AppEnvironment)
						}
					}
				}
			}
		}
	}

	return
}
