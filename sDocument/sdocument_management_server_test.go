package sDocument

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	if _, soteErr := New(); soteErr.ErrCode != nil {
		t.Errorf("New failed:Expected soteErr to be nil:%v ", soteErr.FmtErrMsg)
	}
}
func TestDocumentServerManager_ConvertImageToPDF(t *testing.T) {
	if dsm, soteErr := New(); soteErr.ErrCode == nil {
		if _, soteErr = dsm.ConvertImageFormat("../img/testing_materials/Container Guarantee Form back.jpg",
			"../img/testing_materials/out.png"); soteErr.ErrCode != nil {
			t.Errorf("ConvertImageToPDF failed:Expected soteErr to be nil:%v ", soteErr.FmtErrMsg)
		}
		_ = os.Remove("../img/testing_materials/out.png")

	}
}
