package gaefire

import (
	"golang.org/x/net/context"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"google.golang.org/appengine/log"
	"github.com/eaglesakura/gaefire"
)

type JsonWebTokenVerifierImpl struct {
	service            gaefire.ServiceAccount
	ctx                context.Context
	token              string
	trustedAud         []string
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

/**
 * "有効期限をチェックしない
 */
func (it *JsonWebTokenVerifierImpl)SkipExpireTime() gaefire.JsonWebTokenVerifier {
	it.skipExpireCheck = true
	return it
}

/**
 * 許可対象のAudienceを追加する
 * デフォルトではFirebase Service AccountのIDが登録される。
 */
func (it *JsonWebTokenVerifierImpl)AddTrustedAudience(aud string) gaefire.JsonWebTokenVerifier {
	if it.trustedAud == nil {
		it.trustedAud = []string{aud}
	} else {
		it.trustedAud = append(it.trustedAud, aud)
	}
	return it
}

/**
 * "aud"チェックをスキップする
 *
 * 他のプロジェクトに対して発行されたJWTを許可してしまうので、これを使用する場合は十分にセキュリティに注意を払う
 */
func (it *JsonWebTokenVerifierImpl)SkipProjectId() gaefire.JsonWebTokenVerifier {
	it.skipProjectIdCheck = true
	return it
}

/**
 * 全てのオプションに対し、有効であることが確認できればtrue
 */
func (it *JsonWebTokenVerifierImpl)Valid() (gaefire.VerifiedJsonWebToken, error) {
	// parse & verify
	rawToken, err := jwt.Parse(it.token, func(token *jwt.Token) (interface{}, error) {
		kid := token.Header["kid"].(string)
		if kid == "" {
			return nil, errors.New("NotFound kid")
		}

		publicKey, err := it.service.FindPublicKey(it.ctx, kid)
		if err != nil {
			log.Errorf(it.ctx, "Not found public key")
			return nil, err
		} else {
			return publicKey, nil
		}
	})

	if rawToken == nil {
		// Verify error
		log.Errorf(it.ctx, "Token format error")
		return nil, newTokenError(errors.New("Token format error."))
	}

	if err != nil && strings.Contains(err.Error(), "crypto/rsa") {
		// Verify error
		log.Errorf(it.ctx, "crypt/rsa error")
		return nil, newTokenError(err)
	}

	if rawToken.Claims != nil {
		if err != nil {
			log.Debugf(it.ctx, "Refresh Validate old[%v]", err.Error())
		}
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
		log.Errorf(it.ctx, "Validate error")
		return nil, newTokenError(err)
	}

	result := &VerifiedJsonWebTokenImpl{
		token:rawToken,
	}

	// Token check
	if !it.skipProjectIdCheck {
		trusted := false
		projectId, _ := result.GetAudience()
		if projectId == it.service.GetProjectId() {
			trusted = true
		} else if it.trustedAud != nil {
			for _, value := range it.trustedAud {
				if value == projectId {
					trusted = true
				}
			}
		}


		// 信頼できるIDが登録されていなかった
		if !trusted {
			log.Errorf(it.ctx, "Token replace attack? token[%v] service[%v] ", projectId, it.service.GetProjectId())
			return nil, newTokenError(errors.New("Token project id error."))
		}
	}

	return result, nil

}