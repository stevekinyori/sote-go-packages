package sMigration

import (
	"context"

	// Add imports here
	"gitlab.com/soteapps/packages/v2022/sError"
	"gitlab.com/soteapps/packages/v2022/sLogger"
)

const (
	LOGMESSAGEPREFIX = "packages/sMigration"
)

func init() {
	sLogger.SetLogMessagePrefix(LOGMESSAGEPREFIX)
}

// Migrate migrates all .go & .sql files withing MigrationsSubDir
// if setupDir is empty, then it takes the directory of the calling file
func Migrate(ctx context.Context, environment string, setupDir string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	soteErr = migrationAndSeeding(ctx, environment, MigrationType, ExternalDefaultStackTraceSkips, setupDir)

	return
}
