package sConfigParams

import (
	"testing"

	"gitlab.com/soteapps/packages/v2020/sError"
)

const (
	// Application values
	API  string = "api"
)

func TestGetParameters(t *testing.T) {
	parameters := make(map[string]interface{})
	var soteErr sError.SoteError
	if parameters, soteErr = GetParameters(API, DEVELOPMENT); soteErr.ErrCode != nil {
		t.Errorf("GetParameters failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if len(parameters) == 0 {
		t.Error("GetParameters failed: Expected parameters to have at least one entry")
	}
}
func TestGetDBPassword(t *testing.T) {
	if _, soteErr := GetDBPassword(API, DEVELOPMENT); soteErr.ErrCode != nil {
		t.Errorf("TestGetDBPassword failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
}
func TestGetDBName(t *testing.T) {
	if _, soteErr := GetDBName(API, DEVELOPMENT); soteErr.ErrCode != nil {
		t.Errorf("GetDBName failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
}
func TestGetDBHost(t *testing.T) {
	if _, soteErr := GetDBHost(API, DEVELOPMENT); soteErr.ErrCode != nil {
		t.Errorf("GetDBHost failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
}
func TestGetDBPort(t *testing.T) {
	if _, soteErr := GetDBPort(API, DEVELOPMENT); soteErr.ErrCode != nil {
		t.Errorf("GetDBPort failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
}
func TestGetDBSSLMode(t *testing.T) {
	if _, soteErr := GetDBSSLMode(API, DEVELOPMENT); soteErr.ErrCode != nil {
		t.Errorf("GetDBSSLMode failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
}
func TestGetDBUser(t *testing.T) {
	if _, soteErr := GetDBUser(API, DEVELOPMENT); soteErr.ErrCode != nil {
		t.Errorf("GetDBUser failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
}
