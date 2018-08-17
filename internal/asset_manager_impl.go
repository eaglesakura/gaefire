package gaefire

import (
	"github.com/eaglesakura/gaefire"
	"io/ioutil"
)

type AssetManagerImpl struct {
}

func (it *AssetManagerImpl) LoadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

/*
 * AssetManagerを生成する
 */
func NewAssetManager() gaefire.AssetManager {
	return &AssetManagerImpl{}
}
