package sDatabase

import (
	"testing"

	"gitlab.com/soteapps/packages/v2020/sLogger"
)

func init() {
	sLogger.SetLogMessagePrefix("sconstraininfo_test.go")
}

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

	// TODO This should be len(myConstraints) == 0.  When the sotetest data structures are installed using the test.sh, this must be changed
	var myConstraints []SConstraint
	if myConstraints, soteErr = GetSingleColumnConstraintInfo(TESTSCHEMA, tConnInfo); len(myConstraints) > 0 {
		t.Errorf("GetSingleColumnConstraintInfo Failed: myContraints should be empty")
	}
}
