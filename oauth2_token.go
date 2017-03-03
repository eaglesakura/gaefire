package gaefire

import (
	"net/http"
	"io/ioutil"
	"errors"
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/log"
	"net/url"
)

type OAuth2Token struct {
	/**
	 * UserID
	 */
	Email        string `json:"email,omitempty"`

	/**
	 * Access scopes
	 */
	Scopes       string `json:"scope,omitempty"`

	/**
	 * OAuth2 Access Token
	 */
	AccessToken  string `json:"access_token,omitempty"`

	/**
	 * Token Type "Bearer"
	 */
	TokenType    string `json:"token_type,omitempty"`

	/**
	 * OAuth2 Refresh token
	 */
	RefreshToken string `json:"refresh_token,omitempty"`
}

/**
 * http requestに認証を行なう
 */
func (it *OAuth2Token)Authorize(req *http.Request) {
	req.Header.Set("Authorization", it.TokenType + " " + it.AccessToken)
}

/**
 * トークンが有効であればtrue
 * ただし、有効期限のチェックを行わない。
 */
func (it *OAuth2Token)Valid(ctx context.Context) bool {
	if len(it.AccessToken) == 0 || len(it.TokenType) == 0 {
		return false
	}

	resp, err := urlfetch.Client(ctx).Get("https://www.googleapis.com/oauth2/v1/tokeninfo?access_token=" + it.AccessToken)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	} else if err != nil {
		log.Errorf(ctx, "OAuth2 validate error[%s]", err.Error())
		return false
	}

	tempToken := OAuth2Token{}
	buf, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(buf, &tempToken); err != nil {
		log.Errorf(ctx, "OAuth2 parse error[%s]", err.Error())
		return false
	}

	it.Scopes = tempToken.Scopes
	it.Email = tempToken.Email

	return true
}

/**
 * OAuth2トークンをリフレッシュする。
 *
 * ただし、it.RefreshTokenが含まれていない場合、リフレッシュは行えない。
 * また、ユーザーが明示的にアクセス権限を取り消している場合もトークンを取り出すことはできない。
 */
func (it *OAuth2Token)Refresh(ctx context.Context, clientId string, clientSecret string) error {
	if len(it.RefreshToken) == 0 {
		return errors.New("refresh token empty")
	}

	// fetch
	values := url.Values{}
	values.Add("client_id", clientId)
	values.Add("client_secret", clientSecret)
	values.Add("grant_type", "refresh_token")
	values.Add("refresh_token", it.RefreshToken)
	resp, err := urlfetch.Client(ctx).PostForm("https://www.googleapis.com/oauth2/v4/token", values)
	if (resp != nil && resp.Body != nil) {
		defer resp.Body.Close()
	} else {
		log.Errorf(ctx, "Https error %v", err.Error())
		return err
	}

	tempToken := OAuth2Token{}
	buf, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(buf, &tempToken); err != nil {
		log.Errorf(ctx, "OAuth2 parse error[%s]", err.Error())
		return err
	}

	if len(tempToken.AccessToken) == 0 {
		return errors.New("Access token not found")
	}

	it.AccessToken = tempToken.AccessToken
	if len(it.TokenType) == 0 {
		it.TokenType = "Bearer"
	}
	return nil
}