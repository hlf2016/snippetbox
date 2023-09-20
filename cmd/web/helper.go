package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
	"net/http"
	"runtime/debug"
	"time"
)

// serverError 辅助程序会将错误信息和堆栈跟踪写入 errorLog，然后向用户发送通用的 500 内部服务器错误响应。
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLogger.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError 辅助程序会向用户发送特定的状态代码和相应的描述。
// 在本书的后面部分，当用户发送的请求出现问题时，我们将使用它来发送类似 400 "Bad Request（错误请求）"的响应。
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// 为了保持一致性，我们还将实现一个 notFound 辅助器。这只是客户端错误（clientError）的一个方便封装，它会向用户发送 404 Not Found 响应。
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s doesn't exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	// 将模板写入缓冲区，而不是直接写入 http.ResponseWriter。如果出现错误，则调用我们的 serverError() 辅助程序，然后返回。
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// 如果模板在写入缓冲区时没有出现任何错误，我们就可以继续将 HTTP 状态代码写入 http.ResponseWriter 中。
	w.WriteHeader(status)
	// 将缓冲区的内容写入 http.ResponseWriter。注意：这又是一次我们将 http.ResponseWriter 传递给接收 io.Writer 的函数的情况。
	_, err = buf.WriteTo(w)
	if err != nil {
		app.serverError(w, err)
	}
	// 直接渲染
	//err := ts.ExecuteTemplate(w, "base", data)
	//if err != nil {
	//	app.serverError(w, err)
	//}
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		// 所有页面数据上都加入 CSRFToken 便于每个页面上使用
		CSRFToken: nosurf.Token(r),
	}
}

// 创建一个新的 decodePostForm() 辅助方法。这里的第二个参数 dst 是我们要将表单数据解码到的目标目的地。
func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	// 在我们的 decoder 实例上调用 Decode() ，将目标目的地作为第一个参数传入
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// 如果我们尝试使用无效的目标目的地 如 nil，Decode() 方法将返回一个类型为 form.InvalidDecoderError 的错误。
		// 我们使用 errors.As() 来检查这种情况，并引发恐慌，而不是返回错误
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
