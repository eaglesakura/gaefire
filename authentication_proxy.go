package gaefire

import (
	"net/http"
	"golang.org/x/net/context"
)

const (
	HttpXHeaderUserInfo string = "X-Endpoint-API-UserInfo"
)

/**
 * httpリクエストに対し、ユーザー認証のチェックを行う。
 *
 * 可能な限りGoogle Cloud Endpointsと同様のチェックを行う。
 * 外部から不正なヘッダを受信した場合、エラー扱いとする。
 */
type AuthenticationProxy interface {
	/**
	 * ユーザー認証を行い、必要に応じてhttpリクエストを改変する。
	 * 認証中にエラーが発生した場合はerrorを返却する。
	 * Authorizationヘッダがない場合等、認証自体を行わなかった場合はエラーを返却しない。
	 * 認証サポート用のProxyを生成する。
	 * Google Cloud Endpoints 2.0と同じくユーザーやAPI Keの検証を行うが、APIごとのチェック（認証が必要か、等）は自身で行う必要がある。
	 * Google Cloud Endpoints仕様に従い、下記のAuthorizationヘッダをサポートする
	 * Proxyを通過した時点で、X-Endpoints-UserInfoヘッダが必要に応じて付与される。
	 *
	 * in:[query|header] api_key={your.gcp.api_key}
	 * "Authorization: Bearer {your.oauth2.token}"
	 * "Authorization: Bearer {google.json.web.token}"
	 * "Authorization: Bearer {firebase.json.web.token}"
	 *
	 * このAPIを使用するためには、事前にswagger.json(openapi.json)を `gcloud service-management deploy path/to/swagger.json` でデプロイしておく必要がある。
	 */
	Authentication(ctx context.Context, r *http.Request) (*AuthenticationInfo, error)
}
