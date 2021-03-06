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

func TestGetParametersFound(t *testing.T) {
	parameters := make(map[string]interface{})
	var soteErr sError.SoteError
	if parameters, soteErr = sConfigParams.GetParameters(API, sConfigParams.STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetParameters failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if len(parameters) == 0 {
		t.Error("GetParameters failed: Expected parameters to have at least one entry")
	}
}
func TestGetParametersNotFound(t *testing.T) {
	var soteErr sError.SoteError
	if _, soteErr = sConfigParams.GetParameters(API, "SCOTT"); soteErr.ErrCode != 209110 {
		t.Errorf("GetParameters failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr = sConfigParams.GetParameters("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetParameters failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBPassword(t *testing.T) {
	if _, soteErr := sConfigParams.GetDBPassword(API, sConfigParams.STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBPassword failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBPassword("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBPassword failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBName(t *testing.T) {
	if _, soteErr := sConfigParams.GetDBName(API, sConfigParams.STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBName failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBName("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBName failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBHost(t *testing.T) {
	if _, soteErr := sConfigParams.GetDBHost(API, sConfigParams.STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBHost failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBHost("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBHost failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBPort(t *testing.T) {
	if _, soteErr := sConfigParams.GetDBPort(API, sConfigParams.STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBPort failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBPort("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBPort failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBSSLMode(t *testing.T) {
	if _, soteErr := sConfigParams.GetDBSSLMode(API, sConfigParams.STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBSSLMode failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBSSLMode("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBSSLMode failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBUser(t *testing.T) {
	if _, soteErr := sConfigParams.GetDBUser(API, sConfigParams.STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBUser failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := sConfigParams.GetDBUser("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBUser failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetRegion(t *testing.T) {
	if _, soteErr := sConfigParams.GetRegion(); soteErr.ErrCode != nil {
		t.Errorf("GetRegion failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
}
func TestGetUserPoolId(t *testing.T) {
	if _, soteErr := sConfigParams.GetUserPoolId(sConfigParams.STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetUserPoolId failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
}
func TestValidateEnvironment(t *testing.T) {
	if soteErr := sConfigParams.ValidateEnvironment(sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		t.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if soteErr := sConfigParams.ValidateEnvironment(sConfigParams.STAGING); soteErr.ErrCode != nil {
		t.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if soteErr := sConfigParams.ValidateEnvironment(sConfigParams.DEMO); soteErr.ErrCode != nil {
		t.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if soteErr := sConfigParams.ValidateEnvironment(sConfigParams.PRODUCTION); soteErr.ErrCode != nil {
		t.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if soteErr := sConfigParams.ValidateEnvironment("BAD_ENV"); soteErr.ErrCode != 209110 {
		t.Errorf("ValidateEnvironment failed: Expected soteErr to be 209110: %v", soteErr.ErrCode)
	}
}
func TestGetNATSCredentials(t *testing.T) {
	var (
		credValues func(string, string) (interface{}, sError.SoteError)
		soteErr sError.SoteError
	)
	if credValues = sConfigParams.GetNATSCredentials(); soteErr.ErrCode != nil {
		t.Errorf("GetNATSCredentials failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}

	if _, soteErr = credValues(SYNADIA, sConfigParams.STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetNATSCredentials failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr = credValues("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetNATSCredentials failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr = credValues(SYNADIA, ""); soteErr.ErrCode != 209110 {
		t.Errorf("GetNATSCredentials failed: Expected soteErr to be 209110: %v", soteErr.FmtErrMsg)
	}
}
func TestGetNATSURL(t *testing.T) {
	if _, soteErr := sConfigParams.GetNATSURL(SYNADIA, sConfigParams.STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetNATSURL failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := sConfigParams.GetNATSURL("SCOTT", sConfigParams.STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetNATSURL failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := sConfigParams.GetNATSURL("", sConfigParams.STAGING); soteErr.ErrCode != 200513 {
		t.Errorf("GetNATSURL failed: Expected soteErr to be 200513: %v", soteErr.FmtErrMsg)
	}
}