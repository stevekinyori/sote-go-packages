package sAuthentication

import (
	"context"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	"gitlab.com/soteapps/packages/v2023/sConfigParams"
	"gitlab.com/soteapps/packages/v2023/sError"
	"gitlab.com/soteapps/packages/v2023/sLogger"
)

type Config struct {
	AppEnvironment string
	AwsRegion      string
	UserPoolId     string
	ClientId       string
}

func ValidToken(ctx context.Context, rawToken string, config *Config) (soteErr sError.SoteError) {
	sLogger.DebugMethod()
	if config != nil && config.AppEnvironment != "" && rawToken != "" {
		if len(strings.Split(rawToken, ".")) == 3 {
			token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					soteErr = sError.GetSError(sError.ErrUnexpectedSigningMethod, nil, sError.EmptyMap)
				}

				if soteErr.ErrCode == nil {
					var (
						kid string
						ok  bool
						key jwk.Key
					)
					if kid, ok = token.Header["kid"].(string); !ok {
						soteErr = sError.GetSError(sError.ErrMissingKidHeader, nil, sError.EmptyMap)
					}

					if soteErr.ErrCode == nil {
						if key, soteErr = matchKid(ctx, kid, config); soteErr.ErrCode == nil {
							var raw interface{}
							return raw, key.Raw(&raw)
						}
					}
				}

				return nil, nil
			})

			if err != nil {
				if strings.Contains(err.Error(), "expired") {
					soteErr = sError.GetSError(sError.ErrExpiredToken, nil, sError.EmptyMap)
				}
				if strings.Contains(err.Error(), "invalid type") || strings.Contains(err.Error(), "invalid character") {
					soteErr = sError.GetSError(sError.ErrInvalidToken, nil, sError.EmptyMap)
				}
				if strings.Contains(err.Error(), "invalid number of segments") {
					soteErr = sError.GetSError(sError.ErrMissingTokenSegments, nil, sError.EmptyMap)
				}
			}

			if soteErr.ErrCode == nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					soteErr = validateClaims(ctx, claims, config.AppEnvironment)
				} else {
					soteErr = sError.GetSError(sError.ErrInvalidToken, nil, sError.EmptyMap)
				}
			}
		} else {
			soteErr = sError.GetSError(sError.ErrMissingTokenSegments, nil, sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		}
	} else {
		soteErr = sError.GetSError(sError.ErrExpectedTwoParameters, sError.BuildParams([]string{"Environment", "Token"}), sError.EmptyMap)
		sLogger.Info(soteErr.FmtErrMsg)
	}

	return
}

func matchKid(ctx context.Context, kid string, config *Config) (key jwk.Key, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		keySet jwk.Set
		ok     bool
	)

	if keySet, soteErr = getPublicKey(ctx, config); soteErr.ErrCode != nil {
		return
	}

	key, ok = keySet.LookupKeyID(kid)
	if !ok {
		soteErr = sError.GetSError(sError.ErrMissingKidInPublicKey, sError.BuildParams([]string{kid}), sError.EmptyMap)
	}

	return
}

/*
This will return the public key needed to validate the jwt token

NOTE: If the region or user pool id is not found, getPublicKey will default to the 'eu-west-1' region
and the userPoolId used by development instance of Cognito
*/
func getPublicKey(ctx context.Context, config *Config) (keySet jwk.Set, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if config == nil {
		soteErr = sError.GetSError(sError.ErrExpectedTwoParameters, sError.BuildParams([]string{"Aws Region", "UserPoolId"}), sError.EmptyMap)
		sLogger.Info(soteErr.FmtErrMsg)
		return
	}

	if keySet, soteErr = fetchPublicKey(config.AwsRegion, config.UserPoolId, config.AppEnvironment); soteErr.ErrCode != nil {
		soteErr = sError.GetSError(sError.ErrFetchingJWKError, sError.BuildParams([]string{config.AppEnvironment}), sError.EmptyMap)
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
			soteErr = sError.GetSError(sError.ErrFetchingJWKError, sError.BuildParams([]string{region, tEnvironment}), sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		}
	}

	return
}

/*
This checks if the claim in the token are valid
*/
func validateClaims(ctx context.Context, claims jwt.MapClaims, tEnvironment string) (soteErr sError.SoteError) {
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
					soteErr = sError.GetSError(sError.ErrInvalidClaims, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
					sLogger.Info(soteErr.FmtErrMsg)
				} else {
					sLogger.Info("Claim (scope) was found")
				}
			case "token_use":
				claimCount++
				if claim != "access" {
					soteErr = sError.GetSError(sError.ErrInvalidClaims, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
					sLogger.Info(soteErr.FmtErrMsg)
				} else {
					sLogger.Info("Claim (token_use) was found")
				}
			case "iss":
				claimCount++
				if region, soteErr = sConfigParams.GetRegion(ctx); soteErr.ErrCode == nil {
					if userPoolId, soteErr = sConfigParams.GetUserPoolId(ctx, tEnvironment); soteErr.ErrCode == nil {
						issuerURL := "https://cognito-idp." + region + ".amazonaws.com/" + userPoolId
						if claim != issuerURL {
							soteErr = sError.GetSError(sError.ErrInvalidClaims, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
							sLogger.Info(soteErr.FmtErrMsg)
						} else {
							sLogger.Info("Claim (iss) was found")
						}
					} else {
						soteErr = sError.GetSError(sError.ErrInvalidClaims, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
						sLogger.Info(soteErr.FmtErrMsg)
					}
				} else {
					soteErr = sError.GetSError(sError.ErrInvalidClaims, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
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
		soteErr = sError.GetSError(sError.ErrMissingClaims, nil, sError.EmptyMap)
		sLogger.Info(soteErr.FmtErrMsg)
	}

	return
}

func validateClientId(ctx context.Context, tClientId, tApplication, tEnvironment string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var clientId string

	if clientId, soteErr = sConfigParams.GetClientId(ctx, tApplication, tEnvironment); soteErr.ErrCode == nil {
		if tClientId != clientId {
			soteErr = sError.GetSError(sError.ErrInvalidClientId, nil, sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		}
	} else if soteErr.ErrCode != nil {
		sLogger.Info(soteErr.FmtErrMsg)
	}

	return
}
