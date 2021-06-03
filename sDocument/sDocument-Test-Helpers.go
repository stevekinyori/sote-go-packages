package sDocument

import (
	"os"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

/*  This will return prefix location of tesseracts training data */
func GetTessdataPrefix() (tessdataPrefix string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if tessdataPrefix = os.Getenv("TESSDATA_PREFIX"); tessdataPrefix == "" {
		soteErr = sError.GetSError(209100, sError.BuildParams([]string{"TESSDATA_PREFIX"}), sError.EmptyMap)
		sLogger.Info(soteErr.FmtErrMsg)
	}

	return
}
