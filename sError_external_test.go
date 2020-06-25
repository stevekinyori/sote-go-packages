package packages

import (
	"strings"
	"testing"

	"gitlab.com/soteapps/packages/sError"
)

func TestIncorrectParams(t *testing.T) {
	if x := sError.GetSError(100100, buildParams([]string{""}), sError.EmptyMap); x.ErrCode != 230060 {
		t.Errorf("The wrong error code (%v) was returned.  230060 should have been returned.", x.ErrCode)
	}
}
func TestErrorCodeNotFound(t *testing.T) {
	errCode := 999999999
	if x := sError.GetSError(errCode, nil, sError.EmptyMap); x.ErrCode != 410000 {
		t.Errorf("%v should have returned an error of 410000", errCode)
	}
}
func Test100100Error(t *testing.T) {
	var errCode = 100100
	s := buildParams([]string{"SUPER_USER, EXECUTIVE", "DELETE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test109999Error(t *testing.T) {
	var errCode = 109999
	s := buildParams([]string{"ITEM_NAME"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200200Error(t *testing.T) {
	var errCode = 200200
	s := buildParams([]string{"PARAMETER_NAME", "DATA_TYPE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200250Error(t *testing.T) {
	var errCode = 200250
	s := buildParams([]string{"PARAMETER_NAME", "PARAMETER_VALUE", "LIST_OF_VALUES"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200260Error(t *testing.T) {
	var errCode = 200260
	s := buildParams([]string{"PARAMETER_NAME_MISSING", "PARAMETER_NAME", "PARAMETER_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200500Error(t *testing.T) {
	var errCode = 200500
	s := buildParams([]string{"OBJECT_NAME"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200510Error(t *testing.T) {
	var errCode = 200510
	s := buildParams([]string{"PARAMETER_NAME", "FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200511Error(t *testing.T) {
	var errCode = 200511
	s := buildParams([]string{"PARAMETER_NAME", "PARAMETER_NAME"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test201000Error(t *testing.T) {
	var errCode = 201000
	s := buildParams([]string{"DETAILED_MESSAGE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test202000Error(t *testing.T) {
	var errCode = 202000
	s := buildParams([]string{"ENVIRONMENT"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test230050Error(t *testing.T) {
	var errCode = 230050
	s := buildParams([]string{"NAME", "APP_PACKAGE_NAME"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test230060Error(t *testing.T) {
	var errCode = 230060
	s := buildParams([]string{"PROVIDED_PARAMETER_COUNT", "EXPECTED_PARAMETER_COUNT"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test310000Error(t *testing.T) {
	var errCode = 310000
	s := buildParams([]string{"KEY_NAME"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test310005Error(t *testing.T) {
	var errCode = 310005
	s := buildParams([]string{"ITEM_NAME"})
	params := make([]interface{}, 1)
	params[0] = "KEY_NAME"
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test320000Error(t *testing.T) {
	var errCode = 320000
	s := buildParams([]string{"PARAMETER_LIST"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test400000Error(t *testing.T) {
	var errCode = 400000
	s := buildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test400005Error(t *testing.T) {
	var errCode = 400005
	s := buildParams([]string{"FIELD_NAME", "MINIMAL_LENGTH"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test400010Error(t *testing.T) {
	var errCode = 400010
	s := buildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test400020Error(t *testing.T) {
	var errCode = 400020
	s := buildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test400030Error(t *testing.T) {
	var errCode = 400030
	s := buildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test400040Error(t *testing.T) {
	var errCode = 400040
	s := buildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test400050Error(t *testing.T) {
	var errCode = 400050
	s := buildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test400060Error(t *testing.T) {
	var errCode = 400060
	s := buildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test400065Error(t *testing.T) {
	var errCode = 400065
	s := buildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test400070Error(t *testing.T) {
	var errCode = 400070
	s := buildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test400080Error(t *testing.T) {
	var errCode = 400080
	s := buildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test400090Error(t *testing.T) {
	var errCode = 400090
	s := buildParams([]string{"FIELD_NAME", "FIELD_VALUE", "SMALL_LARGE", "MIN_MAX", "SIZE_EXPECTED", "SIZE_PROVIDED"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test400100Error(t *testing.T) {
	var errCode = 400100
	s := buildParams([]string{"PARAMETER_NAME", "DATA_STRUCTURE_TYPE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test400105Error(t *testing.T) {
	var errCode = 400105
	s := buildParams([]string{"DATA_STRUCTURE_NAME", "DATA_STRUCTURE_TYPE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test400110Error(t *testing.T) {
	var errCode = 400110
	s := buildParams([]string{"PARAMETER_NAME"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test400111Error(t *testing.T) {
	var errCode = 400111
	s := buildParams([]string{"PARAMETER_NAME", "APPLICATION_PACKAGE_NAME"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test405110Error(t *testing.T) {
	var errCode = 405110
	s := buildParams([]string{"OBJECT_TYPE", "SYSTEM_ID"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test405120Error(t *testing.T) {
	var errCode = 405120
	s := buildParams([]string{"JSON_ARRAY", "OBJECT_TYPE", "SYSTEM_ID"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test410000Error(t *testing.T) {
	var errCode = 410000
	s := buildParams([]string{"ERROR_MESSAGE_NUMBER"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test600010Error(t *testing.T) {
	var errCode = 600010
	s := buildParams([]string{"FILE_NAME", "MESSAGE_RETURNED_FROM_OPEN"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test601000Error(t *testing.T) {
	var errCode = 601000
	s := buildParams([]string{"ENVIRONMENT"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test602000Error(t *testing.T) {
	var errCode = 602000
	s := buildParams([]string{"DATABASE_NAME", "DATABASE_DRIVER", "DATABASE_PORT"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test602020Error(t *testing.T) {
	var errCode = 602020
	s := buildParams([]string{"SSL_MODE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test602100Error(t *testing.T) {
	var errCode = 602100
	s := buildParams([]string{"CONNECTION_TYPE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test609999Error(t *testing.T) {
	var errCode = 609999
	s := buildParams([]string{"START_VARIABLE_MISSING"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test700000Error(t *testing.T) {
	var errCode = 700000
	s := buildParams([]string{"LIST_PARAMETERS"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestGenMarkDown(t *testing.T) {
	if x := sError.GenMarkDown(); !strings.Contains(x, sError.MarkDownTitleBar) {
		t.Errorf("GenMarkDown doesn't have the correct header, so there appears something wrong with the code.")
	} else {
		println(x)
	}
}
func TestGenErrorListRequiredParams(t *testing.T) {
	if x := sError.GenErrorLisRequiredParams(); !strings.Contains(x, sError.FuncCommentsHeader) {
		t.Errorf("GenMarkDown doesn't have the correct header, so there appears something wrong with the code.")
	} else {
		println(x)
	}
}
func buildParams(values []string) (s []interface{}){
	s = make([]interface{}, len(values))
	for i, v := range values {
		s[i] = v
	}

	return
}
func validateReply(t *testing.T, errCode int, params []interface{}, x sError.SoteError) {
	if errCode != x.ErrCode {
		t.Errorf("Error Code Tested: %v return %v error code when called.", errCode, x.ErrCode)
		t.Fail()
	}
	for i, _ := range params {
		if !strings.Contains(x.FmtErrMsg, params[i].(string)) {
			t.Errorf("Error Code Tested: %v - The %v was not found in the FmtErrMsg property returned", errCode, params[i])
		}
	}
}
