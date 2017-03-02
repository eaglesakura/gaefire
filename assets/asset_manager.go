package gaefire

import (
	"io/ioutil"
)

/**
 * GAE/Go Appに含まれるアセットを管理する
 */
type AssetManager interface {
	LoadFile(path string) ([]byte, error)
}

type _AssetManagerImpl struct {
}

func (it *_AssetManagerImpl)LoadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (it *_AssetManagerImpl)LoadFileOrNil(path string) ([]byte) {
	buf, err := it.LoadFile(path)
	if err != nil {
		return nil
	} else {
		return buf
	}
}
