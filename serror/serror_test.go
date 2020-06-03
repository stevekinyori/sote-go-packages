package serror

import (
	"strings"
	"testing"
)

func TestIncorrectParams(t *testing.T) {
	var paramValues = []string{""}
	if x := GetSError(100100, paramValues, EmptyMap); x.ErrCode != 230060 {
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
	var paramValues = []string{"SUPER_USER, EXECUTIVE", "DELETE"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test109999Error(t *testing.T) {
	var errCode = 109999
	var paramValues = []string{"ITEM_NAME"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test200200Error(t *testing.T) {
	var errCode = 200200
	var paramValues = []string{"PARAMETER_NAME", "DATA_TYPE"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test200250Error(t *testing.T) {
	var errCode = 200250
	var paramValues = []string{"PARAMETER_NAME", "PARAMETER_VALUE", "LIST_OF_VALUES"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test200260Error(t *testing.T) {
	var errCode = 200260
	var paramValues = []string{"PARAMETER_NAME_MISSING", "PARAMETER_NAME", "PARAMETER_VALUE"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test200500Error(t *testing.T) {
	var errCode = 200500
	var paramValues = []string{"OBJECT_NAME"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test200510Error(t *testing.T) {
	var errCode = 200510
	var paramValues = []string{"PARAMETER_NAME", "FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test200511Error(t *testing.T) {
	var errCode = 200511
	var paramValues = []string{"PARAMETER_NAME", "PARAMETER_NAME"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test201000Error(t *testing.T) {
	var errCode = 201000
	var paramValues = []string{"DETAILED_MESSAGE"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test202000Error(t *testing.T) {
	var errCode = 202000
	var paramValues = []string{"ENVIRONMENT"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test230050Error(t *testing.T) {
	var errCode = 230050
	var paramValues = []string{"NAME", "APP_PACKAGE_NAME"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test230060Error(t *testing.T) {
	var errCode = 230060
	var paramValues = []string{"PROVIDED_PARAMETER_COUNT", "EXPECTED_PARAMETER_COUNT"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test310000Error(t *testing.T) {
	var errCode = 310000
	var paramValues = []string{"KEY_NAME"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test310005Error(t *testing.T) {
	var errCode = 310005
	var paramValues = []string{"KEY_NAME"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test320000Error(t *testing.T) {
	var errCode = 320000
	var paramValues = []string{"PARAMETER_LIST"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test400000Error(t *testing.T) {
	var errCode = 400000
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test400005Error(t *testing.T) {
	var errCode = 400005
	var paramValues = []string{"FIELD_NAME", "MINIMAL_LENGTH"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test400010Error(t *testing.T) {
	var errCode = 400010
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test400020Error(t *testing.T) {
	var errCode = 400020
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test400030Error(t *testing.T) {
	var errCode = 400030
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test400040Error(t *testing.T) {
	var errCode = 400040
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test400050Error(t *testing.T) {
	var errCode = 400050
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test400060Error(t *testing.T) {
	var errCode = 400060
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test400065Error(t *testing.T) {
	var errCode = 400065
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test400070Error(t *testing.T) {
	var errCode = 400070
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test400080Error(t *testing.T) {
	var errCode = 400080
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test400090Error(t *testing.T) {
	var errCode = 400090
	var paramValues = []string{"FIELD_NAME", "FIELD_VALUE", "SMALL_LARGE", "MIN_MAX", "SIZE_EXPECTED", "SIZE_PROVIDED"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test400100Error(t *testing.T) {
	var errCode = 400100
	var paramValues = []string{"FIELD_VALUE"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test405110Error(t *testing.T) {
	var errCode = 405110
	var paramValues = []string{"OBJECT_TYPE", "SYSTEM_ID"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test405120Error(t *testing.T) {
	var errCode = 405120
	var paramValues = []string{"JSON_ARRAY", "OBJECT_TYPE", "SYSTEM_ID"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test410000Error(t *testing.T) {
	var errCode = 410000
	var paramValues = []string{"ERROR_MESSAGE_NUMBER"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test600010Error(t *testing.T) {
	var errCode = 600010
	var paramValues = []string{"FILE_NAME", "MESSAGE_RETURNED_FROM_OPEN"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test601000Error(t *testing.T) {
	var errCode = 601000
	var paramValues = []string{"ENVIRONMENT"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test602000Error(t *testing.T) {
	var errCode = 602000
	var paramValues = []string{"DATABASE_NAME", "DATABASE_DRIVER", "DATABASE_PORT"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test609999Error(t *testing.T) {
	var errCode = 609999
	var paramValues = []string{"START_VARIABLE_MISSING"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func Test700000Error(t *testing.T) {
	var errCode = 700000
	var paramValues = []string{"LIST_PARAMETERS"}
	validateReply(t, errCode, paramValues, GetSError(errCode, paramValues, EmptyMap))
}
func TestGenMarkDown(t *testing.T) {
	if x := GenMarkDown(); !strings.Contains(x, MarkDownTitleBar) {
		t.Errorf("GenMarkDown doesn't have the correct header, so there appears something wrong with the code.")
	} else {
		println(x)
	}
}
func TestGenErrorListRequiredParams(t *testing.T) {
	if x := GenErrorLisRequiredParams(); !strings.Contains(x, FuncCommentsHeader) {
		t.Errorf("GenMarkDown doesn't have the correct header, so there appears something wrong with the code.")
	} else {
		println(x)
	}
}
func validateReply(t *testing.T, errCode int, paramValues []string, x SoteError) {
	if errCode != x.ErrCode {
		t.Errorf("Error Code Tested: %v return v% error code when called.", errCode, x.ErrCode)
		t.Fail()
	}
	for i, _ := range paramValues {
		if !strings.Contains(x.FmtErrMsg, paramValues[i]) {
			t.Errorf("Error Code Tested: %v - The %v was not found in the FmtErrMsg property returned", errCode, paramValues[i])
		}
	}
}
func validateCovertErr(t *testing.T) {

}