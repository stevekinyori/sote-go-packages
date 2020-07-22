package sConfigParams

import (
	"testing"

	"gitlab.com/soteapps/packages/v2020/sError"
)

const (
	// Application values
	API string = "api"
	SDCC string = "sdcc"
)

func TestGetParametersFound(t *testing.T) {
	parameters := make(map[string]interface{})
	var soteErr sError.SoteError
	if parameters, soteErr = GetParameters(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetParameters failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if len(parameters) == 0 {
		t.Error("GetParameters failed: Expected parameters to have at least one entry")
	}
}
func TestGetParametersNotFound(t *testing.T) {
	var soteErr sError.SoteError
	if _, soteErr = GetParameters(API, "SCOTT"); soteErr.ErrCode != 601010 {
		t.Errorf("GetParameters failed: Expected soteErr to be 601010: %v", soteErr.ErrCode)
	}
	if _, soteErr = GetParameters("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetParameters failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBPassword(t *testing.T) {
	if _, soteErr := GetDBPassword(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBPassword failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := GetDBPassword("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBPassword failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBName(t *testing.T) {
	if _, soteErr := GetDBName(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBName failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := GetDBName("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBName failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBHost(t *testing.T) {
	if _, soteErr := GetDBHost(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBHost failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := GetDBHost("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBHost failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBPort(t *testing.T) {
	if _, soteErr := GetDBPort(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBPort failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := GetDBPort("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBPort failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBSSLMode(t *testing.T) {
	if _, soteErr := GetDBSSLMode(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBSSLMode failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := GetDBSSLMode("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBSSLMode failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetDBUser(t *testing.T) {
	if _, soteErr := GetDBUser(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBUser failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if _, soteErr := GetDBUser("SCOTT", STAGING); soteErr.ErrCode != 109999 {
		t.Errorf("GetDBUser failed: Expected soteErr to be 109999: %v", soteErr.ErrCode)
	}
}
func TestGetRegion(t *testing.T) {
	if _, soteErr := GetRegion(); soteErr.ErrCode != nil {
		t.Errorf("GetRegion failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
}
func TestGetUserPoolId(t *testing.T) {
	if _, soteErr := GetUserPoolId(STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetUserPoolId failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
}
func TestGetClientId(t *testing.T) {
	if _, soteErr := GetClientId(SDCC, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetClientId failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
}
