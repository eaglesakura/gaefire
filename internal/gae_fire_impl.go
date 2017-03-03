package internal

import (
	"os"
	"github.com/eaglesakura/gaefire/utils"
	util "github.com/eaglesakura/gaefire/utils"
	fire_auth "github.com/eaglesakura/gaefire/auth"
	fire_auth_internal "github.com/eaglesakura/gaefire/internal/auth"

	fire_assets "github.com/eaglesakura/gaefire/assets"
	fire_assets_internal "github.com/eaglesakura/gaefire/internal/assets"
)

type GaeFireImpl struct {
}

func (it *GaeFireImpl)Initialize() error {
	// 必要に応じてWorkspaceを切り替える
	// 主にUnitTestを行う場合に使う
	{
		workspace := gaefire.GetEnv(util.EnvWorkspace, "");
		if workspace != "" {
			os.Chdir(workspace);
		}
	}

	return nil;
}

/**
 * AssetManagerを生成する
 */
func (it *GaeFireImpl)NewAssetManager() fire_assets.AssetManager {
	return &fire_assets_internal.AssetManagerImpl{}
}

func (it *GaeFireImpl)NewServiceAccount(jsonBuf []byte) fire_auth.FirebaseServiceAccount {
	return fire_auth_internal.NewFirebaseServiceAccount(jsonBuf)
}

/**
 * ユーザーOAuth2認証に利用するWebアプリケーションを生成する
 */

func (it *GaeFireImpl)NewWebApplication(jsonBuf []byte) fire_auth.FirebaseWebApplication {
	return fire_auth_internal.NewFirebaseWebApplication(jsonBuf)
}