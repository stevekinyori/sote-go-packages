package sDocument

import (
	"fmt"
	"testing"
)

func TestNewTesseractServer(t *testing.T) {
	if _, soteError := NewTesseractServer(SGetTessdataPrefix()); soteError.ErrCode != nil {
		t.Errorf("New failed: Expected error code to be %v got %v", "nil", soteError.FmtErrMsg)
	}
}
func TestTesseractServerManager_GetTextFromFile(t *testing.T) {
	filename := "../img/testing_materials/ContainerGuaranteeFormback.jpg"
	var text string

	if tsm, soteErr := NewTesseractServer(SGetTessdataPrefix()); soteErr.ErrCode == nil {
		if text, soteErr = tsm.GetTextFromFile(filename); soteErr.ErrCode != nil {
		} else {
			fmt.Println(text)
		}
	}

	if tsm, soteErr := NewTesseractServer(SGetTessdataPrefix()); soteErr.ErrCode == nil {
		if _, soteErr = tsm.GetTextFromFile(""); soteErr.ErrCode != 209110 {
			t.Errorf("New failed: Expected error code to be %v got  %v", "209110", soteErr.ErrCode)
		}
	}
}
