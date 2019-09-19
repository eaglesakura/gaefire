package gaefire

import (
	"github.com/eaglesakura/gaefire"
	"net/http"
)

/**
 * 通常のRequest用のContextを生成する
 */
func NewContext(request *http.Request) gaefire.Context {
	result := &ContextImpl{
		ctx: request.Context(),
	}

	return result
}
