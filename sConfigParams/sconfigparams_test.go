package sConfigParams

import (
	"testing"

	"gitlab.com/soteapps/packages/v2020/sError"
)

const (
	// Application values
	API string = "api"
)

func TestGetParameters(t *testing.T) {
	parameters := make(map[string]interface{})
	var soteErr sError.SoteError
	if parameters, soteErr = GetParameters(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetParameters failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if len(parameters) == 0 {
		t.Error("GetParameters failed: Expected parameters to have at least one entry")
	}
}
func TestGetDBPassword(t *testing.T) {
	var parameters string
	var soteErr sError.SoteError
	if parameters, soteErr = GetDBPassword(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("TestGetDBPassword failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if len(parameters) == 0 {
		t.Error("GetParameters failed: Expected parameters to have at least one entry")
	}
}
func TestGetDBName(t *testing.T) {
	var parameters string
	var soteErr sError.SoteError
	if parameters, soteErr = GetDBName(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBName failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if len(parameters) == 0 {
		t.Error("GetParameters failed: Expected parameters to have at least one entry")
	}
}
func TestGetDBHost(t *testing.T) {
	var parameters string
	var soteErr sError.SoteError
	if parameters, soteErr = GetDBHost(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBHost failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if len(parameters) == 0 {
		t.Error("GetParameters failed: Expected parameters to have at least one entry")
	}
}
func TestGetDBPort(t *testing.T) {
	var parameters string
	var soteErr sError.SoteError
	if parameters, soteErr = GetDBPort(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBPort failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if len(parameters) == 0 {
		t.Error("GetParameters failed: Expected parameters to have at least one entry")
	}
}
func TestGetDBSSLMode(t *testing.T) {
	var parameters string
	var soteErr sError.SoteError
	if parameters, soteErr = GetDBSSLMode(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBSSLMode failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if len(parameters) == 0 {
		t.Error("GetParameters failed: Expected parameters to have at least one entry")
	}
}
func TestGetDBUser(t *testing.T) {
	var parameters string
	var soteErr sError.SoteError
	if parameters, soteErr = GetDBUser(API, STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetDBUser failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if len(parameters) == 0 {
		t.Error("GetParameters failed: Expected parameters to have at least one entry")
	}
}
func TestGetRegion(t *testing.T) {
	var parameters string
	var soteErr sError.SoteError
	if parameters, soteErr = GetRegion(); soteErr.ErrCode != nil {
		t.Errorf("GetRegion failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if len(parameters) == 0 {
		t.Error("GetParameters failed: Expected parameters to have at least one entry")
	}
}
func TestGetUserPoolId(t *testing.T) {
	var parameters string
	var soteErr sError.SoteError
	if parameters, soteErr = GetUserPoolId(STAGING); soteErr.ErrCode != nil {
		t.Errorf("GetUserPoolId failed: Expected soteErr to be nil: %v", soteErr.ErrCode)
	}
	if len(parameters) == 0 {
		t.Error("GetUserPoolId failed: Expected parameters to have at least one entry")
	}
}
