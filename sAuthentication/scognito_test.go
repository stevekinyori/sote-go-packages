package sAuthentication

import (
	"context"
	"testing"
)

func TestInitiateAuth(tPtr *testing.T) {
	resp, sError := InitiateAuth(context.TODO(), "staging")
	if sError.ErrCode != nil {
		tPtr.Fatalf("Something terrible happened %v", sError.FmtErrMsg)
	}

	tPtr.Logf("Response %v", resp)
}
