package gaefire

import "net/http"

type OAuth2Token struct {
	/**
	 * OAuth2 Access Token
	 */
	AccessToken  string `json:"access_token"`

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