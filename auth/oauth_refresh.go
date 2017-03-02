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

type UserRefreshTokenRequest struct {
	context      context.Context
	service      *FirebaseServiceAccount
	appInfo      WebApplicationInfo
	user         string // user email
	refreshToken string // user refresh token
}

/**
 * User OAuth2 cache,
 */
type _UserTokenCache struct {
	gaestore.BasicStoreData
	Account       string // gmail account
	Token         string // OAuth2 access token
	GoogleIdToken string // google signed jwt
}

func (it*FirebaseServiceAccount)NewUserTokenRefreshRequest(ctx context.Context, appInfo WebApplicationInfo, gmail string, refreshToken string) *UserRefreshTokenRequest {
	return &UserRefreshTokenRequest{
		service:it,
		context:ctx,
		appInfo: appInfo,
		user: gmail,
		refreshToken:refreshToken,
	}
}

func (it *UserRefreshTokenRequest)refreshOauthToken(_auth gaestore.ManagedStoreData) error {
	// https://developers.google.com/identity/protocols/OAuth2ServiceAccount
	ref, _ := _auth.(*_UserTokenCache)

	// fetch
	values := url.Values{}
	values.Add("client_id", it.appInfo.ClientId)
	values.Add("client_secret", it.appInfo.ClientSecret)
	values.Add("grant_type", "refresh_token")
	values.Add("refresh_token", it.refreshToken)
	resp, err := urlfetch.Client(it.context).PostForm("https://www.googleapis.com/oauth2/v4/token", values)
	if (resp != nil && resp.Body != nil) {
		defer resp.Body.Close()
	} else {
		log.Errorf(it.context, "Https error %v", err.Error())
		return err;
	}

	var mapping map[string]interface{}
	if unmarshalJson(resp, &mapping) != nil {
		log.Errorf(it.context, "user=%v", it.user)
		return errors.New("OAuth2 token failed")
	}

	token := FindStringValue(&mapping, "access_token")
	if token == "" {
		log.Errorf(it.context, "user=%v", it.user)
		return errors.New("Token not found.")
	}

	ref.Account = it.user
	ref.Token = token
	ref.GoogleIdToken = FindStringValue(&mapping, "id_token")

	return nil
}


//
// Get OAuth2 token,
//
func (it *UserRefreshTokenRequest)GetToken() (*OAuth2Token, error) {

	req := gaestore.NewMemcacheRequest(it.context)

	auth := _UserTokenCache{
		BasicStoreData:gaestore.BasicStoreData{
			Id:it.user,
		},
		Account:it.user,
	}

	info, err := req.SetExpireDate(time.Now().Add(_OAUTH2_CACHE_DURATION)).
		SetKindInfo(gaestore.KindInfo{Name:"__USER_ACCOUNT_AUTH__", Version:1}).
		LoadOrCreate(&auth, it.refreshOauthToken)
	if err != nil {
		return nil, err
	}

	return &OAuth2Token{
		rawToken:oauth2.Token{
			AccessToken:auth.Token,
			RefreshToken:it.refreshToken,
		},
		byCache:info.DataByCache,
		googleSignedToken:auth.GoogleIdToken,
	}, nil;
}
