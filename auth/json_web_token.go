package gaefire

/**
 * Json Web Tokenの生成オブジェクト
 */
type JsonWebTokenGenerator interface {
	/**
	 * JWT BodyにClaimを追加する。
	 * この値はFirebase AuthのTokenに反映される。
	 *
	 * ex) AddClaim("email", "example@example.como");
	 */
	AddClaim(key string, value interface{}) JsonWebTokenGenerator

	/**
	 * JWT Headerに情報を付与する
	 *
	 * ex) AddClaim("kid", "your.private.key.id");
	 */
	AddHeader(key string, value interface{}) JsonWebTokenGenerator

	/**
	 * Json Web Tokenの仕様に従い、署名付きの文字列を生成する。
	 *
	 * `base64(header).base64(token).base64(sign)`
	 */
	Generate() (string, error)
}

/**
 * Json Web Tokenの検証オブジェクト
 */
type JsonWebTokenVerifier interface {
	/**
	 * "有効期限をチェックしない
	 */
	SkipExpireTime() JsonWebTokenVerifier

	/**
	 * 許可対象のAudienceを追加する
	 * デフォルトではFirebase Service AccountのIDが登録される。
	 */
	AddTrustedAudience(aud string) JsonWebTokenVerifier

	/**
	 * "aud"チェックをスキップする
	 *
	 * 他のプロジェクトに対して発行されたJWTを許可してしまうので、これを使用する場合は十分にセキュリティに注意を払う
	 */
	SkipProjectId() JsonWebTokenVerifier

	/**
	 * 全てのオプションに対し、有効であることが確認できればtrue
	 */
	Valid() (VerifiedJsonWebToken, error)
}

type VerifiedJsonWebToken interface {
	/**
	 * ユーザーID(uid)を取得する
	 * 取得できない場合、errorを返却する
	 *
	 * Firebaseの場合、1-36文字の英数である。
	 */
	GetUserId() (string, error)

	/**
	 * プロジェクトID(aud)を取得する
	 * 取得できない場合、errorを返却する
	 */
	GetProjectId() (string, error)

	/**
	 * 指定したkeyに紐付いた値を取得する。
	 * 取得できない場合、errorを返却する
	 */
	GetClaim(key string) (interface{}, error)

	/**
	 * 指定したkeyに紐付いた値をヘッダから取得する
	 * 取得できない場合、errorを返却する
	 */
	GetHeader(key string) (interface{}, error)
}