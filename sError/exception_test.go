package sError

import (
	"fmt"
	"strings"
	"testing"
)

func TestExceptionParams(t *testing.T) {
	compare(t, fmt.Sprintf("%v"+ItemAlreadyExists.fmtErrMsg, ItemAlreadyExists.errCode, "A"), ItemAlreadyExists.SetParams("A").fmtErrMsg)
	compare(t, fmt.Sprintf("%v"+InvalidParameterCount.fmtErrMsg, InvalidParameterCount.errCode, 1, 2), ItemAlreadyExists.SetParams("A", "B").fmtErrMsg)
}

func TestExceptionGetException(t *testing.T) {
	compare(t, ItemAlreadyExists.fmtErrMsg, GetException(100000).String())
}

func TestExceptionString(t *testing.T) {
	ex := ItemAlreadyExists.SetParams("A").SetDetails("File Already Exists")
	ex.LogInfo()
	compare(t, "100000: A already exists\nFile Already Exists", ex.String())
}

func TestExceptionJson(t *testing.T) {
	exJson := ItemAlreadyExists.SetParams("A").SetDetails("File Already Exists").GenerateJson()
	props := []string{
		`"ErrCode": 100000`,
		`"ErrorDetails": "File Already Exists"`,
		`"FmtErrMsg": "100000: A already exists"`,
		`"ParamDescription": "Item Name"`,
		`"ParamCount": 1`,
	}
	for _, s := range props {
		if !strings.Contains(exJson, s) {
			t.Fatalf("Cannot find '%s' in JSON: %v", s, exJson)
		}
	}
}

func TestExceptionDocumentation(t *testing.T) {
	x := ItemAlreadyExists
	xMarkDown := fmt.Sprintf("| %v | %v | %v |\n", x.errCode, x.paramDescription, x.fmtErrMsg)
	xFuncComments := fmt.Sprintf("\t\t%v\t%v > %v\n", x.errCode, x.paramDescription, x.fmtErrMsg)
	markDown, funcComments := GenerateDoc()
	if !strings.Contains(markDown, xMarkDown) {
		t.Fatal("Cannot find markDown: " + xMarkDown)
	}
	if !strings.Contains(funcComments, xFuncComments) {
		t.Fatal("Cannot find funcComments: " + xFuncComments)
	}
}

func TestValidateExceptions(t *testing.T) {
	validate(t, ItemAlreadyExists)
	validate(t, NotAuthorized)
	validate(t, DirtyRead)
	validate(t, CanceledComplete)
	validate(t, TimeOut)
	validate(t, ItemNotFound)
	validate(t, UnexpectedError)
	validate(t, TableMissing)
	validate(t, InvalidDataType)
	validate(t, RequiredValueMissing)
	validate(t, LinkedParameterValueMissing)
	validate(t, ParameterLockOtherParameterSet)
	validate(t, ParameterMustBeSetOrNull)
	validate(t, ParametersMustBeProvided)
	validate(t, ParameterMustBeSet)
	validate(t, ThreeParametersMustBeSet)
	validate(t, ParameterMustBeEmptyWhenParameterSet)
	validate(t, BadHTTPRequest)
	validate(t, InvalidEnvironmentForAPI)
	validate(t, QuickSightError)
	validate(t, DatabaseError)
	validate(t, SqlError)
	validate(t, CognitoError)
	validate(t, InvalidParameterCount)
	validate(t, AwsSESError)
	validate(t, AwsSTSError)
	validate(t, JetStreamError)
	validate(t, NatsSubscriptionError)
	validate(t, NatsStreamPointerMissing)
	validate(t, NatsStreamCreateError)
	validate(t, NatsConsumerCreateError)
	validate(t, NatsInvalidConsumerSubjectFilter)
	validate(t, ParameterNotNumeric)
	validate(t, ParameterToSmall)
	validate(t, ParameterNotString)
	validate(t, ParameterNotFloat)
	validate(t, ParameterNotArray)
	validate(t, ParameterNotJsonString)
	validate(t, InvalidEmailFormat)
	validate(t, ParameterNotDate)
	validate(t, ParameterNotTimestamp)
	validate(t, ParameterInvalidSize)
	validate(t, JsonConversionError)
	validate(t, ParameterNotMap)
	validate(t, MissingErrorNumber)
	validate(t, InvalidISS)
	validate(t, InvalidSubject)
	validate(t, InvalidToken)
	validate(t, InvalidAppClientId)
	validate(t, TokenExpired)
	validate(t, TokenInvalid)
	validate(t, SegmentsCountInvalid)
	validate(t, InvalidClaim)
	validate(t, MissingClaim)
	validate(t, EnvFileMissing)
	validate(t, FileNotFound)
	validate(t, EnvironmentMissing)
	validate(t, EnvironmentInvalid)
	validate(t, InvalidDBConnection)
	validate(t, InvalidDBAuthentication)
	validate(t, InvalidDBSSLMode)
	validate(t, InvalidConnectionType)
	validate(t, NoDBConnection)
	validate(t, NatsNkeyMissing)
	validate(t, NoNATSConnection)
	validate(t, UnexpectedSign)
	validate(t, KidNotFound)
	validate(t, KidMissingFromToken)
	validate(t, KidDoesNotMatchPublicKeySet)
	validate(t, InvalidRegion)
	validate(t, InvalidURL)
	validate(t, OutOfValidRange)
}

func validate(t *testing.T, ex Exception) {
	var (
		newEx     Exception
		fmtErrMsg = ex.fmtErrMsg
	)
	if ex.paramCount == 0 {
		newEx = ex.SetDetails("Custom Error Message")
		compare(t, "Custom Error Message", newEx.errorDetails)
	} else if ex.paramCount > 0 {
		params := []interface{}{"A", "B", "C", "D", "E", "F"}[:ex.paramCount]
		newEx = ex.SetParams(params...)
		compare(t, fmt.Sprint(ex.errCode)+fmt.Sprintf(fmtErrMsg, params...), newEx.fmtErrMsg)
	}
	compare(t, fmtErrMsg, ex.fmtErrMsg)
	compare(t, "", ex.errorDetails)
	compare(t, ex.GetCode(), newEx.GetCode())
}

func compare(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Fatalf("Not equal:\nexpected: %v\nactual:   %v", expected, actual)
	}
}
