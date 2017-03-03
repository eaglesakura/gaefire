package gaefire

type GaeFire interface {
	/**
	 * AssetManagerを生成する
	 */
	NewAssetManager() AssetManager

	/**
	 * サービスアカウントを生成する
	 */
	NewServiceAccount(jsonBuf []byte) FirebaseServiceAccount

	/**
	 * ユーザーOAuth2認証に利用するWebアプリケーションを生成する
	 */
	NewWebApplication(jsonBuf []byte) FirebaseWebApplication
}

