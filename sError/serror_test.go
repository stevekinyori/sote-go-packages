package sError

import (
	"strings"
	"testing"
)

func TestIncorrectParams(t *testing.T) {
	s := BuildParams([]string{""})
	if x := GetSError(100100, s, EmptyMap); x.ErrCode != 203060 {
		t.Errorf("The wrong error code (%v) was returned.  203060 should have been returned.", x.ErrCode)
	}
}

func TestErrorCodeNotFound(t *testing.T) {
	errCode := 999999999
	if x := GetSError(errCode, nil, EmptyMap); x.ErrCode != 208200 {
		t.Errorf("%v should have returned an error of 208200", errCode)
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
func Test200600Error(t *testing.T) {
	var errCode = 200600
	s := BuildParams([]string{"DETAILED_MESSAGE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test200700Error(t *testing.T) {
	var errCode = 200700
	s := BuildParams([]string{"ENVIRONMENT"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test203050Error(t *testing.T) {
	var errCode = 203050
	s := BuildParams([]string{"NAME", "APP_PACKAGE_NAME"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test203060Error(t *testing.T) {
	var errCode = 203060
	s := BuildParams([]string{"PROVIDED_PARAMETER_COUNT", "EXPECTED_PARAMETER_COUNT"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test206100Error(t *testing.T) {
	var errCode = 206100
	s := BuildParams([]string{"KEY_NAME"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test206105Error(t *testing.T) {
	var errCode = 206105
	s := BuildParams([]string{"KEY_NAME"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test206200Error(t *testing.T) {
	var errCode = 206200
	s := BuildParams([]string{"PARAMETER_LIST"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test206400Error(t *testing.T) {
	var errCode = 206400
	s := BuildParams([]string{"PARAMETER"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test206600Error(t *testing.T) {
	var errCode = 206600
	s := BuildParams([]string{"PARAMETER", "ANOTHER_PARAMETER"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test206700Error(t *testing.T) {
	var errCode = 206700
	s := BuildParams([]string{"PARAMETER", "ANOTHER_PARAMETER"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207000Error(t *testing.T) {
	var errCode = 207000
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207005Error(t *testing.T) {
	var errCode = 207005
	s := BuildParams([]string{"FIELD_NAME", "MINIMAL_LENGTH"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207010Error(t *testing.T) {
	var errCode = 207010
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207020Error(t *testing.T) {
	var errCode = 207020
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207030Error(t *testing.T) {
	var errCode = 207030
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207040Error(t *testing.T) {
	var errCode = 207040
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207050Error(t *testing.T) {
	var errCode = 207050
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207060Error(t *testing.T) {
	var errCode = 207060
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207065Error(t *testing.T) {
	var errCode = 207065
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207070Error(t *testing.T) {
	var errCode = 207070
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207080Error(t *testing.T) {
	var errCode = 207080
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207090Error(t *testing.T) {
	var errCode = 207090
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE", "SMALL_LARGE", "MIN_MAX", "SIZE_EXPECTED", "SIZE_PROVIDED"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207095Error(t *testing.T) {
	var errCode = 207095
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE", "GREATER_THAN", "LESS_THAN"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207100Error(t *testing.T) {
	var errCode = 207100
	s := BuildParams([]string{"PARAMETER_NAME", "DATA_STRUCTURE_TYPE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207200Error(t *testing.T) {
	var errCode = 207200
	s := BuildParams([]string{"PARAMETER_NAME", "DATA_STRUCTURE_TYPE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207105Error(t *testing.T) {
	var errCode = 207105
	s := BuildParams([]string{"DATA_STRUCTURE_NAME", "DATA_STRUCTURE_TYPE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207110Error(t *testing.T) {
	var errCode = 207110
	s := BuildParams([]string{"PARAMETER_NAME"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207111Error(t *testing.T) {
	var errCode = 207111
	s := BuildParams([]string{"PARAMETER_NAME", "APPLICATION_PACKAGE_NAME"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test208110Error(t *testing.T) {
	var errCode = 208110
	s := BuildParams([]string{"OBJECT_TYPE", "SYSTEM_ID"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test208120Error(t *testing.T) {
	var errCode = 208120
	s := BuildParams([]string{"JSON_ARRAY", "OBJECT_TYPE", "SYSTEM_ID"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test208200Error(t *testing.T) {
	var errCode = 208200
	s := BuildParams([]string{"ERROR_MESSAGE_NUMBER"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test209010Error(t *testing.T) {
	var errCode = 209010
	s := BuildParams([]string{"FILE_NAME", "MESSAGE_RETURNED_FROM_OPEN"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test209100Error(t *testing.T) {
	var errCode = 209100
	s := BuildParams([]string{"ENVIRONMENT"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test209110Error(t *testing.T) {
	var errCode = 209110
	s := BuildParams([]string{"ENVIRONMENT"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test209200Error(t *testing.T) {
	var errCode = 209200
	s := BuildParams([]string{"DATABASE_NAME", "DATABASE_DRIVER", "DATABASE_PORT"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test209220Error(t *testing.T) {
	var errCode = 209220
	s := BuildParams([]string{"SSL_MODE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test209230Error(t *testing.T) {
	var errCode = 209230
	s := BuildParams([]string{"CONNECTION_TYPE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test209520Error(t *testing.T) {
	var errCode = 209520
	s := BuildParams([]string{"KID"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test209521Error(t *testing.T) {
	var errCode = 209521
	s := BuildParams([]string{"KID"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test210030Error(t *testing.T) {
	var errCode = 210030
	s := BuildParams([]string{"REGION", "ENVIRONMENT"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test210090Error(t *testing.T) {
	var errCode = 210090
	s := BuildParams([]string{"URL_IS_MISSING"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test210098Error(t *testing.T) {
	var errCode = 210098
	s := BuildParams([]string{"START_VARIABLE_OUT_OF_RANGE"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test210099Error(t *testing.T) {
	var errCode = 210099
	s := BuildParams([]string{"START_VARIABLE_MISSING"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test210100Error(t *testing.T) {
	var errCode = 210100
	s := BuildParams([]string{"LIST_PARAMETERS"})
	validateReply(t, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrorDetails(t *testing.T) {
	var (
		errCode    = 210100
		errDetails = make(map[string]string)
	)
	s := BuildParams([]string{"LIST_PARAMETERS"})
	errDetails["test_1"] = "Test_1_Value"
	errDetails["test_2"] = "Test_2_Value"
	errDetails["test_3"] = "Test_3_Value"
	validateReply(t, errCode, s, GetSError(errCode, s, errDetails))
	errCode = 210400
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
