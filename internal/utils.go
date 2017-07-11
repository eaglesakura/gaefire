package gaefire

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"net/http"
	"os"
	"time"
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
	value, look := os.LookupEnv(key)
	if !look {
		return def
	} else {
		return value
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

func newHttpClient(ctx context.Context) *http.Client {
	// タイムアウトを30秒に延長
	deadline, _ := context.WithTimeout(ctx, 30*time.Second)

	result := &http.Client{
		Transport: &urlfetch.Transport{
			Context: deadline,
		},
	}
	result.Timeout = 30 * time.Second
	return result
}
