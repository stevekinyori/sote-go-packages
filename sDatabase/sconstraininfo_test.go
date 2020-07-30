package sDatabase

import (
	"testing"
)

func TestGetSingleColumnConstraintInfo(t *testing.T) {
	if soteErr = GetAWSParams(); soteErr.ErrCode != nil {
		t.Errorf("getAWSParams Failed: Expected error code to be nil.")
		t.Fatal()
	}

	tConnInfo, soteErr = GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("GetConnection Failed: Please investigate")
		t.Fail()
	}

	var myConstraints []SConstraint
	if myConstraints, soteErr = GetSingleColumnConstraintInfo("sote", tConnInfo); len(myConstraints) == 0 {
		t.Errorf("GetSingleColumnConstraintInfo Failed: myContraints should not be empty")
	}
}
