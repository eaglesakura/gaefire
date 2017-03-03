package gaefire

import (
	gaefire_internal "github.com/eaglesakura/gaefire/internal"
	"net/http"
	"github.com/eaglesakura/gaefire"
)

func NewGaeFire() gaefire.GaeFire {
	result := &gaefire_internal.GaeFireImpl{}
	result.Initialize()
	return result
}


/**
 * ハンドリング用のコンテキストを生成する
 */
func NewContext(request *http.Request) gaefire.Context {
	return gaefire_internal.NewContextImpl(request)
}