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
		setupDir          = ""
	)

	tPtr.Cleanup(func() {
	})

	tPtr.Run("test migrate", func(tPtr *testing.T) {
		if soteErr = Migrate(parentCtx, sConfigParams.DEVELOPMENT, setupDir); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})

	tPtr.Run("test run", func(tPtr *testing.T) {
		if soteErr = Run(parentCtx, sConfigParams.DEVELOPMENT, MigrationType, "setup", setupDir); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code %v got %v", testName, "nil", soteErr.FmtErrMsg)
		}
	})
}
