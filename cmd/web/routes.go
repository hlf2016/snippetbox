package main

import (
	"github.com/hlf2016/snippetbox/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// 创建一个封装 notFound() 辅助函数的处理函数，然后将其指定为 404 Not Found 响应的自定义处理函数。
	// 您还可以通过设置 router.MethodNotAllowed 来为 405 Method Not Allowed 响应设置自定义处理程序。
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	// 当该处理程序接收到一个请求时，它会删除 URL 路径中的前导斜线，然后在 ./ui/static 目录中搜索相应的文件发送给用户。
	// 因此，为了使该处理程序正常工作，我们必须在将 URL 路径传递给 http.FileServer 之前，去掉 URL 路径中以"/static "开头的斜线。
	// 否则，它将寻找一个不存在的文件，用户将收到未找到的 404 页面响应。幸运的是，Go 包含了一个 http.StripPrefix() 助手，专门用于完成这项任务。
	//fileServer := http.FileServer(http.Dir(app.cfg.staticDir))
	//router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", neuter(fileServer)))

	// 改用 embeded files
	// 将 ui.Files 嵌入式文件系统转换为 http.FS 类型，使其满足 http.FileSystem 接口的要求。然后，我们将其传递给 http.FileServer() 函数，创建文件服务器处理程序。
	fileServer := http.FileServer(http.FS(ui.Files))
	// log.Print(cfg)
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)
	// 我们的静态文件包含在ui.Files嵌入式文件系统的“Static”文件夹中。因此，例如，我们的CSS样式表位于“Static/css/main.css”。这意味着我们现在不再需要从请求URL中去掉前缀--任何以静态开头的请求都可以直接传递到文件服务器，并且将提供相应的静态文件(只要它存在)。

	router.HandlerFunc(http.MethodGet, "/ping", ping)
	// 该中间件会在每次 HTTP 请求和响应时自动加载和保存会话数据。
	// 使用 "dynamic "中间件链的无保护应用路由。
	// 在所有 "dynamic"路由上使用 nosurf 中间件
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(app.about))
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// 受保护（仅通过身份验证）的应用路由，使用新的 "protected"中间件链，其中包括 requireAuthentication 中间件。
	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))
	router.Handler(http.MethodGet, "/account/view", protected.ThenFunc(app.accountView))
	router.Handler(http.MethodGet, "/account/password/update", protected.ThenFunc(app.accountPasswordUpdate))
	router.Handler(http.MethodPost, "/account/password/update", protected.ThenFunc(app.accountPasswordUpdatePost))

	// 创建一个中间件链，其中包含我们的 "标准 "中间件，该中间件将用于应用程序收到的每个请求。
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// 将 servemux 作为 "next "参数传递给 secureHeaders 中间件。
	// 因为 secureHeaders 只是一个函数，而函数返回的是 http.Handler，所以我们不需要做其他任何事情。
	return standard.Then(router)
}
