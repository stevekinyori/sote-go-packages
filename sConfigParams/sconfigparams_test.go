package sConfigParams

import (
	"testing"

	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
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
	if _, found := GetDBPassword(tAWS); !found {
		var p = make([]interface{}, 1)
		p[0] = dbPassword
		serror.GetSError(609999, p, serror.EmptyMap)
		t.Errorf("TestGetDBPassword should have found the password: %v", serror.GetSError(609999, p, serror.EmptyMap).FmtErrMsg)
	}

}
