package sDocument

import (
	"testing"
)

func TestNewPreprocessor(t *testing.T) {
	if _, soteErr := NewPreprocessor("../img/testing_materials/Container Guarantee Form back.jpg"); soteErr.ErrCode != nil {
		t.Errorf("NewPreprocessor failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}

	if _, soteErr := NewPreprocessor("img/testing_materials/Container Guarantee Form back.jpg"); soteErr.ErrCode != 210400 {
		t.Errorf("NewPreprocessor failed: Expected error code of 210400: %v", soteErr.FmtErrMsg)
	}
}
func TestCheckIfPathExists(t *testing.T) {
	if _, soteErr := CheckIfPathExists("../img/testing_materials/Container Guarantee Form back.jpg"); soteErr.ErrCode != nil {
		t.Errorf("CheckIfFilePathExists failed: Expected soteErr to be nil:%v", soteErr.FmtErrMsg)
	}

	if _, soteErr := CheckIfPathExists("img/testing_materials/Container Guarantee Form back.jpg"); soteErr.ErrCode != 210400 {
		t.Errorf("CheckIfFilePathExists failed: Expected error code of 210400: %v", soteErr.FmtErrMsg)
	}
}
func TestPreprocessManager_CorrectSkew(t *testing.T) {
	if pm, soteErr := NewPreprocessor("../img/testing_materials/Container Guarantee Form back.jpg"); soteErr.ErrCode == nil {
		if _, soteErr = pm.CorrectSkew(); soteErr.ErrCode != nil {
			t.Errorf("CorrectSkew failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
		}
	}
}
