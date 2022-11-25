package sMigration

import (
	"context"
	"runtime"
	"testing"

	"gitlab.com/soteapps/packages/v2022/sConfigParams"
	"gitlab.com/soteapps/packages/v2022/sError"
)

var (
	parentCtx = context.Background()
)

func TestMigration(tPtr *testing.T) {
	var (
		soteErr           sError.SoteError
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	tPtr.Cleanup(func() {
	})

	if soteErr = Migrate(parentCtx, sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		tPtr.Errorf("%v Failed: Expected error code %v got %v", testName, "nil", soteErr.FmtErrMsg)
	}
}
