package sDocument

import (
	"os"
	"testing"
)

func TestConvertImageFormat(tPtr *testing.T) {
	if _, soteErr := SConvertImageFormat("../img/testing_materials/Container Guarantee Form back.jpg",
		"../img/testing_materials/out.pdf"); soteErr.ErrCode != nil {
		tPtr.Errorf("SConvertImageFormat failed:Expected soteErr to be nil:%v ", soteErr.FmtErrMsg)
	}

	if _, soteErr := SConvertImageFormat("img/testing_materials/Container Guarantee Form back.jpg",
		"../img/testing_materials/out.pdf"); soteErr.ErrCode != 199999 {
		tPtr.Errorf("SConvertImageFormat failed:Expected error code of 199999:%v ", soteErr.FmtErrMsg)
	}

	os.Remove("../img/testing_materials/out.pdf")
}
