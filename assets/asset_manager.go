package gaefire

/**
 * GAE/Go Appに含まれるアセットを管理する
 */
type AssetManager interface {
	LoadFile(path string) ([]byte, error)
}
