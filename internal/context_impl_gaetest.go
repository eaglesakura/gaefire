// +build gaetest

package gaefire

import (
	"net/http"
	"github.com/eaglesakura/gaefire"
	"google.golang.org/appengine/aetest"
)

/**
 * UnitTest用のContextを生成する
 */
func NewContext(request *http.Request) gaefire.Context {
	ctx, delFunc, err := aetest.NewContext();
	if err != nil {
		panic(err);
	}

	result := &ContextImpl{
		ctx:ctx,
		closeFunc:delFunc,
	};

	return result;
}
