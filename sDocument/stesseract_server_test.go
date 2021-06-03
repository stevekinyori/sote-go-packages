package sDocument

import (
	"testing"

	"gitlab.com/soteapps/packages/v2021/sError"
)

func TestNewTesseractServer(t *testing.T) {
	var (
		soteErr        sError.SoteError
		tessdataPrefix string
	)

	if tessdataPrefix, soteErr = GetTessdataPrefix(); soteErr.ErrCode != nil {
		t.Errorf("NewTesseractServer failed: Expected error code to be %v got %v", "nil", soteErr.FmtErrMsg)
	}

	if soteErr.ErrCode == nil {
		if _, soteErr = NewTesseractServer(tessdataPrefix); soteErr.ErrCode != nil {
			t.Errorf("NewTesseractServer failed: Expected error code to be %v got %v", "nil", soteErr.FmtErrMsg)
		}
	}
}
func TestTesseractServerManager_GetTextFromFile(t *testing.T) {

	var (
		soteErr        sError.SoteError
		tessdataPrefix string
		tsm            *TesseractServerManager
	)
	// filename := "../img/testing_materials/ContainerGuaranteeFormback.jpg"
	filename := "/Users/fionamurie/Desktop/sote/golang/src/packages/img/testing_materials/Invoice.jpeg"

	if tessdataPrefix, soteErr = GetTessdataPrefix(); soteErr.ErrCode != nil {
		t.Errorf("GetTextFromFile failed: Expected error code to be %v got %v", "nil", soteErr.FmtErrMsg)
	}

	if tsm, soteErr = NewTesseractServer(tessdataPrefix); soteErr.ErrCode == nil {
		if _, soteErr = tsm.GetTextFromFile(filename); soteErr.ErrCode != nil {
			t.Errorf("GetTextFromFile failed: Expected error code to be %v got %v", "nil", soteErr.FmtErrMsg)
		}
	}

	if tsm, soteErr = NewTesseractServer(tessdataPrefix); soteErr.ErrCode == nil {
		if _, soteErr = tsm.GetTextFromFile(""); soteErr.ErrCode != 209110 {
			t.Errorf("GetTextFromFile failed: Expected error code to %v got  %v", "209110", soteErr.FmtErrMsg)
		}
	}

}
