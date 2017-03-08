package gaefire

/**
 * X-Endpoint-API-UserInfoにbase64エンコードされて書き込まれる認証情報
 */
type UserInfo struct {
	Issuer *string        `json:"issuer"`
	Id     *string        `json:"id"`
	Email  *string        `json:"email"`
}

/**
 * 認証が正しく行われた場合の結果情報
 */
type AuthenticationInfo struct {
	/**
	 * 妥当なAPI Keyが使われた場合、!=nilとなる。
	 * nil以外の場合、len()は必ず1以上となる。
	 */
	ApiKey        *string

	/**
	 * Firebase Tokenが使用された場合にsetされる
	 */
	FirebaseToken *string

	/**
	 * Google Tokenが使用された場合にsetされる
	 */
	GoogleIdToken *string

	/**
	 * OAuth2 Tokenが使用された場合にsetされる
	 */
	OAuth2Token   *string

	/**
	 * 認証されたユーザー情報
	 * 認証されなければとなる。
	 */
	User          *UserInfo
}
