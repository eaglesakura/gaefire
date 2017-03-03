package gaefire

import (
	"fmt"
	"os"
	"errors"
	"appengine_internal/gopkg.in/yaml.v2"
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

// find string by map
// NotFound => ""
func FindStringValue(values *map[string]interface{}, key string) string {
	result, ok := (*values)[key]
	if !ok {
		return ""
	} else {
		return fmt.Sprintf("%v", result)
	}
}

func UnmarshalJson(resp *http.Response, result interface{}) error {
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, result)
}

func UnmarshalYaml(asset gaefire.AssetManager, path string, result interface{}) error {
	buf, _ := asset.LoadFile(path)
	if buf == nil {
		return errors.New(fmt.Sprintf("Asset[%s] load failed", path))
	}
	return yaml.Unmarshal(buf, result)
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
