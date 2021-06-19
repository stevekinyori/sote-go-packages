package sConfigParams

import (
	"testing"

	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

const (
	// Application values
	API       string = "api"
	SDCC      string = "sdcc"
	SYNADIA   string = "synadia"
	DOCUMENTS string = "documents"
)

func init() {
	sLogger.SetLogMessagePrefix("sconfigparams_test.go")
}

//  THIS MUST BE THE FIRST TEST RUN - DO NOT MOVE OR PLACE A TEST BEFORE TestGetRegion
func TestGetRegion(tPtr *testing.T) {
	if _, soteErr := GetRegion(); soteErr.ErrCode != nil {
		tPtr.Errorf("GetRegion failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
		tPtr.Fatalf("RUN AWS_CONFIG or verify that ~/.aws/config exists and the region is set: %v", soteErr.FmtErrMsg)
	}
}
func TestGetParametersFound(tPtr *testing.T) {
	parameters := make(map[string]interface{})
	var soteErr sError.SoteError
	if parameters, soteErr = GetParameters(API, STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetParameters failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if len(parameters) == 0 {
		tPtr.Error("GetParameters failed: Expected parameters to have at least one entry")
	}
}
func TestGetParametersNotFound(tPtr *testing.T) {
	var soteErr sError.SoteError
	if _, soteErr = GetParameters(API, "SCOTT"); soteErr.ErrCode != 209110 {
		tPtr.Errorf("GetParameters failed: Expected soteErr to be 209110: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr = GetParameters("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetParameters failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr = GetParameters("", STAGING); soteErr.ErrCode != 200513 {
		tPtr.Errorf("GetParameters failed: Expected soteErr to be 200513: %v", soteErr.FmtErrMsg)
	}
}
func TestGetSmtpUsername(tPtr *testing.T) {
	if _, soteErr := GetSmtpUsername(API, STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetSmtpUsername failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetSmtpUsername("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetSmtpUsername failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetSmtpUsername("", STAGING); soteErr.ErrCode != 200513 {
		tPtr.Errorf("GetSmtpUsername failed: Expected soteErr to be 200513: %v", soteErr.FmtErrMsg)
	}
}
func TestGetSmtpPassword(tPtr *testing.T) {
	if _, soteErr := GetSmtpPassword(API, STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetSmtpPassword failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetSmtpPassword("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetSmtpPassword failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetSmtpPassword("", STAGING); soteErr.ErrCode != 200513 {
		tPtr.Errorf("GetSmtpPassword failed: Expected soteErr to be 200513: %v", soteErr.FmtErrMsg)
	}
}
func TestGetDBPassword(tPtr *testing.T) {
	if _, soteErr := GetDBPassword(API, STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBPassword failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBPassword("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBPassword failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBPassword("", STAGING); soteErr.ErrCode != 200513 {
		tPtr.Errorf("GetDBPassword failed: Expected soteErr to be 200513: %v", soteErr.FmtErrMsg)
	}
}
func TestGetDBName(tPtr *testing.T) {
	if _, soteErr := GetDBName(API, STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBName failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBName("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBName failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBName("", STAGING); soteErr.ErrCode != 200513 {
		tPtr.Errorf("GetDBName failed: Expected soteErr to be 200513: %v", soteErr.FmtErrMsg)
	}
}
func TestGetDBHost(tPtr *testing.T) {
	if _, soteErr := GetDBHost(API, STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBHost failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBHost("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBHost failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBHost("", STAGING); soteErr.ErrCode != 200513 {
		tPtr.Errorf("GetDBHost failed: Expected soteErr to be 200513: %v", soteErr.FmtErrMsg)
	}
}
func TestGetDBPort(tPtr *testing.T) {
	if _, soteErr := GetDBPort(API, STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBPort failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBPort("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBPort failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBPort("", STAGING); soteErr.ErrCode != 200513 {
		tPtr.Errorf("GetDBPort failed: Expected soteErr to be 200513: %v", soteErr.FmtErrMsg)
	}
}
func TestGetDBSSLMode(tPtr *testing.T) {
	if _, soteErr := GetDBSSLMode(API, STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBSSLMode failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBSSLMode("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBSSLMode failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBSSLMode("", STAGING); soteErr.ErrCode != 200513 {
		tPtr.Errorf("GetDBSSLMode failed: Expected soteErr to be 200513: %v", soteErr.FmtErrMsg)
	}
}
func TestGetDBUser(tPtr *testing.T) {
	if _, soteErr := GetDBUser(API, STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBUser failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBUser("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetDBUser failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBUser("", STAGING); soteErr.ErrCode != 200513 {
		tPtr.Errorf("GetDBUser failed: Expected soteErr to be 200513: %v", soteErr.FmtErrMsg)
	}
}
func TestGetUserPoolId(tPtr *testing.T) {
	if _, soteErr := GetUserPoolId(STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetUserPoolId failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetUserPoolId(""); soteErr.ErrCode != 209110 {
		tPtr.Errorf("GetUserPoolId failed: Expected soteErr to be 209110: %v", soteErr.FmtErrMsg)
	}
}
func TestGetClientId(tPtr *testing.T) {
	if _, soteErr := GetClientId(SDCC, STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetClientId failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetClientId("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetClientId failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetClientId("", STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetClientId failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
}
func TestValidateEnvironment(tPtr *testing.T) {
	if soteErr := ValidateEnvironment(DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if soteErr := ValidateEnvironment(STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if soteErr := ValidateEnvironment(DEMO); soteErr.ErrCode != nil {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if soteErr := ValidateEnvironment(PRODUCTION); soteErr.ErrCode != nil {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if soteErr := ValidateEnvironment("BAD_ENV"); soteErr.ErrCode != 209110 {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be 209110: %v", soteErr.FmtErrMsg)
	}
	if soteErr := ValidateEnvironment(""); soteErr.ErrCode != 209110 {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be 209110: %v", soteErr.FmtErrMsg)
	}
}
func TestGetEnvironmentVariable(tPtr *testing.T) {
	if _, soteErr := GetEnvironmentVariable(APPENV); soteErr.ErrCode != nil {
		tPtr.Errorf("GetEnvironmentVariable failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetEnvironmentAppEnvironment(); soteErr.ErrCode != nil {
		tPtr.Errorf("GetEnvironmentAppEnvironment failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
}
func TestGetNATSCredentials(tPtr *testing.T) {
	var (
		credValues func(string, string) (interface{}, sError.SoteError)
		soteErr    sError.SoteError
	)
	if credValues = GetNATSCredentials(); soteErr.ErrCode != nil {
		tPtr.Errorf("GetNATSCredentials failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}

	if _, soteErr = credValues(SYNADIA, STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetNATSCredentials failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr = credValues("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetNATSCredentials failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr = credValues(SYNADIA, ""); soteErr.ErrCode != 209110 {
		tPtr.Errorf("GetNATSCredentials failed: Expected soteErr to be 209110: %v", soteErr.FmtErrMsg)
	}
}
func TestGetNATSURL(tPtr *testing.T) {
	if _, soteErr := GetNATSURL(SYNADIA, STAGING); soteErr.ErrCode != nil {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetNATSURL("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetNATSURL("", STAGING); soteErr.ErrCode != 200513 {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be 200513: %v", soteErr.FmtErrMsg)
	}
}
func TestGetNATSTLSURLMask(tPtr *testing.T) {
	if _, soteErr := GetNATSTLSURLMask(SYNADIA); soteErr.ErrCode != nil {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetNATSTLSURLMask("SCOTT"); soteErr.ErrCode != 109999 {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be 109999: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetNATSTLSURLMask(""); soteErr.ErrCode != 200513 {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be 200513: %v", soteErr.FmtErrMsg)
	}
}
func TestSGetS3BucketURL(t *testing.T) {
	if _, soteErr := SGetS3BucketURL(DOCUMENTS, STAGING, PROCESSEDDOCUMENTSKEY); soteErr.ErrCode != nil {
		t.Errorf("SGetS3BucketURL failed: Expected soteErr to be nil got %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := SGetS3BucketURL(DOCUMENTS, STAGING, ""); soteErr.ErrCode != 109999 {
		t.Errorf("SGetS3BucketURL failed: Expected error code of 109999 got %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := SGetS3BucketURL("SCOTT", STAGING, UNPROCESSEDDOCUMENTSKEY); soteErr.ErrCode != 109999 {
		t.Errorf("SGetS3BucketURL failed: Expected error code of 109999 got %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := SGetS3BucketURL("", STAGING, UNPROCESSEDDOCUMENTSKEY); soteErr.ErrCode != 200513 {
		t.Errorf("SGetS3BucketURL failed: Expected error code of 200513 got %v", soteErr.FmtErrMsg)
	}

}
