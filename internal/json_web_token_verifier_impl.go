package internal

import (
	"github.com/eaglesakura/gaefire/auth"
	"golang.org/x/net/context"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"google.golang.org/appengine/log"
)

type JsonWebTokenVerifierImpl struct {
	service            gaefire.FirebaseServiceAccount
	ctx                context.Context
	token              string
	skipExpireCheck    bool
	skipProjectIdCheck bool
}

type VerifyError struct {
	internalError *jwt.ValidationError
	rawError      error
}

func newTokenError(raw error) (it *VerifyError) {
	internal, _ := raw.(*jwt.ValidationError)

	return &VerifyError{
		internalError: internal,
		rawError:raw,
	}
}

func (it *VerifyError)Error() string {
	if it.internalError != nil {
		return it.internalError.Error()
	} else {
		return it.rawError.Error()
	}
}


//
// Verify user JWT
// check public key, by Header["kid"]. firebase default CustomToken
// ex. Android => user.getToken(true)
//
func (it *JsonWebTokenVerifierImpl)Verify(jwtToken string) (gaefire.VerifiedJsonWebToken, error) {
	// parse & verify
	rawToken, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		kid := token.Header["kid"].(string)
		if kid == "" {
			return nil, errors.New("NotFound kid")
		}

		publicKey, err := it.service.FindPublicKey(it.ctx, kid)
		if err != nil {
			return nil, err
		} else {
			return publicKey, nil
		}
	})

	if strings.Contains(err.Error(), "crypto/rsa") {
		// Verify error
		return nil, newTokenError(err)
	}

	if rawToken != nil && rawToken.Claims != nil {
		log.Errorf(it.ctx, "Refresh Validate old[%v]", err.Error())
		err = rawToken.Claims.Valid()
	}

	// update error
	if err != nil {
		internalError, _ := err.(*jwt.ValidationError)

		if internalError != nil {
			log.Debugf(it.ctx, "Start Errors. %v", internalError.Errors)
			if it.skipExpireCheck {
				internalError.Errors &= ^(jwt.ValidationErrorExpired | jwt.ValidationErrorIssuedAt)
				log.Debugf(it.ctx, "Skip Expire check")
			}

			// remove error?
			if internalError.Errors == 0 {
				err = nil
			} else {
				log.Debugf(it.ctx, "Updated Errors. %v", internalError.Errors)
			}
		}
	}

	// check error
	if err != nil {
		return nil, newTokenError(err)
	}

	result := &VerifiedJsonWebTokenImpl{
		token:rawToken,
	}

	// Token check
	projectId, err := result.GetProjectId()
	if !it.skipProjectIdCheck && projectId != it.service.GetProjectId() {
		log.Errorf(it.ctx, "Token replace attack? token[%v] require[%v] ", projectId, it.service.GetProjectId())
		return nil, newTokenError(errors.New("Token project id error."))
	}

	return result, nil

}