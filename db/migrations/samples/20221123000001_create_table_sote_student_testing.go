package migrations

import (
	"context"

	"gitlab.com/soteapps/packages/v2022/sDatabase"
	"gitlab.com/soteapps/packages/v2022/sError"
)

func (config Config) CreateTableSoteStudentTesting(ctx context.Context) (soteErr sError.SoteError) {
	var (
		tableName = "sote_student_test"
		exists    bool
	)

	if exists, soteErr = config.DBConnInfo.HasTable(ctx, tableName, sDatabase.MigrationType); exists {
		return
	}

	table := sDatabase.TableInfo{
		Name: tableName,
		PrimaryKey: &sDatabase.PrimaryKeyInfo{
			Columns: []string{"sote_student_test_id"},
			AutoIncrementInfo: &sDatabase.AutoIncrementInfo{
				IsAutoIncrement:       true,
				AutoIncrementStartBy:  10000,
				AutoIncrementInterval: 3,
			},
			Description: "This is the unique identifier for each student",
		},
		Description: nil,
	}

	columns := []sDatabase.ColumnInfo{
		{
			Name:        "name",
			DataType:    sDatabase.TEXT,
			IsNullable:  false,
			Description: "This is the name of the student",
		},
		{
			Name:        "age",
			DataType:    sDatabase.INTEGER,
			IsNullable:  false,
			Description: "This is the age of the student",
		},
		{
			Name:        "class",
			DataType:    sDatabase.TEXT,
			IsNullable:  false,
			Description: "This is the class or grade of the student",
		},
	}

	soteErr = config.DBConnInfo.CreateTable(ctx, table, columns, sDatabase.MigrationType)
	if soteErr.ErrCode != nil {
		config.DBConnInfo.DropTable(ctx, "sote_student_test", sDatabase.MigrationType) // rollback
	}

	return
}
