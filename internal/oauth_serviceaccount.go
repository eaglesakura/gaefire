package gaefire

import (
	"strings"
	"serverutil"
	"time"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"net/url"
	"golang.org/x/oauth2"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/log"
	"gaestore"
	"github.com/eaglesakura/gaefire/internal"
)

var (
	_OAUTH2_CACHE_DURATION, _ = time.ParseDuration("3500s")
)

type ServiceAccountOAuthTokenCache struct {
	OAuth2Token
	Scopes string
}

type ServiceAccountAuthRequest struct {
	service FirebaseServiceAccount
	context context.Context
	scope   string
}

/**
 * Memcache用データ
 */
type _ServiceAccountAuthCache struct {
	gaestore.BasicStoreData
	Token     string
	GoogleJwt string
}

func (it*FirebaseServiceAccount)NewServiceAccountAuthRequest(ctx context.Context) *ServiceAccountAuthRequest {
	return &ServiceAccountAuthRequest{
		service:it,
		context:ctx,
	}
}

//
// Add access scopes
//
func (it *ServiceAccountAuthRequest)AddScope(scope string) *ServiceAccountAuthRequest {
	if strings.Contains(it.scope, scope) {
		return it
	} else {
		if len(it.scope) > 0 {
			it.scope += (" " + scope)
		} else {
			it.scope = scope
		}
	}
	return it
}

func (it *ServiceAccountAuthRequest)newOauthToken(_auth gaestore.ManagedStoreData) error {
	// https://developers.google.com/identity/protocols/OAuth2ServiceAccount
	ref, _ := _auth.(*_ServiceAccountAuthCache)

	gen := &TokenGenerator{
		service:it.service,
		source:internal.TokenSourceModel{
			StandardClaims:jwt.StandardClaims{
				ExpiresAt: time.Now().Unix() + 3600,
				IssuedAt:time.Now().Unix(),
				Audience:"https://www.googleapis.com/oauth2/v4/token",
				Issuer:it.service.GetAccountEmail(),
				Subject:it.service.GetAccountEmail(),
			},
			Scope:it.scope,
			Claims:map[string]string{},
		},
		headers:map[string]string{},
	}

	jwtToken, _ := gen.Generate()
	if jwtToken == "" {
		return errors.New("JwtToken Generate failed")
	}

	// fetch
	values := url.Values{}
	values.Add("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	values.Add("assertion", jwtToken)
	resp, err := urlfetch.Client(it.context).PostForm("https://www.googleapis.com/oauth2/v4/token", values)
	if (resp != nil && resp.Body != nil) {
		defer resp.Body.Close()
	} else {
		log.Errorf(it.context, "Https error %v", err.Error())
		return err;
	}

	var mapping map[string]interface{}
	if unmarshalJson(resp, &mapping) != nil {
		log.Errorf(it.context, "jwt=%v", jwtToken)
		return errors.New("OAuth2 token failed")
	}

	token := FindStringValue(&mapping, "access_token")
	if token == "" {
		log.Errorf(it.context, "jwt=%v", jwtToken)
		return errors.New("Token not found.")
	}

	ref.Token = token
	ref.GoogleJwt = FindStringValue(&mapping, "id_token")

	return nil
}

//
// Get OAuth2 token,
//
func (it *ServiceAccountAuthRequest)GetToken() (*OAuth2Token, error) {

	req := gaestore.NewMemcacheRequest(it.context)

	auth := _ServiceAccountAuthCache{
		BasicStoreData:gaestore.BasicStoreData{
			Id:it.service.serviceAccount.ClientEmail + "-" + serverutil.ToMD5(it.scope),
		},
	}

	info, err := req.SetExpireDate(time.Now().Add(_OAUTH2_CACHE_DURATION)).
		SetKindInfo(gaestore.KindInfo{Name:"__SERVICE_ACCOUNT_AUTH__", Version:1}).
		LoadOrCreate(&auth, it.newOauthToken)
	if err != nil {
		return nil, err
	}

	return &OAuth2Token{
		rawToken:oauth2.Token{
			AccessToken:auth.Token,
		},
		byCache:info.DataByCache,
	}, nil;
}
