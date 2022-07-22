package sDatabase

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"gitlab.com/soteapps/packages/v2022/sError"
	"gitlab.com/soteapps/packages/v2022/sLogger"
)

// FilterCommon describes the  format for common filter expression in the request json payload
type FilterCommon struct {
	Operator string      `json:"operator" mod:"ucase" validate:"required,oneof=> < = >= <= != <> IN 'NOT IN'"`
	Value    interface{} `json:"value" validate:"required_unless=Operator != Operator ="`
}

// FilterFields describes the  format for all filter expression in the request json payload
type FilterFields struct {
	FilterCommon
	FieldName string `json:"field-name"`
}

// ArrFilterParam filter parameters for a slice/array
type ArrFilterParam struct {
	FilterCommon
	FieldName         string
	Prefix            string
	InitialParamCount int
	CaseInsensitive   bool
}

type ArrFilterResponse struct {
	QueryStr   string
	Params     []interface{}
	ParamCount int
}

type FormatConditionParams struct {
	InitialParamCount int
	RecordLimitCount  int
	TblPrefixes       []string // e.g. tbl.'the prefix must have a dot at the end'
	SortOrderStr      string
	ColName           string
	Operator          string
	Filters           map[string][]FilterFields
	SortOrderKeysMap  map[string]SortOrder
}

type SortOrder struct {
	ColumnName      string
	CaseInsensitive bool
}

type FormatConditionsResp struct {
	Where      string
	Limit      string
	Order      string
	Params     []interface{}
	ParamCount int
}

// FormatArrayFilterCondition formats slice/array filter conditions for a get/list request
func FormatArrayFilterCondition(ctx context.Context, sortOrderKeysMap map[string]SortOrder,
	reqParams *ArrFilterParam) (arrFilterResp *ArrFilterResponse, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		paramStart string
		paramEnd   string
	)

	arrFilterResp = &ArrFilterResponse{}
	s := reflect.ValueOf(reqParams.Value)
	if kind := s.Kind(); kind == reflect.Slice || kind == reflect.Array {
		reqParamLen := s.Len()
		if reqParamLen > 0 {
			arrFilterResp.ParamCount = reqParams.InitialParamCount
			if reqParams.CaseInsensitive {
				paramStart = "UPPER("
				paramEnd = ")"
			}

			arrFilterResp.QueryStr = fmt.Sprintf(" %v%v%v%v %v (", paramStart, reqParams.Prefix, sortOrderKeysMap[reqParams.FieldName].ColumnName,
				paramEnd,
				reqParams.Operator)
			arrFilterResp.Params = make([]interface{}, reqParamLen)

			for i := 0; i < reqParamLen; i++ {
				arrFilterResp.ParamCount++
				arrFilterResp.QueryStr += fmt.Sprintf("%v$%v%v,", paramStart, arrFilterResp.ParamCount, paramEnd)
				arrFilterResp.Params[i] = s.Index(i).Interface()
			}

			arrFilterResp.QueryStr = strings.TrimSuffix(arrFilterResp.QueryStr, ",") + ")"
		}
	} else {
		soteErr = sError.GetSError(207030, sError.BuildParams([]string{reqParams.FieldName, fmt.Sprint(reqParams.Value)}), sError.EmptyMap)
	}

	return
}

func formatFilterCondition(ctx context.Context, fmtConditionParams *FormatConditionParams) (fmtConditionResp FormatConditionsResp,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		col           string
		val           string
		join          string
		paramCount    int
		queryStr      string
		tQueryStr     string
		params        []interface{}
		arrFilterResp = &ArrFilterResponse{}
		prefix        string
	)
	if len(fmtConditionParams.Filters) > 0 {
		paramCount = fmtConditionParams.InitialParamCount
		join = " AND "
		if fmtConditionParams.InitialParamCount > 0 {
			queryStr = join
		} else {
			queryStr = " WHERE "
		}
		if len(fmtConditionParams.TblPrefixes) > 0 {
			prefix = fmtConditionParams.TblPrefixes[0]
			fmtConditionParams.TblPrefixes = fmtConditionParams.TblPrefixes[1:]
		}

	firstLoop:
		for operand, filterValues := range fmtConditionParams.Filters {
			tQueryStr += "("
			for _, field := range filterValues {
				if field.Operator == "IN" || field.Operator == "NOT IN" {
					if arrFilterResp, soteErr = FormatArrayFilterCondition(ctx, fmtConditionParams.SortOrderKeysMap, &ArrFilterParam{
						FieldName:         field.FieldName,
						FilterCommon:      FilterCommon{Operator: field.Operator, Value: field.Value},
						Prefix:            prefix,
						InitialParamCount: paramCount,
						CaseInsensitive:   fmtConditionParams.SortOrderKeysMap[field.FieldName].CaseInsensitive,
					}); soteErr.ErrCode == nil {
						tQueryStr += arrFilterResp.QueryStr + " " + operand
						params = append(params, arrFilterResp.Params...)
						paramCount = arrFilterResp.ParamCount
					} else {
						break firstLoop
					}
				} else {
					paramCount++
					if fmtConditionParams.SortOrderKeysMap[field.FieldName].CaseInsensitive {
						col = fmt.Sprintf("UPPER(%v%v)", fmtConditionParams.TblPrefixes,
							fmtConditionParams.SortOrderKeysMap[field.FieldName].ColumnName)
						val = fmt.Sprintf("UPPER($%v)", paramCount)
					} else {
						col = prefix + fmtConditionParams.SortOrderKeysMap[field.FieldName].ColumnName
						val = fmt.Sprintf("$%v", paramCount)
					}

					// filter by is not null or is null
					if field.Value == nil {
						subQuery := "NULL"
						switch field.Operator {
						case "=":
							subQuery = "IS " + subQuery
						case "!=":
							subQuery = "IS NOT " + subQuery
						}

						tQueryStr += fmt.Sprintf(" %v %v %v", col, subQuery, operand)
					} else {
						tQueryStr += fmt.Sprintf(" %v %v %v %v", col, field.Operator, val, operand)
						params = append(params, field.Value)
					}
				}
			}
			tQueryStr = fmt.Sprintf("%v)%v", strings.TrimSuffix(tQueryStr, operand), join)
		}

		if soteErr.ErrCode == nil {
			ttQueryStr := strings.TrimSuffix(tQueryStr, join)
			for _, p := range fmtConditionParams.TblPrefixes {
				ttQueryStr += " OR " + strings.ReplaceAll(ttQueryStr, prefix, p)
			}

			queryStr += "(" + ttQueryStr + ")"
		}

		fmtConditionResp.Where = queryStr
		fmtConditionResp.Params = params
		fmtConditionResp.ParamCount = paramCount
	}

	return
}

// FormatListQueryConditions parses the query list for a /list endpoints and list nats action types to form relevant sql queries
func FormatListQueryConditions(ctx context.Context, fmtConditionParams *FormatConditionParams) (fmtConditionResp FormatConditionsResp,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		wg             sync.WaitGroup
		limitChan      = make(chan string)
		whereChan      = make(chan string)
		orderChan      = make(chan string)
		paramsChan     = make(chan []interface{})
		paramCountChan = make(chan int, 1)
		soteErrChan    = make(chan sError.SoteError, 1)
	)

	runtime.GOMAXPROCS(runtime.NumCPU())
	wg.Add(3)

	go func() {
		wg.Wait()
		close(limitChan)
		close(orderChan)
		close(whereChan)
		close(paramsChan)
		close(paramCountChan)
	}()

	go func() {
		defer wg.Done()
		if fmtConditionParams.RecordLimitCount > 0 {
			limitChan <- fmt.Sprintf("LIMIT %v ", fmtConditionParams.RecordLimitCount)
		} else {
			limitChan <- ""
		}
	}()

	go func() {
		defer wg.Done()

		if fmtConditionParams.SortOrderStr == "" {
			orderChan <- ""
		} else {
			orderChan <- fmt.Sprintf("ORDER BY %v%v ", fmtConditionParams.TblPrefixes, fmtConditionParams.SortOrderStr)
		}
	}()

	go func() {
		defer wg.Done()
		var (
			tSoteErr          sError.SoteError
			tFmtConditionResp FormatConditionsResp
		)

		if len(fmtConditionParams.Filters) > 0 {
			if tFmtConditionResp, tSoteErr = formatFilterCondition(ctx, fmtConditionParams); tSoteErr.ErrCode != nil {
				soteErrChan <- tSoteErr
				whereChan <- tFmtConditionResp.Where // where clause string
				paramsChan <- tFmtConditionResp.Params
				paramCountChan <- tFmtConditionResp.ParamCount
				return
			}
			soteErrChan <- tSoteErr
			whereChan <- tFmtConditionResp.Where // where clause string
			paramsChan <- tFmtConditionResp.Params
			paramCountChan <- tFmtConditionResp.ParamCount
		} else {
			soteErrChan <- tSoteErr
			whereChan <- ""
			paramsChan <- []interface{}{}
			paramCountChan <- fmtConditionParams.InitialParamCount
		}
	}()

	fmtConditionResp.Limit = <-limitChan
	fmtConditionResp.Where = <-whereChan
	fmtConditionResp.Order = <-orderChan
	fmtConditionResp.Params = <-paramsChan
	fmtConditionResp.ParamCount = <-paramCountChan
	soteErr = <-soteErrChan

	return
}

// FormatGenericFilterArray formats params from slice/array for additional filters that are not supported by the filters list. (i.e for summary endpoints)
func FormatGenericFilterArray(ctx context.Context, fmtConditionParams *FormatConditionParams, args []string) (queryStr string, params []interface{},
	paramCount int) {
	paramStart := "UPPER("
	paramEnd := ")"

	reqParamLen := len(args)

	// we are being paranoid here! We already know that this function will never be called with an empty slice:-)
	if reqParamLen > 0 {
		paramCount = fmtConditionParams.InitialParamCount
		queryStr = fmt.Sprintf(" %v%v%v%v %v (", paramStart, fmtConditionParams.TblPrefixes, fmtConditionParams.ColName, paramEnd,
			fmtConditionParams.Operator)

		params = make([]interface{}, reqParamLen)
		for i := 0; i < reqParamLen; i++ {
			paramCount++
			queryStr += fmt.Sprintf("%v$%v%v,", paramStart, paramCount, paramEnd)
			params[i] = args[i]
		}
		queryStr = strings.TrimSuffix(queryStr, ",") + ")"
	}
	return
}
