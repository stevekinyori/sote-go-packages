package migrations

import (
	"context"
	"log"

	"gitlab.com/soteapps/packages/v2022/sDatabase"
	"gitlab.com/soteapps/packages/v2022/sError"
	"gitlab.com/soteapps/packages/v2022/sLogger"
)

func (config Config) DropTableSoteStudentTesting(ctx context.Context) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	config.DBConnInfo.DropTable(ctx, "sote_student_test", sDatabase.MigrationType)
	if soteErr.ErrCode != nil {
		log.Fatalln(soteErr.FmtErrMsg) // terminate migration process
	}

	return
}
