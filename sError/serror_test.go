package sError

import (
	"fmt"
	"strings"
	"testing"
)

func TestIncorrectParams(tPtr *testing.T) {
	s := BuildParams([]string{""})
	if x := GetSError(ErrAuthorized, s, EmptyMap); x.ErrCode != ErrInvalidParameterCount {
		tPtr.Errorf("The wrong error code (%v) was returned.  %v should have been returned.", x.ErrCode, ErrInvalidParameterCount)
	}
}
func TestErrorCodeNotFound(tPtr *testing.T) {
	errCode := 999999999
	if x := GetSError(errCode, nil, EmptyMap); x.ErrCode != ErrMissingErrorMessage {
		tPtr.Errorf("%v should have returned an error of %v", errCode, ErrMissingErrorMessage)
	}
}
func TestErrAuthorizedError(tPtr *testing.T) {
	var errCode = ErrAuthorized
	s := BuildParams([]string{"SUPER_USER, EXECUTIVE", "DELETE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test100200Error(tPtr *testing.T) {
	var errCode = 100200
	validateReply(tPtr, errCode, nil, GetSError(errCode, nil, EmptyMap))
}
func TestErrCancelOrCompleteError(tPtr *testing.T) {
	var errCode = ErrCancelOrComplete
	s := BuildParams([]string{"OBJECT"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test100600Error(tPtr *testing.T) {
	var errCode = 100600
	s := BuildParams([]string{"OBJECT"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrTimeoutError(tPtr *testing.T) {
	var errCode = ErrTimeout
	s := BuildParams([]string{"SERVICE_NAME"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrItemNotFoundError(tPtr *testing.T) {
	var errCode = ErrItemNotFound
	s := BuildParams([]string{"ITEM_NAME"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidParameterTypeError(tPtr *testing.T) {
	var errCode = ErrInvalidParameterType
	s := BuildParams([]string{"PARAMETER_NAME", "DATA_TYPE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidParameterValueError(tPtr *testing.T) {
	var errCode = ErrInvalidParameterValue
	s := BuildParams([]string{"PARAMETER_NAME", "PARAMETER_VALUE", "LIST_OF_VALUES"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test200260Error(tPtr *testing.T) {
	var errCode = 200260
	s := BuildParams([]string{"PARAMETER_NAME_MISSING", "PARAMETER_NAME", "PARAMETER_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test200510Error(tPtr *testing.T) {
	var errCode = 200510
	s := BuildParams([]string{"PARAMETER_NAME", "FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test200511Error(tPtr *testing.T) {
	var errCode = 200511
	s := BuildParams([]string{"PARAMETER_NAME", "PARAMETER_NAME"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrExpectedTwoParametersError(tPtr *testing.T) {
	var errCode = ErrExpectedTwoParameters
	s := BuildParams([]string{"PARAMETER_NAME", "PARAMETER_NAME"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrMissingParametersError(tPtr *testing.T) {
	var errCode = ErrMissingParameters
	s := BuildParams([]string{"PARAMETER_NAME"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrExpectedThreeParametersError(tPtr *testing.T) {
	var errCode = ErrExpectedThreeParameters
	s := BuildParams([]string{"PARAMETER_NAME", "SECOND_PARAMETER_NAME", "THIRD_PARAMETER_NAME"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test200515Error(tPtr *testing.T) {
	var errCode = 200515
	s := BuildParams([]string{"PARAMETER_NAME", "ANOTHER_PARAMETER_NAME"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrBadHTTPRequestError(tPtr *testing.T) {
	var errCode = ErrBadHTTPRequest
	s := BuildParams([]string{"DETAILED_MESSAGE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test200700Error(tPtr *testing.T) {
	var errCode = 200700
	s := BuildParams([]string{"ENVIRONMENT"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test203050Error(tPtr *testing.T) {
	var errCode = 203050
	s := BuildParams([]string{"NAME", "APP_PACKAGE_NAME"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidParameterCountError(tPtr *testing.T) {
	var errCode = ErrInvalidParameterCount
	s := BuildParams([]string{"PROVIDED_PARAMETER_COUNT", "EXPECTED_PARAMETER_COUNT"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test206100Error(tPtr *testing.T) {
	var errCode = 206100
	s := BuildParams([]string{"KEY_NAME"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test206105Error(tPtr *testing.T) {
	var errCode = 206105
	s := BuildParams([]string{"KEY_NAME"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidMsgSignatureError(tPtr *testing.T) {
	var errCode = ErrInvalidMsgSignature
	s := BuildParams([]string{"PARAMETER_LIST"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test206400Error(tPtr *testing.T) {
	var errCode = 206400
	s := BuildParams([]string{"PARAMETER"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test206600Error(tPtr *testing.T) {
	var errCode = 206600
	s := BuildParams([]string{"PARAMETER", "ANOTHER_PARAMETER"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test206700Error(tPtr *testing.T) {
	var errCode = 206700
	s := BuildParams([]string{"PARAMETER", "ANOTHER_PARAMETER"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrNotNumericError(tPtr *testing.T) {
	var errCode = ErrNotNumeric
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207005Error(tPtr *testing.T) {
	var errCode = 207005
	s := BuildParams([]string{"FIELD_NAME", "MINIMAL_LENGTH"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrNotStringError(tPtr *testing.T) {
	var errCode = ErrNotString
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207020Error(tPtr *testing.T) {
	var errCode = 207020
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrNotArrayError(tPtr *testing.T) {
	var errCode = ErrNotArray
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207040Error(tPtr *testing.T) {
	var errCode = 207040
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207050Error(tPtr *testing.T) {
	var errCode = 207050
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207060Error(tPtr *testing.T) {
	var errCode = 207060
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207065Error(tPtr *testing.T) {
	var errCode = 207065
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207070Error(tPtr *testing.T) {
	var errCode = 207070
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207080Error(tPtr *testing.T) {
	var errCode = 207080
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207090Error(tPtr *testing.T) {
	var errCode = 207090
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE", "SMALL_LARGE", "MIN_MAX", "SIZE_EXPECTED", "SIZE_PROVIDED"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test207095Error(tPtr *testing.T) {
	var errCode = 207095
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE", "GREATER_THAN", "LESS_THAN"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrCustomJSONConversionError(tPtr *testing.T) {
	var errCode = ErrCustomJSONConversionError
	s := BuildParams([]string{"PARAMETER_NAME", "DATA_STRUCTURE_TYPE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrConversionErrorError(tPtr *testing.T) {
	var errCode = ErrConversionError
	s := BuildParams([]string{"PARAMETER_NAME", "DATA_STRUCTURE_TYPE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrJSONConversionError(tPtr *testing.T) {
	var errCode = ErrJSONConversionError
	s := BuildParams([]string{"DATA_STRUCTURE_NAME", "DATA_STRUCTURE_TYPE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidJSONError(tPtr *testing.T) {
	var errCode = ErrInvalidJSON
	s := BuildParams([]string{"PARAMETER_NAME"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrMapConversionError(tPtr *testing.T) {
	var errCode = ErrMapConversionError
	s := BuildParams([]string{"PARAMETER_NAME", "APPLICATION_PACKAGE_NAME"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test208110Error(tPtr *testing.T) {
	var errCode = 208110
	s := BuildParams([]string{"OBJECT_TYPE", "SYSTEM_ID"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test208120Error(tPtr *testing.T) {
	var errCode = 208120
	s := BuildParams([]string{"JSON_ARRAY", "OBJECT_TYPE", "SYSTEM_ID"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrMissingErrorMessageError(tPtr *testing.T) {
	var errCode = ErrMissingErrorMessage
	s := BuildParams([]string{"ERROR_MESSAGE_NUMBER"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrMissingFileError(tPtr *testing.T) {
	var errCode = ErrMissingFile
	s := BuildParams([]string{"FILE_NAME", "MESSAGE_RETURNED_FROM_OPEN"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrMissingEnvVariableError(tPtr *testing.T) {
	var errCode = ErrMissingEnvVariable
	s := BuildParams([]string{"ENVIRONMENT"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidEnvValueError(tPtr *testing.T) {
	var errCode = ErrInvalidEnvValue
	s := BuildParams([]string{"ENVIRONMENT"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test209200Error(tPtr *testing.T) {
	var errCode = 209200
	s := BuildParams([]string{"DATABASE_NAME", "DATABASE_DRIVER", "DATABASE_PORT"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidSSLModeError(tPtr *testing.T) {
	var errCode = ErrInvalidSSLMode
	s := BuildParams([]string{"SSL_MODE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test209230Error(tPtr *testing.T) {
	var errCode = 209230
	s := BuildParams([]string{"CONNECTION_TYPE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test209520Error(tPtr *testing.T) {
	var errCode = 209520
	s := BuildParams([]string{"KID"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrMissingKidInPublicKeyError(tPtr *testing.T) {
	var errCode = ErrMissingKidInPublicKey
	s := BuildParams([]string{"KID"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrFetchingJWKError(tPtr *testing.T) {
	var errCode = ErrFetchingJWKError
	s := BuildParams([]string{"REGION", "ENVIRONMENT"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrMissingURLError(tPtr *testing.T) {
	var errCode = ErrMissingURL
	s := BuildParams([]string{"URL_IS_MISSING"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test210098Error(tPtr *testing.T) {
	var errCode = 210098
	s := BuildParams([]string{"START_VARIABLE_OUT_OF_RANGE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test210099Error(tPtr *testing.T) {
	var errCode = 210099
	s := BuildParams([]string{"START_VARIABLE_MISSING"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func Test210100Error(tPtr *testing.T) {
	var errCode = 210100
	s := BuildParams([]string{"LIST_PARAMETERS"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrorDetails(tPtr *testing.T) {
	var (
		errCode    = 210100
		errDetails = make(map[string]string)
	)
	s := BuildParams([]string{"LIST_PARAMETERS"})
	errDetails["test_1"] = "Test_1_Value"
	errDetails["test_2"] = "Test_2_Value"
	errDetails["test_3"] = "Test_3_Value"
	validateReply(tPtr, errCode, s, GetSError(errCode, s, errDetails))
	errCode = ErrSynadiaError
	validateReply(tPtr, errCode, nil, GetSError(errCode, nil, errDetails))
	validateReply(tPtr, errCode, nil, GetSError(errCode, nil, EmptyMap))
}
func TestGenerateDocumentation(tPtr *testing.T) {
	x, y := GenerateDocumentation()

	if !strings.Contains(x, MARKDOWNTITLEBAR) {
		tPtr.Errorf("TestGenerateDocumentation doesn't have the correct header, so there appears something wrong with the code.")
	} else {
		// This outputs the documentation so it can be cut/paste to the wiki.
		fmt.Print(x)
	}

	if !strings.Contains(y, FUNCCOMMENTSHEADER) {
		tPtr.Errorf("TestGenerateDocumentation doesn't have the correct header, so there appears something wrong with the code.")
	} else {
		// This outputs the documentation so it can be cut/paste to the source code.
		fmt.Print(y)
	}
}
func TestOutputErrorJSON(tPtr *testing.T) {
	var (
		errCode               = 999999999
		s       []interface{} = BuildParams([]string{""})
		x       []byte
	)
	if x = OutputErrorJSON(GetSError(errCode, nil, EmptyMap)); string(x) == "" {
		tPtr.Errorf(GetSError(ErrCustomJSONConversionError, s, EmptyMap).FmtErrMsg)
	}
}
func validateReply(tPtr *testing.T, errCode int, params []interface{}, x SoteError) {
	if errCode != x.ErrCode {
		tPtr.Errorf("Error Code Tested: %v return %v error code when called.", errCode, x.ErrCode)
		tPtr.Fail()
	}
	for i := range params {
		if !strings.Contains(x.FmtErrMsg, params[i].(string)) {
			tPtr.Errorf("Error Code Tested: %v - The %v was not found in the FmtErrMsg property returned", errCode, params[i])
		}
	}
}
