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
func TestErrItemAlreadyUpdatedError(tPtr *testing.T) {
	var errCode = ErrItemAlreadyUpdated
	validateReply(tPtr, errCode, nil, GetSError(errCode, nil, EmptyMap))
}
func TestErrCancelOrCompleteError(tPtr *testing.T) {
	var errCode = ErrCancelOrComplete
	s := BuildParams([]string{"OBJECT"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInactiveItemError(tPtr *testing.T) {
	var errCode = ErrInactiveItem
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
func TestErrMissingConditionalParameterError(tPtr *testing.T) {
	var errCode = ErrMissingConditionalParameter
	s := BuildParams([]string{"PARAMETER_NAME_MISSING", "PARAMETER_NAME", "PARAMETER_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrConditionalNonUpdatableParameterError(tPtr *testing.T) {
	var errCode = ErrConditionalNonUpdatableParameter
	s := BuildParams([]string{"PARAMETER_NAME", "FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrExpectedTwoParametersOrNullError(tPtr *testing.T) {
	var errCode = ErrExpectedTwoParametersOrNull
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
func TestErrConditionalEmptyFieldError(tPtr *testing.T) {
	var errCode = ErrConditionalEmptyFieldError
	s := BuildParams([]string{"PARAMETER_NAME", "ANOTHER_PARAMETER_NAME"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrBadHTTPRequestError(tPtr *testing.T) {
	var errCode = ErrBadHTTPRequest
	s := BuildParams([]string{"DETAILED_MESSAGE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidAPIError(tPtr *testing.T) {
	var errCode = ErrInvalidAPI
	s := BuildParams([]string{"ENVIRONMENT"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrUnsupportedParameterCountError(tPtr *testing.T) {
	var errCode = ErrUnsupportedParameterCount
	s := BuildParams([]string{"NAME", "APP_PACKAGE_NAME"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidParameterCountError(tPtr *testing.T) {
	var errCode = ErrInvalidParameterCount
	s := BuildParams([]string{"PROVIDED_PARAMETER_COUNT", "EXPECTED_PARAMETER_COUNT"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrMissingKeyError(tPtr *testing.T) {
	var errCode = ErrMissingKey
	s := BuildParams([]string{"KEY_NAME"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrMissingKeyValueError(tPtr *testing.T) {
	var errCode = ErrMissingKeyValue
	s := BuildParams([]string{"KEY_NAME"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidMsgSignatureError(tPtr *testing.T) {
	var errCode = ErrInvalidMsgSignature
	s := BuildParams([]string{"PARAMETER_LIST"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrStreamCreationFailureError(tPtr *testing.T) {
	var errCode = ErrStreamCreationFailure
	s := BuildParams([]string{"PARAMETER"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrConsumerCreationFailureError(tPtr *testing.T) {
	var errCode = ErrConsumerCreationFailure
	s := BuildParams([]string{"PARAMETER", "ANOTHER_PARAMETER"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidConsumerSubjectFilterError(tPtr *testing.T) {
	var errCode = ErrInvalidConsumerSubjectFilter
	s := BuildParams([]string{"PARAMETER", "ANOTHER_PARAMETER"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrNotNumericError(tPtr *testing.T) {
	var errCode = ErrNotNumeric
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrNotMinimumLengthError(tPtr *testing.T) {
	var errCode = ErrNotMinimumLength
	s := BuildParams([]string{"FIELD_NAME", "MINIMAL_LENGTH"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrNotStringError(tPtr *testing.T) {
	var errCode = ErrNotString
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrNotFloatError(tPtr *testing.T) {
	var errCode = ErrNotFloat
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrNotArrayError(tPtr *testing.T) {
	var errCode = ErrNotArray
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrNotJSONStringError(tPtr *testing.T) {
	var errCode = ErrNotJSONString
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidEmailError(tPtr *testing.T) {
	var errCode = ErrInvalidEmail
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrForbiddenSpecialCharactersError(tPtr *testing.T) {
	var errCode = ErrForbiddenSpecialCharacters
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrOnlyUnderscoreSpecialCharacterAllowedError(tPtr *testing.T) {
	var errCode = ErrOnlyUnderscoreSpecialCharacterAllowed
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidDateError(tPtr *testing.T) {
	var errCode = ErrInvalidDate
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidTimestampError(tPtr *testing.T) {
	var errCode = ErrInvalidTimestamp
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidSizeError(tPtr *testing.T) {
	var errCode = ErrInvalidSize
	s := BuildParams([]string{"FIELD_NAME", "FIELD_VALUE", "SMALL_LARGE", "MIN_MAX", "SIZE_EXPECTED", "SIZE_PROVIDED"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrMustBeGreaterOrLessThanError(tPtr *testing.T) {
	var errCode = ErrMustBeGreaterOrLessThan
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
func TestErrUnnecessaryRecordUpdateError(tPtr *testing.T) {
	var errCode = ErrUnnecessaryRecordUpdate
	s := BuildParams([]string{"OBJECT_TYPE", "SYSTEM_ID"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrMissingJSONRecordKeyValueError(tPtr *testing.T) {
	var errCode = ErrMissingJSONRecordKeyValue
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
func TestErrInvalidDBConnectionPortError(tPtr *testing.T) {
	var errCode = ErrInvalidDBConnectionPort
	s := BuildParams([]string{"DATABASE_NAME", "DATABASE_DRIVER", "DATABASE_PORT"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidSSLModeError(tPtr *testing.T) {
	var errCode = ErrInvalidSSLMode
	s := BuildParams([]string{"SSL_MODE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidConnectionTypeError(tPtr *testing.T) {
	var errCode = ErrInvalidConnectionType
	s := BuildParams([]string{"CONNECTION_TYPE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrMissingTokenKeyError(tPtr *testing.T) {
	var errCode = ErrMissingTokenKey
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
func TestErrInvalidStartUpValueRangeError(tPtr *testing.T) {
	var errCode = ErrInvalidStartUpValueRange
	s := BuildParams([]string{"START_VARIABLE_OUT_OF_RANGE"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrMissingStartUpParameterError(tPtr *testing.T) {
	var errCode = ErrMissingStartUpParameter
	s := BuildParams([]string{"START_VARIABLE_MISSING"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrInvalidAPICallSignatureError(tPtr *testing.T) {
	var errCode = ErrInvalidAPICallSignature
	s := BuildParams([]string{"LIST_PARAMETERS"})
	validateReply(tPtr, errCode, s, GetSError(errCode, s, EmptyMap))
}
func TestErrorDetails(tPtr *testing.T) {
	var (
		errCode    = ErrInvalidAPICallSignature
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
