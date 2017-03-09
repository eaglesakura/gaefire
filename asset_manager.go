package gaefire

/**
 * GAE/Go Appに含まれるアセットを管理する
 *
 * 読み込み時点で環境変数 `WORKSPACE` がルートになるようにカレントディレクトリが変更される。
 */
type AssetManager interface {
	LoadFile(path string) ([]byte, error)
}
