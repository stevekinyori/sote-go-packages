package sAuthorize

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

var holdSoteErr sError.SoteError

func ValidToken(tApplication, tEnvironment, rawToken string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	if tApplication != "" && tEnvironment != "" && rawToken != "" {
		if len(strings.Split(rawToken, ".")) == 3 {
			token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					holdSoteErr = sError.GetSError(605000, nil, sError.EmptyMap)
				}

				if holdSoteErr.ErrCode == nil {
					var (
						kid  string
						ok   bool
						keys []jwk.Key
					)
					if kid, ok = token.Header["kid"].(string); !ok {
						holdSoteErr = sError.GetSError(605010, nil, sError.EmptyMap)
					}

					if holdSoteErr.ErrCode == nil {
						if keys, holdSoteErr = matchKid(tEnvironment, kid); holdSoteErr.ErrCode == nil {
							var raw interface{}
							return raw, keys[0].Raw(&raw)
						}
					}
				}

				return nil, nil
			})

			if err != nil {
				if strings.Contains(err.Error(), "expired") {
					holdSoteErr = sError.GetSError(500050, nil, sError.EmptyMap)
				}
				if strings.Contains(err.Error(), "invalid type") || strings.Contains(err.Error(), "invalid character") {
					holdSoteErr = sError.GetSError(500055, nil, sError.EmptyMap)
				}
				if strings.Contains(err.Error(), "invalid number of segments") {
					holdSoteErr = sError.GetSError(500056, nil, sError.EmptyMap)
				}
			}

			if holdSoteErr.ErrCode == nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					holdSoteErr = validateClaims(claims, tApplication, tEnvironment)
				} else {
					holdSoteErr = sError.GetSError(500055, nil, sError.EmptyMap)
				}
			}
			soteErr = holdSoteErr
		} else {
			soteErr = sError.GetSError(500056, nil, sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		}
	} else {
		soteErr = sError.GetSError(200514, sError.BuildParams([]string{"tApplication", "tEnvironment", "rawToken"}), sError.EmptyMap)
		sLogger.Info(soteErr.FmtErrMsg)
	}

	return
}

func matchKid(tEnvironment, kid string) (keys []jwk.Key, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		keySet *jwk.Set
	)
	keySet, soteErr = getPublicKey(tEnvironment)
	keys = keySet.LookupKeyID(kid)
	if len(keys) == 0 {
		soteErr = sError.GetSError(605021, sError.BuildParams([]string{kid}), sError.EmptyMap)
	}

	return
}

/*
This will return the public key needed to validate the jwt token

NOTE: If the region or user pool id is not found, getPublicKey will default to the 'eu-west-1' region
and the userPoolId used by development instance of Cognito
*/
func getPublicKey(tEnvironment string) (keySet *jwk.Set, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if region, soteErr := sConfigParams.GetRegion(); soteErr.ErrCode == nil {
		if userPoolId, soteErr := sConfigParams.GetUserPoolId(tEnvironment); soteErr.ErrCode == nil {
			if keySet, soteErr = fetchPublicKey(region, userPoolId, tEnvironment); soteErr.ErrCode != nil {
				soteErr = sError.GetSError(605030, sError.BuildParams([]string{tEnvironment}), sError.EmptyMap)
				panic(soteErr.FmtErrMsg)
			}
		}
	}

	return
}

/*
This will pull the public key from the internet. The URL should not be output to logs.
*/
func fetchPublicKey(region, userPoolId, tEnvironment string) (keySet *jwk.Set, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err            error
		keyPubLocation = "https://cognito-idp." + region + ".amazonaws.com/" + userPoolId + "/.well-known/jwks.json"
	)

	keySet, err = jwk.Fetch(keyPubLocation)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "remote JWK") {
			soteErr = sError.GetSError(605030, sError.BuildParams([]string{region, tEnvironment}), sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
This checks if the claim in the token are valid
*/
func validateClaims(claims jwt.MapClaims, tApplication, tEnvironment string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		claimCount         = 0
		region, userPoolId string
	)

	for key, claim := range claims {
		if soteErr.ErrCode == nil {
			switch key {
			case "scope":
				claimCount++
				if claim != "aws.cognito.signin.user.admin" {
					soteErr = sError.GetSError(500060, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
					sLogger.Info(soteErr.FmtErrMsg)
				} else {
					sLogger.Info("Claim (scope) was found")
				}
			case "token_use":
				claimCount++
				if claim != "access" {
					soteErr = sError.GetSError(500060, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
					sLogger.Info(soteErr.FmtErrMsg)
				} else {
					sLogger.Info("Claim (token_use) was found")
				}
			case "iss":
				claimCount++
				if region, soteErr = sConfigParams.GetRegion(); soteErr.ErrCode == nil {
					if userPoolId, soteErr = sConfigParams.GetUserPoolId(tEnvironment); soteErr.ErrCode == nil {
						issuerURL := "https://cognito-idp." + region + ".amazonaws.com/" + userPoolId
						if claim != issuerURL {
							soteErr = sError.GetSError(500060, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
							sLogger.Info(soteErr.FmtErrMsg)
						} else {
							sLogger.Info("Claim (iss) was found")
						}
					} else {
						soteErr = sError.GetSError(500060, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
						sLogger.Info(soteErr.FmtErrMsg)
					}
				} else {
					soteErr = sError.GetSError(500060, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
					sLogger.Info(soteErr.FmtErrMsg)
				}
			case "client_id":
				claimCount++
				if soteErr = validateClientId(claim.(string), tApplication, tEnvironment); soteErr.ErrCode == nil {
					sLogger.Info("Claim (client_id) was found")
				}
			}
		} else {
			break
		}
	}

	if claimCount != 4 && soteErr.ErrCode == nil {
		soteErr = sError.GetSError(500070, nil, sError.EmptyMap)
		sLogger.Info(soteErr.FmtErrMsg)
	}

	return
}

func validateClientId(tClientId, tApplication, tEnvironment string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var clientId string

	if clientId, soteErr = sConfigParams.GetClientId(tApplication, tEnvironment); soteErr.ErrCode == nil {
		if tClientId != clientId {
			soteErr = sError.GetSError(500040, nil, sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		}
	} else if soteErr.ErrCode != nil {
		sLogger.Info(soteErr.FmtErrMsg)
	}

	return
}
