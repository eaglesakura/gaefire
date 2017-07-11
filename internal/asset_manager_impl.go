package gaefire

import (
	"github.com/eaglesakura/gaefire"
	"io/ioutil"
	"os"
)

type AssetManagerImpl struct {
}

func (it *AssetManagerImpl) LoadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

/**
 * AssetManagerを生成する
 */
func NewAssetManager() gaefire.AssetManager {
	// 必要に応じてWorkspaceを切り替える
	// 主にUnitTestを行う場合に使う
	workspace := GetEnv(EnvWorkspace, "")
	if workspace != "" {
		os.Chdir(workspace)
	}
	return &AssetManagerImpl{}
}
