package sDatabase

import (
	"context"
	"runtime"
	"testing"
)

func TestFormatArrayFilterCondition(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)

	tPtr.Run("multiple prefixes", func(tPtr *testing.T) {
		if _, soteErr := FormatListQueryConditions(context.Background(), &FormatConditionParams{
			InitialParamCount: 1,
			RecordLimitCount:  0,
			TblPrefixes:       []string{"table1.", "table2.", "table3."},
			Filters: map[string][]FilterFields{
				"AND": {
					FilterFields{
						FilterCommon: FilterCommon{
							Operator: "IN",
							Value:    []string{"name1", "name2"},
						},
						FieldName: "column-name",
					},
					FilterFields{
						FilterCommon: FilterCommon{
							Operator: "=",
							Value:    1,
						},
						FieldName: "column-id",
					},
				},
			},
			SortOrderKeysMap: map[string]TableColumn{
				"column-name": {
					ColumnName:      "column_name",
					CaseInsensitive: true,
				},
				"column-id": {
					ColumnName:      "column_id",
					CaseInsensitive: false,
				},
			},
			SortOrder: SortOrder{
				TblPrefix: "table.",
				Fields:    map[string]string{"column-name": "DESC", "column-id": "DESC"},
			},
		}); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
		}
	})

}

func TestSetSortOrder(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function).Name()
	)
	tPtr.Run("order by", func(tPtr *testing.T) {
		if sortOrderStr := setSortOrder(SortOrder{
			TblPrefix: "table.",
			Fields:    map[string]string{"column-name": "DESC", "column-id": "DESC"},
		}, map[string]TableColumn{
			"column-name": {
				ColumnName:      "column_name",
				CaseInsensitive: true,
			},
			"column-id": {
				ColumnName:      "column_id",
				CaseInsensitive: false,
			},
		}); sortOrderStr == "" {
			tPtr.Errorf("%v Failed: Expected sort order string not to be empty", testName)
		}
	})
}
