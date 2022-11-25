package seeds

import (
	"context"

	"gitlab.com/soteapps/packages/v2022/sDatabase"
	"gitlab.com/soteapps/packages/v2022/sError"
)

func (config Config) SoteStudentTestData(ctx context.Context) (soteErr sError.SoteError) {
	var (
		tableName = "sote_student_test"
		data      = []map[string]any{
			{
				"name":  "Mary N M",
				"age":   17,
				"class": "12th Grade",
			},
			{
				"name":  "Fiona Q M",
				"age":   14,
				"class": "8th Grade",
			},
		}
	)

	soteErr = config.DBConnInfo.InsertTableData(ctx, tableName, data, sDatabase.SeedingType)

	return
}
