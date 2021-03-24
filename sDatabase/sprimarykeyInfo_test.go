package sDatabase

import (
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

const (
	SOTETESTSCHEMA   = "sotetest"
	REFERENCETABLE   = "referencetable"
	REFTBLCOLUMNNAME = "reference_name"
	PARENTCHILDTABLE = "parentchildtable"
	PCTBLCOLUMNNAME  = "reference_name"
	EMPTYVALUE       = ""
	SELF             = "self"
)

func init() {
	sLogger.SetLogMessagePrefix("sprimarykeyInfo_test.go")
}

// func TestPKPrimer(tPtr *testing.T) {
	// if soteErr := GetAWSParams(); soteErr.ErrCode != nil {
	// 	tPtr.Errorf("getAWSParams Failed: Expected error code to be nil.")
	// 	t.Fatal()
	// }
	//
	// tConnInfo, soteErr := GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
	// if soteErr.ErrCode != nil {
	// 	tPtr.Errorf("GetConnection Failed: Please investigate")
	// 	tPtr.Fail()
	// }
	//
	// if pkPrimer("sotetest", tConnInfo); len(pkList) == 0 {
	// 	tPtr.Errorf("pkPrimer Failed: Expected the pkList to have at least one entry.")
	// }
// }
// func TestPkPrimer(tPtr *testing.T) {
	// if soteErr := GetAWSParams(); soteErr.ErrCode != nil {
	// 	tPtr.Errorf("getAWSParams Failed: Expected error code to be nil.")
	// 	t.Fatal()
	// }
	//
	// tConnInfo, soteErr := GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
	// if soteErr.ErrCode != nil {
	// 	tPtr.Errorf("GetConnection Failed: Please investigate")
	// 	tPtr.Fail()
	// }
	//
	// if pkPrimer("sotetest", tConnInfo); len(pkList) == 0 {
	// 	tPtr.Errorf("pkPrimer Failed: Expected the pkList to have at least one entry.")
	// }
// }
// func TestPKLookup(tPtr *testing.T) {
	// if soteErr := GetAWSParams(); soteErr.ErrCode != nil {
	// 	tPtr.Errorf("getAWSParams Failed: Expected error code to be nil.")
	// 	t.Fatal()
	// }
	//
	// tConnInfo, soteErr := GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
	// if soteErr.ErrCode != nil {
	// 	tPtr.Errorf("GetConnection Failed: Please investigate")
	// 	tPtr.Fail()
	// }
	//
	// if tbName, soteErr := PKLookup(SOTETESTSCHEMA, REFERENCETABLE, REFTBLCOLUMNNAME, tConnInfo); soteErr.ErrCode != nil && tbName != SELF {
	// 	tPtr.Errorf("pkLookup Failed: Expected error code to be nil and tbName should be self.")
	// }
	//
	// if tbName, soteErr := PKLookup(SOTETESTSCHEMA, PARENTCHILDTABLE, PCTBLCOLUMNNAME, tConnInfo); soteErr.ErrCode != nil && tbName != REFERENCETABLE {
	// 	tPtr.Errorf("pkLookup Failed: Expected error code to be nil and tbName should be referencetable.")
	// }
// }
// func TestPKLookupEmptyValues(tPtr *testing.T) {
	// if soteErr := GetAWSParams(); soteErr.ErrCode != nil {
	// 	tPtr.Errorf("getAWSParams Failed: Expected error code to be nil.")
	// 	t.Fatal()
	// }
	//
	// tConnInfo, soteErr := GetConnection(DBName, DBUser, DBPassword, DBHost, DBSSLMode, DBPort, 3)
	// if soteErr.ErrCode != nil {
	// 	tPtr.Errorf("GetConnection Failed: Please investigate")
	// 	tPtr.Fail()
	// }
	//
	// if _, soteErr := PKLookup(EMPTYVALUE, PARENTCHILDTABLE, PCTBLCOLUMNNAME, tConnInfo); soteErr.ErrCode != 200513 {
	// 	tPtr.Errorf("pkLookup Failed: Expected error code to be 200513.")
	// }
	// if _, soteErr := PKLookup(SOTETESTSCHEMA, EMPTYVALUE, PCTBLCOLUMNNAME, tConnInfo); soteErr.ErrCode != 200513 {
	// 	tPtr.Errorf("pkLookup Failed: Expected error code to be 200513.")
	// }
	// if _, soteErr := PKLookup(SOTETESTSCHEMA, PARENTCHILDTABLE, EMPTYVALUE, tConnInfo); soteErr.ErrCode != 200513 {
	// 	tPtr.Errorf("pkLookup Failed: Expected error code to be 200513.")
	// }
// }
// func TestNoConnection(tPtr *testing.T) {
// 	var tConnInfo = ConnInfo{nil, ConnValues{}}
// 	if _, soteErr := PKLookup(SOTETESTSCHEMA, PARENTCHILDTABLE, PCTBLCOLUMNNAME, tConnInfo); soteErr.ErrCode != 209299 {
// 		tPtr.Errorf("pkLookup Failed: Expected error code to be 209299.")
// 	}
// }
