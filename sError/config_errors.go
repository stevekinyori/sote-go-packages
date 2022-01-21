package sError

import "fmt"

const (
	// Error Category
	CONFIGURATION_ERROR = "config-error"

	FAILED_DB_CONNECTION               = 209299
	ENV_FILE_MISSING                   = 209000
	FILE_NOT_FOUND                     = 209010
	ENVIRONMENT_MISSING                = 209100
	ENVIRONMENT_INVALID                = 209110
	INVALID_DB_CONNECTION              = 209200
	INVALID_DB_AUTHENTICATION          = 209210
	INVALID_DB_SSL_MODE                = 209220
	INVALID_CONNECTION_TYPE            = 209230
	NATS_NKEY_MISSING                  = 209398
	NO_NATS_CONNECTION                 = 209499
	UNEXPECTED_SIGN                    = 209500
	K_ID_NOT_FOUND                     = 209510
	K_ID_MISSING_FROM_TOKEN            = 209520
	K_ID_DOES_NOT_MATCH_PUBLIC_KEY_SET = 209521
	INVALID_REGION                     = 210030
	INVALID_URL                        = 210090
	OUT_OF_VALID_RANGE                 = 210098
)

func ErrDBConnectionFailed(message string, operation string, Err error) error {
	return &SoteError{
		ErrCode:   FAILED_DB_CONNECTION,
		ErrType:   CONFIGURATION_ERROR,
		FmtErrMsg: message,
		Loc:       operation,
		Err:       Err,
	}
}

func ErrEnvFileMissing(operation string, Err error) error {
	return &SoteError{
		ErrCode:   ENV_FILE_MISSING,
		ErrType:   CONFIGURATION_ERROR,
		FmtErrMsg: ": .env files are missing",
		Loc:       operation,
		Err:       Err,
	}
}

func ErrFileNotFound(fileName, message interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode:   FILE_NOT_FOUND,
		ErrType:   CONFIGURATION_ERROR,
		FmtErrMsg: fmt.Sprintf(": %v file was not found. Message return: %v", fileName, message),
		Loc:       operation,
		Err:       Err,
	}
}

func ErrEnvironmentMissing(envName interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode:   ENVIRONMENT_MISSING,
		ErrType:   CONFIGURATION_ERROR,
		FmtErrMsg: fmt.Sprintf(": environment variable is missing (%v)", envName),
		Loc:       operation,
		Err:       Err,
	}
}

func ErrEnvironmentInvalid(envName interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode:   ENVIRONMENT_INVALID,
		ErrType:   CONFIGURATION_ERROR,
		FmtErrMsg: fmt.Sprintf(": environment value (%v) is invalid", envName),
		Loc:       operation,
		Err:       Err,
	}
}

func ErrInvalidDBConnection(DBName, DBDriverName, Port interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode: INVALID_DB_CONNECTION,
		ErrType: CONFIGURATION_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": Unable to connect to database %v using driver %v on port %v",
			DBName, DBDriverName, Port),
		Loc: operation,
		Err: Err,
	}
}

func ErrInvalidDBAuthentication(operation string, Err error) error {
	return &SoteError{
		ErrCode:   INVALID_DB_AUTHENTICATION,
		ErrType:   CONFIGURATION_ERROR,
		FmtErrMsg: ": Unable to pass database authentication",
		Loc:       operation,
		Err:       Err,
	}
}

func ErrInvalidDBSSLMode(operation string, Err error) error {
	return &SoteError{
		ErrCode:   INVALID_DB_SSL_MODE,
		ErrType:   CONFIGURATION_ERROR,
		FmtErrMsg: ": Only disable, allow, prefer and required are supported.",
		Loc:       operation,
		Err:       Err,
	}
}

func ErrInvalidConnectionType(operation string, Err error) error {
	return &SoteError{
		ErrCode:   INVALID_CONNECTION_TYPE,
		ErrType:   CONFIGURATION_ERROR,
		FmtErrMsg: ": Only single or pool are supported.",
		Loc:       operation,
		Err:       Err,
	}
}

func ErrNATSNKeyMissing(operation string, Err error) error {
	return &SoteError{
		ErrCode:   NATS_NKEY_MISSING,
		ErrType:   CONFIGURATION_ERROR,
		FmtErrMsg: ": no nkey seed found",
		Loc:       operation,
		Err:       Err,
	}
}

func ErrNoNATSConnection(operation string, Err error) error {
	return &SoteError{
		ErrCode:   NO_NATS_CONNECTION,
		ErrType:   CONFIGURATION_ERROR,
		FmtErrMsg: ": No nats connection has been established",
		Loc:       operation,
		Err:       Err,
	}
}

func ErrUnexpectedSign(operation string, Err error) error {
	return &SoteError{
		ErrCode:   UNEXPECTED_SIGN,
		ErrType:   CONFIGURATION_ERROR,
		FmtErrMsg: ": Unexpected signing method",
		Loc:       operation,
		Err:       Err,
	}
}

func ErrKIDNotFound(operation string, Err error) error {
	return &SoteError{
		ErrCode:   K_ID_NOT_FOUND,
		ErrType:   CONFIGURATION_ERROR,
		FmtErrMsg: ": Kid header not found",
		Loc:       operation,
		Err:       Err,
	}
}

func ErrKIDMissingFromToken(kID, operation string, Err error) error {
	return &SoteError{
		ErrCode:   K_ID_MISSING_FROM_TOKEN,
		ErrType:   CONFIGURATION_ERROR,
		FmtErrMsg: fmt.Sprintf(": key (%v) was not found in token", kID),
		Loc:       operation,
		Err:       Err,
	}
}

func ErrKIDDoesNotMatchPublicKeySet(kID, operation string, Err error) error {
	return &SoteError{
		ErrCode:   K_ID_DOES_NOT_MATCH_PUBLIC_KEY_SET,
		ErrType:   CONFIGURATION_ERROR,
		FmtErrMsg: fmt.Sprintf(": Kid (%v) was not found in public key set", kID),
		Loc:       operation,
		Err:       Err,
	}
}

func ErrInvalidRegion(region, environment, operation string, Err error) error {
	return &SoteError{
		ErrCode: INVALID_REGION,
		ErrType: CONFIGURATION_ERROR,
		FmtErrMsg: fmt.Sprintf(
			": Failed to fetch remote JWK (status = 404) for %v region %v environment",
			region, environment),
		Loc: operation,
		Err: Err,
	}
}

func ErrInvalidURL(param, operation string, Err error) error {
	return &SoteError{
		ErrCode:   INVALID_URL,
		ErrType:   CONFIGURATION_ERROR,
		FmtErrMsg: fmt.Sprintf(": URL is missing (%v)", param),
		Loc:       operation,
		Err:       Err,
	}
}

func ErrOutofValidRange(param, operation string, Err error) error {
	return &SoteError{
		ErrCode:   OUT_OF_VALID_RANGE,
		ErrType:   CONFIGURATION_ERROR,
		FmtErrMsg: fmt.Sprintf(": Start up parameter is out of value range (%v)", param),
		Loc:       operation,
		Err:       Err,
	}
}
