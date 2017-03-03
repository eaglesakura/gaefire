package gaefire

import (
	"github.com/eaglesakura/gaefire/internal"
	fire_context "github.com/eaglesakura/gaefire/context"
	fire_context_internal "github.com/eaglesakura/gaefire/internal/context"
	"net/http"
)

func NewGaeFire() GaeFire {
	result := &internal.GaeFireImpl{}
	result.Initialize()
	return result
}


/**
 * ハンドリング用のコンテキストを生成する
 */
func NewContext(request *http.Request) fire_context.Context {
	return fire_context_internal.NewContext(request)
}