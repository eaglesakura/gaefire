package internal

import (
	"os"
	"github.com/eaglesakura/gaefire/utils"
	"encoding/json"
	"errors"
	util "github.com/eaglesakura/gaefire/utils"
	fire_auth "github.com/eaglesakura/gaefire/auth"
	"github.com/eaglesakura/gaefire/internal"
)

type GaeFireImpl struct {
	/**
	 * Google公開鍵署名
	 */
	googleKeystore *internal.PublicKeystore
}

func (it *GaeFireImpl)Initialize() error {
	it.googleKeystore = internal.NewPublicKeystore(util.GooglePublicKeystoreAccount)

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
	if jsonBuf == nil {
		return errors.New("NotFound")
	}

	raw := ServiceAccountModel{}
	if json.Unmarshal(jsonBuf, &raw) != nil {
		return errors.New("Json parse failed")
	}

	service, err := NewServiceAccount(&raw)
	if err != nil {
		return err
	}

	it.firebaseAuthService = service
	return nil
}