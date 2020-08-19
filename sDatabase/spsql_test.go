package sDatabase

import (
	"testing"

	"gitlab.com/soteapps/packages/v2020/sLogger"
)

func init() {
	sLogger.SetLogMessagePrefix("spsql_test.go")
}

func TestExecutePSQLWithFile(t *testing.T) {
	if soteErr := ExecutePSQLWithFile("empty.file"); soteErr.ErrCode != nil {
		t.Errorf("ExecutePSQLWithFile Failed: Expected error code of nil")
		t.Fail()
	}
}
