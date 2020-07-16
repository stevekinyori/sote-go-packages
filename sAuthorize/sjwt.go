package sAuthorize

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	"gitlab.com/soteapps/packages/v2020/sConfigParams"
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

var holdSoteErr sError.SoteError

func ValidToken(tEnvironment, rawToken string) (soteErr sError.SoteError) {
	sLogger.DebugMethod()

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
	}

	if holdSoteErr.ErrCode == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims["scope"], claims["token_use"], claims["iss"], claims["client_id"])
		} else {
			fmt.Println(err)
		}
	}
	soteErr = holdSoteErr

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
			if keySet, soteErr = fetchPublicKey(region, userPoolId); soteErr.ErrCode != nil {
				soteErr = sError.GetSError(605030, sError.BuildParams([]string{tEnvironment}), sError.EmptyMap)
				panic(soteErr.FmtErrMsg)
			}
		}
	} else {
		// The following is only needed until the implementation of the region and user pool id System Manager parameters.
		// If this has been completed, remove this code!
		soteErr = sError.GetSError(605030, sError.BuildParams([]string{tEnvironment + ". RUNNING FALLBACK CODE"}), sError.EmptyMap)
		if keySet, soteErr = fetchPublicKey("eu-west-1", "eu-west-1_QVPwwCg2c"); soteErr.ErrCode != nil {
			soteErr = sError.GetSError(605030, sError.BuildParams([]string{"Fallback for getPublicKey failed!"}), sError.EmptyMap)
			panic("Fallback for getPublicKey failed!" + soteErr.FmtErrMsg)
		}
		// End of fallback code.
	}

	return
}

/*
This will pull the public key from the internet. The URL should not be output to logs.
*/
func fetchPublicKey(region, userPoolId string) (keySet *jwk.Set, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	var (
		err            error
		keyPubLocation = "https://cognito-idp." + region + ".amazonaws.com/" + userPoolId + "/.well-known/jwks.json"
	)

	keySet, err = jwk.Fetch(keyPubLocation)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			soteErr = sError.GetSError(605030, nil, sError.EmptyMap)
			sLogger.Info(soteErr.FmtErrMsg)
		}
	}

	return
}
