package sMigration

import (
	"context"

	// Add imports here
	"gitlab.com/soteapps/packages/v2022/sError"
	"gitlab.com/soteapps/packages/v2022/sLogger"
)

func init() {
	sLogger.SetLogMessagePrefix(LOGMESSAGEPREFIX)
}

// Seed seeds all .go & .sql files withing SeedsSubDir
func Seed(ctx context.Context, environment string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	soteErr = migrationAndSeeding(ctx, environment, SeedingType)

	return
}
