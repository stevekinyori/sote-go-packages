package sHelper

import (
	"context"
	"errors"
	"fmt"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

type panicServiceFunc func(context.Context, sError.SoteError) sError.SoteError

type FormatConditionParams struct {
	InitialParamCount int
	RecordLimitCount  int
	TblPrefix         string //e.g. tbl.'the prefix must have a dot at the end'
	SortOrderStr      string
	ColName           string
	Operator          string
	Filters           map[string][]FilterFields
	SortOrderKeysMap  map[string]map[string]interface{}
}

type FormatConditionsResp struct {
	Where      string
	Limit      string
	Order      string
	Params     []interface{}
	ParamCount int
}

//FormatListQueryConditions returns a formatted list of all query params parsed in from the request payload.
func FormatListQueryConditions(ctx context.Context, fmtConditionParams *FormatConditionParams, panicService panicServiceFunc) (fmtConditionResp FormatConditionsResp, soteErr sError.SoteError) {
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
			orderChan <- fmt.Sprintf("ORDER BY shipmentfintrans.%v ", fmtConditionParams.SortOrderStr)
		}
	}()

	go func() {
		defer wg.Done()
		var (
			tWhere      string
			tParams     []interface{}
			tParamCount int
			tSoteErr    sError.SoteError
		)

		if len(fmtConditionParams.Filters) > 0 {
			if tWhere, tParams, tParamCount, tSoteErr = FormatFilterCondition(ctx, fmtConditionParams, panicService); tSoteErr.ErrCode != nil {
				soteErrChan <- tSoteErr
				whereChan <- tWhere
				paramsChan <- tParams
				paramCountChan <- tParamCount
				return
			}
			soteErrChan <- tSoteErr
			whereChan <- tWhere
			paramsChan <- tParams
			paramCountChan <- tParamCount
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

func FormatSummaryListQueryConditions(ctx context.Context, fmtConditionParams *FormatConditionParams, panicService panicServiceFunc) (where string, params []interface{}, paramCount int, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		wg sync.WaitGroup
		// limitChan      = make(chan string)
		whereChan = make(chan string)
		// orderChan      = make(chan string)
		paramsChan     = make(chan []interface{})
		paramCountChan = make(chan int, 1)
		soteErrChan    = make(chan sError.SoteError, 1)
	)

	runtime.GOMAXPROCS(runtime.NumCPU())
	wg.Add(1)

	go func() {
		wg.Wait()
		close(whereChan)
		close(paramsChan)
		close(paramCountChan)
	}()

	go func() {
		defer wg.Done()
		var (
			tWhere      string
			tParams     []interface{}
			tParamCount int
			tSoteErr    sError.SoteError
		)

		if len(fmtConditionParams.Filters) > 0 {
			if tWhere, tParams, tParamCount, tSoteErr = FormatFilterCondition(ctx, fmtConditionParams, panicService); tSoteErr.ErrCode != nil {
				soteErrChan <- tSoteErr
				whereChan <- tWhere
				paramsChan <- tParams
				paramCountChan <- tParamCount
				return
			}
			soteErrChan <- tSoteErr
			whereChan <- tWhere
			paramsChan <- tParams
			paramCountChan <- tParamCount
		} else {
			soteErrChan <- tSoteErr
			whereChan <- ""
			paramsChan <- []interface{}{}
			paramCountChan <- fmtConditionParams.InitialParamCount
		}
	}()

	where = <-whereChan
	params = <-paramsChan
	paramCount = <-paramCountChan
	soteErr = <-soteErrChan

	return
}

func FormatFilterCondition(ctx context.Context, fmtConditionParams *FormatConditionParams, panicService panicServiceFunc) (queryStr string, params []interface{},
	paramCount int, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		col         string
		val         string
		join        string
		tQueryStr   string
		tParams     []interface{}
		tParamCount int
	)
	if len(fmtConditionParams.Filters) > 0 {
		paramCount = fmtConditionParams.InitialParamCount
		join = " AND "
		if fmtConditionParams.InitialParamCount > 0 {
			queryStr += join
		}
	firstLoop:
		for operand, filterValues := range fmtConditionParams.Filters {
			queryStr += "("
			for _, field := range filterValues {
				if field.Operator == "IN" || field.Operator == "NOT IN" {
					if tQueryStr, tParams, tParamCount, soteErr = formatArrayFilterCondition(ctx, &ArrFilterParam{
						FilterFields:      field,
						Prefix:            fmtConditionParams.TblPrefix,
						InitialParamCount: paramCount,
						CaseInsensitive:   fmtConditionParams.SortOrderKeysMap[field.FieldName]["case-insensitive"].(bool),
					}, fmtConditionParams.SortOrderKeysMap, panicService); soteErr.ErrCode == nil {
						queryStr += tQueryStr + operand
						params = append(params, tParams...)
						paramCount = tParamCount
					} else {
						break firstLoop
					}
				} else {
					paramCount++
					if fmtConditionParams.SortOrderKeysMap[field.FieldName]["case-insensitive"].(bool) {
						col = fmt.Sprintf("UPPER(%v%v)", fmtConditionParams.TblPrefix, fmtConditionParams.SortOrderKeysMap[field.FieldName]["field"])
						val = fmt.Sprintf("UPPER($%v)", paramCount)
					} else {
						col = fmtConditionParams.TblPrefix + fmtConditionParams.SortOrderKeysMap[field.FieldName]["field"].(string)
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
						queryStr += fmt.Sprintf(" %v %v %v", col, subQuery, operand)
					} else {
						queryStr += fmt.Sprintf(" %v %v %v %v", col, field.Operator, val, operand)
						params = append(params, field.Value)
					}
				}
			}
			queryStr = fmt.Sprintf("%v)%v", strings.TrimSuffix(queryStr, operand), join)
		}

		if soteErr.ErrCode == nil {
			queryStr = strings.TrimSuffix(queryStr, join)
		}
	}

	return

}

// formatArrayFilterCondition formats slice/array filter conditions for a get/list request
func formatArrayFilterCondition(ctx context.Context, reqParams *ArrFilterParam, sortOrderKeysMap map[string]map[string]interface{}, panicService panicServiceFunc) (queryStr string, params []interface{}, paramCount int,
	soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		paramStart string
		paramEnd   string
	)

	s := reflect.ValueOf(reqParams.Value)
	if kind := s.Kind(); kind == reflect.Slice || kind == reflect.Array {
		reqParamLen := s.Len()
		if reqParamLen > 0 {
			paramCount = reqParams.InitialParamCount

			if reqParams.CaseInsensitive {
				paramStart = "UPPER("
				paramEnd = ")"
			}

			queryStr = fmt.Sprintf(" %v%v%v%v %v (", paramStart, reqParams.Prefix, sortOrderKeysMap[reqParams.FieldName]["field"], paramEnd,
				reqParams.Operator)
			params = make([]interface{}, reqParamLen)

			for i := 0; i < reqParamLen; i++ {
				paramCount++
				queryStr += fmt.Sprintf("%v$%v%v,", paramStart, paramCount, paramEnd)
				params[i] = s.Index(i).Interface()
			}

			queryStr = strings.TrimSuffix(queryStr, ",") + ")"
		}
	} else {
		soteErr = panicService(ctx,
			sError.GetSError(207030, sError.BuildParams([]string{reqParams.FieldName, fmt.Sprint(reqParams.Value)}), sError.EmptyMap))
	}

	return
}

//FormatGenericFilterArray formats params from slice/array for additional filters that are not supported by the filters list. (i.e for summary endpoints)
func FormatGenericFilterArray(ctx context.Context, conditionParams *FormatConditionParams, shipmentStatuses []string) (queryStr string, params []interface{}, paramCount int, soteErr sError.SoteError) {
	paramStart := "UPPER("
	paramEnd := ")"

	reqParamLen := len(shipmentStatuses)
	//Check if the prefix has a dot suffix
	if ok := strings.HasSuffix(conditionParams.TblPrefix, "."); ok {
		//we are being paranoid here! We already know that this function will never be called with an empty slice:-)
		if reqParamLen > 0 {
			paramCount = conditionParams.InitialParamCount
			queryStr = fmt.Sprintf(" %v%v%v%v %v (", paramStart, conditionParams.TblPrefix, conditionParams.ColName, paramEnd,
				conditionParams.Operator)

			params = make([]interface{}, reqParamLen)
			for i := 0; i < reqParamLen; i++ {
				paramCount++
				queryStr += fmt.Sprintf("%v$%v%v,", paramStart, paramCount, paramEnd)
				params[i] = shipmentStatuses[i]
			}
			queryStr = strings.TrimSuffix(queryStr, ",") + ")"
		}
	} else {
		err := errors.New(fmt.Sprintf("%v missing suffix,dot", conditionParams.TblPrefix))
		soteErr = sError.GetSError(109999, sError.BuildParams([]string{err.Error()}), sError.EmptyMap)
	}

	return
}
