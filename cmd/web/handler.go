package main

import (
	"errors"
	"fmt"
	"github.com/hlf2016/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// 因为Httprouter与“/”路径完全匹配，所以我们现在可以从此处理程序中删除对r.URL.Path！=“/”的手动检查。
	//if r.URL.Path != "/" {
	//	app.notFound(w)
	//	return
	//}

	// 故意制造错误 查看 recoverPanic 中间件的反应
	// panic("oops! something went wrong")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)
}
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// 当 httprouter 解析请求时，任何已命名参数的值都将存储在请求上下文中。关于请求上下文，我们将在本书后面的章节中详细讨论，
	// 但现在只要知道可以使用 ParamsFromContext() 函数检索包含这些参数名称和值的片段就足够了，就像下面这样：
	params := httprouter.ParamsFromContext(r.Context())
	// 然后，我们就可以使用 ByName() 方法从片段中获取名为 "id "的参数值，并按常规进行验证
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		// 这是因为 Go 1.13 引入了通过封装错误为错误添加附加信息的功能。
		// 如果一个错误碰巧被封装，就会创建一个全新的错误值--这反过来又意味着无法使用常规的 == 平等运算符来检查原始底层错误的值
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)
	// 将片段数据写成纯文本 HTTP 响应体。
	//fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	//Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
