package gaefire

import (
	"context"
	"github.com/eaglesakura/gaefire"
	"net/http"
)

/**
 * 通常のRequest用のContextを生成する
 */
func NewContext(request *http.Request) gaefire.Context {
	if request == nil {
		return &ContextImpl{
			ctx: context.Background(),
		}
	} else {
		return &ContextImpl{
			ctx: request.Context(),
		}
	}
}
