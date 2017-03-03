package gaefire

import (
	"github.com/eaglesakura/gaefire"
	fire_utils "github.com/eaglesakura/gaefire/utils"
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"io/ioutil"
	"os"
)

func newTestWebApp() gaefire.FirebaseWebApplication {
	fire := fire_utils.NewGaeFire()
	if json, err := fire.NewAssetManager().LoadFile("assets/firebase-web.json"); err != nil {
		panic(err)
	} else {
		return fire.NewWebApplication(json)
	}
}

func newOAuthTestData() UserOAuthTestData {
	result := UserOAuthTestData{}
	fire := fire_utils.NewGaeFire()
	if jsonBuf, err := fire.NewAssetManager().LoadFile("private/oauth-test-token.json"); err != nil {
		panic(err)
	} else {
		json.Unmarshal(jsonBuf, &result)
	}
	return result
}

/**
 * サービスアカウントの生成が行える
 */
func TestNewFirebaseWebApp(t *testing.T) {
	webApp := newTestWebApp()
	assert.NotNil(t, webApp)
}

type UserOAuthTestData struct {
	AccessCode string `json:"accessCode"`
}

/**
 * アクセスコードからOAuth2トークンを生成する
 */
func TestOAuth2TokenNew(t *testing.T) {
	ctx := fire_utils.NewContext(nil)
	defer ctx.Close()

	webApp := newTestWebApp()
	testData := newOAuthTestData()
	if len(testData.AccessCode) == 0 {
		// skip testing
		return
	}

	token0, err := webApp.NewUserAccountToken(ctx.GetAppengineContext(), testData.AccessCode)
	assert.Nil(t, err)
	assert.NotEqual(t, token0.Email, "")
	assert.NotEqual(t, token0.Scopes, "")
	assert.Equal(t, token0.TokenType, "Bearer")
	assert.NotEqual(t, token0.AccessToken, "")
	assert.NotEqual(t, token0.RefreshToken, "")
	ioutil.WriteFile("private/user-token0.txt", []byte(token0.AccessToken), os.ModePerm)
	ioutil.WriteFile("private/user-token0-refresh.txt", []byte(token0.RefreshToken), os.ModePerm)

	// 一度作ったトークンは無効である
	_, err = webApp.NewUserAccountToken(ctx.GetAppengineContext(), testData.AccessCode)
	assert.NotNil(t, err)

	// キャッシュ済みなので、キャッシュから取り出せなければならない
	token1, err := webApp.GetUserAccountToken(ctx.GetAppengineContext(), token0.RefreshToken)
	assert.Equal(t, token0.AccessToken, token1.AccessToken)
	ioutil.WriteFile("private/user-token1.txt", []byte(token1.AccessToken), os.ModePerm)
	ioutil.WriteFile("private/user-token1-refresh.txt", []byte(token1.RefreshToken), os.ModePerm)
}