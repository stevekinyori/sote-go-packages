package sConfigParams

import (
	"testing"

	"gitlab.com/soteapps/packages/serror"
	"gitlab.com/soteapps/packages/slogger"
)

const (
	tAWS       string = "staging"
	dbPassword string = "password"
)

func TestGetDBPassword(t *testing.T) {
	slogger.DebugMethod()

	if _, found := GetDBPassword(tAWS); !found {
		var p = make([]interface{}, 1)
		p[0] = dbPassword
		serror.GetSError(609999, p, serror.EmptyMap)
		t.Errorf("TestGetDBPassword should have found the password: %v", serror.GetSError(609999, p, serror.EmptyMap).FmtErrMsg)
	}

}
