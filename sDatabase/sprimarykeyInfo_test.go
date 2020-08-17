package sDatabase

import (
	"fmt"
	"testing"
)

func TestPKPrimer(t *testing.T) {
	if soteErr = GetAWSParams(); soteErr.ErrCode != nil {
		t.Errorf("getAWSParams Failed: Expected error code to be nil.")
		t.Fatal()
	}

	tConnInfo, soteErr = GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("GetConnection Failed: Please investigate")
		t.Fail()
	}

	if pkPrimer("sotetest", tConnInfo); len(pkList) == 0 {
		t.Errorf("pkPrimer Failed: Expected the pkList to have at least one entry.")
	}
}
func TestPKLookup(t *testing.T) {
	if soteErr = GetAWSParams(); soteErr.ErrCode != nil {
		t.Errorf("getAWSParams Failed: Expected error code to be nil.")
		t.Fatal()
	}

	tConnInfo, soteErr = GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("GetConnection Failed: Please investigate")
		t.Fail()
	}

	if pkPrimer("sotetest", tConnInfo); len(pkList) == 0 {
		t.Errorf("pkPrimer Failed: Expected the pkList to have at least one entry.")
	}

	if tbName, soteErr := pkLookup("sotetest", "referencetable", "reference_name", tConnInfo, false); soteErr.ErrCode != nil {
		t.Errorf("pkLookup Failed: Expected error code to be nil.")
	} else {
		fmt.Println(tbName)
	}
}