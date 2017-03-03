package internal

import (
	"strings"
	"time"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"net/url"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/log"
	fire_datastore "github.com/eaglesakura/gaefire/datastore"
	fire_auth "github.com/eaglesakura/gaefire/auth"
	fire_util "github.com/eaglesakura/gaefire/utils"
)

var (
	_OAUTH2_CACHE_DURATION = time.Duration(55 * time.Minute)
	_OAUTH2_KIND_INFO = fire_datastore.KindInfo{Name:"oauth2-token-cache", Version:1}
)

type OAuth2RefreshRequest struct {
	ctx            context.Context

	serviceAccount fire_auth.FirebaseServiceAccount // for Service Account
	scope          string                           // for Service Account

	webApplication fire_auth.FirebaseWebApplication // for User Account
	accessCode     string                           // for User Account
	refreshToken   string                           // for User Account
}


/**
 * アクセススコープを追加する
 */
func (it *OAuth2RefreshRequest)AddScope(scope string) *OAuth2RefreshRequest {
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


/**
 * サービスアカウントのOAuth2情報をリフレッシュする
 */
func (it *OAuth2RefreshRequest)_newServiceOauth2Token() (fire_auth.OAuth2Token, error) {
	// https://developers.google.com/identity/protocols/OAuth2ServiceAccount
	gen := &JsonWebTokenGeneratorImpl{
		service:it.serviceAccount,
		source:TokenSourceModel{
			StandardClaims:jwt.StandardClaims{
				ExpiresAt: time.Now().Unix() + 3600,
				IssuedAt:time.Now().Unix(),
				Audience:"https://www.googleapis.com/oauth2/v4/token",
				Issuer:it.serviceAccount.GetAccountEmail(),
				Subject:it.serviceAccount.GetAccountEmail(),
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
	resp, err := urlfetch.Client(it.ctx).PostForm("https://www.googleapis.com/oauth2/v4/token", values)
	if (resp != nil && resp.Body != nil) {
		defer resp.Body.Close()
	} else {
		log.Errorf(it.ctx, "Https error %v", err.Error())
		return fire_auth.OAuth2Token{}, err;
	}

	token := fire_auth.OAuth2Token{}
	if err := fire_util.UnmarshalJson(resp, &token); err != nil {
		log.Errorf(it.ctx, "jwt=%v", jwtToken)
		return fire_auth.OAuth2Token{}, err
	}

	if !token.Valid(it.ctx) {
		errors.New("OAuth2 token validate error")
	}
	//token.UserId = it.serviceAccount.GetAccountEmail()
	//token.Scopes = it.scope
	return token, nil
}

/**
 * ユーザー用のOAuth2トークンを取得する
 */
func (it *OAuth2RefreshRequest)_newUserOauth2Token() (fire_auth.OAuth2Token, error) {
	token := fire_auth.OAuth2Token{}

	if len(it.refreshToken) > 0 {
		// リフレッシュトークンがあるならば、リフレッシュが可能である
		token.RefreshToken = it.refreshToken
		err := token.Refresh(it.ctx, it.webApplication.GetClientId(), it.webApplication.GetClientSecret())
		if err != nil {
			// リフレッシュに失敗した。恐らくExpireされている。
			return fire_auth.OAuth2Token{}, err
		}
	} else {
		// 新規にトークンを取得する
		// fetch
		values := url.Values{}
		values.Add("client_id", it.webApplication.GetClientId())
		values.Add("client_secret", it.webApplication.GetClientSecret())
		values.Add("grant_type", "authorization_code")
		values.Add("code", it.accessCode)
		resp, err := urlfetch.Client(it.ctx).PostForm("https://www.googleapis.com/oauth2/v4/token", values)
		if (resp != nil && resp.Body != nil) {
			defer resp.Body.Close()
		} else {
			log.Errorf(it.ctx, "Https error %v", err.Error())
			return fire_auth.OAuth2Token{}, err;
		}

		if err := fire_util.UnmarshalJson(resp, &token); err != nil {
			return fire_auth.OAuth2Token{}, err
		}
	}

	if !token.Valid(it.ctx) {
		// 何らかの原因でToken検証に失敗した
		return errors.New("OAuth2 token validate error.")
	}

	return token, nil
}

/**
 * OAuth2トークンを取得する
 */
func (it *OAuth2RefreshRequest)GetToken() (fire_auth.OAuth2Token, error) {

	var keyId string
	if it.serviceAccount != nil {
		keyId = it.serviceAccount.GetAccountEmail() + "-" + fire_util.GenMD5(it.scope)
	} else {
		keyId = "user-" + fire_util.GenMD5(it.refreshToken)
	}

	req := fire_datastore.NewMemcacheRequest(it.ctx).
		SetKindInfo(_OAUTH2_KIND_INFO).
		SetExpireDate(time.Now().
		Add(_OAUTH2_CACHE_DURATION)).
		SetId(keyId)
	token := fire_auth.OAuth2Token{}

	// Memcacheを優先ロードし、データが見つからなければ新規に取得する
	err := req.Load(&token, func(ref interface{}) error {
		tokenRef, _ := ref.(*fire_auth.OAuth2Token)
		var err error
		*tokenRef, err = it._newServiceOauth2Token()
		return err
	})

	if it.webApplication != nil {
		// ユーザー権限の場合、検証を行なう
		if !token.Valid(it.ctx) {
			// 権限が不正のため、リフレッシュする
			err := token.Refresh(it.ctx, it.webApplication.GetClientId(), it.webApplication.GetClientSecret())
			if err != nil {
				// リフレッシュにも失敗したため、恐らくユーザーはExpireしている
				return fire_auth.OAuth2Token{}, err
			}
		}
	}

	return token, err
}
