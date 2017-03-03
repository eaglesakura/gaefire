package internal

import "io/ioutil"

type AssetManagerImpl struct {
}

func (it *AssetManagerImpl)LoadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

