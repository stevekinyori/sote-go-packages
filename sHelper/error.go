package sHelper

import (
	"gitlab.com/soteapps/packages/v2021/sError"
)

type sErrorHelper struct {
	errorDetails map[string]string
}

func NewError(errorDetails ...map[string]string) sErrorHelper {
	if len(errorDetails) > 0 {
		return sErrorHelper{errorDetails: errorDetails[0]}
	}
	return sErrorHelper{errorDetails: sError.EmptyMap}
}

func (err *sErrorHelper) factory(code int, params ...interface{}) sError.SoteError {
	return sError.GetSError(code, params, err.errorDetails)
}

// User_Error's
func (err sErrorHelper) AlreadyExists(itemName interface{}) sError.SoteError {
	return err.factory(100000, itemName) //"100000: %v already exists"
}

func (err sErrorHelper) ItemNotFound(itemName interface{}) sError.SoteError {
	return err.factory(109999, itemName) //"109999: %v was/were not found"
}

// Process_Error's
func (err sErrorHelper) SqlError(error string) sError.SoteError {
	err.errorDetails = map[string]string{"SQL ERROR": error}
	return err.factory(200999) //"200999: SQL error - see Details"
}

func (err sErrorHelper) MustBePopulated(param interface{}) sError.SoteError {
	return err.factory(200513, param) //"200513: %v must be populated"
}

func (err sErrorHelper) MustBeType(param string, types interface{}) sError.SoteError {
	return err.factory(200200, param, types) //"200200: %v must be of type %v"
}

func (err sErrorHelper) AllowValues(param string, name interface{}, values interface{}) sError.SoteError {
	return err.factory(200250, param, name, values) //"200250: %v (%v) must contain one of these values: %v"
}

// Configuration_Issue's
func (err sErrorHelper) NoDbConnection() sError.SoteError {
	return err.factory(209299) //"209299: No database connection has been established"
}

func (err sErrorHelper) FileNotFound(fileName string, exception string) sError.SoteError {
	return err.factory(209010, fileName, exception) //"209010: %v file was not found. Message return: %v"
}

// Content_Error's
func (err sErrorHelper) InvalidJson(fileName string) sError.SoteError {
	return err.factory(207110, fileName) //"207110: %v couldn't be parsed - Invalid JSON error"
}

func (err sErrorHelper) InvalidEmailAddress(fieldName string, value interface{}) sError.SoteError {
	return err.factory(207050, value, fieldName) //"207050: %v (%v) is not a valid email address"
}

// NATS_Error's
func (err sErrorHelper) InvalidParameters(params ...interface{}) sError.SoteError {
	return err.factory(206200, params...) //"206200: Message doesn't match signature. Sender must provide the following parameter names: %v"
}

// General_Error's
func (err sErrorHelper) InternalError() sError.SoteError {
	return err.factory(210599) //"210599: Business Service error has occurred that is not expected."
}

// Permission_Error's
func (err sErrorHelper) InvalidToken() sError.SoteError {
	return err.factory(208355) //"208355: Token is invalid"
}

func (err sErrorHelper) ExpiredToken() sError.SoteError {
	return err.factory(208350) //"208350: Token is expired"
}
