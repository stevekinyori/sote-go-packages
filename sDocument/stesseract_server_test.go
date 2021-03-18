package sDocument

import "testing"

//func TestGetTessdataLocation(t *testing.T) {
//	if st := GetTessdataLocation(); st == ""{
//		fmt.Println(st)
//	}
//}

func TestNew(t *testing.T) {
	if _, soteErr := New(); soteErr.ErrCode != nil{
		t.Errorf("New failed: Expected error code to be nil")
	}
}