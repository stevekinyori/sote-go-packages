package sAuthentication

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

const (
	/*
		Before completing testing, make sure to put a non-expired token value in the EXPIRED_TOKEN const, so you have tested a successful case.
	*/
	// Expired Token
	EXPIRED_TOKEN = "eyJraWQiOiJjbmhHanQwVXNqZFNHNzFvUWQ1cThTRjNTb29mOHBPNU1qTThMaDdNWDlrPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJkZDg5NmVhNi03NmI2LTQ1OGYtYWYyNC0zMTAyN2JiOGQzODMiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xXzc3MTViMjlmLTVmZmMtNGE4MC04MjkxLTI0YzhjYzAwNjdkMSIsImV2ZW50X2lkIjoiZDFlODg2NDktMTI0Ni00YTIxLThjZjYtYjZmMGQ4ZmMwNDVmIiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU5ODAyNDU4NCwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfUVZQd3dDZzJjIiwiZXhwIjoxNTk4MDI4MTg0LCJpYXQiOjE1OTgwMjQ1ODQsImp0aSI6ImIwYjgzNGE3LWY4MjQtNGMwMi1hZGRiLTU5ZmUwODY2YzJkOSIsImNsaWVudF9pZCI6InR0c21yNjUzcXRma3VjZ2dvZnBtcWJic3AiLCJ1c2VybmFtZSI6ImRkODk2ZWE2LTc2YjYtNDU4Zi1hZjI0LTMxMDI3YmI4ZDM4MyJ9.R0h6otGxbqX3Novw5qHEpOvkSvzSPVHBxtXd-eb6Zq0Oq8IXMkHrgfmBZH0CiaGUzjaCHF5OyZEfqq9vfbJ0iu7YF9hGbBVHDVbJ1mywbpwQ-5l53rOEFeaXU-jshVXnNb_VbBpPrxC0Cha_w-_MJ1JF5c5_jw-93iT9azU57mwcwfc95Ro9P4mWsQ7i6Wurk1Mw7ijhTaguJceB8cRcfoYCt2xx5BuGeYBLe-5QuTbUebpkMM6WoGschwpUiZDsPhhsHXO9Tu_Exk58Ad3BSaHVOgrquF1qq6KNoCObiUIPjq2z4BELev0jH1B0KrY_0kK77IkdAokte1kmniAHuQ"
	// Invalid Signature error
	INVALID_SIGNED_TOKEN = "eyJraWQiOiJjbmhHanQwVXNqZFNHNzFvUWQ1cThTRjNTb29mOHBPNU1qTThMaDdNWDlrPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJkZDg5NmVhNi03NmI2LTQ1OGYtYWYyNC0zMTAyN2JiOGQzODMiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xXzU1MDg3NThmLWM4NmMtNDYwYS05NjMyLTBiMzBjYjE0NzBhNSIsImV2ZW50X2lkIjoiYWFkNTdhMjMtZDRjMC00M2E2LWIyZGQtNzIzZjA3ZjkxNzk5IiwiDG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU5NTQ2MTE4NSwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfUVZQd3dDZzJjIiwiZXhwIjoxNTk1NDY0Nzg1LCJpYXQiOjE1OTU0NjExODUsImp0aSI6IjFiNTYyMWNlLTk5YzAtNDczOC05YjA0LTZkOTk0ODUxYTM2MiIsImNsaWVudF9pZCI6InR0c21yNjUzcXRma3VjZ2dvZnBtcWJic3AiLCJ1c2VybmFtZSI6ImRkODk2ZWE2LTc2YjYtNDU4Zi1hZjI0LTMxMDI3YmI4ZDM4MyJ9.Y2ypGn1cVPvjoOMxE5jbm7HwAWghN4zX2RJ8UxGwfwYLsHEPpGgJHpCqHJal1jufe-ciM_XQM7QvPFYFO5BL0wzDtHZmx0ntCP26Tt6lnwi7a_XflWlhb48CPE4WIk_1TcgKXVwquIhf437NvmsfXo_ctoSCJ4EGPYN4BUQCugYWmsMdh5aFzXVS3nz9DEHJVAh5IB7C3N9TTYOmplUVIoRLLfCyk16eMhO-I3zv2T3PKTiM01vUe_7zxaXqPLdG52GQ_U-wmJueMhYABWkKDJtVdWqYn7RS-dJckbEozbdalMqwyIe9ejMz8MlMthVTq6qaDMD8-n26WlIAA09VUw"
	// Invalid type error
	TOKEN_INVALID = "eyJraWQiOiJlOCt4TW4rOGYrZmlIXC9OZDNDZGNxOVRvU3FPKzdZYldcL1wvSUxCYVJyTElNPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiI1ZDUxNDdlMi01N2ZjLTQ4YTYtYjQ5My0xNzgzOTMxYWU5YzAiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xX2IyODdkMTQ0LTJhNjEtNDAzZC04MmNjLTkyYWY3ZmZhYmZjNiIsImV2ZW50X2lkIjoiZjFhNTYwMzItYjYyZi00OWFhLWI0MmEtMmI3YzIwZmVhY2VkIiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU5NTQ2MTQyMSwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfZnBaQ3lxbFFOIiwiZXhwIjoxNTk1NDY1MDIxLCJpYXQiOjE1OTU0NjE0MjEsImp0aSI6ImU5ODNkNDJjLWM4YWUtNDZmNi05MmFlLTE1YzczMGEwODRhZCIsImNsaWVudF9pZCI6IjQ5ZzUwOGgwanIyczJ0bmJpM2toZ3Rkc3AyIiwidXNlcm5hbWUiOiI1ZDUxNDdlMi01N2ZjLTQ4YTYtYjQ5My0xNzgzOTMxYWU5YzAifQ.T6IOtBsU0QLzOplbEWqa1QRAS7nMTknWP-meYaE4WQybLzFr-9-dPGMA3spFoQfgjD40Mxl5CFtkVmDSn1W8Yo3wvVATt7t1220YV2WAIQLf4SKghS-dwans10BpAnC5BwLgEE_sDPaJ064mV56xbOO9R-ePmNws-_qYp_R615RtHwJPtQedVwFRH7W63wS6ATs0wBaQ6McIAu1QSyoOj6ePegSzhJd_bhrDD4i42GqC2rb0rca_IRYxd3Ev44Rjx9QGZHni-BYw04jBKSFuYUvtoIpPl9gD4PzxSU5d_4lRX264uBY9F8cgeVmV97JvzVSdoAqaxfpj_drX0789Sg"
	// Only two segments
	TOKEN_MISSING_SEGMENT = "eyJraWQiOiJjbmhHanQwVXNqZFNHNzFvUWQ1cThTRjNTb29mOHBPNU1qTThMaDdNWDlrPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJkZDg5NmVhNi03NmI2LTQ1OGYtYWYyNC0zMTAyN2JiOGQzODMiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xX2Q3MmQ5NzJiLTRiNGQtNGJjMi05YTU0LWNmZmJlOTU1YTExMiIsImV2ZW50X2lkIjoiNWRmMGUwMWMtMTAyMy00OGRmLThjOTgtM2Y1MWI4N2Y1Y2E0IiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU5Nzg1ODE5NCwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfUVZQd3dDZzJjIiwiZXhwIjoxNTk3ODYxNzk0LCJpYXQiOjE1OTc4NTgxOTQsImp0aSI6IjI3NTkxMWYxLWJhYmYtNDNmNS04YjE4LWMxZWFkMTljMmY1YiIsImNsaWVudF9pZCI6InR0c21yNjUzcXRma3VjZ2dvZnBtcWJic3AiLCJ1c2VybmFtZSI6ImRkODk2ZWE2LTc2YjYtNDU4Zi1hZjI0LTMxMDI3YmI4ZDM4MyJ9ECr1bgVGRNYscxP10HCvu3WE6RQKtwOvs3yoGUGsDr9TDi_hWHc7VUxcDXvjxYy5C9zm75Jg4CIjbp42GhA2L82h0cAyI3eBiWq_GSIMdKo1ZNWVrIrksYDLgHJUaRqAuO4ui2AEsN5P2fXlLPmOngMOEQbVo22ybUeZpM83Q0R08A92z3cAxW2MPjCvYzdf1MPO3OCS2OpK54SzRQftR_mRqWhSMYM6NHn_hlE33quhCQ1AMFrQsoZOT4ujwA_h0WHWhxikPnJHJIuqDatWQ4PPwF4DcbGFSZfOm3fDOY9GBXXypjywUIn7cRCTLEm3YN1-HzYtBoLq9VZF9uxLKA"
	GOOD_KID              = "cnhGjt0UsjdSG71oQd5q8SF3Soof8pO5MjM8Lh7MX9k="
	FAKE_TOKEN            = "TOKEN"
	// KIDBAD          = "e8+xMn+8f+fiH/Nd3Cdcq9ToSqO+7YbW//"
	// Application
	SDCC = "sdcc"
)

func init() {
	sLogger.SetLogMessagePrefix("sjwt_test.go")
}

func TestValidToken(t *testing.T) {
	type args struct {
		tEnvironment string
		rawToken     string
	}
	tests := []struct {
		name              string
		args              args
		wantErr           bool
		ExpectedErrorCode interface{}
	}{

		{
			name: "positive case: valid token",
			args: args{
				tEnvironment: sConfigParams.DEVELOPMENT,
				rawToken:     EXPIRED_TOKEN,
			},
			// We expect an error here because we have no programmatic way of generating a jwt token
			// Hence a token expired error is always returned
			wantErr:           true,
			ExpectedErrorCode: sError.TOKEN_EXPIRED,
		},
		{
			name: "negative case: token missing segment",
			args: args{
				tEnvironment: sConfigParams.DEVELOPMENT,
				rawToken:     TOKEN_MISSING_SEGMENT,
			},
			wantErr:           true,
			ExpectedErrorCode: sError.SEGEMNTS_COUNT_INVALID,
		},
		{
			name: "negative case: fake token",
			args: args{
				tEnvironment: sConfigParams.DEVELOPMENT,
				rawToken:     FAKE_TOKEN,
			},
			wantErr:           true,
			ExpectedErrorCode: sError.SEGEMNTS_COUNT_INVALID,
		},
		{
			name: "negative case: invalidly signed token",
			args: args{
				tEnvironment: sConfigParams.DEVELOPMENT,
				rawToken:     INVALID_SIGNED_TOKEN,
			},
			wantErr:           true,
			ExpectedErrorCode: sError.INVALID_TOKEN,
		},
		{
			name: "negative case: valid token missing params",
			args: args{
				tEnvironment: sConfigParams.DEVELOPMENT,
				rawToken:     "",
			},
			wantErr:           true,
			ExpectedErrorCode: sError.PARAMS_MUST_BE_PROVIDED,
		},
		{
			name: "negative case: invalid token",
			args: args{
				tEnvironment: sConfigParams.DEVELOPMENT,
				rawToken:     TOKEN_INVALID,
			},
			wantErr:           true,
			ExpectedErrorCode: sError.SEGEMNTS_COUNT_INVALID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidToken(tt.args.tEnvironment, tt.args.rawToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if e, ok := err.(*sError.SoteError); ok {
				if e.ErrCode != tt.ExpectedErrorCode {
					t.Errorf("ValidToken() error-code = %v, expected-error-code = %v", e.ErrCode, tt.ExpectedErrorCode)
				}
			} else {
				t.Errorf("error returned was not successfully cast")
			}

		})
	}
}

func Test_matchKid(t *testing.T) {
	type args struct {
		tEnvironment string
		kid          string
	}
	tests := []struct {
		name    string
		args    args
		wantKey jwk.Key
		wantErr bool
	}{
		{
			name: "positive case: KID successfully matched",
			args: args{
				tEnvironment: sConfigParams.DEVELOPMENT,
				kid:          GOOD_KID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, err := matchKid(tt.args.tEnvironment, tt.args.kid)
			if (err != nil) != tt.wantErr {
				t.Errorf("matchKid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (gotKey == nil) != tt.wantErr {
				t.Errorf("wantKey is nil, want non-nil Key: %v", !tt.wantErr)
			}
		})
	}
}

func Test_getPublicKey(t *testing.T) {
	type args struct {
		tEnvironment string
	}
	tests := []struct {
		name    string
		args    args
		want    jwk.Set
		wantErr bool
	}{

		{
			name: "positive case: successfully got public key",
			args: args{
				tEnvironment: sConfigParams.DEVELOPMENT,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPublicKey(tt.args.tEnvironment)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("result of type jwk.Set is nil, want non-nil Key: %v", !tt.wantErr)
			}
		})
	}
}

func Test_fetchPublicKey(t *testing.T) {

	region, soteErr := sConfigParams.GetRegion()
	if soteErr.ErrCode != nil {
		t.Fatalf("could not get region: %v", soteErr.FmtErrMsg)
	}

	userPoolId, soteErr := sConfigParams.GetUserPoolId(sConfigParams.DEVELOPMENT)
	if soteErr.ErrCode != nil {
		t.Fatalf("could not get user pool id: %v", soteErr.FmtErrMsg)
	}

	type args struct {
		region       string
		userPoolId   string
		tEnvironment string
	}
	tests := []struct {
		name       string
		args       args
		wantKeySet jwk.Set
		wantErr    bool
	}{

		{
			name: "positive case: successfully fetched public key",
			args: args{
				region:       region,
				userPoolId:   userPoolId,
				tEnvironment: sConfigParams.DEVELOPMENT,
			},
			wantErr: false,
		},

		{
			name: "negative case: could not fetch public key",
			args: args{
				region:       "SCOTT_LAND",
				userPoolId:   userPoolId,
				tEnvironment: sConfigParams.DEVELOPMENT,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKeySet, err := fetchPublicKey(tt.args.region, tt.args.userPoolId, tt.args.tEnvironment)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (gotKeySet == nil) != tt.wantErr {
				t.Errorf("gotKeySet is nil, want non-nil Key: %v", !tt.wantErr)
			}
		})
	}
}

func Test_validateClaims(t *testing.T) {

	var invalidClaims jwt.MapClaims = make(map[string]interface{})
	invalidClaims["scope"] = "scott.fake.scope"

	var invalidAppClientIDClaim jwt.MapClaims = make(map[string]interface{})
	invalidAppClientIDClaim["client_id"] = "scott.fake.client_id"

	var invalidTokenUseIDClaim jwt.MapClaims = make(map[string]interface{})
	invalidAppClientIDClaim["token_use"] = "scott.fake.use"

	var invalidISSClaim jwt.MapClaims = make(map[string]interface{})
	invalidAppClientIDClaim["iss"] = "scott.fake.iss"

	type args struct {
		claims       jwt.MapClaims
		tEnvironment string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "negative case: missing claims",
			args: args{
				claims:       nil,
				tEnvironment: sConfigParams.STAGING,
			},
			wantErr: true,
		},
		{
			name: "negative case: invalid claims",
			args: args{
				claims:       invalidClaims,
				tEnvironment: sConfigParams.STAGING,
			},
			wantErr: true,
		},
		{
			name: "negative case: invalid client id claims",
			args: args{
				claims:       invalidAppClientIDClaim,
				tEnvironment: sConfigParams.STAGING,
			},
			wantErr: true,
		},
		{
			name: "negative case: invalid iss claims",
			args: args{
				claims:       invalidISSClaim,
				tEnvironment: sConfigParams.STAGING,
			},
			wantErr: true,
		},
		{
			name: "negative case: invalid token use claims",
			args: args{
				claims:       invalidTokenUseIDClaim,
				tEnvironment: sConfigParams.STAGING,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateClaims(tt.args.claims, tt.args.tEnvironment); (err != nil) != tt.wantErr {
				t.Errorf("validateClaims() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
