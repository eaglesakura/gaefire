package gaefire

/*
X-Endpoint-API-UserInfoにbase64エンコードされて書き込まれる認証情報
Google Cloud Endpoints 2.0互換のModelとなる。
*/
type UserInfo struct {
	Issuer *string `json:"issuer"`
	Id     *string `json:"id"`
	Email  *string `json:"email"`
}

/*
認証が正しく行われた場合の結果情報
*/
type AuthenticationInfo struct {
	/*
		"iss"に相当する要素が含まれている場合にsetされる
	 	nil以外の場合、len()は必ず1以上となる。
	*/
	Issuer *string

	/*
		"aud","audience"に相当する要素が含まれている場合にsetされる
		nil以外の場合、len()は必ず1以上となる。
	*/
	Audience *string

	/*
		API Keyが使われた場合にsetされる
		nil以外の場合、len()は必ず1以上となる。
	*/
	ApiKey *string

	/*
		Firebaseが署名したTokenが使用された場合にsetされる
		nil以外の場合、len()は必ず1以上となる。
	*/
	FirebaseToken *string

	/*
		Firebase Service Accountが署名したTokenが使用された場合にsetされる
		nil以外の場合、len()は必ず1以上となる。
	*/
	ServiceAccountToken *string

	/*
		Google Tokenが使用された場合にsetされる
		nil以外の場合、len()は必ず1以上となる。
	*/
	GoogleIdToken *string

	/*
		OAuth2 Tokenが使用された場合にsetされる
		nil以外の場合、len()は必ず1以上となる。
	*/
	OAuth2Token *string

	/*
		認証されたユーザー情報
		認証されなければnil(匿名)となる。
	*/
	User *UserInfo
}
