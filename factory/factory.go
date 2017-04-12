package factory

import (
	gaefire_internal "github.com/eaglesakura/gaefire/internal"
	"net/http"
	"github.com/eaglesakura/gaefire"
)


/**
 * ハンドリング用のコンテキストを生成する
 */
func NewContext(request *http.Request) gaefire.Context {
	return gaefire_internal.NewContext(request)
}

/**
 * 内蔵アセット管理クラスを取得する
 *
 * 環境変数"WORKSPACE"が指定されている場合、そのパスをカレントディレクトリとして扱うようになる。
 */
func NewAssetManager() gaefire.AssetManager {
	return gaefire_internal.NewAssetManager()
}

/**
 * Firebase(GCP)のサービスアカウントを生成する
 *
 * Service Accountには適切なパーミッションを指定する必要がある。
 */
func NewServiceAccount(serviceAccountJson []byte) gaefire.ServiceAccount {
	return gaefire_internal.NewServiceAccount(serviceAccountJson)
}

/**
 * Firebase(GCP)のOAuth2認証用Applicationを生成する
 *
 * Google Authにてaccess_codeが発行された場合、このインターフェースを利用してOAuth2トークンを生成することができる。
 */
func NewWebApplication(webAppJson []byte) gaefire.WebApplication {
	return gaefire_internal.NewWebApplication(webAppJson)
}

/**
 * 認証サポート用のProxyを生成する。
 * 認証情報はswagger.jsonを元にパースされる。
 * パースに失敗した場合はnilが返却される
 */
func NewAuthenticationProxy(serviceAccount gaefire.ServiceAccount, swaggerJson []byte) gaefire.AuthenticationProxy {
	return NewAuthenticationProxyWithOption(serviceAccount, gaefire.AuthenticationProxyOption{}, swaggerJson)
}

/**
 * 認証サポート用のProxyを生成する。
 * 認証情報はswagger.jsonを元にパースされる。
 * パースに失敗した場合はnilが返却される
 */
func NewAuthenticationProxyWithOption(serviceAccount gaefire.ServiceAccount, options gaefire.AuthenticationProxyOption, swaggerJson []byte) gaefire.AuthenticationProxy {
	return gaefire_internal.NewAuthenticationProxy(serviceAccount, options, swaggerJson)
}