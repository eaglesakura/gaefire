// +build !gaetest

package gaefire

import (
	"net/http"
	"google.golang.org/appengine"
	"github.com/eaglesakura/gaefire"
)

/**
 * 通常のRequest用のContextを生成する
 */
func NewContext(request *http.Request) gaefire.Context {
	result := &ContextImpl{
		ctx:appengine.NewContext(request),
	};

	return result;
}
