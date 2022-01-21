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

func ValidToken(tEnvironment, rawToken string) error {
	sLogger.DebugMethod()

	if tEnvironment == "" || rawToken == "" {
		err := sError.ErrParamMustBeSet("", nil)
		sLogger.Info(err.Error())

		return err
	}

	if len(strings.Split(rawToken, ".")) != 3 {
		err := sError.ErrSegmentsCountInvalid("", nil)
		sLogger.Info(err.Error())

		return err
	}

	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			err := sError.ErrUnexpectedSign("", nil)
			return nil, err
		}

		var (
			kid string
			ok  bool
		)
		if kid, ok = token.Header["kid"].(string); !ok {
			err := sError.ErrKIDNotFound("", nil)
			return nil, err
		}

		key, err := matchKid(tEnvironment, kid)
		if err != nil {
			return nil, err
		}

		var raw interface{}
		return raw, key.Raw(&raw)
	})

	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			err := sError.ErrTokenExpired("", err)
			return err
		}
		if strings.Contains(err.Error(), "invalid type") || strings.Contains(err.Error(), "invalid character") {
			err := sError.ErrInvalidToken("", err)
			return err
		}
		if strings.Contains(err.Error(), "invalid number of segments") {
			err := sError.ErrSegmentsCountInvalid("", err)
			return err
		}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		err := validateClaims(claims, tEnvironment)
		if err != nil {
			return err
		}
	} else {
		err := sError.ErrSegmentsCountInvalid("", nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func matchKid(tEnvironment, kid string) (key jwk.Key, soteErr error) {
	sLogger.DebugMethod()

	var (
		keySet jwk.Set
		ok     bool
	)
	keySet, soteErr = getPublicKey(tEnvironment)
	key, ok = keySet.LookupKeyID(kid)
	if !ok {
		soteErr = sError.ErrKIDMissingFromToken(kid, "", nil)
	}

	return
}

/*
This will return the public key needed to validate the jwt token

NOTE: If the region or user pool id is not found, getPublicKey will default to the 'eu-west-1' region
and the userPoolId used by development instance of Cognito
*/
func getPublicKey(tEnvironment string) (jwk.Set, error) {
	sLogger.DebugMethod()

	region, soteErr := sConfigParams.GetRegion()
	if soteErr.ErrCode != nil {
		return nil, soteErr
	}

	userPoolId, soteErr := sConfigParams.GetUserPoolId(tEnvironment)
	if soteErr.ErrCode != nil {
		return nil, soteErr
	}

	keySet, err := fetchPublicKey(region, userPoolId, tEnvironment)
	if err != nil {
		err = sError.ErrInvalidRegion(region, tEnvironment, "", err)
		panic(err.Error())
	}

	return keySet, nil
}

/*
This will pull the public key from the internet. The URL should not be output to logs.
*/
func fetchPublicKey(region, userPoolId, tEnvironment string) (keySet jwk.Set, soteErr error) {
	sLogger.DebugMethod()

	var (
		err            error
		keyPubLocation = "https://cognito-idp." + region + ".amazonaws.com/" + userPoolId + "/.well-known/jwks.json"
	)

	keySet, err = jwk.Fetch(context.Background(), keyPubLocation)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "remote JWK") {
			soteErr = sError.ErrInvalidRegion(region, tEnvironment, "", nil)
			sLogger.Info(soteErr.Error())
		}
	}

	return
}

/*
This checks if the claim in the token are valid
*/
func validateClaims(claims jwt.MapClaims, tEnvironment string) error {
	sLogger.DebugMethod()

	var (
		claimCount = 0
	)

	for key, claim := range claims {
		switch key {
		case "scope":

			claimCount++
			if claim != "aws.cognito.signin.user.admin" {
				err := sError.ErrMissingClaim("", nil)
				sLogger.Info(err.Error())
				return err

			} else {
				sLogger.Info("Claim (scope) was found")
			}

		case "token_use":

			claimCount++
			if claim != "access" {
				err := sError.ErrMissingClaim("", nil)
				sLogger.Info(err.Error())
				return err

			} else {
				sLogger.Info("Claim (token_use) was found")
			}

		case "iss":

			claimCount++
			region, soteErr := sConfigParams.GetRegion()
			if soteErr.ErrCode != nil {
				soteErr = sError.GetSError(208360, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
				sLogger.Info(soteErr.Error())
				return soteErr
			}

			userPoolId, soteErr := sConfigParams.GetUserPoolId(tEnvironment)
			if soteErr.ErrCode != nil {
				soteErr = sError.GetSError(208360, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
				sLogger.Info(soteErr.Error())
				return soteErr
			}

			issuerURL := "https://cognito-idp." + region + ".amazonaws.com/" + userPoolId
			if claim != issuerURL {
				err := sError.GetSError(208360, sError.BuildParams([]string{claim.(string)}), sError.EmptyMap)
				sLogger.Info(err.Error())
				return err
			} else {
				sLogger.Info("Claim (iss) was found")
			}

		}
	}

	if claimCount != 3 {
		err := sError.ErrMissingClaim("", nil)
		sLogger.Info(err.Error())

		return err
	}

	return nil
}
