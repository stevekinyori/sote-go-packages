package packages

import (
	"context"
	"runtime"
	"testing"

	"gitlab.com/soteapps/packages/v2022/sConfigParams"
	"gitlab.com/soteapps/packages/v2022/sError"
)

const (
	// Application values
	API       string = "api"
	SYNADIA   string = "synadia"
	DOCUMENTS string = "documents"
)

var parentCtx = context.Background()

func TestGetParametersFound(tPtr *testing.T) {
	parameters := make(map[string]interface{})
	var soteErr sError.SoteError
	if parameters, soteErr = sConfigParams.GetParameters(parentCtx, API, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetParameters failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if len(parameters) == 0 {
		tPtr.Error("GetParameters failed: Expected parameters to have at least one entry")
	}
}
func TestGetParametersNotFound(tPtr *testing.T) {
	var soteErr sError.SoteError
	if _, soteErr = sConfigParams.GetParameters(parentCtx, API, "SCOTT"); soteErr.ErrCode != 209110 {
		tPtr.Errorf("GetParameters failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr = sConfigParams.GetParameters(parentCtx, "SCOTT", sConfigParams.DEVELOPMENT); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetParameters failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBPassword(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetDBPassword(parentCtx, API, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBPassword failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBPassword(parentCtx, "SCOTT", sConfigParams.DEVELOPMENT); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBPassword failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBName(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetDBName(parentCtx, API, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBName failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBName(parentCtx, "SCOTT", sConfigParams.DEVELOPMENT); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBName failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBHost(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetDBHost(parentCtx, API, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBHost failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBHost(parentCtx, "SCOTT", sConfigParams.DEVELOPMENT); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBHost failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBPort(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetDBPort(parentCtx, API, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBPort failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBPort(parentCtx, "SCOTT", sConfigParams.DEVELOPMENT); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBPort failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBSSLMode(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetDBSSLMode(parentCtx, API, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBSSLMode failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBSSLMode(parentCtx, "SCOTT", sConfigParams.DEVELOPMENT); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBSSLMode failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBUser(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetDBUser(parentCtx, API, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBUser failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBUser(parentCtx, "SCOTT", sConfigParams.DEVELOPMENT); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBUser failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetRegion(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetRegion(parentCtx); soteErr.ErrCode != nil {
		tPtr.Errorf("GetRegion failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
}
func TestGetUserPoolId(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetUserPoolId(parentCtx, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetUserPoolId failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
}
func TestValidateEnvironment(tPtr *testing.T) {
	if soteErr := sConfigParams.ValidateEnvironment(sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if soteErr := sConfigParams.ValidateEnvironment(sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if soteErr := sConfigParams.ValidateEnvironment(sConfigParams.DEMO); soteErr.ErrCode != nil {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if soteErr := sConfigParams.ValidateEnvironment(sConfigParams.PRODUCTION); soteErr.ErrCode != nil {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if soteErr := sConfigParams.ValidateEnvironment("BAD_ENV"); soteErr.ErrCode != 209110 {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be 209110: %v", soteErr.ErrCode)
	}
}
func TestGetNATSCredentials(tPtr *testing.T) {
	var (
		credValues func(string, string) (interface{}, sError.SoteError)
		soteErr    sError.SoteError
	)
	if credValues = sConfigParams.GetNATSCredentials(parentCtx); soteErr.ErrCode != nil {
		tPtr.Errorf("GetNATSCredentials failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}

	if _, soteErr = credValues(SYNADIA, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetNATSCredentials failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr = credValues("SCOTT", sConfigParams.DEVELOPMENT); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetNATSCredentials failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr = credValues(SYNADIA, ""); soteErr.ErrCode != 209110 {
		tPtr.Errorf("GetNATSCredentials failed: Expected soteErr to be 209110: %v", soteErr.FmtErrMsg)
	}
}
func TestGetNATSURL(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetNATSURL(parentCtx, SYNADIA, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := sConfigParams.GetNATSURL(parentCtx, "SCOTT", sConfigParams.DEVELOPMENT); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := sConfigParams.GetNATSURL(parentCtx, "", sConfigParams.DEVELOPMENT); soteErr.ErrCode != 200513 {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be 200513: %v", soteErr.FmtErrMsg)
	}
}
func TestGetAWSClientId(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	tPtr.Run("Get AWS Client ID", func(tPtr *testing.T) {
		if _, soteErr = sConfigParams.GetAWSAccountId(parentCtx); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})
}
func TestGetAWSS3Bucket(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	tPtr.Run("Get AWS S3 Bucket", func(tPtr *testing.T) {
		if _, soteErr = sConfigParams.GetAWSS3Bucket(parentCtx, DOCUMENTS); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})
}
