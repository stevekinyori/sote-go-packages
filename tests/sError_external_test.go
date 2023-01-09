package tests

import (
	"strings"
	"testing"

	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
)

func init() {
	sLogger.SetLogMessagePrefix("packages")
}

func TestIncorrectParams(tPtr *testing.T) {
	if x := sError.GetSError(sError.ErrAuthorized, sError.BuildParams([]string{""}),
		sError.EmptyMap); x.ErrCode != sError.ErrInvalidParameterCount {
		tPtr.Errorf("The wrong error code (%v) was returned.  %v should have been returned.", x.ErrCode, sError.ErrInvalidParameterCount)
	}
}
func TestErrorCodeNotFound(tPtr *testing.T) {
	errCode := 999999999
	if x := sError.GetSError(errCode, nil, sError.EmptyMap); x.ErrCode != sError.ErrMissingErrorMessage {
		tPtr.Errorf("%v should have returned an error of %v", errCode, sError.ErrMissingErrorMessage)
	}
}
func TestErrAuthorizedError(tPtr *testing.T) {
	var errCode = sError.ErrAuthorized
	s := sError.BuildParams([]string{"SUPER_USER, EXECUTIVE", "DELETE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrItemNotFoundError(tPtr *testing.T) {
	var errCode = sError.ErrItemNotFound
	s := sError.BuildParams([]string{"ITEM_NAME"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrInvalidParameterTypeError(tPtr *testing.T) {
	var errCode = sError.ErrInvalidParameterType
	s := sError.BuildParams([]string{"PARAMETER_NAME", "DATA_TYPE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrInvalidParameterValueError(tPtr *testing.T) {
	var errCode = sError.ErrInvalidParameterValue
	s := sError.BuildParams([]string{"PARAMETER_NAME", "PARAMETER_VALUE", "LIST_OF_VALUES"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200260Error(tPtr *testing.T) {
	var errCode = 200260
	s := sError.BuildParams([]string{"PARAMETER_NAME_MISSING", "PARAMETER_NAME", "PARAMETER_VALUE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200510Error(tPtr *testing.T) {
	var errCode = 200510
	s := sError.BuildParams([]string{"PARAMETER_NAME", "FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200511Error(tPtr *testing.T) {
	var errCode = 200511
	s := sError.BuildParams([]string{"PARAMETER_NAME", "PARAMETER_NAME"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrExpectedTwoParametersError(tPtr *testing.T) {
	var errCode = sError.ErrExpectedTwoParameters
	s := sError.BuildParams([]string{"PARAMETER_NAME", "PARAMETER_NAME"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrMissingParametersError(tPtr *testing.T) {
	var errCode = sError.ErrMissingParameters
	s := sError.BuildParams([]string{"PARAMETER_NAME"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrBadHTTPRequestError(tPtr *testing.T) {
	var errCode = sError.ErrBadHTTPRequest
	s := sError.BuildParams([]string{"DETAILED_MESSAGE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test200700Error(tPtr *testing.T) {
	var errCode = 200700
	s := sError.BuildParams([]string{"ENVIRONMENT"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test203050Error(tPtr *testing.T) {
	var errCode = 203050
	s := sError.BuildParams([]string{"NAME", "APP_PACKAGE_NAME"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrInvalidParameterCountError(tPtr *testing.T) {
	var errCode = sError.ErrInvalidParameterCount
	s := sError.BuildParams([]string{"PROVIDED_PARAMETER_COUNT", "EXPECTED_PARAMETER_COUNT"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test206100Error(tPtr *testing.T) {
	var errCode = 206100
	s := sError.BuildParams([]string{"KEY_NAME"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test206105Error(tPtr *testing.T) {
	var errCode = 206105
	s := sError.BuildParams([]string{"KEY_NAME"})
	params := make([]interface{}, 1)
	params[0] = "KEY_NAME"
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrInvalidMsgSignatureError(tPtr *testing.T) {
	var errCode = sError.ErrInvalidMsgSignature
	s := sError.BuildParams([]string{"PARAMETER_LIST"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test206400Error(tPtr *testing.T) {
	var errCode = 206400
	s := sError.BuildParams([]string{"PARAMETER"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test335699Error(tPtr *testing.T) {
	var errCode = 206600
	s := sError.BuildParams([]string{"PARAMETER", "ANOTHER_PARAMETER"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test206700Error(tPtr *testing.T) {
	var errCode = 206700
	s := sError.BuildParams([]string{"PARAMETER", "ANOTHER_PARAMETER"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrNotNumericError(tPtr *testing.T) {
	var errCode = sError.ErrNotNumeric
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207005Error(tPtr *testing.T) {
	var errCode = 207005
	s := sError.BuildParams([]string{"FIELD_NAME", "MINIMAL_LENGTH"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrNotStringError(tPtr *testing.T) {
	var errCode = sError.ErrNotString
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207020Error(tPtr *testing.T) {
	var errCode = 207020
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrNotArrayError(tPtr *testing.T) {
	var errCode = sError.ErrNotArray
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207040Error(tPtr *testing.T) {
	var errCode = 207040
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207050Error(tPtr *testing.T) {
	var errCode = 207050
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207060Error(tPtr *testing.T) {
	var errCode = 207060
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207065Error(tPtr *testing.T) {
	var errCode = 207065
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207070Error(tPtr *testing.T) {
	var errCode = 207070
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207080Error(tPtr *testing.T) {
	var errCode = 207080
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test207090Error(tPtr *testing.T) {
	var errCode = 207090
	s := sError.BuildParams([]string{"FIELD_NAME", "FIELD_VALUE", "SMALL_LARGE", "MIN_MAX", "SIZE_EXPECTED", "SIZE_PROVIDED"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrCustomJSONConversionError(tPtr *testing.T) {
	var errCode = sError.ErrCustomJSONConversionError
	s := sError.BuildParams([]string{"PARAMETER_NAME", "DATA_STRUCTURE_TYPE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrConversionError(tPtr *testing.T) {
	var errCode = sError.ErrConversionError
	s := sError.BuildParams([]string{"PARAMETER_NAME", "DATA_STRUCTURE_TYPE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrJSONConversionError(tPtr *testing.T) {
	var errCode = sError.ErrJSONConversionError
	s := sError.BuildParams([]string{"DATA_STRUCTURE_NAME", "DATA_STRUCTURE_TYPE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrInvalidJSONError(tPtr *testing.T) {
	var errCode = sError.ErrInvalidJSON
	s := sError.BuildParams([]string{"PARAMETER_NAME"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrMapConversionError(tPtr *testing.T) {
	var errCode = sError.ErrMapConversionError
	s := sError.BuildParams([]string{"PARAMETER_NAME", "APPLICATION_PACKAGE_NAME"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test208110Error(tPtr *testing.T) {
	var errCode = 208110
	s := sError.BuildParams([]string{"OBJECT_TYPE", "SYSTEM_ID"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test208120Error(tPtr *testing.T) {
	var errCode = 208120
	s := sError.BuildParams([]string{"JSON_ARRAY", "OBJECT_TYPE", "SYSTEM_ID"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrMissingErrorMessageError(tPtr *testing.T) {
	var errCode = sError.ErrMissingErrorMessage
	s := sError.BuildParams([]string{"ERROR_MESSAGE_NUMBER"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrMissingFileError(tPtr *testing.T) {
	var errCode = sError.ErrMissingFile
	s := sError.BuildParams([]string{"FILE_NAME", "MESSAGE_RETURNED_FROM_OPEN"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrMissingEnvVariableError(tPtr *testing.T) {
	var errCode = sError.ErrMissingEnvVariable
	s := sError.BuildParams([]string{"ENVIRONMENT"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test209200Error(tPtr *testing.T) {
	var errCode = 209200
	s := sError.BuildParams([]string{"DATABASE_NAME", "DATABASE_DRIVER", "DATABASE_PORT"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrInvalidSSLModeError(tPtr *testing.T) {
	var errCode = sError.ErrInvalidSSLMode
	s := sError.BuildParams([]string{"SSL_MODE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test209230Error(tPtr *testing.T) {
	var errCode = 209230
	s := sError.BuildParams([]string{"CONNECTION_TYPE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test209520Error(tPtr *testing.T) {
	var errCode = 209520
	s := sError.BuildParams([]string{"KID"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrMissingKidInPublicKeyError(tPtr *testing.T) {
	var errCode = sError.ErrMissingKidInPublicKey
	s := sError.BuildParams([]string{"KID"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrFetchingJWKError(tPtr *testing.T) {
	var errCode = sError.ErrFetchingJWKError
	s := sError.BuildParams([]string{"REGION", "ENVIRONMENT"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrMissingURLError(tPtr *testing.T) {
	var errCode = sError.ErrMissingURL
	s := sError.BuildParams([]string{"URL_IS_MISSING"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test210098Error(tPtr *testing.T) {
	var errCode = 210098
	s := sError.BuildParams([]string{"START_VARIABLE_OUT_OF_RANGE"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test210099Error(tPtr *testing.T) {
	var errCode = 210099
	s := sError.BuildParams([]string{"START_VARIABLE_MISSING"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func Test210100Error(tPtr *testing.T) {
	var errCode = 210100
	s := sError.BuildParams([]string{"LIST_PARAMETERS"})
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, sError.EmptyMap))
}
func TestErrorDetails(tPtr *testing.T) {
	var (
		errCode    = 210100
		errDetails = make(map[string]string)
	)
	s := sError.BuildParams([]string{"LIST_PARAMETERS"})
	errDetails["test_1"] = "Test_1_Value"
	errDetails["test_2"] = "Test_2_Value"
	errDetails["test_3"] = "Test_3_Value"
	validateReply(tPtr, errCode, s, sError.GetSError(errCode, s, errDetails))
	errCode = sError.ErrSynadiaError
	validateReply(tPtr, errCode, nil, sError.GetSError(errCode, nil, errDetails))
	validateReply(tPtr, errCode, nil, sError.GetSError(errCode, nil, sError.EmptyMap))
}
func TestGenerateDocumentation(tPtr *testing.T) {
	if x, y := sError.GenerateDocumentation(); !strings.Contains(x, sError.MARKDOWNTITLEBAR) {
		tPtr.Errorf("TestGenerateDocumentation doesn't have the correct header, so there appears something wrong with the code.")
	} else {
		println(x)
		println(y)
	}
}
func validateReply(tPtr *testing.T, errCode int, params []interface{}, x sError.SoteError) {
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
