package gaefire

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

/*
Firebaseのサービスアカウントを定義する。

* 下記の機能を提供する:
	* JWT生成
 	* JWT検証
 	* サービスアカウントのOAuthトークン生成・リフレッシュ
	* ユーザー認証(OAuth2, Firebase/JWT, ServiceAccount/JWT, GoogleIdToken/JWT)
*/
type ServiceAccount interface {
	/*
		GCP Project IDを取得する
	*/
	GetProjectId() string

	/*
		Service Accountのメールアドレスを取得する
	*/
	GetClientEmail() string

	/*
		サービスアカウント識別IDを取得する
		oauth2トークンを生成した場合、audに対応される
	*/
	GetClientId() string

	/*
		署名のためのPrivate Keyを取得する
	*/
	GetPrivateKey() *rsa.PrivateKey

	/**
	 * 署名用のKeyIDを取得する
	 * これは公開鍵チェックの際 "kid"として利用できる
	 */
	GetPrivateKeyId() string

	/*
		JWT署名検証のために公開鍵を検索する。
		デフォルトではServiceAccountsの公開鍵、もしくはGoogleの公開鍵を検索する
	*/
	FindPublicKey(ctx context.Context, kid string) (*rsa.PublicKey, error)

	/*
		ユーザー名を指定し、JWTの生成を開始する
		Firebase用に生成する場合、userUniqueIdは[1-36文字の英数]である必要がある。
	*/
	NewFirebaseAuthTokenGenerator(userUniqueId string) JsonWebTokenGenerator

	/**
	 * Json Web TokenのVerifyオブジェクトを生成する
	 */
	NewFirebaseAuthTokenVerifier(ctx context.Context, jwt string) JsonWebTokenVerifier

	/*
			Json Web TokenのVerifyオブジェクトを生成する
		 	Google Play Service:Authによって認証されたトークンはGoogle経由でVerifyを行なうほうが効率的
	*/
	NewGoogleAuthTokenVerifier(ctx context.Context, jwt string) JsonWebTokenVerifier

	/*
		Service Accountとして認証するためのOAuth2トークンを取得する
		OAuth2トークンはMemcacheにキャッシュされ、再取得は最低限となるよう実装される。
	*/
	GetServiceAccountToken(ctx context.Context, scope string, addScopes ...string) (OAuth2Token, error)

	/*
		Service Accountとして認証するためのOAuth2トークンを取得する
		この生成結果はキャッシュされず、必ず新しいものとなる
	*/
	NewServiceAccountToken(ctx context.Context, scope string, addScopes ...string) (OAuth2Token, error)
}

/*
	データをSHA256で署名する.
	これはFirebase StorageのURL生成用の関数として利用できる.
*/
func SignSHA256(account ServiceAccount, buffer []byte) ([]byte, error) {
	sum := sha256.Sum256(buffer)
	return rsa.SignPKCS1v15(
		rand.Reader,
		account.GetPrivateKey(),
		crypto.SHA256,
		sum[:],
	)
}
