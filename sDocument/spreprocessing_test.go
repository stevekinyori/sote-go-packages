package sDocument

import (
	"testing"
)

func TestNewPreprocessor(t *testing.T) {
	if _, soteErr := NewPreprocessor(); soteErr.ErrCode != nil {
		t.Errorf("NewPreprocessor failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
}
func TestPreprocessManager_CorrectSkew(t *testing.T) {
	if pm, soteErr := NewPreprocessor(); soteErr.ErrCode == nil {
		if _, soteErr = pm.CorrectSkew("../img/testing_materials/Container Guarantee Form back.jpg"); soteErr.ErrCode != nil {
			t.Errorf("CorrectSkew failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
		}
	}

	if pm, soteErr := NewPreprocessor(); soteErr.ErrCode == nil {
		if _, soteErr = pm.CorrectSkew("img/testing_materials/Container Guarantee Form back.jpg"); soteErr.ErrCode != 210400 {
			t.Errorf("CorrectSkew failed: Expected error code of 210400: %v", soteErr.FmtErrMsg)
		}
	}
}




