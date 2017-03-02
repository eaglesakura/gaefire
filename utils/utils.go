package gaefire

import (
	"fmt"
	"os"
	"errors"
	"appengine_internal/gopkg.in/yaml.v2"
	assets "github.com/eaglesakura/gaefire/assets"
)

const (
	EnvWorkspace = "WORKSPACE"
	GooglePublicKeystoreAccount = "securetoken@system.gserviceaccount.com"
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

func UnmarshalYaml(asset assets.AssetManager, path string, result interface{}) error {
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