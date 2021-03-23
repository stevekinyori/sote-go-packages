package sDocument

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	if _, soteError := NewTesseractServer(SGetTessdataPrefix()); soteError.ErrCode != nil {
		t.Errorf("New failed: Expected error code to be nil")
	}

	if _, soteError := NewTesseractServer(""); soteError.ErrCode != 209100 {
		t.Errorf("New failed: Expected error code to be 209100")
	}
}
func TestTesseractServerManager_GetTextFromFile(t *testing.T) {
	var stext string
	sfilename := "../img/testing_materials/Container Guarantee Form back.jpg"

	if tsm, soteError := NewTesseractServer(SGetTessdataPrefix()); soteError.ErrCode == nil {
		stext, soteError = tsm.GetTextFromFile(sfilename)
		fmt.Println(stext)
	}

}
