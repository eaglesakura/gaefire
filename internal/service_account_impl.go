package gaefire

import (
	"crypto/rsa"
	"time"
	"fmt"
	"errors"
	"encoding/json"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/log"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
	"github.com/eaglesakura/gaefire"
)

const (
	GooglePublicKeystoreAccount = "securetoken@system.gserviceaccount.com"
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

// impl FirebaseServiceAccount
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

func NewServiceAccount(jsonBuf []byte) gaefire.ServiceAccount {
	if jsonBuf == nil {
		panic(errors.New("NotFound"))
	}

	result := &FirebaseServiceAccountImpl{

	}

	if json.Unmarshal(jsonBuf, &result.rawServiceAccount) != nil {
		panic(errors.New("Json parse failed"))
	}

	if privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(result.rawServiceAccount.PrivateKey)); err == nil {
		result.privateKey = privateKey
	} else {
		panic(errors.New("Private key error"))
	}

	if keystore := NewPublicKeystore(GooglePublicKeystoreAccount); keystore != nil {
		result.googlePublicKeys = keystore
	} else {
		panic(errors.New("Google Public key error"))
	}
	if keystore := NewPublicKeystore(result.rawServiceAccount.ClientEmail); keystore != nil {
		result.firebasePublicKeys = keystore
	} else {
		panic(errors.New("Firebase PUblic key error"))
	}
	return result
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
func (it *FirebaseServiceAccountImpl)GetClientEmail() string {
	return it.rawServiceAccount.ClientEmail
}

/**
 * クライアント識別IDを取得する
 */
func (it *FirebaseServiceAccountImpl)GetClientId() string {
	return it.rawServiceAccount.ClientId
}

/**
 * 署名のためのPrivate Keyを取得する
 */
func (it *FirebaseServiceAccountImpl)GetPrivateKey() *rsa.PrivateKey {
	return it.privateKey
}

/**
 * 署名用のKeyIDを取得する
 * これは公開鍵チェックの際 "kid"として利用できる
 */
func (it *FirebaseServiceAccountImpl)GetPrivateKeyId() string {
	return it.rawServiceAccount.PrivateKeyId
}

/**
 * JWT署名検証のために公開鍵を検索する。
 *
 * デフォルトではServiceAccountsの公開鍵、もしくはGoogleの公開鍵を検索する
 */
func (it *FirebaseServiceAccountImpl)FindPublicKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {

	// Google公開鍵を探す
	if key, err := it.googlePublicKeys.FindPublicKey(kid); err == nil {
		return key, nil
	}

	// Firebase公開鍵を探す
	if key, err := it.firebasePublicKeys.FindPublicKey(kid); err == nil {
		return key, nil
	}

	// Google公開鍵をリフレッシュして探す
	if err := it.googlePublicKeys.Refresh(ctx); err != nil {
		return nil, err
	}
	if key, err := it.googlePublicKeys.FindPublicKey(kid); err == nil {
		return key, nil
	}


	// Firebase公開鍵をリフレッシュして再探索
	if err := it.firebasePublicKeys.Refresh(ctx); err != nil {
		return nil, err
	}
	if key, err := it.firebasePublicKeys.FindPublicKey(kid); err == nil {
		return key, nil
	}

	log.Errorf(ctx, fmt.Sprintf("Not Found Keystore[%v]", kid))
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
				Issuer:it.GetClientEmail(),
				Subject:it.GetClientEmail(),
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

/**
 * Json Web TokenのVerifyオブジェクトを生成する
 * Google Play Service:Authによって認証されたトークンはGoogle経由でVerifyを行なうほうが効率的
 */
func (it *FirebaseServiceAccountImpl)NewGoogleAuthTokenVerifier(ctx context.Context, jwt string) gaefire.JsonWebTokenVerifier {
	client := urlfetch.Client(ctx)
	resp, err := client.Head("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + jwt)
	defer resp.Body.Close()

	if err != nil {
		panic(err)
		return nil
	}

	if resp.StatusCode != 200 {
		// 少なくともGoogleが発行したトークンでは無いので、無効である
		return &JsonWebTokenVerifierStubImpl{
		}
	}

	return &JsonWebTokenVerifierImpl{
		service:it,
		ctx:ctx,
		token:jwt,
	}
}

/**
 * Service Accountとして認証するためのOAuth2トークンを取得する
 */
func (it *FirebaseServiceAccountImpl)GetServiceAccountToken(ctx context.Context, scope string, addScopes ...string) (gaefire.OAuth2Token, error) {
	token := &OAuth2RefreshRequest{ctx:ctx, serviceAccount:it }
	token.AddScope(scope)
	for _, value := range addScopes {
		token.AddScope(value)
	}
	return token.GetToken()
}

/**
 * Service Accountとして認証するためのOAuth2トークンを生成する
 * この生成結果はキャッシュされず、必ず新しいものとなる
 */
func (it *FirebaseServiceAccountImpl)NewServiceAccountToken(ctx context.Context, scope string, addScopes ...string) (gaefire.OAuth2Token, error) {
	token := &OAuth2RefreshRequest{ctx:ctx, serviceAccount:it }
	token.AddScope(scope)
	for _, value := range addScopes {
		token.AddScope(value)
	}
	return token._newServiceOauth2Token()
}
