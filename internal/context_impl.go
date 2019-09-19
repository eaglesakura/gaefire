package gaefire

import (
	"context"
	"fmt"
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
func (it *ContextImpl) LogError(format string, args ...interface{}) {
	logError(fmt.Sprintf(format, args...))
}

/**
 * デバッグログ出力を行なう
 */
func (it *ContextImpl) LogDebug(format string, args ...interface{}) {
	logDebug(fmt.Sprintf(format, args...))
}

/**
 * インフォログ出力を行なう
 */
func (it *ContextImpl) LogInfo(format string, args ...interface{}) {
	logInfo(fmt.Sprintf(format, args...))
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
