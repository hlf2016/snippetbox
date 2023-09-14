package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	// 当该处理程序接收到一个请求时，它会删除 URL 路径中的前导斜线，然后在 ./ui/static 目录中搜索相应的文件发送给用户。
	// 因此，为了使该处理程序正常工作，我们必须在将 URL 路径传递给 http.FileServer 之前，去掉 URL 路径中以"/static "开头的斜线。
	// 否则，它将寻找一个不存在的文件，用户将收到未找到的 404 页面响应。幸运的是，Go 包含了一个 http.StripPrefix() 助手，专门用于完成这项任务。
	fileServer := http.FileServer(http.Dir(app.cfg.staticDir))
	// log.Print(cfg)
	mux.Handle("/static/", http.StripPrefix("/static", neuter(fileServer)))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// 将 servemux 作为 "next "参数传递给 secureHeaders 中间件。
	// 因为 secureHeaders 只是一个函数，而函数返回的是 http.Handler，所以我们不需要做其他任何事情。
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
