// +build gaetest

package gaefire

import "github.com/eaglesakura/gaefire/context"


/**
 * UnitTest用のContextを生成する
 */
func NewContextImpl(request *http.Request) gaefire.Context {
	ctx, delFunc, err := aetest.NewContext();
	if err != nil {
		panic(err);
	}

	result := &ContextImpl{
		ctx:ctx,
		deleteFunc:delFunc,
	};

	return result;
}
