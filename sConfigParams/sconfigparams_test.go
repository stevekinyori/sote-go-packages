package sConfigParams

import (
	"testing"

	"gitlab.com/soteapps/packages/v2020/sError"
)

const (
	// Application values
	API  string = "api"
	SDCC string = "sdcc"
)

func TestGetParametersFound(t *testing.T) {
	parameters := make(map[string]interface{})
	var soteErr sError.SoteError
	if parameters, soteErr = GetParameters(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetParameters failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if len(parameters) == 0 {
		t.Error("GetParameters failed: Expected parameters to have at least one entry")
	}
}
func TestGetParametersNotFound(t *testing.T) {
	var soteErr sError.SoteError
	if _, soteErr = GetParameters(API, "SCOTT"); soteErr.ErrCode != 601010 {
		t.Errorf("GetParameters failed: Expected soteErr to be 601010: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr = GetParameters("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetParameters failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr = GetParameters("", STAGING); soteErr.ErrCode != 200512 {
		t.Errorf("GetParameters failed: Expected soteErr to be 200512: %v", soteErr.FmtErrMsg)
	}
}
func TestGetDBPassword(t *testing.T) {
	if _, soteErr := GetDBPassword(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBPassword failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBPassword("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBPassword failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBPassword("", STAGING); soteErr.ErrCode != 200512 {
		t.Errorf("GetDBPassword failed: Expected soteErr to be 200512: %v", soteErr.FmtErrMsg)
	}
}
func TestGetDBName(t *testing.T) {
	if _, soteErr := GetDBName(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBName failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBName("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBName failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBName("", STAGING); soteErr.ErrCode != 200512 {
		t.Errorf("GetDBName failed: Expected soteErr to be 200512: %v", soteErr.FmtErrMsg)
	}
}
func TestGetDBHost(t *testing.T) {
	if _, soteErr := GetDBHost(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBHost failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBHost("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBHost failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBHost("", STAGING); soteErr.ErrCode != 200512 {
		t.Errorf("GetDBHost failed: Expected soteErr to be 200512: %v", soteErr.FmtErrMsg)
	}
}
func TestGetDBPort(t *testing.T) {
	if _, soteErr := GetDBPort(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBPort failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBPort("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBPort failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBPort("", STAGING); soteErr.ErrCode != 200512 {
		t.Errorf("GetDBPort failed: Expected soteErr to be 200512: %v", soteErr.FmtErrMsg)
	}
}
func TestGetDBSSLMode(t *testing.T) {
	if _, soteErr := GetDBSSLMode(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBSSLMode failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBSSLMode("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBSSLMode failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBSSLMode("", STAGING); soteErr.ErrCode != 200512 {
		t.Errorf("GetDBSSLMode failed: Expected soteErr to be 200512: %v", soteErr.FmtErrMsg)
	}
}
func TestGetDBUser(t *testing.T) {
	if _, soteErr := GetDBUser(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBUser failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBUser("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBUser failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBUser("", STAGING); soteErr.ErrCode != 200512 {
		t.Errorf("GetDBUser failed: Expected soteErr to be 200512: %v", soteErr.FmtErrMsg)
	}
}
func TestGetRegion(t *testing.T) {
	if _, soteErr := GetRegion(); soteErr.ErrCode != nil {
		t.Errorf("GetRegion failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
}
func TestGetUserPoolId(t *testing.T) {
	if _, soteErr := GetUserPoolId(STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetUserPoolId failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetUserPoolId(""); soteErr.ErrCode != 200513 {
		t.Errorf("GetUserPoolId failed: Expected soteErr to be 200513: %v", soteErr.FmtErrMsg)
	}
}
func TestGetClientId(t *testing.T) {
	if _, soteErr := GetClientId(SDCC, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetClientId failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetClientId("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetClientId failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetClientId("", STAGING); soteErr.ErrCode != 200512 {
		t.Errorf("GetClientId failed: Expected soteErr to be 200512: %v", soteErr.FmtErrMsg)
	}
}
func TestValidateEnvironment(t *testing.T) {
	if soteErr := ValidateEnvironment(DEVELOPMENT); soteErr.ErrCode != nil {
		t.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if soteErr := ValidateEnvironment(STAGING); soteErr.ErrCode != nil {
		t.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if soteErr := ValidateEnvironment(DEMO); soteErr.ErrCode != nil {
		t.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if soteErr := ValidateEnvironment(PRODUCTION); soteErr.ErrCode != nil {
		t.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if soteErr := ValidateEnvironment("BAD_ENV"); soteErr.ErrCode != 601010 {
		t.Errorf("ValidateEnvironment failed: Expected soteErr to be 601010: %v", soteErr.FmtErrMsg)
	}
	if soteErr := ValidateEnvironment(""); soteErr.ErrCode != 601010 {
		t.Errorf("ValidateEnvironment failed: Expected soteErr to be 601010: %v", soteErr.FmtErrMsg)
	}
}
func TestGetEnvironmentVariable(t *testing.T) {
	if _, soteErr := GetEnvironmentVariable(APPENV); soteErr.ErrCode != nil {
		t.Errorf("GetEnvironmentVariable failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetEnvironmentVariable(AWSREGION); soteErr.ErrCode != nil {
		t.Errorf("GetEnvironmentVariable failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetEnvironmentAWSRegion(); soteErr.ErrCode != nil {
		t.Errorf("GetEnvironmentAWSRegion failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetEnvironmentAppEnvironment(); soteErr.ErrCode != nil {
		t.Errorf("GetEnvironmentAppEnvironment failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
}
