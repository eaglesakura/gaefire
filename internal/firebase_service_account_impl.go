package internal

import (
	"crypto/rsa"
	"golang.org/x/net/context"
	"github.com/eaglesakura/gaefire/auth"
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
	"errors"
)
/**
 * service-account.jsonからロードしたサービスアカウント情報を定義する
 */
type ServiceAccountModel struct {
	ProjectId         string `json:"project_id,omitempty"`
	PrivateKeyId      string `json:"private_key_id,omitempty"`
	PrivateKey        string `json:"private_key,omitempty"`
	ClientEmail       string `json:"client_email,omitempty"`
	ClientId          string `json:"client_id,omitempty"`
	ClientX509CertUrl string `json:"client_x509_cert_url,omitempty"`
}

type FirebaseServiceAccountImpl struct {
	/**
	 * JSONをデコードしたそのままのデータ
	 */
	rawServiceAccount  ServiceAccountModel

	/**
	 * サービスアカウントの秘密鍵
	 */
	privateKey         *rsa.PrivateKey

	/**
	 * Google公開鍵キャッシュ
	 */
	googlePublicKeys   *PublicKeystore

	/**
	 * Firebase公開鍵キャッシュ
	 */
	firebasePublicKeys *PublicKeystore
}

/**
 * GCP Project IDを取得する
 */
func (it *FirebaseServiceAccountImpl)GetProjectId() string {
	return it.rawServiceAccount.ProjectId
}

/**
 * Service Accountのメールアドレスを取得する
 */
func (it *FirebaseServiceAccountImpl)GetAccountEmail() string {
	return it.rawServiceAccount.ClientEmail
}

/**
 * 署名のためのPrivate Keyを取得する
 */
func (it *FirebaseServiceAccountImpl)GetPrivateKey() *rsa.PrivateKey {
	return it.privateKey
}

/**
 * JWT署名検証のために公開鍵を検索する。
 *
 * デフォルトではServiceAccountsの公開鍵、もしくはGoogleの公開鍵を検索する
 */
func (it *FirebaseServiceAccountImpl)FindPublicKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {

	// Google公開鍵を探す
	if key, err := it.googlePublicKeys.FindPublicKey(kid); err == nil {
		return key
	}

	// Firebase公開鍵を探す
	if key, err := it.firebasePublicKeys.FindPublicKey(kid); err == nil {
		return key
	}

	// Google公開鍵をリフレッシュして探す
	if err := it.googlePublicKeys.Refresh(ctx); err != nil {
		return nil, err
	}
	if key, err := it.googlePublicKeys.FindPublicKey(kid); err == nil {
		return key
	}


	// Firebase公開鍵をリフレッシュして再探索
	if err := it.firebasePublicKeys.Refresh(ctx); err != nil {
		return nil, err
	}
	if key, err := it.firebasePublicKeys.FindPublicKey(kid); err == nil {
		return key
	}

	return nil, errors.New(fmt.Sprintf("Not Found Keystore[%v]", kid))
}

/**
 * ユーザー名を指定し、JWTの生成を開始する
 *
 * Firebase用に生成する場合、userUniqueIdは[1-36文字の英数]である必要がある。
 */
func (it *FirebaseServiceAccountImpl)NewFirebaseAuthTokenGenerator(userUniqueId string) gaefire.JsonWebTokenGenerator {
	return &JsonWebTokenGeneratorImpl{
		service:it,
		source:TokenSourceModel{
			StandardClaims:jwt.StandardClaims{
				ExpiresAt: time.Now().Unix() + 3600,
				IssuedAt:time.Now().Unix(),
				Audience:"https://identitytoolkit.googleapis.com/google.identity.identitytoolkit.v1.IdentityToolkit",
				Issuer:it.GetAccountEmail(),
				Subject:it.GetAccountEmail(),
			},
			Uid:userUniqueId,
			Claims:map[string]string{},
		},
		headers:map[string]string{},
	}
}

/**
 * Json Web TokenのVerifyオブジェクトを生成する
 */
func (it *FirebaseServiceAccountImpl)NewFirebaseAuthTokenVerifier(ctx context.Context, jwt string) gaefire.JsonWebTokenVerifier {
	return &JsonWebTokenVerifierImpl{
		service:it,
		ctx:ctx,
		token:jwt,
	}
}