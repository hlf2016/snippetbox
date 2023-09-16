package main

import (
	"errors"
	"fmt"
	"github.com/hlf2016/snippetbox/internal/models"
	"github.com/hlf2016/snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// 因为httprouter与“/”路径完全匹配，所以我们现在可以从此处理程序中删除对r.URL.Path！=“/”的手动检查。
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

	// 使用 PopString() 方法获取 "flash "键的值。PopString() 还会从会话数据中删除键和值，因此它的作用类似于一次性获取。如果会话数据中没有匹配的键，该方法将返回空字符串。
	// 如果只想从会话数据中获取一个值（并将其保留在其中），可以使用 GetString() 方法。scs 软件包还提供了检索其他常见数据类型的方法，包括 GetInt()、GetBool()、GetBytes() 和 GetTime()。
	flash := app.sessionManager.PopString(r.Context(), "flash")

	data := app.newTemplateData(r)
	data.Snippet = snippet
	data.Flash = flash

	app.render(w, http.StatusOK, "view.tmpl", data)
	// 将片段数据写成纯文本 HTTP 响应体。
	//fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl", data)
}

// 定义一个 snippetCreateForm 结构，用于表示表单数据和表单字段的验证错误。
// 请注意，所有结构字段都是特意导出的（即以大写字母开头）。这是因为结构字段必须导出，才能在渲染模板时被 html/template 包读取
// 更新我们的 snippetCreateForm 结构，使其包含 struct 标记，告诉解码器如何将 HTML 表单值映射到不同的结构字段中。例如，在这里我们告诉解码器将 HTML 表单输入的名称为 "title "的值存储在 Title 字段中。结构标记 `form:"-"` 会告诉解码器在解码时完全忽略某个字段。
type snippetCreateForm struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`
	// 删除显式 FieldErrors 结构字段，转而嵌入 Validator 类型。嵌入 Validator 类型意味着我们的片段创建表格 "继承 "了 Validator 类型的所有字段和方法（包括 FieldErrors 字段）。
	validator.Validator `form:"-"`
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// 将请求正文大小限制为 4096 字节 如果超出大小 那么 r.ParseForm() 将会报错
	//r.Body = http.MaxBytesReader(w, r.Body, 4096)

	// 首先，我们调用 r.ParseForm()，将 POST 请求体中的任何数据添加到 r.PostForm 映射中。
	// 这也同样适用于 PUT 和 PATCH 请求。如果出现任何错误，我们将使用 app.ClientError() 助手向用户发送 400 Bad Request 响应
	// ************** 已经封装到 helper 中 decodePostForm *************

	var form snippetCreateForm
	// 调用表单解码器的 Decode() 方法，传入当前请求和指向我们的 snippetCreateForm 结构的指针。这主要是用 HTML 表单中的相关值填充我们的结构。如果出现问题，我们将向客户端返回 400 Bad Request 响应
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// 处理 多值字段 如 多选
	//for i, item := range r.PostForm["items"] {
	//	fmt.Fprintf(w, "%d: Item %s\n", i, item)
	//}

	// 由于 Validator 类型已嵌入到 snippetCreateForm 结构中，因此我们可以直接调用 CheckField() 来执行验证检查。
	// 如果检查结果不为 true，CheckField() 将把提供的键和错误信息添加到 FieldErrors 映射中。例如，在第一行中，我们 "检查 form.Title 字段是否为空"。
	// 在第二行中，我们 "检查 form.Title 字段的最大字符长度是否为 100"，以此类推。
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// 使用 Put() 方法将字符串值（"片段创建成功！"）和相应的键（"flash"）添加到会话数据中。
	// r.Context 在处理程序处理请求时，将其作为会话管理器临时存储信息的地方
	// 第二个参数（在我们的例子中是字符串 "flash"）是我们要添加到会话数据中的特定消息的密钥。随后，我们也将使用该键从会话数据中获取信息
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	//Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
