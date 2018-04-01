package gaefire

import "golang.org/x/net/context"

/*
	通常request時のContextと、aetestを使ったUnitTestのContextのラップを行なう
*/
type Context interface {
	/*
		AppEngineのハンドラに紐付いたコンテキストを取得する
	*/
	GetAppengineContext() context.Context

	/*
		エラーログ出力を行なう
	*/
	LogError(fmt string, args ...interface{})

	/*
		デバッグログ出力を行なう
	*/
	LogDebug(fmt string, args ...interface{})

	/*
		インフォログ出力を行なう
	*/
	LogInfo(fmt string, args ...interface{})

	/*
		関連するリソースの解放を行なう
	*/
	Close()
}
