package internal

import (
	"github.com/eaglesakura/gaefire/auth"
	"errors"
)

type JsonWebTokenVerifierStubImpl struct {
}


/**
 * "有効期限をチェックしない
 */
func (it *JsonWebTokenVerifierStubImpl)SkipExpireTime() gaefire.JsonWebTokenVerifier {
	return it
}

/**
 * 許可対象のAudienceを追加する
 * デフォルトではFirebase Service AccountのIDが登録される。
 */
func (it *JsonWebTokenVerifierStubImpl)AddTrustedAudience(aud string) gaefire.JsonWebTokenVerifier {
	return it
}

/**
 * "aud"チェックをスキップする
 *
 * 他のプロジェクトに対して発行されたJWTを許可してしまうので、これを使用する場合は十分にセキュリティに注意を払う
 */
func (it *JsonWebTokenVerifierStubImpl)SkipProjectId() gaefire.JsonWebTokenVerifier {
	return it
}

/**
 * 全てのオプションに対し、有効であることが確認できればtrue
 */
func (it *JsonWebTokenVerifierStubImpl)Verify(jwtToken string) (gaefire.VerifiedJsonWebToken, error) {
	return nil, errors.New("Error Verified")

}