package sDatabase

import (
	"testing"

	"gitlab.com/soteapps/packages/v2020/sLogger"
)

const (
	SOTETESTSCHEMA = "sotetest"
	REFERENCETABLE = "referencetable"
	REFTBLCOLUMNNAME = "reference_name"
	PARENTCHILDTABLE = "parentchildtable"
	PCTBLCOLUMNNAME = "reference_name"
	SELF = "self"
)
func init() {
	sLogger.SetLogMessagePrefix("sprimarykeyInfo_test.go")
}

func TestPKPrimer(t *testing.T) {
	if soteErr := GetAWSParams(); soteErr.ErrCode != nil {
		t.Errorf("getAWSParams Failed: Expected error code to be nil.")
		t.Fatal()
	}

	tConnInfo, soteErr := GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("GetConnection Failed: Please investigate")
		t.Fail()
	}

	if pkPrimer("sotetest", tConnInfo); len(pkList) == 0 {
		t.Errorf("pkPrimer Failed: Expected the pkList to have at least one entry.")
	}
}
func TestPKLookup(t *testing.T) {
	if soteErr := GetAWSParams(); soteErr.ErrCode != nil {
		t.Errorf("getAWSParams Failed: Expected error code to be nil.")
		t.Fatal()
	}

	tConnInfo, soteErr := GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
	if soteErr.ErrCode != nil {
		t.Errorf("GetConnection Failed: Please investigate")
		t.Fail()
	}

	if pkPrimer("sotetest", tConnInfo); len(pkList) == 0 {
		t.Errorf("pkPrimer Failed: Expected the pkList to have at least one entry.")
	}

	if tbName, soteErr := PKLookup(SOTETESTSCHEMA, REFERENCETABLE, REFTBLCOLUMNNAME, tConnInfo, false); soteErr.ErrCode != nil && tbName != SELF {
		t.Errorf("pkLookup Failed: Expected error code to be nil and tbName should be self.")
	}

	if tbName, soteErr := PKLookup(SOTETESTSCHEMA, PARENTCHILDTABLE, PCTBLCOLUMNNAME, tConnInfo, false); soteErr.ErrCode != nil && tbName != REFERENCETABLE {
		t.Errorf("pkLookup Failed: Expected error code to be nil and tbName should be referencetable.")
	}
}