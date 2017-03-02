package gaefire

import (
	"gaestore"
	"golang.org/x/net/context"
	"net/url"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/log"
	"errors"
	"time"
	"golang.org/x/oauth2"
)

type UserTokenRequest struct {
	context    context.Context
	service    *FirebaseServiceAccount
	appInfo    WebApplicationInfo
	user       string // user email
	serverCode string // user init code
}

type UserAccountInfo struct {
	AccessToken   string // OAuth2 access token
	RefreshToken  string // OAuth2 access token
	GoogleIdToken string // google signed jwt
}

func (it*FirebaseServiceAccount)NewUserAuthRequest(ctx context.Context, appInfo WebApplicationInfo, gmail string, code string) *UserTokenRequest {
	return &UserTokenRequest{
		service:it,
		context:ctx,
		appInfo: appInfo,
		user:gmail,
		serverCode: code,
	}
}

func (it *UserTokenRequest)newOauthToken(ref *UserAccountInfo) error {
	// fetch
	values := url.Values{}
	values.Add("client_id", it.appInfo.ClientId)
	values.Add("client_secret", it.appInfo.ClientSecret)
	values.Add("grant_type", "authorization_code")
	values.Add("code", it.serverCode)
	resp, err := urlfetch.Client(it.context).PostForm("https://www.googleapis.com/oauth2/v4/token", values)
	if (resp != nil && resp.Body != nil) {
		defer resp.Body.Close()
	} else {
		log.Errorf(it.context, "Https error %v", err.Error())
		return err;
	}

	var mapping map[string]interface{}
	if unmarshalJson(resp, &mapping) != nil {
		return errors.New("OAuth2 token failed")
	}

	token := FindStringValue(&mapping, "access_token")
	refresh_token := FindStringValue(&mapping, "refresh_token")
	if token == "" || refresh_token == "" {
		log.Errorf(it.context, "Token not found")
		return errors.New("Token not found.")
	}

	ref.AccessToken = token
	ref.RefreshToken = refresh_token
	ref.GoogleIdToken = FindStringValue(&mapping, "id_token")

	return nil
}


//
// Get OAuth2 token,
//
func (it *UserTokenRequest)GetToken() (*OAuth2Token, error) {

	info := UserAccountInfo{}
	err := it.newOauthToken(&info)
	if err != nil {
		return nil, err
	}

	// put memcache
	auth := _ServiceAccountAuthCache{
		BasicStoreData:gaestore.BasicStoreData{
			Id:it.service.serviceAccount.ClientEmail,
		},
		Token:info.AccessToken,
		GoogleJwt:info.GoogleIdToken,
	}

	err = gaestore.NewMemcacheRequest(it.context).
		SetExpireDate(time.Now().Add(_OAUTH2_CACHE_DURATION)).
		SetKindInfo(gaestore.KindInfo{Name:"__USER_ACCOUNT_AUTH__", Version:1}).
		Save(&auth)

	return &OAuth2Token{
		rawToken:oauth2.Token{
			AccessToken:info.AccessToken,
			RefreshToken:info.RefreshToken,
		},
		byCache:false,
		googleSignedToken:info.GoogleIdToken,
	}, err;
}
