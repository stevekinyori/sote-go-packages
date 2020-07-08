package sDatabase

import (
	"testing"
)

func TestGetColumnConstraintInfo(t *testing.T) {
	soteErr := GetConnection("single", "sote_development", "sote", "password", "localhost", "disable", 5442, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("Get Connection Failed: Should have returned a pointer to the single database connection")
		t.Fail()
	}

	var myConstraints []sConstraint
	if myConstraints, soteErr = getColumnConstraintInfo("sote"); len(myConstraints) == 0 {
		t.Errorf("Test Get Column Constraint Info Failed: myContraints should not be empty")
	}
}
