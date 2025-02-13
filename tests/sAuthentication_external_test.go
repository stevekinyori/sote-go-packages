package tests

import (
	"context"
	"testing"

	"gitlab.com/soteapps/packages/v2023/sAuthentication"
	"gitlab.com/soteapps/packages/v2023/sConfigParams"
	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
)

const (
	/*
		Before completing testing, make sure to put a non-expired token value in the TOKENEXPIRED const, so you have tested a successful case.
	*/
	// Expired Token
	TOKENEXPIRED = "eyJraWQiOiJjbmhHanQwVXNqZFNHNzFvUWQ1cThTRjNTb29mOHBPNU1qTThMaDdNWDlrPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJkZDg5NmVhNi03NmI2LTQ1OGYtYWYyNC0zMTAyN2JiOGQzODMiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xXzc3MTViMjlmLTVmZmMtNGE4MC04MjkxLTI0YzhjYzAwNjdkMSIsImV2ZW50X2lkIjoiZDFlODg2NDktMTI0Ni00YTIxLThjZjYtYjZmMGQ4ZmMwNDVmIiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU5ODAyNDU4NCwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfUVZQd3dDZzJjIiwiZXhwIjoxNTk4MDI4MTg0LCJpYXQiOjE1OTgwMjQ1ODQsImp0aSI6ImIwYjgzNGE3LWY4MjQtNGMwMi1hZGRiLTU5ZmUwODY2YzJkOSIsImNsaWVudF9pZCI6InR0c21yNjUzcXRma3VjZ2dvZnBtcWJic3AiLCJ1c2VybmFtZSI6ImRkODk2ZWE2LTc2YjYtNDU4Zi1hZjI0LTMxMDI3YmI4ZDM4MyJ9.R0h6otGxbqX3Novw5qHEpOvkSvzSPVHBxtXd-eb6Zq0Oq8IXMkHrgfmBZH0CiaGUzjaCHF5OyZEfqq9vfbJ0iu7YF9hGbBVHDVbJ1mywbpwQ-5l53rOEFeaXU-jshVXnNb_VbBpPrxC0Cha_w-_MJ1JF5c5_jw-93iT9azU57mwcwfc95Ro9P4mWsQ7i6Wurk1Mw7ijhTaguJceB8cRcfoYCt2xx5BuGeYBLe-5QuTbUebpkMM6WoGschwpUiZDsPhhsHXO9Tu_Exk58Ad3BSaHVOgrquF1qq6KNoCObiUIPjq2z4BELev0jH1B0KrY_0kK77IkdAokte1kmniAHuQ"
	// Invalid Signature error
	TOKENINVALIDSIG = "eyJraWQiOiJjbmhHanQwVXNqZFNHNzFvUWQ1cThTRjNTb29mOHBPNU1qTThMaDdNWDlrPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJkZDg5NmVhNi03NmI2LTQ1OGYtYWYyNC0zMTAyN2JiOGQzODMiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xXzU1MDg3NThmLWM4NmMtNDYwYS05NjMyLTBiMzBjYjE0NzBhNSIsImV2ZW50X2lkIjoiYWFkNTdhMjMtZDRjMC00M2E2LWIyZGQtNzIzZjA3ZjkxNzk5IiwiDG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU5NTQ2MTE4NSwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfUVZQd3dDZzJjIiwiZXhwIjoxNTk1NDY0Nzg1LCJpYXQiOjE1OTU0NjExODUsImp0aSI6IjFiNTYyMWNlLTk5YzAtNDczOC05YjA0LTZkOTk0ODUxYTM2MiIsImNsaWVudF9pZCI6InR0c21yNjUzcXRma3VjZ2dvZnBtcWJic3AiLCJ1c2VybmFtZSI6ImRkODk2ZWE2LTc2YjYtNDU4Zi1hZjI0LTMxMDI3YmI4ZDM4MyJ9.Y2ypGn1cVPvjoOMxE5jbm7HwAWghN4zX2RJ8UxGwfwYLsHEPpGgJHpCqHJal1jufe-ciM_XQM7QvPFYFO5BL0wzDtHZmx0ntCP26Tt6lnwi7a_XflWlhb48CPE4WIk_1TcgKXVwquIhf437NvmsfXo_ctoSCJ4EGPYN4BUQCugYWmsMdh5aFzXVS3nz9DEHJVAh5IB7C3N9TTYOmplUVIoRLLfCyk16eMhO-I3zv2T3PKTiM01vUe_7zxaXqPLdG52GQ_U-wmJueMhYABWkKDJtVdWqYn7RS-dJckbEozbdalMqwyIe9ejMz8MlMthVTq6qaDMD8-n26WlIAA09VUw"
	// Invalid type error
	TOKENINVALID = "eyJraWQiOiJlOCt4TW4rOGYrZmlIXC9OZDNDZGNxOVRvU3FPKzdZYldcL1wvSUxCYVJyTElNPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiI1ZDUxNDdlMi01N2ZjLTQ4YTYtYjQ5My0xNzgzOTMxYWU5YzAiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xX2IyODdkMTQ0LTJhNjEtNDAzZC04MmNjLTkyYWY3ZmZhYmZjNiIsImV2ZW50X2lkIjoiZjFhNTYwMzItYjYyZi00OWFhLWI0MmEtMmI3YzIwZmVhY2VkIiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU5NTQ2MTQyMSwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfZnBaQ3lxbFFOIiwiZXhwIjoxNTk1NDY1MDIxLCJpYXQiOjE1OTU0NjE0MjEsImp0aSI6ImU5ODNkNDJjLWM4YWUtNDZmNi05MmFlLTE1YzczMGEwODRhZCIsImNsaWVudF9pZCI6IjQ5ZzUwOGgwanIyczJ0bmJpM2toZ3Rkc3AyIiwidXNlcm5hbWUiOiI1ZDUxNDdlMi01N2ZjLTQ4YTYtYjQ5My0xNzgzOTMxYWU5YzAifQ.T6IOtBsU0QLzOplbEWqa1QRAS7nMTknWP-meYaE4WQybLzFr-9-dPGMA3spFoQfgjD40Mxl5CFtkVmDSn1W8Yo3wvVATt7t1220YV2WAIQLf4SKghS-dwans10BpAnC5BwLgEE_sDPaJ064mV56xbOO9R-ePmNws-_qYp_R615RtHwJPtQedVwFRH7W63wS6ATs0wBaQ6McIAu1QSyoOj6ePegSzhJd_bhrDD4i42GqC2rb0rca_IRYxd3Ev44Rjx9QGZHni-BYw04jBKSFuYUvtoIpPl9gD4PzxSU5d_4lRX264uBY9F8cgeVmV97JvzVSdoAqaxfpj_drX0789Sg"
	// Only two segments
	TOKENMISSINGSEGMENT = "eyJraWQiOiJjbmhHanQwVXNqZFNHNzFvUWQ1cThTRjNTb29mOHBPNU1qTThMaDdNWDlrPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiJkZDg5NmVhNi03NmI2LTQ1OGYtYWYyNC0zMTAyN2JiOGQzODMiLCJkZXZpY2Vfa2V5IjoiZXUtd2VzdC0xX2Q3MmQ5NzJiLTRiNGQtNGJjMi05YTU0LWNmZmJlOTU1YTExMiIsImV2ZW50X2lkIjoiNWRmMGUwMWMtMTAyMy00OGRmLThjOTgtM2Y1MWI4N2Y1Y2E0IiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU5Nzg1ODE5NCwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMS5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTFfUVZQd3dDZzJjIiwiZXhwIjoxNTk3ODYxNzk0LCJpYXQiOjE1OTc4NTgxOTQsImp0aSI6IjI3NTkxMWYxLWJhYmYtNDNmNS04YjE4LWMxZWFkMTljMmY1YiIsImNsaWVudF9pZCI6InR0c21yNjUzcXRma3VjZ2dvZnBtcWJic3AiLCJ1c2VybmFtZSI6ImRkODk2ZWE2LTc2YjYtNDU4Zi1hZjI0LTMxMDI3YmI4ZDM4MyJ9ECr1bgVGRNYscxP10HCvu3WE6RQKtwOvs3yoGUGsDr9TDi_hWHc7VUxcDXvjxYy5C9zm75Jg4CIjbp42GhA2L82h0cAyI3eBiWq_GSIMdKo1ZNWVrIrksYDLgHJUaRqAuO4ui2AEsN5P2fXlLPmOngMOEQbVo22ybUeZpM83Q0R08A92z3cAxW2MPjCvYzdf1MPO3OCS2OpK54SzRQftR_mRqWhSMYM6NHn_hlE33quhCQ1AMFrQsoZOT4ujwA_h0WHWhxikPnJHJIuqDatWQ4PPwF4DcbGFSZfOm3fDOY9GBXXypjywUIn7cRCTLEm3YN1-HzYtBoLq9VZF9uxLKA"
	KIDGOOD             = "cnhGjt0UsjdSG71oQd5q8SF3Soof8pO5MjM8Lh7MX9k="
	FAKETOKEN           = "TOKEN"
	// KIDBAD          = "e8+xMn+8f+fiH/Nd3Cdcq9ToSqO+7YbW//"
	// Application
	SDCC = "sdcc"
)

var (
	sAuthCtx      = context.Background()
	awsRegion, _  = sConfigParams.GetRegion(sAuthCtx)
	userPoolId, _ = sConfigParams.GetUserPoolId(sAuthCtx, sConfigParams.DEVELOPMENT)
)

func init() {
	sLogger.SetLogMessagePrefix("packages")
}

func TestValidToken(tPtr *testing.T) {
	var soteErr sError.SoteError
	if soteErr = sAuthentication.ValidToken(sAuthCtx, TOKENEXPIRED, &sAuthentication.Config{
		AppEnvironment: sConfigParams.DEVELOPMENT,
		AwsRegion:      awsRegion,
		UserPoolId:     userPoolId,
		ClientId:       "",
	}); soteErr.ErrCode != sError.ErrExpiredToken && soteErr.ErrCode != nil {
		tPtr.Errorf("ValidToken failed: Expected soteErr to be %v or nil: %v", sError.ErrExpiredToken, soteErr.ErrCode)
	}
}
func TestValidFakeToken(tPtr *testing.T) {
	var soteErr sError.SoteError
	if soteErr = sAuthentication.ValidToken(sAuthCtx, TOKENMISSINGSEGMENT, &sAuthentication.Config{
		AppEnvironment: sConfigParams.DEVELOPMENT,
		AwsRegion:      awsRegion,
		UserPoolId:     userPoolId,
		ClientId:       "",
	}); soteErr.ErrCode != sError.ErrMissingTokenSegments && soteErr.ErrCode != nil {
		tPtr.Errorf("ValidToken failed: Expected soteErr to be nil or nil: %v", soteErr.FmtErrMsg)
	}
}
func TestValidMissingSegmentToken(tPtr *testing.T) {
	var soteErr sError.SoteError
	if soteErr = sAuthentication.ValidToken(sAuthCtx, FAKETOKEN, &sAuthentication.Config{
		AppEnvironment: sConfigParams.DEVELOPMENT,
		AwsRegion:      awsRegion,
		UserPoolId:     userPoolId,
		ClientId:       "",
	}); soteErr.ErrCode != sError.ErrMissingTokenSegments && soteErr.ErrCode != nil {
		tPtr.Errorf("ValidToken failed: Expected soteErr to be %v or nil: %v", sError.ErrInvalidToken, soteErr.FmtErrMsg)
	}
}
func TestInValidSignatureToken(tPtr *testing.T) {
	var soteErr sError.SoteError
	if soteErr = sAuthentication.ValidToken(sAuthCtx, TOKENINVALIDSIG, &sAuthentication.Config{
		AppEnvironment: sConfigParams.DEVELOPMENT,
		AwsRegion:      awsRegion,
		UserPoolId:     userPoolId,
		ClientId:       "",
	}); soteErr.ErrCode != sError.ErrExpiredToken && soteErr.ErrCode != sError.ErrInvalidToken && soteErr.ErrCode != sError.ErrMissingTokenSegments {
		tPtr.Errorf("ValidToken failed: Expected soteErr to be %v, %v or %v: %v", sError.ErrExpiredToken,
			sError.ErrMissingTokenSegments, sError.ErrInvalidToken, soteErr.ErrCode)
	}
}
func TestInValidToken(tPtr *testing.T) {
	var soteErr sError.SoteError
	if soteErr = sAuthentication.ValidToken(sAuthCtx, TOKENINVALID, &sAuthentication.Config{
		AppEnvironment: sConfigParams.DEVELOPMENT,
		AwsRegion:      awsRegion,
		UserPoolId:     userPoolId,
		ClientId:       "",
	}); soteErr.ErrCode != sError.ErrInvalidToken {
		tPtr.Errorf("ValidToken failed: Expected soteErr to be %v: %v", sError.ErrInvalidToken, soteErr.ErrCode)
	}
}
