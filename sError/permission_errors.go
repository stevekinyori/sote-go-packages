package sError

import "fmt"

const (
	// Error Category
	PERMISSION_ERROR = "permission-error"

	INVALID_ISS            = 208300
	INVALID_SUBJECT        = 208310
	INVALID_TOKEN          = 208320
	INVALID_APP_CLIENT_ID  = 208340
	TOKEN_EXPIRED          = 208350
	TOKEN_INVALID          = 208355
	SEGEMNTS_COUNT_INVALID = 208356
	INVALID_CLAIM          = 208360
	MISSING_CLAIM          = 208370
)

func ErrInvalidISS(operation string, Err error) error {
	return &SoteError{
		ErrCode:   INVALID_ISS,
		ErrType:   PERMISSION_ERROR,
		FmtErrMsg: ": iss (Issuer) is not valid",
		Loc:       operation,
		Err:       Err,
	}
}

func ErrInvalidSubject(subject interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode:   INVALID_SUBJECT,
		ErrType:   PERMISSION_ERROR,
		FmtErrMsg: fmt.Sprintf(": sub (Subject: %v) was not present", subject),
		Loc:       operation,
		Err:       Err,
	}
}

// TODO? :- Similar Error Below [ASK Scott]
func ErrInvalidToken(operation string, Err error) error {
	return &SoteError{
		ErrCode:   INVALID_TOKEN,
		ErrType:   PERMISSION_ERROR,
		FmtErrMsg: ": token_use is not valid",
		Loc:       operation,
		Err:       Err,
	}
}

func ErrInvalidAppClientID(operation string, Err error) error {
	return &SoteError{
		ErrCode:   INVALID_APP_CLIENT_ID,
		ErrType:   PERMISSION_ERROR,
		FmtErrMsg: ": client id is not valid for this application",
		Loc:       operation,
		Err:       Err,
	}
}

func ErrTokenExpired(operation string, Err error) error {
	return &SoteError{
		ErrCode:   TOKEN_EXPIRED,
		ErrType:   PERMISSION_ERROR,
		FmtErrMsg: ": Token is expired",
		Loc:       operation,
		Err:       Err,
	}
}

func ErrSegmentsCountInvalid(operation string, Err error) error {
	return &SoteError{
		ErrCode:   SEGEMNTS_COUNT_INVALID,
		ErrType:   PERMISSION_ERROR,
		FmtErrMsg: ": Token contains an invalid number of segments",
		Loc:       operation,
		Err:       Err,
	}
}

func ErrInvalidClaim(claimNames interface{}, operation string, Err error) error {
	return &SoteError{
		ErrCode:   SEGEMNTS_COUNT_INVALID,
		ErrType:   PERMISSION_ERROR,
		FmtErrMsg: fmt.Sprintf(": These claims are invalid: %v", claimNames),
		Loc:       operation,
		Err:       Err,
	}
}

func ErrMissingClaim(operation string, Err error) error {
	return &SoteError{
		ErrCode:   MISSING_CLAIM,
		ErrType:   PERMISSION_ERROR,
		FmtErrMsg: ": Required claim(s) is/are missing",
		Loc:       operation,
		Err:       Err,
	}
}
