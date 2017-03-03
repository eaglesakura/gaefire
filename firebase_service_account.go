package gaefire

import (
	"crypto/rsa"
	"golang.org/x/net/context"
)

/**
 * Firebaseのサービスアカウントを定義する。
 *
 * 下記の機能を提供する:
 * * JWT生成
 * * JWT検証
 * * サービスアカウントのOAuthトークン生成・リフレッシュ
 */
type FirebaseServiceAccount interface {
	/**
	 * GCP Project IDを取得する
	 */
	GetProjectId() string

	/**
	 * Service Accountのメールアドレスを取得する
	 */
	GetAccountEmail() string

	/**
	 * 署名のためのPrivate Keyを取得する
	 */
	GetPrivateKey() *rsa.PrivateKey

	/**
	 * 署名用のKeyIDを取得する
	 * これは公開鍵チェックの際 "kid"として利用できる
	 */
	GetPrivateKeyId() string

	/**
	 * JWT署名検証のために公開鍵を検索する。
	 *
	 * デフォルトではServiceAccountsの公開鍵、もしくはGoogleの公開鍵を検索する
	 */
	FindPublicKey(ctx context.Context, kid string) (*rsa.PublicKey, error)

	/**
	 * ユーザー名を指定し、JWTの生成を開始する
	 *
	 * Firebase用に生成する場合、userUniqueIdは[1-36文字の英数]である必要がある。
	 */
	NewFirebaseAuthTokenGenerator(userUniqueId string) JsonWebTokenGenerator

	/**
	 * Json Web TokenのVerifyオブジェクトを生成する
	 */
	NewFirebaseAuthTokenVerifier(ctx context.Context, jwt string) JsonWebTokenVerifier

	/**
	 * Json Web TokenのVerifyオブジェクトを生成する
	 * Google Play Service:Authによって認証されたトークンはGoogle経由でVerifyを行なうほうが効率的
	 */
	NewGoogleAuthTokenVerifier(ctx context.Context, jwt string) JsonWebTokenVerifier

	/**
	 * Service Accountとして認証するためのOAuth2トークンを取得する
	 *
	 * OAuth2トークンはMemcacheにキャッシュされ、再取得は最低限となるよう実装される。
	 */
	GetServiceAccountToken(ctx context.Context, scope string, addScopes ...string) (OAuth2Token, error)
}