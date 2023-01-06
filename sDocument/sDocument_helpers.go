package sDocument

/*
General description of the purpose of the go file.

RESTRICTIONS:
    AWS functions:
    * Program must have access to a .aws/credentials file in the default location.
    * This will only access system parameters that start with '/sote' (ROOTPATH).
    * {Enter other restrictions here for AWS

    {Other categories of restrictions}
    * AWS S3 Bucket must be mounted.
    * Mount point environment variable must be defined (Check constant.go for the name).

NOTES:
    {Enter any additional notes that you believe will help the next developer.}
*/

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gitlab.com/soteapps/packages/v2022/sError"
	"gitlab.com/soteapps/packages/v2022/sLogger"
)

func init() {
	sLogger.SetLogMessagePrefix(LOGMESSAGEPREFIX)
	testMode = true
}

// AmazonTextractErrorHandler converts error from Amazon Textract to a Sote Error
func AmazonTextractErrorHandler(ctx context.Context, err error) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	sLogger.Info(err.Error())
	// todo: build appropriate Sote error
	if strings.Contains(err.Error(), "InvalidParameterException") {
		soteErr = sError.GetSError(200514, []interface{}{"filename", "bucket name", "S3 bucket mount point"}, sError.EmptyMap)
	} else if strings.Contains(err.Error(), "InvalidS3ObjectException") || strings.Contains(strings.ToLower(err.Error()),
		"no such file or directory") || strings.Contains(strings.ToLower(err.Error()), "404") {
		soteErr = sError.GetSError(109999, []interface{}{"document"}, sError.EmptyMap)
	} else if strings.Contains(err.Error(), "input member Bucket must not be empty") {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"bucket-name"}), sError.EmptyMap)
	} else {
		soteErr = sError.GetSError(210599, nil, sError.EmptyMap)
	}

	return
}

// GetDocumentsMountPoint  returns location where documents s3 Bucket has been mounted*/
func GetDocumentsMountPoint(ctx context.Context, mountPointEnvName string) (documentsMountPoint string, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if documentsMountPoint = os.Getenv(mountPointEnvName); documentsMountPoint == "" {
		soteErr = panicService(ctx, sError.GetSError(209100, sError.BuildParams([]string{DOCUMENTSMOUNTPOINTENVIRONMENTVARNAME}), sError.EmptyMap))

	}

	return
}

/*GetFullDirectoryPath Will return current working directory */
func GetFullDirectoryPath() (dirPath string) {
	sLogger.DebugMethod()

	var (
		_, b, _, _ = runtime.Caller(0)
	)

	dirPath = filepath.Dir(b)

	return
}

// RemoveFile  will delete file in specified path and returns Sote Error
func RemoveFile(ctx context.Context, filepath string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	var (
		err error
	)

	sLogger.Info("Removing " + filepath)
	if err = os.Remove(filepath); err != nil {
		soteErr = AmazonTextractErrorHandler(ctx, err)
	}

	return
}

/*ValidateFilepath Checks if filepath exists */
func ValidateFilepath(filepath string) (pathExists bool, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err error
	)

	_, err = os.Stat(filepath)
	if os.IsNotExist(err) {
		pathExists = false
		soteErr = sError.GetSError(109999, sError.BuildParams([]string{fmt.Sprintf("document in %v path", filepath)}), sError.EmptyMap)
	} else {
		pathExists = true
		sLogger.Info(fmt.Sprintf("Document %v was found", filepath))
	}

	return
}

// panicService panic when not in test mode/production/demo
func panicService(ctx context.Context, inSoteErr sError.SoteError) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	sLogger.Info(inSoteErr.FmtErrMsg)
	soteErr = inSoteErr
	if !testMode {
		if inSoteErr.ErrCode != 199999 {
			soteErr = sError.GetSError(199999, sError.BuildParams([]string{""}), sError.EmptyMap)
		}

		if appEnvironment == "development" || appEnvironment == "staging" {
			defer func() {
				if r := recover(); r != nil {
					sLogger.Info(fmt.Sprintf("Recovered from panic %v", r))
				}
			}()

			panic(inSoteErr.FmtErrMsg)
		}
	}

	return
}
