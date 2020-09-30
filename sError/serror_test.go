package sError

import (
	"strings"
	"testing"
)

func TestIncorrectParams(t *testing.T) {
	s := BuildParams([]string{""})
	if x := GetSError(100100, s, EmptyMap); x.ErrCode != 230060 {
		t.Errorf("The wrong error code (%v) was returned.  230060 should have been returned.", x.ErrCode)
	}
}

func TestErrorCodeNotFound(t *testing.T) {
	errCode := 999999999
	if x := GetSError(errCode, nil, EmptyMap); x.ErrCode != 410000 {
		t.Errorf("%v should have returned an error of 410000", errCode)
	}
}
func Test100100Error(t *testing.T) {
	var errCode = 100100
	s := BuildParams([]string{"SUPER_USER, EXECUTIVE", "DELETE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test109999Error(t *testing.T) {
	var errCode = 109999
	s := BuildParams([]string{"ITEM_NAME"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test200200Error(t *testing.T) {
	var errCode = 200200
	s := BuildParams([]string{"PARAMETER_NAME", "DATA_TYPE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test200250Error(t *testing.T) {
	var errCode = 200250
	s := BuildParams([]string{"PARAMETER_NAME", "PARAMETER_VALUE", "LIST_OF_VALUES"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test200260Error(t *testing.T) {
	var errCode = 200260
	s := BuildParams([]string{"PARAMETER_NAME_MISSING", "PARAMETER_NAME", "PARAMETER_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test200500Error(t *testing.T) {
	var errCode = 200500
	s := BuildParams([]string{"OBJECT_NAME"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test200510Error(t *testing.T) {
	var errCode = 200510
	s := BuildParams([]string{"PARAMETER_NAME", "FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test200511Error(t *testing.T) {
	var errCode = 200511
	s := BuildParams([]string{"PARAMETER_NAME", "PARAMETER_NAME"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test200512Error(t *testing.T) {
	var errCode = 200512
	s := BuildParams([]string{"PARAMETER_NAME", "PARAMETER_NAME"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test200513Error(t *testing.T) {
	var errCode = 200513
	s := BuildParams([]string{"PARAMETER_NAME"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test200514Error(t *testing.T) {
	var errCode = 200514
	s := BuildParams([]string{"PARAMETER_NAME", "SECOND_PARAMETER_NAME", "THIRD_PARAMETER_NAME"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test200515Error(t *testing.T) {
	var errCode = 200515
	s := BuildParams([]string{"PARAMETER_NAME", "ANOTHER_PARAMETER_NAME"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test201000Error(t *testing.T) {
	var errCode = 201000
	s := BuildParams([]string{"DETAILED_MESSAGE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test202000Error(t *testing.T) {
	var errCode = 202000
	s := BuildParams([]string{"ENVIRONMENT"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test230050Error(t *testing.T) {
	var errCode = 230050
	s := BuildParams([]string{"NAME", "APP_PACKAGE_NAME"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test230060Error(t *testing.T) {
	var errCode = 230060
	s := BuildParams([]string{"PROVIDED_PARAMETER_COUNT", "EXPECTED_PARAMETER_COUNT"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test310000Error(t *testing.T) {
	var errCode = 310000
	s := BuildParams([]string{"KEY_NAME"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test310005Error(t *testing.T) {
	var errCode = 310005
	s := BuildParams([]string{"KEY_NAME"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test320000Error(t *testing.T) {
	var errCode = 320000
	s := BuildParams([]string{"PARAMETER_LIST"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test335299Error(t *testing.T) {
	var errCode = 335299
	s := BuildParams([]string{"PARAMETER"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test335699Error(t *testing.T) {
	var errCode = 335599
	s := BuildParams([]string{"PARAMETER", "ANOTHER_PARAMETER"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test400000Error(t *testing.T) {
	var errCode = 400000
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test400005Error(t *testing.T) {
	var errCode = 400005
	s := BuildParams([]string{"FIELD_NAME", "MINIMAL_LENGTH"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test400010Error(t *testing.T) {
	var errCode = 400010
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test400020Error(t *testing.T) {
	var errCode = 400020
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test400030Error(t *testing.T) {
	var errCode = 400030
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test400040Error(t *testing.T) {
	var errCode = 400040
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test400050Error(t *testing.T) {
	var errCode = 400050
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test400060Error(t *testing.T) {
	var errCode = 400060
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test400065Error(t *testing.T) {
	var errCode = 400065
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test400070Error(t *testing.T) {
	var errCode = 400070
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test400080Error(t *testing.T) {
	var errCode = 400080
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test400090Error(t *testing.T) {
	var errCode = 400090
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE", "SMALL_LARGE", "MIN_MAX", "SIZE_EXPECTED", "SIZE_PROVIDED"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test400095Error(t *testing.T) {
	var errCode = 400095
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE", "GREATER_THAN", "LESS_THAN"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test400100Error(t *testing.T) {
	var errCode = 400100
	s := BuildParams([]string{"PARAMETER_NAME", "DATA_STRUCTURE_TYPE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test400105Error(t *testing.T) {
	var errCode = 400105
	s := BuildParams([]string{"DATA_STRUCTURE_NAME", "DATA_STRUCTURE_TYPE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test400110Error(t *testing.T) {
	var errCode = 400110
	s := BuildParams([]string{"PARAMETER_NAME"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test400111Error(t *testing.T) {
	var errCode = 400111
	s := BuildParams([]string{"PARAMETER_NAME", "APPLICATION_PACKAGE_NAME"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test405110Error(t *testing.T) {
	var errCode = 405110
	s := BuildParams([]string{"OBJECT_TYPE", "SYSTEM_ID"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test405120Error(t *testing.T) {
	var errCode = 405120
	s := BuildParams([]string{"JSON_ARRAY", "OBJECT_TYPE", "SYSTEM_ID"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test410000Error(t *testing.T) {
	var errCode = 410000
	s := BuildParams([]string{"ERROR_MESSAGE_NUMBER"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test600010Error(t *testing.T) {
	var errCode = 600010
	s := BuildParams([]string{"FILE_NAME", "MESSAGE_RETURNED_FROM_OPEN"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test601000Error(t *testing.T) {
	var errCode = 601000
	s := BuildParams([]string{"ENVIRONMENT"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test602000Error(t *testing.T) {
	var errCode = 602000
	s := BuildParams([]string{"DATABASE_NAME", "DATABASE_DRIVER", "DATABASE_PORT"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test602020Error(t *testing.T) {
	var errCode = 602020
	s := BuildParams([]string{"SSL_MODE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test602030Error(t *testing.T) {
	var errCode = 602030
	s := BuildParams([]string{"CONNECTION_TYPE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test605020Error(t *testing.T) {
	var errCode = 605020
	s := BuildParams([]string{"KID"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test605021Error(t *testing.T) {
	var errCode = 605021
	s := BuildParams([]string{"KID"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test605030Error(t *testing.T) {
	var errCode = 605030
	s := BuildParams([]string{"REGION", "ENVIRONMENT"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test609998Error(t *testing.T) {
	var errCode = 609998
	s := BuildParams([]string{"START_VARIABLE_OUT_OF_RANGE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test609999Error(t *testing.T) {
	var errCode = 609999
	s := BuildParams([]string{"START_VARIABLE_MISSING"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test700000Error(t *testing.T) {
	var errCode = 700000
	s := BuildParams([]string{"LIST_PARAMETERS"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrorDetails(t *testing.T) {
	var (
		errCode    = 700000
		errDetails = make(map[string]string)
	)
	s := BuildParams([]string{"LIST_PARAMETERS"})
	errDetails["test_1"] = "Test_1_Value"
	errDetails["test_2"] = "Test_2_Value"
	errDetails["test_3"] = "Test_3_Value"
	validateReply(t, errCode, s, GetSError(errCode, s, errDetails))
	errCode = 805000
	validateReply(t, errCode, nil, GetSError(errCode, nil, errDetails))
	validateReply(t, errCode, nil, GetSError(errCode, nil, EmptyMap))
}
func TestGenMarkDown(t *testing.T) {
	if x := GenMarkDown(); !strings.Contains(x, MARKDOWNTITLEBAR) {
		t.Errorf("GenMarkDown doesn't have the correct header, so there appears something wrong with the code.")
	} else {
		println(x)
	}
}
func TestGenErrorListRequiredParams(t *testing.T) {
	if x := GenErrorListRequiredParams(); !strings.Contains(x, FUNCCOMMENTSHEADER) {
		t.Errorf("GenMarkDown doesn't have the correct header, so there appears something wrong with the code.")
	} else {
		println(x)
	}
}
func validateReply(t *testing.T, errCode int, params []interface{}, x SoteError) {
	if errCode != x.ErrCode {
		t.Errorf("Error Code Tested: %v return %v error code when called.", errCode, x.ErrCode)
		t.Fail()
	}
	for i := range params {
		if !strings.Contains(x.FmtErrMsg, params[i].(string)) {
			t.Errorf("Error Code Tested: %v - The %v was not found in the FmtErrMsg property returned", errCode, params[i])
		}
	}
}
