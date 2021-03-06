package gaefire

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type OAuth2Token struct {
	/**
	 * mail
	 */
	Email string `json:"email,omitempty"`

	/**
	 * Access scopes
	 */
	Scopes string `json:"scope,omitempty"`

	/**
	 * OAuth2 Access Token
	 */
	AccessToken string `json:"access_token,omitempty"`

	/**
	 * Token Type "Bearer"
	 */
	TokenType string `json:"token_type,omitempty"`

	/**
	 * OAuth2 Refresh token
	 */
	RefreshToken string `json:"refresh_token,omitempty"`

	/**
	 * OAuth2 aud
	 */
	Audience string `json:"aud,omitempty"`
}

/**
 * http requestに認証を行なう
 */
func (it *OAuth2Token) Authorize(req *http.Request) {
	req.Header.Set("Authorization", it.TokenType+" "+it.AccessToken)
}

/**
 * トークンが有効であればtrue
 * ただし、有効期限のチェックを行わない。
 */
func (it *OAuth2Token) Valid(ctx context.Context) bool {
	if len(it.AccessToken) == 0 || len(it.TokenType) == 0 {
		return false
	}

	resp, err := newHttpClient(ctx).Get("https://www.googleapis.com/oauth2/v3/tokeninfo?access_token=" + it.AccessToken)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		logError(fmt.Sprintf("OAuth2 validate error[%s]", err.Error()))
		return false
	}

	if resp.StatusCode != 200 {
		logError("OAuth2 invalid_token")
		return false
	}

	tempToken := OAuth2Token{}
	buf, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(buf, &tempToken); err != nil {
		logError(fmt.Sprintf("OAuth2 parse error[%s]", err.Error()))
		return false
	}

	if len(tempToken.Scopes) > 0 {
		it.Scopes = tempToken.Scopes
	}
	if len(tempToken.Email) > 0 {
		it.Email = tempToken.Email
	}
	if len(tempToken.Audience) > 0 {
		it.Audience = tempToken.Audience
	}

	return true
}

func newHttpClient(_ context.Context) *http.Client {
	// タイムアウトを30秒に延長
	result := &http.Client{
		Timeout: 30 * time.Second,
	}
	return result
}

/**
 * OAuth2トークンをリフレッシュする。
 *
 * ただし、it.RefreshTokenが含まれていない場合、リフレッシュは行えない。
 * また、ユーザーが明示的にアクセス権限を取り消している場合もトークンを取り出すことはできない。
 */
func (it *OAuth2Token) Refresh(ctx context.Context, clientId string, clientSecret string) error {
	if len(it.RefreshToken) == 0 {
		return errors.New("refresh token empty")
	}

	// fetch
	values := url.Values{}
	values.Add("client_id", clientId)
	values.Add("client_secret", clientSecret)
	values.Add("grant_type", "refresh_token")
	values.Add("refresh_token", it.RefreshToken)
	resp, err := newHttpClient(ctx).PostForm("https://www.googleapis.com/oauth2/v4/token", values)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	} else {
		logError(fmt.Sprintf("Https error %v", err))
		return err
	}

	tempToken := OAuth2Token{}
	buf, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(buf, &tempToken); err != nil {
		logError(fmt.Sprintf("OAuth2 parse error[%s]", err))
		return err
	}

	if len(tempToken.AccessToken) == 0 {
		return errors.New("access token not found")
	}

	it.AccessToken = tempToken.AccessToken
	if len(it.TokenType) == 0 {
		it.TokenType = "Bearer"
	}
	return nil
}
