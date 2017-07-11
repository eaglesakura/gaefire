// +build !gaetest

package gaefire

import (
	"github.com/eaglesakura/gaefire"
	"google.golang.org/appengine"
	"net/http"
)

/**
 * 通常のRequest用のContextを生成する
 */
func NewContext(request *http.Request) gaefire.Context {
	result := &ContextImpl{
		ctx: appengine.NewContext(request),
	}

	return result
}
