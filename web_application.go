package gaefire

import "context"

/**
 * Firebase経由で生成されたWebアプリケーションの管理を行なう
 *
 * 下記の機能を提供する:
 * * ユーザーOAuth2トークン管理
 */
type WebApplication interface {
	/**
	 * GCP Project IDを取得する
	 */
	GetProjectId() string

	/**
	 * OAuth2 Client IDを取得する
	 */
	GetClientId() string

	/**
	 * OAuth2 Client Secretを取得する
	 */
	GetClientSecret() string

	/**
	 * 一般ユーザーがOAuth2認証を行なうためのトークンを取得する
	 *
	 * 取得したOAuth2トークンはMemcacheに登録される。
	 * ユーザーが明示的に権限を取り消している場合、エラーを返却する。
	 */
	NewUserAccountToken(ctx context.Context, accessCode string) (OAuth2Token, error)

	/**
	 * 一般ユーザーがOAuth2認証を行なうためのトークンを取得する
	 *
	 * OAuth2トークンはMemcacheにキャッシュされ、再取得は最低限となるよう実装される。
	 * キャッシュが存在しないかExpireされている場合、自動的にRefreshTokenを用いてリフレッシュした結果を返す。
	 * ユーザーが明示的に権限を取り消している場合、エラーを返却する。
	 */
	GetUserAccountToken(ctx context.Context, refreshToken string) (OAuth2Token, error)
}
