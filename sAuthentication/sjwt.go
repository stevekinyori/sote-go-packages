package sAuthentication

import (
	"context"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	"gitlab.com/soteapps/packages/v2021/sConfigParams"
	"gitlab.com/soteapps/packages/v2021/sError"
	"gitlab.com/soteapps/packages/v2021/sLogger"
)

func ValidToken(tEnvironment, rawToken string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	if tEnvironment != "" && rawToken != "" {
		if len(strings.Split(rawToken, ".")) == 3 {
			token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					soteErr = sError.GetSError(209500, nil, sError.EmptyMap)
				}

				if soteErr.ErrCode == nil {
					var (
						kid string
						ok  bool
						key jwk.Key
					)
					if kid, ok = token.Header["kid"].(string); !ok {
						soteErr = sError.GetSError(209510, nil, sError.EmptyMap)
					}

					if soteErr.ErrCode == nil {
						if key, soteErr = matchKid(tEnvironment, kid); soteErr.ErrCode == nil {
							var raw interface{}
							return raw, key.Raw(&raw)
						}
					}
				}

				return nil, nil
			})

			if err != nil {
				if strings.Contains(err.Error(), "expired") {
					soteErr = sError.GetSError(208350, nil, sError.EmptyMap)
				}
				if strings.Contains(err.Error(), "invalid type") || strings.Contains(err.Error(), "invalid character") {
					soteErr = sError.GetSError(208355, nil, sError.EmptyMap)
				}
				if strings.Contains(err.Error(), "invalid number of segments") {
					soteErr = sError.GetSError(208356, nil, sError.EmptyMap)
				}
			}

			if soteErr.ErrCode == nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					soteErr = validateClaims(claims, tEnvironment)
				} else {
					soteErr = sError.GetSError(208355, nil, sError.EmptyMap)
				}
			}
		} else {
			soteErr = sError.GetSError(208356, nil, sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		}
	} else {
		soteErr = sError.GetSError(200512, sError.BuildParams([]string{"Environment", "Token"}), sError.EmptyMap)
		sLogger.Info(soteErr.FmtErrMsg)
	}

	return
}

func matchKid(tEnvironment, kid string) (key jwk.Key, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		keySet jwk.Set
		ok     bool
	)
	keySet, soteErr = getPublicKey(tEnvironment)
	key, ok = keySet.LookupKeyID(kid)
	if !ok {
		soteErr = sError.GetSError(209521, sError.BuildParams([]string{kid}), sError.EmptyMap)
	}

	return
}

/*
This will return the public key needed to validate the jwt token

NOTE: If the region or user pool id is not found, getPublicKey will default to the 'eu-west-1' region
and the userPoolId used by development instance of Cognito
*/
func getPublicKey(tEnvironment string) (keySet jwk.Set, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if region, soteErr := sConfigParams.GetRegion(); soteErr.ErrCode == nil {
		if userPoolId, soteErr := sConfigParams.GetUserPoolId(tEnvironment); soteErr.ErrCode == nil {
			if keySet, soteErr = fetchPublicKey(region, userPoolId, tEnvironment); soteErr.ErrCode != nil {
				soteErr = sError.GetSError(210030, sError.BuildParams([]string{tEnvironment}), sError.EmptyMap)
				panic(soteErr.FmtErrMsg)
			}
		}
	}

	return
}

/*
This will pull the public key from the internet. The URL should not be output to logs.
*/
func fetchPublicKey(region, userPoolId, tEnvironment string) (keySet jwk.Set, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err            error
		keyPubLocation = "https://cognito-idp." + region + ".amazonaws.com/" + userPoolId + "/.well-known/jwks.json"
	)

	keySet, err = jwk.Fetch(context.Background(), keyPubLocation)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "remote JWK") {
			soteErr = sError.GetSError(210030, sError.BuildParams([]string{region, tEnvironment}), sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
This checks if the claim in the token are valid
*/
func validateClaims(claims jwt.MapClaims, tEnvironment string) (soteErr sError.SoteError) {
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
					soteErr = sError.GetSError(208360, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
					sLogger.Info(soteErr.FmtErrMsg)
				} else {
					sLogger.Info("Claim (scope) was found")
				}
			case "token_use":
				claimCount++
				if claim != "access" {
					soteErr = sError.GetSError(208360, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
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
							soteErr = sError.GetSError(208360, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
							sLogger.Info(soteErr.FmtErrMsg)
						} else {
							sLogger.Info("Claim (iss) was found")
						}
					} else {
						soteErr = sError.GetSError(208360, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
						sLogger.Info(soteErr.FmtErrMsg)
					}
				} else {
					soteErr = sError.GetSError(208360, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
					sLogger.Info(soteErr.FmtErrMsg)
				}
				/*case "client_id":
				claimCount++
				if soteErr = validateClientId(claim.(string), tApplication, tEnvironment); soteErr.ErrCode == nil {
					sLogger.Info("Claim (client_id) was found")
				}*/
			}
		} else {
			break
		}
	}

	if claimCount != 3 && soteErr.ErrCode == nil {
		soteErr = sError.GetSError(208370, nil, sError.EmptyMap)
		sLogger.Info(soteErr.FmtErrMsg)
	}

	return
}

func validateClientId(tClientId, tApplication, tEnvironment string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var clientId string

	if clientId, soteErr = sConfigParams.GetClientId(tApplication, tEnvironment); soteErr.ErrCode == nil {
		if tClientId != clientId {
			soteErr = sError.GetSError(208340, nil, sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		}
	} else if soteErr.ErrCode != nil {
		sLogger.Info(soteErr.FmtErrMsg)
	}

	return
}