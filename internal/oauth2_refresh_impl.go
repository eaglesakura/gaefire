package gaefire

import (
	"strings"
	"time"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"net/url"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/log"
	"github.com/eaglesakura/gaefire"
)

var (
	_OAUTH2_CACHE_DURATION = time.Duration(55 * time.Minute)
	_OAUTH2_KIND_INFO = gaefire.KindInfo{Name:"oauth2-token-cache", Version:1}
)

type OAuth2RefreshRequest struct {
	ctx            context.Context

	serviceAccount gaefire.FirebaseServiceAccount // for Service Account
	scope          string                         // for Service Account

	webApplication gaefire.FirebaseWebApplication // for User Account
	accessCode     string                         // for User Account
	refreshToken   string                         // for User Account
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
func (it *OAuth2RefreshRequest)_newServiceOauth2Token() (gaefire.OAuth2Token, error) {
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
		return gaefire.OAuth2Token{}, errors.New("JwtToken Generate failed")
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
		return gaefire.OAuth2Token{}, err;
	}

	token := gaefire.OAuth2Token{}
	if err := UnmarshalJson(resp, &token); err != nil {
		log.Errorf(it.ctx, "jwt=%v", jwtToken)
		return gaefire.OAuth2Token{}, err
	}

	if !token.Valid(it.ctx) {
		return gaefire.OAuth2Token{}, errors.New("OAuth2 token validate error")
	}
	//token.UserId = it.serviceAccount.GetAccountEmail()
	//token.Scopes = it.scope
	return token, nil
}

/**
 * ユーザー用のOAuth2トークンを取得する
 */
func (it *OAuth2RefreshRequest)_newUserOauth2Token() (gaefire.OAuth2Token, error) {
	token := gaefire.OAuth2Token{}

	if len(it.refreshToken) > 0 {
		// リフレッシュトークンがあるならば、リフレッシュが可能である
		token.RefreshToken = it.refreshToken
		err := token.Refresh(it.ctx, it.webApplication.GetClientId(), it.webApplication.GetClientSecret())
		if err != nil {
			// リフレッシュに失敗した。恐らくExpireされている。
			return gaefire.OAuth2Token{}, err
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
			return gaefire.OAuth2Token{}, err;
		}

		if err := UnmarshalJson(resp, &token); err != nil {
			return gaefire.OAuth2Token{}, err
		}
	}

	if !token.Valid(it.ctx) {
		// 何らかの原因でToken検証に失敗した
		return gaefire.OAuth2Token{}, errors.New("OAuth2 token validate error.")
	}

	return token, nil
}

/**
 * OAuth2トークンを取得する
 */
func (it *OAuth2RefreshRequest)GetToken() (gaefire.OAuth2Token, error) {

	var keyId string
	if it.serviceAccount != nil {
		keyId = it.serviceAccount.GetAccountEmail() + "-" + GenMD5(it.scope)
	} else {
		keyId = "user-" + GenMD5(it.refreshToken)
	}

	req := gaefire.NewMemcacheRequest(it.ctx).
		SetKindInfo(_OAUTH2_KIND_INFO).
		SetExpireDate(time.Now().
		Add(_OAUTH2_CACHE_DURATION)).
		SetId(keyId)
	token := gaefire.OAuth2Token{}

	// Memcacheを優先ロードし、データが見つからなければ新規に取得する
	if err := req.Load(&token, func(ref interface{}) error {
		tokenRef, _ := ref.(*gaefire.OAuth2Token)
		var err error
		*tokenRef, err = it._newServiceOauth2Token()
		return err
	}); err != nil {
		return gaefire.OAuth2Token{}, err
	}

	if it.webApplication != nil {
		// ユーザー権限の場合、検証を行なう
		if !token.Valid(it.ctx) {
			// 権限が不正のため、リフレッシュする
			err := token.Refresh(it.ctx, it.webApplication.GetClientId(), it.webApplication.GetClientSecret())
			if err != nil {
				// リフレッシュにも失敗したため、恐らくユーザーはExpireしている
				return gaefire.OAuth2Token{}, err
			}
		}
	}

	return token, nil
}
