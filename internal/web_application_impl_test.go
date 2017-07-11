package gaefire

import (
	"encoding/json"
	"github.com/eaglesakura/gaefire"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func newTestWebApp() gaefire.WebApplication {
	if json, err := NewAssetManager().LoadFile("assets/firebase-web.json"); err != nil {
		panic(err)
	} else {
		return NewWebApplication(json)
	}
}

func newOAuthTestData() UserOAuthTestData {
	result := UserOAuthTestData{}
	if jsonBuf, err := NewAssetManager().LoadFile("private/oauth-test-token.json"); err == nil {
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
	AccessCode    string `json:"accessCode,omitempty"`
	RefreshToken  string `json:"refreshToken,omitempty"`
	GoogleIdToken string `json:"googleIdToken,omitempty"`
}

/**
 * アクセスコードからOAuth2トークンを生成する
 */
func TestOAuth2TokenNew(t *testing.T) {
	ctx := NewContext(nil)
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

/**
 * リフレッシュトークンからOAuthTokenを再取得する
 */
func TestOAuth2TokenRefresh(t *testing.T) {
	ctx := NewContext(nil)
	defer ctx.Close()

	webApp := newTestWebApp()
	testData := newOAuthTestData()
	if len(testData.RefreshToken) == 0 {
		// skip testing
		return
	}

	token2, err := webApp.GetUserAccountToken(ctx.GetAppengineContext(), testData.RefreshToken)
	assert.Nil(t, err)
	assert.NotEqual(t, token2.Email, "")
	assert.NotEqual(t, token2.Scopes, "")
	assert.Equal(t, token2.TokenType, "Bearer")
	assert.NotEqual(t, token2.AccessToken, "")
	assert.NotEqual(t, token2.RefreshToken, "")
	ioutil.WriteFile("private/user-token2.txt", []byte(token2.AccessToken), os.ModePerm)

	token3, err := webApp.GetUserAccountToken(ctx.GetAppengineContext(), testData.RefreshToken)
	assert.Nil(t, err)
	assert.Equal(t, token2.AccessToken, token3.AccessToken)
	assert.Equal(t, token2.RefreshToken, token3.RefreshToken)
	ioutil.WriteFile("private/user-token3.txt", []byte(token3.AccessToken), os.ModePerm)
}
