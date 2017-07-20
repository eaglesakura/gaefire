package gaefire

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

type ContextImpl struct {
	/**
	 * Appengine Context
	 */
	ctx context.Context

	/**
	 * destroy func
	 */
	closeFunc func()
}

/**
 * AppEngineのハンドラに紐付いたコンテキストを取得する
 */
func (it *ContextImpl) GetAppengineContext() context.Context {
	return it.ctx
}

/**
 * エラーログ出力を行なう
 */
func (it *ContextImpl) LogError(fmt string, args ...interface{}) {
	log.Errorf(it.ctx, fmt, args...)
}

/**
 * デバッグログ出力を行なう
 */
func (it *ContextImpl) LogDebug(fmt string, args ...interface{}) {
	log.Debugf(it.ctx, fmt, args...)
}

/**
 * インフォログ出力を行なう
 */
func (it *ContextImpl) LogInfo(fmt string, args ...interface{}) {
	log.Infof(it.ctx, fmt, args...)
}

/**
 * 関連するリソースの解放を行なう
 */
func (it *ContextImpl) Close() {
	if it.closeFunc != nil {
		it.closeFunc()
	}
	it.closeFunc = nil
	it.ctx = nil
}
