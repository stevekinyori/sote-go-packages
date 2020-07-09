package sDatabase

import (
	"testing"
)

func TestGetSingleColumnConstraintInfo(t *testing.T) {
	tConnInfo, soteErr := GetConnection("sote_development", "sote", "password", "localhost", "disable", 5442, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("GetConnection Failed: Please investigate")
		t.Fail()
	}

	var myConstraints []sConstraint
	if myConstraints, soteErr = GetSingleColumnConstraintInfo("sote", tConnInfo); len(myConstraints) == 0 {
		t.Errorf("GetSingleColumnConstraintInfo Failed: myContraints should not be empty")
	}
}
