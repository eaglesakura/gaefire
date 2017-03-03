package gaefire

import (
	"os"
	"github.com/eaglesakura/gaefire"
)

type GaeFireImpl struct {
}

func (it *GaeFireImpl)Initialize() error {
	// 必要に応じてWorkspaceを切り替える
	// 主にUnitTestを行う場合に使う
	{
		workspace := GetEnv(EnvWorkspace, "");
		if workspace != "" {
			os.Chdir(workspace);
		}
	}

	return nil;
}

/**
 * AssetManagerを生成する
 */
func (it *GaeFireImpl)NewAssetManager() gaefire.AssetManager {
	return &AssetManagerImpl{}
}

func (it *GaeFireImpl)NewServiceAccount(jsonBuf []byte) gaefire.FirebaseServiceAccount {
	return NewFirebaseServiceAccount(jsonBuf)
}

/**
 * ユーザーOAuth2認証に利用するWebアプリケーションを生成する
 */
func (it *GaeFireImpl)NewWebApplication(jsonBuf []byte) gaefire.FirebaseWebApplication {
	return NewFirebaseWebApplication(jsonBuf)
}