package sDocument

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	if _, soteError := New(GetTessdataPrefix()); soteError.ErrCode != nil {
		t.Errorf("New failed: Expected error code to be nil")
	}

	if _, soteError := New(""); soteError.ErrCode != 209100 {
		t.Errorf("New failed: Expected error code to be 209100")
	}
}
func TestDocumentManager_GetTextFromDocument(t *testing.T) {
	var stext string
	sfilename := "/Users/fionamurie/Desktop/sote/practice/testing_materials/Container Guarantee Form back.jpg"

	if dm, soteError := New(GetTessdataPrefix()); soteError.ErrCode == nil{
		stext, soteError = dm.GetTextFromDocument(sfilename)
		fmt.Println(stext)
	}


}