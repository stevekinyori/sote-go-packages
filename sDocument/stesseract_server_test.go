package sDocument

import (
	"testing"
)

func TestNewTesseractServer(t *testing.T) {
	if _, soteError := NewTesseractServer(SGetTessdataPrefix()); soteError.ErrCode != nil {
		t.Errorf("New failed: Expected error code to be nil")
	}

	if _, soteError := NewTesseractServer(""); soteError.ErrCode != 209100 {
		t.Errorf("New failed: Expected error code to be 209100")
	}
}
func TestTesseractServerManager_GetTextFromFile(t *testing.T) {
	sfilename := "../img/testing_materials/Container Guarantee Form back.jpg"

	if tsm, soteErr := NewTesseractServer(SGetTessdataPrefix()); soteErr.ErrCode == nil {
		if _, soteErr = tsm.GetTextFromFile(sfilename); soteErr.ErrCode != nil {

		}
	}
	if tsm, soteErr := NewTesseractServer(""); soteErr.ErrCode == nil {
		if _, soteErr = tsm.GetTextFromFile(sfilename); soteErr.ErrCode != 209100 {
			t.Errorf("New failed: Expected error code to be 209100 %v", soteErr.ErrCode)
		}
	}

}
