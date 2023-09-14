package main

import (
	"fmt"
	"net/http"
	"strings"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLogger.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 创建一个defer函数（在出现panic时，Go 会释放堆栈，并始终运行该函数）。
		defer func() {
			// 使用内置的recover方法检查是否发生了panic 如果发生了
			if err := recover(); err != nil {
				// 在响应中设置 Connection: Close 标头作为触发器，使 Go 的 HTTP 服务器在发送响应后自动关闭当前连接。它还通知用户连接将被关闭。
				// 注意：如果使用的协议是 HTTP/2，Go 会自动从响应中剥离 Connection: Close 标头（因此它没有格式错误）并发送一个 GOAWAY 帧。
				w.Header().Set("Connection", "close")
				// 内置的recover()函数返回的值的类型为any，其底层类型可以是字符串、错误或其他类型——无论传递给panic()的参数是什么。在我们的例子中，它是字符串“oops! something went wrong”。
				// 在上面的代码中，我们通过使用 fmt.Errorf() 函数创建一个包含 any 值的默认文本表示形式的新错误对象，将其标准化为错误，然后将此错误传递给 app.server Error() 帮助程序方法。
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
