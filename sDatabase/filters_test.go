package sDatabase

import (
	"context"
	"fmt"
	"runtime"
	"testing"
)

func TestFormatArrayFilterCondition(tPtr *testing.T) {
	var (
		function, _, _, _ = runtime.Caller(0)
		testName          = runtime.FuncForPC(function)
	)

	tPtr.Run("multiple prefixes", func(tPtr *testing.T) {
		if resp, soteErr := FormatFilterCondition(context.Background(), &FormatConditionParams{
			InitialParamCount: 1,
			RecordLimitCount:  0,
			TblPrefixes:       []string{"table1.", "table2.", "table3."},
			SortOrderStr:      "",
			ColName:           "",
			Operator:          "=",
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
			SortOrderKeysMap: map[string]SortOrder{
				"column-name": {
					ColumnName:      "column_name",
					CaseInsensitive: true,
				},
				"column-id": {
					ColumnName:      "column_id",
					CaseInsensitive: false,
				},
			},
		}); soteErr.ErrCode != nil {
			tPtr.Errorf("%v Failed: Expected error code to be nil got %v", testName, soteErr.FmtErrMsg)
		} else {
			fmt.Println(resp.Where)
		}
	})

}
