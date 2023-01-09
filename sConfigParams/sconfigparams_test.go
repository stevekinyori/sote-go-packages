package sConfigParams

import (
	"context"
	"runtime"
	"testing"
	"time"

	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
)

var parentCtx = context.Background()

func init() {
	sLogger.SetLogMessagePrefix("sconfigparams_test.go")
}

//  THIS MUST BE THE FIRST TEST RUN - DO NOT MOVE OR PLACE A TEST BEFORE TestGetRegion
func TestGetRegion(tPtr *testing.T) {
	if _, soteErr := GetRegion(parentCtx); soteErr.ErrCode != nil {
		tPtr.Errorf("GetRegion failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
		tPtr.Fatalf("RUN AWS_CONFIG or verify that ~/.aws/config exists and the region is set: %v", soteErr.FmtErrMsg)
	}
}
func TestGetParametersFound(tPtr *testing.T) {
	parameters := make(map[string]interface{})
	var soteErr sError.SoteError
	if parameters, soteErr = GetParameters(parentCtx, API, DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetParameters failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if len(parameters) == 0 {
		tPtr.Error("GetParameters failed: Expected parameters to have at least one entry")
	}
}
func TestGetParametersNotFound(tPtr *testing.T) {
	var soteErr sError.SoteError
	if _, soteErr = GetParameters(parentCtx, API, "SCOTT"); soteErr.ErrCode != sError.ErrInvalidEnvValue {
		tPtr.Errorf("GetParameters failed: Expected soteErr to be %v: %v", sError.ErrInvalidEnvValue, soteErr.FmtErrMsg)
	}
	if _, soteErr = GetParameters(parentCtx, "SCOTT", DEVELOPMENT); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetParameters failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if _, soteErr = GetParameters(parentCtx, "", DEVELOPMENT); soteErr.ErrCode != sError.ErrMissingParameters {
		tPtr.Errorf("GetParameters failed: Expected soteErr to be %v: %v", sError.ErrMissingParameters,
			soteErr.FmtErrMsg)
	}
}
func TestGetSMTPConfig(tPtr *testing.T) {
	if _, soteErr := GetSMTPConfig(parentCtx, SMTP, DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetSmtpUsername failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetSMTPConfig(parentCtx, "MARY", DEVELOPMENT); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetSmtpUsername failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if _, soteErr := GetSMTPConfig(parentCtx, "", DEVELOPMENT); soteErr.ErrCode != sError.ErrMissingParameters {
		tPtr.Errorf("GetSmtpUsername failed: Expected soteErr to be %v: %v", sError.ErrMissingParameters, soteErr.FmtErrMsg)
	}
}
func TestGetQuickbooksConfig(tPtr *testing.T) {
	if _, soteErr := GetQuickbooksConfig(parentCtx, QUICKBOOKS, DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetQuickbooksConfig failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetQuickbooksConfig(parentCtx, "MARY", DEVELOPMENT); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetQuickbooksConfig failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound,
			soteErr.FmtErrMsg)
	}
	if _, soteErr := GetQuickbooksConfig(parentCtx, "", DEVELOPMENT); soteErr.ErrCode != sError.ErrMissingParameters {
		tPtr.Errorf("GetQuickbooksConfig failed: Expected soteErr to be %v: %v", sError.ErrMissingParameters, soteErr.FmtErrMsg)
	}
}
func TestGetCognitoConfig(tPtr *testing.T) {
	if _, soteErr := GetCognitoConfig(parentCtx, COGNITO, DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("TestGetCognitoConfig failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetCognitoConfig(parentCtx, "MARY", DEVELOPMENT); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("TestGetCognitoConfig failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound,
			soteErr.FmtErrMsg)
	}
	if _, soteErr := GetCognitoConfig(parentCtx, "", DEVELOPMENT); soteErr.ErrCode != sError.ErrMissingParameters {
		tPtr.Errorf("TestGetCognitoConfig failed: Expected soteErr to be %v: %v", sError.ErrMissingParameters, soteErr.FmtErrMsg)
	}
}
func TestGetSmtpUsername(tPtr *testing.T) {
	if _, soteErr := GetSmtpUsername(parentCtx, SMTP, DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetSmtpUsername failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetSmtpUsername(parentCtx, "SCOTT", DEVELOPMENT); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetSmtpUsername failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if _, soteErr := GetSmtpUsername(parentCtx, "", DEVELOPMENT); soteErr.ErrCode != sError.ErrMissingParameters {
		tPtr.Errorf("GetSmtpUsername failed: Expected soteErr to be %v: %v", sError.ErrMissingParameters, soteErr.FmtErrMsg)
	}
}
func TestGetSmtpPassword(tPtr *testing.T) {
	if _, soteErr := GetSmtpPassword(parentCtx, SMTP, DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetSmtpPassword failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetSmtpPassword(parentCtx, "SCOTT", DEVELOPMENT); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetSmtpPassword failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if _, soteErr := GetSmtpPassword(parentCtx, "", DEVELOPMENT); soteErr.ErrCode != sError.ErrMissingParameters {
		tPtr.Errorf("GetSmtpPassword failed: Expected soteErr to be %v: %v", sError.ErrMissingParameters, soteErr.FmtErrMsg)
	}
}
func TestGetDBPassword(tPtr *testing.T) {
	if _, soteErr := GetDBPassword(parentCtx, API, DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBPassword failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBPassword(parentCtx, "SCOTT", DEVELOPMENT); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetDBPassword failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBPassword(parentCtx, "", DEVELOPMENT); soteErr.ErrCode != sError.ErrMissingParameters {
		tPtr.Errorf("GetDBPassword failed: Expected soteErr to be %v: %v", sError.ErrMissingParameters, soteErr.FmtErrMsg)
	}
}
func TestGetDBName(tPtr *testing.T) {
	if _, soteErr := GetDBName(parentCtx, API, DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBName failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBName(parentCtx, "SCOTT", DEVELOPMENT); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetDBName failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBName(parentCtx, "", DEVELOPMENT); soteErr.ErrCode != sError.ErrMissingParameters {
		tPtr.Errorf("GetDBName failed: Expected soteErr to be %v: %v", sError.ErrMissingParameters, soteErr.FmtErrMsg)
	}
}
func TestGetDBSchema(tPtr *testing.T) {
	if _, soteErr := GetDBSchema(parentCtx, API, DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBSchema failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBSchema(parentCtx, "MARY", DEVELOPMENT); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetDBSchema failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBSchema(parentCtx, "", DEVELOPMENT); soteErr.ErrCode != sError.ErrMissingParameters {
		tPtr.Errorf("GetDBSchema failed: Expected soteErr to be %v: %v", sError.ErrMissingParameters, soteErr.FmtErrMsg)
	}
}
func TestGetDBHost(tPtr *testing.T) {
	if _, soteErr := GetDBHost(parentCtx, API, DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBHost failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBHost(parentCtx, "SCOTT", DEVELOPMENT); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetDBHost failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBHost(parentCtx, "", DEVELOPMENT); soteErr.ErrCode != sError.ErrMissingParameters {
		tPtr.Errorf("GetDBHost failed: Expected soteErr to be %v: %v", sError.ErrMissingParameters, soteErr.FmtErrMsg)
	}
}
func TestGetDBPort(tPtr *testing.T) {
	if _, soteErr := GetDBPort(parentCtx, API, DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBPort failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBPort(parentCtx, "SCOTT", DEVELOPMENT); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetDBPort failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBPort(parentCtx, "", DEVELOPMENT); soteErr.ErrCode != sError.ErrMissingParameters {
		tPtr.Errorf("GetDBPort failed: Expected soteErr to be %v: %v", sError.ErrMissingParameters, soteErr.FmtErrMsg)
	}
}
func TestGetDBSSLMode(tPtr *testing.T) {
	if _, soteErr := GetDBSSLMode(parentCtx, API, DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBSSLMode failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBSSLMode(parentCtx, "SCOTT", DEVELOPMENT); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetDBSSLMode failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBSSLMode(parentCtx, "", DEVELOPMENT); soteErr.ErrCode != sError.ErrMissingParameters {
		tPtr.Errorf("GetDBSSLMode failed: Expected soteErr to be %v: %v", sError.ErrMissingParameters, soteErr.FmtErrMsg)
	}
}
func TestGetDBUser(tPtr *testing.T) {
	if _, soteErr := GetDBUser(parentCtx, API, DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetDBUser failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBUser(parentCtx, "SCOTT", DEVELOPMENT); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetDBUser failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if _, soteErr := GetDBUser(parentCtx, "", DEVELOPMENT); soteErr.ErrCode != sError.ErrMissingParameters {
		tPtr.Errorf("GetDBUser failed: Expected soteErr to be %v: %v", sError.ErrMissingParameters, soteErr.FmtErrMsg)
	}
}
func TestGetUserPoolId(tPtr *testing.T) {
	if _, soteErr := GetUserPoolId(parentCtx, DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetUserPoolId failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetUserPoolId(parentCtx, ""); soteErr.ErrCode != sError.ErrInvalidEnvValue {
		tPtr.Errorf("GetUserPoolId failed: Expected soteErr to be %v: %v", sError.ErrInvalidEnvValue, soteErr.FmtErrMsg)
	}
}
func TestGetClientId(tPtr *testing.T) {
	if _, soteErr := GetClientId(parentCtx, SDCC, DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetClientId failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetClientId(parentCtx, "SCOTT", DEVELOPMENT); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetClientId failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if _, soteErr := GetClientId(parentCtx, "", DEVELOPMENT); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetClientId failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
}
func TestValidateEnvironment(tPtr *testing.T) {
	if soteErr := ValidateEnvironment(DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if soteErr := ValidateEnvironment(DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if soteErr := ValidateEnvironment(DEMO); soteErr.ErrCode != nil {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if soteErr := ValidateEnvironment(PRODUCTION); soteErr.ErrCode != nil {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if soteErr := ValidateEnvironment("BAD_ENV"); soteErr.ErrCode != sError.ErrInvalidEnvValue {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be %v: %v", sError.ErrInvalidEnvValue, soteErr.FmtErrMsg)
	}
	if soteErr := ValidateEnvironment(""); soteErr.ErrCode != sError.ErrInvalidEnvValue {
		tPtr.Errorf("ValidateEnvironment failed: Expected soteErr to be %v: %v", sError.ErrInvalidEnvValue, soteErr.FmtErrMsg)
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
	if credValues = GetNATSCredentials(parentCtx); soteErr.ErrCode != nil {
		tPtr.Errorf("GetNATSCredentials failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}

	if _, soteErr = credValues(SYNADIA, DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetNATSCredentials failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr = credValues("SCOTT", DEVELOPMENT); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetNATSCredentials failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if _, soteErr = credValues(SYNADIA, ""); soteErr.ErrCode != sError.ErrInvalidEnvValue {
		tPtr.Errorf("GetNATSCredentials failed: Expected soteErr to be %v: %v", sError.ErrInvalidEnvValue, soteErr.FmtErrMsg)
	}
}
func TestGetNATSURL(tPtr *testing.T) {
	if _, soteErr := GetNATSURL(parentCtx, SYNADIA, DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetNATSURL(parentCtx, "SCOTT", DEVELOPMENT); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if _, soteErr := GetNATSURL(parentCtx, "", DEVELOPMENT); soteErr.ErrCode != sError.ErrMissingParameters {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be %v: %v", sError.ErrMissingParameters, soteErr.FmtErrMsg)
	}
}
func TestGetNATSTLSURLMask(tPtr *testing.T) {
	if _, soteErr := GetNATSTLSURLMask(parentCtx, SYNADIA); soteErr.ErrCode != nil {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := GetNATSTLSURLMask(parentCtx, "SCOTT"); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if _, soteErr := GetNATSTLSURLMask(parentCtx, ""); soteErr.ErrCode != sError.ErrMissingParameters {
		tPtr.Errorf("GetNATSURL failed: Expected soteErr to be %v: %v", sError.ErrMissingParameters, soteErr.FmtErrMsg)
	}
}
func TestSGetS3BucketURL(t *testing.T) {
	if _, soteErr := SGetS3BucketURL(parentCtx, DOCUMENTS, DEVELOPMENT, PROCESSEDDOCUMENTSKEY); soteErr.ErrCode != nil {
		t.Errorf("SGetS3BucketURL failed: Expected soteErr to be nil got %v", soteErr.FmtErrMsg)
	}
	if _, soteErr := SGetS3BucketURL(parentCtx, DOCUMENTS, DEVELOPMENT, ""); soteErr.ErrCode != sError.ErrItemNotFound {
		t.Errorf("SGetS3BucketURL failed: Expected error code of %v got %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if _, soteErr := SGetS3BucketURL(parentCtx, "SCOTT", DEVELOPMENT, UNPROCESSEDDOCUMENTSKEY); soteErr.ErrCode != sError.ErrItemNotFound {
		t.Errorf("SGetS3BucketURL failed: Expected error code of %v got %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if _, soteErr := SGetS3BucketURL(parentCtx, "", DEVELOPMENT, UNPROCESSEDDOCUMENTSKEY); soteErr.ErrCode != sError.ErrMissingParameters {
		t.Errorf("SGetS3BucketURL failed: Expected error code of %v got %v", sError.ErrMissingParameters, soteErr.FmtErrMsg)
	}

}
func TestGetAWSS3Bucket(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
		soteErr           sError.SoteError
	)

	tPtr.Run("Get AWS S3 Bucket", func(tPtr *testing.T) {
		if _, soteErr = GetAWSS3Bucket(parentCtx, DOCUMENTS); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be %v but got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})
}
func TestUpdateQuickbooksRefreshToken(tPtr *testing.T) {
	var soteErr sError.SoteError

	param := QuickBooksRefreshToken{
		Token:      "test2",
		ExpiryDate: time.Now(),
	}
	if soteErr = UpdateQuickbooksRefreshToken(parentCtx, QUICKBOOKS, "MARY", param); soteErr.ErrCode != sError.ErrInvalidEnvValue {
		tPtr.Errorf("GetParameters failed: Expected soteErr to be %v: %v", sError.ErrInvalidEnvValue, soteErr.FmtErrMsg)
	}
	if soteErr = UpdateQuickbooksRefreshToken(parentCtx, "MARY", DEVELOPMENT, param); soteErr.ErrCode != sError.ErrItemNotFound {
		tPtr.Errorf("GetParameters failed: Expected soteErr to be %v: %v", sError.ErrItemNotFound, soteErr.FmtErrMsg)
	}
	if soteErr = UpdateQuickbooksRefreshToken(parentCtx, "", DEVELOPMENT, param); soteErr.ErrCode != sError.ErrMissingParameters {
		tPtr.Errorf("GetParameters failed: Expected soteErr to be %v: %v", sError.ErrMissingParameters, soteErr.FmtErrMsg)
	}
}
