package gaefire

import (
	"golang.org/x/oauth2"
	"net/http"
	"io/ioutil"
	"encoding/json"
)
//
//type OAuth2Token struct {
//	rawToken          oauth2.Token
//	byCache           bool
//	googleSignedToken string
//}
//
//func (it *OAuth2Token)GetRawToken() *oauth2.Token {
//	return &it.rawToken
//}
//
//func (it *OAuth2Token)Sign(req *http.Request) {
//	//req.Header.Add("Authorization", "Bearer " + it.GetBearer())
//	it.rawToken.SetAuthHeader(req)
//}
//
//func (it *OAuth2Token)GetToken() string {
//	return it.rawToken.AccessToken
//}
//
//func (it *OAuth2Token)GetRefreshToken() string {
//	return it.rawToken.RefreshToken
//}
//
//func (it *OAuth2Token)ByCache() bool {
//	return it.byCache
//}

//
// get JWT signed by google
//
func (it *OAuth2Token)GetSignedToken() string {
	return it.googleSignedToken
}

type WebApplicationInfo struct {
	ClientId     string        `json:"client_id"`
	ClientSecret string        `json:"client_secret"`
}

type _WebApplicationJson struct {
	Web WebApplicationInfo `json:web`
}

/**
 * WebApplication情報をアセットからロードする
 * JSONはcloud-consoleから取得できる
 */
func NewWebApplicationFromFile(path string) (WebApplicationInfo, error) {
	buf, err := ioutil.ReadFile(path);
	if err != nil {
		return WebApplicationInfo{}, err
	} else {
		return NewWebApplicationFromJson(buf)
	}
}

/**
 * WebApplication情報をJSONからロードする
 * JSONはcloud-consoleから取得できる
 */
func NewWebApplicationFromJson(jsonText []byte) (WebApplicationInfo, error) {
	root := _WebApplicationJson{}
	err := json.Unmarshal(jsonText, &root)
	if err != nil {
		return WebApplicationInfo{}, err
	} else {
		return root.Web, nil
	}
}