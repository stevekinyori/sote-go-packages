package sDocument

import (
	"testing"
)

func TestNewPreprocessor(t *testing.T) {
	if _, soteErr := NewPreprocessor("../img/testing_materials/Container Guarantee Form back.jpg"); soteErr.ErrCode != nil {
		t.Errorf("NewPreprocessor failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}

	if _, soteErr := NewPreprocessor("img/testing_materials/Container Guarantee Form back.jpg"); soteErr.ErrCode != 199999 {
		t.Errorf("NewPreprocessor failed: Expected error code of %v got  %v", "199999", soteErr.FmtErrMsg)
	}
}
func TestCheckIfPathExists(t *testing.T) {
	if _, soteErr := CheckIfPathExists("../img/testing_materials/Container Guarantee Form back.jpg"); soteErr.ErrCode != nil {
		t.Errorf("CheckIfFilePathExists failed: Expected soteErr to be nil:%v",     soteErr.FmtErrMsg)
	}

	if _, soteErr := CheckIfPathExists("img/testing_materials/Container Guarantee Form back.jpg"); soteErr.ErrCode != 199999 {
		t.Errorf("CheckIfFilePathExists failed: Expected error code of %v got  %v", "199999", soteErr.FmtErrMsg)
	}
}
func TestPreprocessManager_CorrectSkew(t *testing.T) {
	if pm, soteErr := NewPreprocessor("../img/testing_materials/Delivery Order.jpg"); soteErr.ErrCode == nil {
		if _, soteErr = pm.CorrectSkew(); soteErr.ErrCode != nil {
			t.Errorf("CorrectSkew failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
		}
	}
}
