package sDocument

import (
	"testing"

	"gitlab.com/soteapps/packages/v2020/sLogger"
)

func init() {
	sLogger.SetLogMessagePrefix("sconnection_test.go")
}

func TestSetConnectionValue(t *testing.T) {
	New()

	// if soteErr.ErrCode != 602020 {
	// 	t.Errorf("setConnectionValues Failed: Error code is not for an invalid sslMode.")
	// 	t.Fail()
	// }
	// _, soteErr = setConnectionValues("dbName", "User", "Password", "Host", "disable", 1, 1)
	// if soteErr.ErrCode != nil {
	// 	t.Errorf("setConnectionValues Failed: Expected a nil error code.")
	// 	t.Fail()
	// }
}
