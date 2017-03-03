package gaefire

import (
	"github.com/eaglesakura/gaefire"
	fire_utils "github.com/eaglesakura/gaefire/utils"
	"testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
)

func newTestServiceAccount() gaefire.FirebaseServiceAccount {
	fire := fire_utils.NewGaeFire()
	if json, err := fire.NewAssetManager().LoadFile("assets/firebase-admin.json"); err != nil {
		panic(err)
	} else {
		return fire.NewServiceAccount(json)
	}
}

/**
 * サービスアカウントの生成が行える
 */
func TestNewFirebaseServiceAccount(t *testing.T) {
	account := newTestServiceAccount()
	assert.NotNil(t, account)
}

/**
 * OAuth2トークンが生成される
 */
func TestServiceAccountAuthGen(t *testing.T) {
	ctx := fire_utils.NewContext(nil)
	defer ctx.Close()

	service := newTestServiceAccount()

	// トークンを取得する
	token1, err := service.GetServiceAccountToken(ctx.GetAppengineContext(), "https://www.googleapis.com/auth/firebase", "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/cloud-platform")
	assert.Nil(t, err)
	assert.NotEqual(t, token1.AccessToken, "")
	assert.NotEqual(t, token1.Scopes, "")
	assert.NotEqual(t, token1.Email, "")
	assert.NotEqual(t, token1.TokenType, "")

	ioutil.WriteFile("private/service-token0.txt", []byte(token1.AccessToken), os.ModePerm)
}

/**
 * OAuth2トークンがキャッシュされる
 */
func TestServiceAccountAuthCache(t *testing.T) {
	ctx := fire_utils.NewContext(nil)
	defer ctx.Close()

	service := newTestServiceAccount()

	// トークンを取得する
	token1, _ := service.GetServiceAccountToken(ctx.GetAppengineContext(), "https://www.googleapis.com/auth/firebase", "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/cloud-platform")
	ioutil.WriteFile("private/service-token1.txt", []byte(token1.AccessToken), os.ModePerm)

	// 再度取得する。
	// 2度目はキャッシュされているはずである
	token2, err := service.GetServiceAccountToken(ctx.GetAppengineContext(), "https://www.googleapis.com/auth/firebase", "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/cloud-platform")
	assert.Nil(t, err)
	ioutil.WriteFile("private/service-token2.txt", []byte(token1.AccessToken), os.ModePerm)

	// 2つのトークンは一致する
	assert.Equal(t, token1.AccessToken, token2.AccessToken)
}

/**
 * 認証用JWTが生成できる
 */
func TestServiceAccountJwtGen(t *testing.T) {
	ctx := fire_utils.NewContext(nil)
	defer ctx.Close()

	service := newTestServiceAccount()
	jwt, err := service.NewFirebaseAuthTokenGenerator("TestUserId").AddClaim("email", "example@example.com").AddHeader("kid", service.GetPrivateKeyId()).Generate()
	assert.Nil(t, err)
	assert.NotEqual(t, jwt, "")

	// 生成した鍵は、当然デコードを行えるはずである
	// audはFirebaseAuthTokenの仕様に従って "https://identitytoolkit.googleapis.com/google.identity.identitytoolkit.v1.IdentityToolkit" 固定
	validToken, err := service.NewFirebaseAuthTokenVerifier(ctx.GetAppengineContext(), jwt).AddTrustedAudience("https://identitytoolkit.googleapis.com/google.identity.identitytoolkit.v1.IdentityToolkit").Valid()
	assert.Nil(t, err)
	assert.NotNil(t, validToken)

	ioutil.WriteFile("private/firebase-auth-token.txt", []byte(jwt), os.ModePerm)
}