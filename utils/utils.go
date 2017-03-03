package gaefire

import (
	"os"
	"github.com/eaglesakura/gaefire/assets"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"crypto/md5"
	"encoding/hex"
)

const (
	EnvWorkspace = "WORKSPACE"
)

func UnmarshalJson(resp *http.Response, result interface{}) error {
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, result)
}

// 環境変数を取得する
func GetEnv(key string, def string) string {
	value, look := os.LookupEnv(key);
	if !look {
		return def;
	} else {
		return value;
	}
}

/**
 * 文字列をMD5に変換する
 */
func GenMD5(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

func LoadFileOrNil(assets gaefire.AssetManager, path string) []byte {
	buf, _ := assets.LoadFile(path)
	return buf
}
