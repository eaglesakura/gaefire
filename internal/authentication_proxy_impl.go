package gaefire

import (
	"github.com/eaglesakura/gaefire"
	"net/http"
	"github.com/eaglesakura/swagger-go-core/errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
	"time"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/appengine"
	"encoding/base64"
	"google.golang.org/appengine/log"
	"github.com/eaglesakura/swagger-go-core/swag-port"
)

var (
	// セキュリティ上許されないヘッダ
	SecurityCheckHeaders []string = []string{gaefire.HttpXHeaderUserInfo }
)

/**
 * Swaggerで定義した情報を元に、セキュリティ設定を行う。
 *
 * https://servicecontrol.googleapis.com/v1/services/{ENDPOINTS_SERVICE_NAME}:check
 * https://servicemanagement.googleapis.com/v1/services/{ENDPOINTS_SERVICE_NAME}/configs/{ENDPOINTS_SERVICE_VERSION}
 */
type SwaggerJsonModel struct {
	Host                string `json:"host"`

	SecurityDefinitions struct {
				    ApiKey        *struct {
					    Name string `json:"name"`
					    In   string `json:"in"`
				    }`json:"api_key,omitempty"`
				    GoogleIdToken *struct {
					    Issuer    string `json:"x-google-issuer"`
					    Audiences []string `json:"x-google-audiences"`
				    }`json:"google_id_token,omitempty"`
			    }`json:"securityDefinitions"`
}

/**
 * https://servicecontrol.googleapis.com/v1/services/{host}:check にてバリデーションを行う
 */
type ServiceCheckModel struct {
	Operation struct {
			  OperationId   string `json:"operationId"`
			  OperationName string `json:"operationName"`
			  ConsumerId    string `json:"consumerId"`
			  StartTime     string `json:"startTime"`
		  }`json:"operation"`
}

type ServiceCheckResultModel struct {
	OperationId     string `json:"operationId"`
	ServiceConfigId string `json:"serviceConfigId"`
	CheckErrors     *[]struct {
		Code   string `json:"code"`
		Detail string `json:"detail"`
	}`json:"checkErrors,omitempty"`
}

type AuthenticationProxyImpl struct {
	/**
	 * 処理対象のサービスアカウント情報
	 */
	ServiceAccount gaefire.FirebaseServiceAccount

	/**
	 * セキュリティチェックソースとなるSwagger情報
	 */
	Swagger        SwaggerJsonModel
}

/**
 * API Keyが入力されている場合、APIキーの妥当性をチェックする
 */
func (it *AuthenticationProxyImpl)validApiKey(ctx context.Context, r *http.Request, result *gaefire.AuthenticationInfo) error {
	if it.Swagger.SecurityDefinitions.ApiKey == nil {
		// API Keyはチェック対象ではない
		return nil
	}

	var apiKey string
	switch it.Swagger.SecurityDefinitions.ApiKey.In {
	case "query":
		apiKey = r.URL.Query().Get(it.Swagger.SecurityDefinitions.ApiKey.Name)
	case "header":
		apiKey = r.Header.Get(it.Swagger.SecurityDefinitions.ApiKey.Name)
	}

	if len(apiKey) == 0 {
		// APIキーが指定されていない
		return nil
	}

	// Api Keyが使用されているので、妥当性をチェックする
	model := ServiceCheckModel{}
	model.Operation.OperationId = appengine.RequestID(ctx)
	model.Operation.OperationName = "check:" + model.Operation.OperationId
	model.Operation.ConsumerId = "api_key:" + apiKey
	model.Operation.StartTime = fmt.Sprintf("%v", time.Now())

	buf, _ := json.Marshal(model)

	resp, err := urlfetch.Client(ctx).Post("https://servicecontrol.googleapis.com/v1/services/" + it.Swagger.Host + ":check", "application/json", bytes.NewReader(buf))
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil || resp.StatusCode != 200 {
		log.Errorf(ctx, "User security service check error api_key[%v]", apiKey)
		return errors.New(http.StatusBadRequest, "ApiKey check error")
	}

	// バッファを読み取る
	buf, _ = ioutil.ReadAll(resp.Body)

	validModel := ServiceCheckResultModel{}
	if err := json.Unmarshal(buf, &validModel); err != nil {
		return errors.New(http.StatusBadRequest, "ApiKey parse error")
	}

	if len(*validModel.CheckErrors) > 0 {
		// エラーが発生している
		log.Errorf(ctx, "User api_key check error error[%v]", (*validModel.CheckErrors)[0].Detail)
		return errors.New(http.StatusBadRequest, (*validModel.CheckErrors)[0].Detail)
	}

	// API Keyは問題ない
	result.ApiKey = &apiKey
	return nil
}

func (it *AuthenticationProxyImpl)validOAuth2(ctx context.Context, authorization string, result *gaefire.AuthenticationInfo) error {
	// OAuth2 Tokenである
	token := gaefire.OAuth2Token{TokenType:"Bearer", AccessToken:authorization}
	if !token.Valid(ctx) {
		log.Errorf(ctx, "Invalid OAuth2 token[%v]", authorization)
		return errors.New(http.StatusForbidden, "Invalid oauth2 token")
	}

	// Audienceチェックを行う
	validAudience := false
	for _, aud := range it.Swagger.SecurityDefinitions.GoogleIdToken.Audiences {
		// 許可済みのAudienceを見つけた
		if len(aud) > 0 && aud == token.Audience {
			validAudience = true
		}
	}

	// サービスアカウントのOAuth2トークンも許可する
	if !validAudience && (token.Audience == it.ServiceAccount.GetClientId()) {
		validAudience = true
	}

	if !validAudience {
		// 許可済のaudを見つけられなかった
		// 恐らく、このトークンはこのプロジェクトのために用意されたものでは無いだろう
		log.Errorf(ctx, "User OAuth2 check error aud[%v]", token.Audience)
		return errors.New(http.StatusForbidden, "Not supported oauth2 aud :: " + token.Audience)
	}

	result.OAuth2Token = &authorization
	result.User = &gaefire.UserInfo{
	}

	// ユーザー情報を書き出す
	if len(token.Email) > 0 {
		result.User.Email = &token.Email
	}

	return nil
}

func (it *AuthenticationProxyImpl)validJsonWebToken(ctx context.Context, jwtString string, result *gaefire.AuthenticationInfo) error {
	token, _ := jwt.Parse(jwtString, nil)
	if token == nil || token.Claims == nil {
		return errors.New(http.StatusForbidden, "Token not supported format")
	}

	if iss, ok := token.Claims.(jwt.MapClaims)["iss"]; !ok {
		return errors.New(http.StatusForbidden, "Invalid token aud")
	} else {
		issuer := fmt.Sprintf("%v", iss)

		var verifier gaefire.JsonWebTokenVerifier
		var googleJwt *string
		var firebaseJwt *string
		var serviceJwt *string

		if issuer == it.Swagger.SecurityDefinitions.GoogleIdToken.Issuer {
			// Google IdTokenとして検証する
			verifier = it.ServiceAccount.NewGoogleAuthTokenVerifier(ctx, jwtString)
			for _, aud := range it.Swagger.SecurityDefinitions.GoogleIdToken.Audiences {
				verifier.AddTrustedAudience(aud)
			}
			googleJwt = &jwtString
		} else if issuer == ("https://securetoken.google.com/" + it.ServiceAccount.GetProjectId()) {
			// Firebase Userとして扱う
			verifier = it.ServiceAccount.NewFirebaseAuthTokenVerifier(ctx, jwtString)
			firebaseJwt = &jwtString
		} else if issuer == it.ServiceAccount.GetClientEmail() {
			// 自身が発行したFirebase JsonWebTokenとして扱う
			verifier = it.ServiceAccount.NewFirebaseAuthTokenVerifier(ctx, jwtString)
			verifier.AddTrustedAudience("https://identitytoolkit.googleapis.com/google.identity.identitytoolkit.v1.IdentityToolkit")
			serviceJwt = &jwtString
		} else {
			log.Errorf(ctx, "User JWT check error issuer[%v]", issuer)
			return errors.New(http.StatusForbidden, "Issuer error")
		}

		if validToken, err := verifier.Valid(); err != nil {
			log.Errorf(ctx, "User JWT check error err[%v]", err.Error())
			return err
		} else {
			result.FirebaseToken = firebaseJwt
			result.GoogleIdToken = googleJwt
			result.ServiceAccountToken = serviceJwt
			result.User = &gaefire.UserInfo{}
			if email, err := validToken.GetClaim("email"); err == nil {
				result.User.Email = swag_port.String(fmt.Sprintf("%v", email))
			}

			if user_id, err := validToken.GetClaim("user_id"); err == nil {
				result.User.Id = swag_port.String(fmt.Sprintf("%v", user_id))
			} else if uid, err := validToken.GetClaim("uid"); err == nil {
				result.User.Id = swag_port.String(fmt.Sprintf("%v", uid))
			}
		}
	}

	return nil
}

/**
 * 認証情報をチェックする
 */
func (it *AuthenticationProxyImpl)validAuthentication(ctx context.Context, r *http.Request, result *gaefire.AuthenticationInfo) error {
	authorization := r.Header.Get("Authorization")
	if len(authorization) == 0 {
		// 認証は必要ない
		return nil
	}

	if strings.Index(authorization, "Bearer ") != 0 {
		return errors.New(http.StatusForbidden, "Authorization not bearer")
	}

	authorization = authorization[len("Bearer "):]
	// OAuth2 or ID Token?
	if strings.Index(authorization, "ey") == 0 {
		// JWTを検証する
		return it.validJsonWebToken(ctx, authorization, result)
	} else {
		// OAuth2チェック
		return it.validOAuth2(ctx, authorization, result)
	}
}


/**
 * ユーザー認証を行い、必要に応じてhttpリクエストを改変する。
 */
func (it *AuthenticationProxyImpl)Authentication(ctx context.Context, r *http.Request) (*gaefire.AuthenticationInfo, error) {
	for _, key := range SecurityCheckHeaders {
		value := r.Header.Get(key)
		// セキュリティ上許されないヘッダを見つけた
		if len(value) > 0 {
			log.Errorf(ctx, "User security header key[%v] vakue[%v]", key, value)
			return nil, errors.New(http.StatusForbidden, fmt.Sprintf("SecurityValue %v", value))
		}
	}

	auth := &gaefire.AuthenticationInfo{}
	// ApiKeyをチェック
	if err := it.validApiKey(ctx, r, auth); err != nil {
		// API Keyが妥当ではない
		return nil, err
	}

	// 認証を行う
	if err := it.validAuthentication(ctx, r, auth); err != nil {
		return nil, err
	}

	if auth.User != nil {
		// ユーザー情報が設定されているなら、カスタムヘッダを設定する
		buf, _ := json.Marshal(*auth.User)
		r.Header.Add(gaefire.HttpXHeaderUserInfo, base64.StdEncoding.EncodeToString(buf))
	}

	return auth, nil
}