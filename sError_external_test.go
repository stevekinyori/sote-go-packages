package packages

import (
	"strings"
	"testing"

	"gitlab.com/soteapps/packages/v2021/sError"
)

func TestIncorrectParams(t *testing.T) {
	if x := sError.GetSError(100100, sError.BuildParams([]string{""}), sError.EmptyMap); x.ErrCode != 203060 {
		t.Errorf("The wrong error code (%v) was returned.  203060 should have been returned.", x.ErrCode)
	}
}
func TestErrorCodeNotFound(t *testing.T) {
	errCode := 999999999
	if x := sError.GetSError(errCode, nil, sError.EmptyMap); x.ErrCode != 208200 {
		t.Errorf("%v should have returned an error of 208200", errCode)
	}
}
func Test100100Error(t *testing.T) {
	var errCode = 100100
	s := sError.BuildParams([]string{"SUPER_USER, EXECUTIVE", "DELETE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test109999Error(t *testing.T) {
	var errCode = 109999
	s := sError.BuildParams([]string{"ITEM_NAME"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200200Error(t *testing.T) {
	var errCode = 200200
	s := sError.BuildParams([]string{"PARAMETER_NAME", "DATA_TYPE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200250Error(t *testing.T) {
	var errCode = 200250
	s := sError.BuildParams([]string{"PARAMETER_NAME", "PARAMETER_VALUE", "LIST_OF_VALUES"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200260Error(t *testing.T) {
	var errCode = 200260
	s := sError.BuildParams([]string{"PARAMETER_NAME_MISSING", "PARAMETER_NAME", "PARAMETER_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200510Error(t *testing.T) {
	var errCode = 200510
	s := sError.BuildParams([]string{"PARAMETER_NAME", "FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200511Error(t *testing.T) {
	var errCode = 200511
	s := sError.BuildParams([]string{"PARAMETER_NAME", "PARAMETER_NAME"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200512Error(t *testing.T) {
	var errCode = 200512
	s := sError.BuildParams([]string{"PARAMETER_NAME", "PARAMETER_NAME"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200513Error(t *testing.T) {
	var errCode = 200513
	s := sError.BuildParams([]string{"PARAMETER_NAME"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200600Error(t *testing.T) {
	var errCode = 200600
	s := sError.BuildParams([]string{"DETAILED_MESSAGE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200700Error(t *testing.T) {
	var errCode = 200700
	s := sError.BuildParams([]string{"ENVIRONMENT"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test203050Error(t *testing.T) {
	var errCode = 203050
	s := sError.BuildParams([]string{"NAME", "APP_PACKAGE_NAME"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test203060Error(t *testing.T) {
	var errCode = 203060
	s := sError.BuildParams([]string{"PROVIDED_PARAMETER_COUNT", "EXPECTED_PARAMETER_COUNT"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test206100Error(t *testing.T) {
	var errCode = 206100
	s := sError.BuildParams([]string{"KEY_NAME"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test206105Error(t *testing.T) {
	var errCode = 206105
	s := sError.BuildParams([]string{"KEY_NAME"})
	params := make([]interface{}, 1)
	params[0] = "KEY_NAME"
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test206200Error(t *testing.T) {
	var errCode = 206200
	s := sError.BuildParams([]string{"PARAMETER_LIST"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test206400Error(t *testing.T) {
	var errCode = 206400
	s := sError.BuildParams([]string{"PARAMETER"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test335699Error(t *testing.T) {
	var errCode = 206600
	s := sError.BuildParams([]string{"PARAMETER", "ANOTHER_PARAMETER"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test206700Error(t *testing.T) {
	var errCode = 206700
	s := sError.BuildParams([]string{"PARAMETER", "ANOTHER_PARAMETER"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207000Error(t *testing.T) {
	var errCode = 207000
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207005Error(t *testing.T) {
	var errCode = 207005
	s := sError.BuildParams([]string{"FIELD_NAME", "MINIMAL_LENGTH"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207010Error(t *testing.T) {
	var errCode = 207010
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207020Error(t *testing.T) {
	var errCode = 207020
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207030Error(t *testing.T) {
	var errCode = 207030
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207040Error(t *testing.T) {
	var errCode = 207040
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207050Error(t *testing.T) {
	var errCode = 207050
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207060Error(t *testing.T) {
	var errCode = 207060
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207065Error(t *testing.T) {
	var errCode = 207065
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207070Error(t *testing.T) {
	var errCode = 207070
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207080Error(t *testing.T) {
	var errCode = 207080
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207090Error(t *testing.T) {
	var errCode = 207090
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE", "SMALL_LARGE", "MIN_MAX", "SIZE_EXPECTED", "SIZE_PROVIDED"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207100Error(t *testing.T) {
	var errCode = 207100
	s := sError.BuildParams([]string{"PARAMETER_NAME", "DATA_STRUCTURE_TYPE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207200Error(t *testing.T) {
	var errCode = 207200
	s := sError.BuildParams([]string{"PARAMETER_NAME", "DATA_STRUCTURE_TYPE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207105Error(t *testing.T) {
	var errCode = 207105
	s := sError.BuildParams([]string{"DATA_STRUCTURE_NAME", "DATA_STRUCTURE_TYPE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207110Error(t *testing.T) {
	var errCode = 207110
	s := sError.BuildParams([]string{"PARAMETER_NAME"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207111Error(t *testing.T) {
	var errCode = 207111
	s := sError.BuildParams([]string{"PARAMETER_NAME", "APPLICATION_PACKAGE_NAME"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test208110Error(t *testing.T) {
	var errCode = 208110
	s := sError.BuildParams([]string{"OBJECT_TYPE", "SYSTEM_ID"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test208120Error(t *testing.T) {
	var errCode = 208120
	s := sError.BuildParams([]string{"JSON_ARRAY", "OBJECT_TYPE", "SYSTEM_ID"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test208200Error(t *testing.T) {
	var errCode = 208200
	s := sError.BuildParams([]string{"ERROR_MESSAGE_NUMBER"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test209010Error(t *testing.T) {
	var errCode = 209010
	s := sError.BuildParams([]string{"FILE_NAME", "MESSAGE_RETURNED_FROM_OPEN"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test209100Error(t *testing.T) {
	var errCode = 209100
	s := sError.BuildParams([]string{"ENVIRONMENT"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test209200Error(t *testing.T) {
	var errCode = 209200
	s := sError.BuildParams([]string{"DATABASE_NAME", "DATABASE_DRIVER", "DATABASE_PORT"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test209220Error(t *testing.T) {
	var errCode = 209220
	s := sError.BuildParams([]string{"SSL_MODE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test209230Error(t *testing.T) {
	var errCode = 209230
	s := sError.BuildParams([]string{"CONNECTION_TYPE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test209520Error(t *testing.T) {
	var errCode = 209520
	s := sError.BuildParams([]string{"KID"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test209521Error(t *testing.T) {
	var errCode = 209521
	s := sError.BuildParams([]string{"KID"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test210030Error(t *testing.T) {
	var errCode = 210030
	s := sError.BuildParams([]string{"REGION","ENVIRONMENT"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test210090Error(t *testing.T) {
	var errCode = 210090
	s := sError.BuildParams([]string{"URL_IS_MISSING"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test210098Error(t *testing.T) {
	var errCode = 210098
	s := sError.BuildParams([]string{"START_VARIABLE_OUT_OF_RANGE"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test210099Error(t *testing.T) {
	var errCode = 210099
	s := sError.BuildParams([]string{"START_VARIABLE_MISSING"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test210100Error(t *testing.T) {
	var errCode = 210100
	s := sError.BuildParams([]string{"LIST_PARAMETERS"})
	validateReply(t, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrorDetails(t *testing.T) {
	var (
		errCode    = 210100
		errDetails = make(map[string]string)
	)
	s := sError.BuildParams([]string{"LIST_PARAMETERS"})
	errDetails["test_1"] = "Test_1_Value"
	errDetails["test_2"] = "Test_2_Value"
	errDetails["test_3"] = "Test_3_Value"
	validateReply(t, errCode, s, sError.GetSError(errCode, s, errDetails))
	errCode = 210400
	validateReply(t, errCode, nil, sError.GetSError(errCode, nil, errDetails))
	validateReply(t, errCode, nil, sError.GetSError(errCode, nil, sError.EmptyMap))
}
func TestGenMarkDown(t *testing.T) {
	if x := sError.GenMarkDown(); !strings.Contains(x, sError.MARKDOWNTITLEBAR) {
		t.Errorf("GenMarkDown doesn't have the correct header, so there appears something wrong with the code.")
	} else {
		println(x)
	}
}
func TestGenErrorListRequiredParams(t *testing.T) {
	if x := sError.GenErrorListRequiredParams(); !strings.Contains(x, sError.FUNCCOMMENTSHEADER) {
		t.Errorf("GenMarkDown doesn't have the correct header, so there appears something wrong with the code.")
	} else {
		println(x)
	}
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
