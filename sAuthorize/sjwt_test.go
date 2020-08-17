package sAuthorize

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	"gitlab.com/soteapps/packages/v2020/sConfigParams"
	"gitlab.com/soteapps/packages/v2020/sError"
)

const (
	// Token
	TOKENEXPIRED    = "eyJraWQiOiJjbmhHanQwVXNqZFNHNzFvUWQ1cThTRjNTb29mOHBPNU1qTThMaDdNWDlrPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJkZDg5NmVhNi03NmI2LTQ1OGYtYWYyNC0zMTAyN2JiOGQzODMiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xX2M0NWYwOWFhLTYwNDEtNDdkMC1hNTNhLWYyYTg3NTE0M2U1YiIsImV2ZW50X2lkIjoiZWQxN2NmZGMtZGMyOC00ZmFhLTk4NjctZTIxYzJmMDc3ZTU2IiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU5NzQzODc5MSwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfUVZQd3dDZzJjIiwiZXhwIjoxNTk3NDQyMzkxLCJpYXQiOjE1OTc0Mzg3OTEsImp0aSI6IjE0YTY5NTczLTVhYjQtNGQyYS04NzM1LTZjYjExZGNkMGZkNSIsImNsaWVudF9pZCI6InR0c21yNjUzcXRma3VjZ2dvZnBtcWJic3AiLCJ1c2VybmFtZSI6ImRkODk2ZWE2LTc2YjYtNDU4Zi1hZjI0LTMxMDI3YmI4ZDM4MyJ9.XHGq7CPs2kcoX4HpnuM5TF7O6_ids5XTjlU5k7C1mNeQBfZgHhMpWvdNqc3S_PbWozVI927xZVtzRB6rG6O1Uutem7CzUpr38mSjCf6yrWqaV2NLWQN9lSPc7MgvtihACl082emgxHMKi-Exw_N5Ikk60MeC1-ChsCFas2h5_HHT-WLNVS16Nn2X_UyRZIc3CADFr88PsmZYvKC_Lnjh3CpnExJl3p88suVti2_pR_GdVsc_11agfnrtWFh52uri4rkE6i9yCLusEb0OhPTqmU8avACGlzl6TAbSTwftNA1Hfu5efG5eqH3tk5-OQY1P0IKDPLekCGnpGctp48B8Fg"
	TOKENINVALIDSIG = "eyJraWQiOiJjbmhHanQwVXNqZFNHNzFvUWQ1cThTRjNTb29mOHBPNU1qTThMaDdNWDlrPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJkZDg5NmVhNi03NmI2LTQ1OGYtYWYyNC0zMTAyN2JiOGQzODMiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xXzU1MDg3NThmLWM4NmMtNDYwYS05NjMyLTBiMzBjYjE0NzBhNSIsImV2ZW50X2lkIjoiYWFkNTdhMjMtZDRjMC00M2E2LWIyZGQtNzIzZjA3ZjkxNzk5IiwiDG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU5NTQ2MTE4NSwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfUVZQd3dDZzJjIiwiZXhwIjoxNTk1NDY0Nzg1LCJpYXQiOjE1OTU0NjExODUsImp0aSI6IjFiNTYyMWNlLTk5YzAtNDczOC05YjA0LTZkOTk0ODUxYTM2MiIsImNsaWVudF9pZCI6InR0c21yNjUzcXRma3VjZ2dvZnBtcWJic3AiLCJ1c2VybmFtZSI6ImRkODk2ZWE2LTc2YjYtNDU4Zi1hZjI0LTMxMDI3YmI4ZDM4MyJ9.Y2ypGn1cVPvjoOMxE5jbm7HwAWghN4zX2RJ8UxGwfwYLsHEPpGgJHpCqHJal1jufe-ciM_XQM7QvPFYFO5BL0wzDtHZmx0ntCP26Tt6lnwi7a_XflWlhb48CPE4WIk_1TcgKXVwquIhf437NvmsfXo_ctoSCJ4EGPYN4BUQCugYWmsMdh5aFzXVS3nz9DEHJVAh5IB7C3N9TTYOmplUVIoRLLfCyk16eMhO-I3zv2T3PKTiM01vUe_7zxaXqPLdG52GQ_U-wmJueMhYABWkKDJtVdWqYn7RS-dJckbEozbdalMqwyIe9ejMz8MlMthVTq6qaDMD8-n26WlIAA09VUw"
	TOKENINVALID    = "eyJraWQiOiJlOCt4TW4rOGYrZmlIXC9OZDNDZGNxOVRvU3FPKzdZYldcL1wvSUxCYVJyTElNPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiI1ZDUxNDdlMi01N2ZjLTQ4YTYtYjQ5My0xNzgzOTMxYWU5YzAiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xX2IyODdkMTQ0LTJhNjEtNDAzZC04MmNjLTkyYWY3ZmZhYmZjNiIsImV2ZW50X2lkIjoiZjFhNTYwMzItYjYyZi00OWFhLWI0MmEtMmI3YzIwZmVhY2VkIiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU5NTQ2MTQyMSwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfZnBaQ3lxbFFOIiwiZXhwIjoxNTk1NDY1MDIxLCJpYXQiOjE1OTU0NjE0MjEsImp0aSI6ImU5ODNkNDJjLWM4YWUtNDZmNi05MmFlLTE1YzczMGEwODRhZCIsImNsaWVudF9pZCI6IjQ5ZzUwOGgwanIyczJ0bmJpM2toZ3Rkc3AyIiwidXNlcm5hbWUiOiI1ZDUxNDdlMi01N2ZjLTQ4YTYtYjQ5My0xNzgzOTMxYWU5YzAifQ.T6IOtBsU0QLzOplbEWqa1QRAS7nMTknWP-meYaE4WQybLzFr-9-dPGMA3spFoQfgjD40Mxl5CFtkVmDSn1W8Yo3wvVATt7t1220YV2WAIQLf4SKghS-dwans10BpAnC5BwLgEE_sDPaJ064mV56xbOO9R-ePmNws-_qYp_R615RtHwJPtQedVwFRH7W63wS6ATs0wBaQ6McIAu1QSyoOj6ePegSzhJd_bhrDD4i42GqC2rb0rca_IRYxd3Ev44Rjx9QGZHni-BYw04jBKSFuYUvtoIpPl9gD4PzxSU5d_4lRX264uBY9F8cgeVmV97JvzVSdoAqaxfpj_drX0789Sg"
	KIDGOOD         = "cnhGjt0UsjdSG71oQd5q8SF3Soof8pO5MjM8Lh7MX9k="
	// KIDBAD          = "e8+xMn+8f+fiH/Nd3Cdcq9ToSqO+7YbW//"
	// Application
	SDCC = "sdcc"
)

func TestValidToken(t *testing.T) {
	var soteErr sError.SoteError
	if soteErr = ValidToken(SDCC, sConfigParams.DEVELOPMENT, TOKENEXPIRED); soteErr.ErrCode != 500050 && soteErr.ErrCode != nil {
		t.Errorf("ValidToken failed: Expected soteErr to be 500050 or nil: %v", soteErr.FmtErrMsg)
	}
}
func TestValidTokenMissingParams(t *testing.T) {
	var soteErr sError.SoteError
	if soteErr = ValidToken("", sConfigParams.DEVELOPMENT, TOKENEXPIRED); soteErr.ErrCode != 200514 {
		t.Errorf("ValidToken failed: Expected soteErr to be 200514 %v", soteErr.FmtErrMsg)
	}
}
func TestInValidSignatureToken(t *testing.T) {
	var soteErr sError.SoteError
	if soteErr = ValidToken(SDCC, sConfigParams.DEVELOPMENT, TOKENINVALIDSIG); soteErr.ErrCode != 500050 && soteErr.ErrCode != 500055 {
		t.Errorf("ValidToken failed: Expected soteErr to be 500050 or 500055: %v", soteErr.FmtErrMsg)
	}
}
func TestInValidToken(t *testing.T) {
	var soteErr sError.SoteError
	if soteErr = ValidToken(SDCC, sConfigParams.DEVELOPMENT, TOKENINVALID); soteErr.ErrCode != 500055 {
		t.Errorf("ValidToken failed: Expected soteErr to be 500055: %v", soteErr.FmtErrMsg)
	}
}
func TestMatchKid(t *testing.T) {
	if keys, soteErr := matchKid(sConfigParams.DEVELOPMENT, KIDGOOD); soteErr.ErrCode == nil {
		if len(keys) == 0 {
			t.Errorf("matchKid failed: Expected keys count to be greater than zero: %v", len(keys))
		}
	} else {
		t.Errorf("matchKid failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
}
func TestGetPublicKey(t *testing.T) {
	if keySet, soteErr := getPublicKey(sConfigParams.DEVELOPMENT); soteErr.ErrCode == nil {
		if len(keySet.Keys) == 0 {
			t.Errorf("matchKid failed: Expected keys count to be greater than zero: %v", len(keySet.Keys))
		}
	} else {
		t.Errorf("matchKid failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}
}
func TestFetchPublicKey(t *testing.T) {
	var (
		region, userPoolId string
		soteErr            sError.SoteError
		keySet             *jwk.Set
	)

	if region, soteErr = sConfigParams.GetRegion(); soteErr.ErrCode != nil {
		t.Fatalf("matchKid failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}

	if userPoolId, soteErr = sConfigParams.GetUserPoolId(sConfigParams.DEVELOPMENT); soteErr.ErrCode != nil {
		t.Fatalf("matchKid failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}

	if keySet, soteErr = fetchPublicKey(region, userPoolId, sConfigParams.DEVELOPMENT); soteErr.ErrCode == nil {
		if len(keySet.Keys) == 0 {
			t.Errorf("matchKid failed: Expected keys count to be greater than zero: %v", len(keySet.Keys))
		}
	} else {
		t.Errorf("matchKid failed: Expected soteErr to be nil: %v", soteErr.FmtErrMsg)
	}

	if keySet, soteErr = fetchPublicKey("SCOTT_LAND", userPoolId, sConfigParams.DEVELOPMENT); soteErr.ErrCode != 605030 {
		t.Errorf("matchKid failed: Expected soteErr to be 605030: %v", soteErr.FmtErrMsg)
	}
}
func TestValidateClaims(t *testing.T) {
	var claims jwt.MapClaims
	if soteErr := validateClaims(claims, SDCC, sConfigParams.STAGING); soteErr.ErrCode != 500070 {
		t.Errorf("validateClaims failed: Expected soteErr to be 500070")
	}
	claims = make(map[string]interface{})
	claims["scope"] = "scott.fake.scope"
	if soteErr := validateClaims(claims, SDCC, sConfigParams.STAGING); soteErr.ErrCode != 500060 {
		t.Errorf("validateClaims failed: Expected soteErr to be 500060")
	}
	claims = make(map[string]interface{})
	claims["token_use"] = "scott.fake.use"
	if soteErr := validateClaims(claims, SDCC, sConfigParams.STAGING); soteErr.ErrCode != 500060 {
		t.Errorf("validateClaims failed: Expected soteErr to be 500060")
	}
	claims = make(map[string]interface{})
	claims["iss"] = "scott.fake.iss"
	if soteErr := validateClaims(claims, SDCC, sConfigParams.STAGING); soteErr.ErrCode != 500060 {
		t.Errorf("validateClaims failed: Expected soteErr to be 500060")
	}
	claims = make(map[string]interface{})
	claims["client_id"] = "scott.fake.client_id"
	if soteErr := validateClaims(claims, SDCC, sConfigParams.STAGING); soteErr.ErrCode != 500040 {
		t.Errorf("validateClaims failed: Expected soteErr to be 500040")
	}
}
func TestValidateClientId(t *testing.T) {
	clientId, soteErr := sConfigParams.GetClientId(SDCC, sConfigParams.STAGING)
	if soteErr.ErrCode == nil {
		if soteErr = validateClientId(clientId, SDCC, sConfigParams.STAGING); soteErr.ErrCode != nil {
			t.Errorf("validateClientId failed: Expected clientId to match")
		}
	}
	clientId = "FAKE_CLIENT_ID"
	if soteErr = validateClientId(clientId, SDCC, sConfigParams.STAGING); soteErr.ErrCode == nil {
		t.Errorf("validateClientId failed: Expected soteErr to be other than nil")
	}
}
