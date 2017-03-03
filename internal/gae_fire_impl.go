package internal

import (
	"os"
	"github.com/eaglesakura/gaefire/utils"
	util "github.com/eaglesakura/gaefire/utils"
	fire_auth "github.com/eaglesakura/gaefire/auth"
	"github.com/eaglesakura/gaefire/internal/auth"
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

func (it *GaeFireImpl)NewServiceAccount(jsonBuf []byte) fire_auth.FirebaseServiceAccount {
	return internal.NewFirebaseServiceAccount(jsonBuf)
}

/**
 * ユーザーOAuth2認証に利用するWebアプリケーションを生成する
 */

func (it *GaeFireImpl)NewWebApplication(jsonBuf []byte) fire_auth.FirebaseWebApplication {
	return internal.NewFirebaseWebApplication(jsonBuf)
}