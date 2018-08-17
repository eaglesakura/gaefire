package gaefire

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/eaglesakura/gaefire"
)

type FirebaseWebApplicationImpl struct {
	rawWebApp WebApplicationModel
}

func NewWebApplication(jsonBuf []byte) gaefire.WebApplication {

	result := &FirebaseWebApplicationImpl{}

	if err := json.Unmarshal(jsonBuf, &result.rawWebApp); err != nil {
		panic(err)
		return nil
	}

	if len(result.rawWebApp.Web.ClientId) == 0 || len(result.rawWebApp.Web.ClientSecret) == 0 {
		panic(errors.New("Client data not found."))
		return nil
	}

	return result
}

/*
 * GCP Project IDを取得する
 */
func (it *FirebaseWebApplicationImpl) GetProjectId() string {
	return it.rawWebApp.Web.ProjectId
}

/*
 * OAuth2 Client IDを取得する
 */
func (it *FirebaseWebApplicationImpl) GetClientId() string {
	return it.rawWebApp.Web.ClientId
}

/*
 * OAuth2 Client Secretを取得する
 */
func (it *FirebaseWebApplicationImpl) GetClientSecret() string {
	return it.rawWebApp.Web.ClientSecret
}

/*
 * 一般ユーザーがOAuth2認証を行なうためのトークンを取得する
 *
 * 取得したOAuth2トークンはMemcacheに登録される。
 * ユーザーが明示的に権限を取り消している場合、エラーを返却する。
 */
func (it *FirebaseWebApplicationImpl) NewUserAccountToken(ctx context.Context, accessCode string) (gaefire.OAuth2Token, error) {
	tokenGen := &OAuth2RefreshRequest{
		ctx:            ctx,
		webApplication: it,
		accessCode:     accessCode,
	}

	return tokenGen.GetToken()
}

/*
 * 一般ユーザーがOAuth2認証を行なうためのトークンを取得する
 *
 * OAuth2トークンはMemcacheにキャッシュされ、再取得は最低限となるよう実装される。
 * キャッシュが存在しないかExpireされている場合、自動的にRefreshTokenを用いてリフレッシュした結果を返す。
 * ユーザーが明示的に権限を取り消している場合、エラーを返却する。
 */
func (it *FirebaseWebApplicationImpl) GetUserAccountToken(ctx context.Context, refreshToken string) (gaefire.OAuth2Token, error) {
	tokenGen := &OAuth2RefreshRequest{
		ctx:            ctx,
		webApplication: it,
		refreshToken:   refreshToken,
	}
	return tokenGen.GetToken()
}
