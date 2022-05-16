package packages

import (
	"testing"

	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
)

const (
	// Application values
	API string = "api"
	SYNADIA string = "synadia"
)

func TestGetParametersFound(tPtr *testing.T) {
	parameters := make(map[string]interface{})
	var soteErr sError.SoteError
	if parameters, soteErr = sConfigParams.GetParameters(API, sConfigParams.STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetParameters failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if len(parameters) == 0 {
		tPtr.Error("GetParameters failed: Expected parameters to have at least one entry")
	}
}
func TestGetParametersNotFound(tPtr *testing.T) {
	var soteErr sError.SoteError
	if _, soteErr = sConfigParams.GetParameters(API, "SCOTT"); soteErr.ErrCode != 209110 {
		tPtr.Errorf("GetParameters failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr = sConfigParams.GetParameters("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetParameters failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBPassword(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetDBPassword(API, sConfigParams.STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBPassword failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBPassword("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBPassword failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBName(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetDBName(API, sConfigParams.STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBName failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBName("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBName failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBHost(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetDBHost(API, sConfigParams.STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBHost failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBHost("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBHost failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBPort(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetDBPort(API, sConfigParams.STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBPort failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBPort("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBPort failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBSSLMode(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetDBSSLMode(API, sConfigParams.STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBSSLMode failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBSSLMode("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBSSLMode failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBUser(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetDBUser(API, sConfigParams.STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBUser failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBUser("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBUser failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetRegion(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetRegion(); soteErr.ErrCode != nil {
		tPtr.Errorf("GetRegion failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
}
func TestGetUserPoolId(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetUserPoolId(sConfigParams.STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetUserPoolId failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
}
func TestValidateEnvironment(tPtr *testing.T) {
	if soteErr := sConfigParams.ValidateEnvironment(sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if soteErr := sConfigParams.ValidateEnvironment(sConfigParams.STAGING); soteErr.ErrCode != nil {
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
		soteErr sError.SoteError
	)
	if credValues = sConfigParams.GetNATSCredentials(); soteErr.ErrCode != nil {
		tPtr.Errorf("GetNATSCredentials failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}

	if _, soteErr = credValues(SYNADIA, sConfigParams.STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetNATSCredentials failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr = credValues("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetNATSCredentials failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr = credValues(SYNADIA, ""); soteErr.ErrCode != 209110 {
		tPtr.Errorf("GetNATSCredentials failed: Expected soteErr to be 209110: %v", soteErr.FmtErrMsg)
	}
}
func TestGetNATSURL(tPtr *testing.T) {
	if _, soteErr := sConfigParams.GetNATSURL(SYNADIA, sConfigParams.STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := sConfigParams.GetNATSURL("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := sConfigParams.GetNATSURL("", sConfigParams.STAGING); soteErr.ErrCode != 200513 {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be 200513: %v", soteErr.FmtErrMsg)
	}
}