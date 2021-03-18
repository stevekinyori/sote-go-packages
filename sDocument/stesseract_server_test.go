package sDocument

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	if _, soteError := New(SGetTessdataPrefix()); soteError.ErrCode != nil {
		t.Errorf("New failed: Expected error code to be nil")
	}

	if _, soteError := New(""); soteError.ErrCode != 209100 {
		t.Errorf("New failed: Expected error code to be 209100")
	}
}
func TestDocumentManager_GetTextFromDocument(t *testing.T) {
	var stext string
	sfilename := "../img/testing_materials/Container Guarantee Form back.jpg"

	if dm, soteError := New(SGetTessdataPrefix()); soteError.ErrCode == nil {
		stext, soteError = dm.GetTextFromFile(sfilename)
		fmt.Println(stext)
	}

}
