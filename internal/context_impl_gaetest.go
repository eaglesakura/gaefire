// +build gaetest

package gaefire

import (
	"net/http"
	"github.com/eaglesakura/gaefire"
	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine"
)

func testNewContext() (context.Context, func(), error) {
	inst, err := aetest.NewInstance(&aetest.Options{StronglyConsistentDatastore:true})
	if err != nil {
		return nil, nil, err
	}
	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		inst.Close()
		return nil, nil, err
	}
	ctx := appengine.NewContext(req)
	return ctx, func() {
		inst.Close()
	}, nil
}

/**
 * UnitTest用のContextを生成する
 */
func NewContext(request *http.Request) gaefire.Context {
	ctx, delFunc, err := testNewContext()
	if err != nil {
		panic(err);
	}
	result := &ContextImpl{
		ctx:ctx,
		closeFunc:delFunc,
	};

	return result;
}
