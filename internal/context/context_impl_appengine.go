// +build !gaetest

package internal

import (
	"net/http"
	"google.golang.org/appengine"
	"github.com/eaglesakura/gaefire/context"
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
