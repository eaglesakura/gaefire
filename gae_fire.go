package gaefire

import (
	fire_assets "github.com/eaglesakura/gaefire/assets"
	fire_auth"github.com/eaglesakura/gaefire/auth"
)

type GaeFire interface {
	/**
	 * 初期化を行なう
	 */
	Initialize() error

	/**
	 * AssetManagerを生成する
	 */
	NewAssetManager() fire_assets.AssetManager

	/**
	 * サービスアカウントを生成する
	 */
	NewServiceAccount(jsonBuf []byte) fire_auth.FirebaseServiceAccount

	/**
	 * ユーザーOAuth2認証に利用するWebアプリケーションを生成する
	 */
	NewWebApplication(jsonBuf []byte) fire_auth.FirebaseWebApplication
}

