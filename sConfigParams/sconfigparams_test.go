package sConfigParams

import (
	"testing"

	"gitlab.com/soteapps/packages/v2020/sLogger"
)

const (
	tAWS       string = "staging"
	dbPassword string = "password"
)

func TestGetDBPassword(t *testing.T) {
	sLogger.DebugMethod()

	if _, found := GetDBPassword(tAWS); !found {
		var p = make([]interface{}, 1)
		p[0] = dbPassword
		t.Errorf("TestGetDBPassword should have found the password: %v", "asdf")
	}
	if _, found := GetDBPassword(tAWS); !found {
		var p = make([]interface{}, 1)
		p[0] = dbPassword
		t.Errorf("TestGetDBPassword should have found the password: %v", "asdf")
	}

}
