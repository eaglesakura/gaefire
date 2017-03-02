package gaefire

import (
	fire_assets "github.com/eaglesakura/gaefire/assets"
	fire_auth"github.com/eaglesakura/gaefire/auth"
)

const (
	EnvWorkspace = "WORKSPACE"
	GooglePublicKeystoreAccount = "securetoken@system.gserviceaccount.com"
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
}

