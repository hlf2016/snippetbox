package main

import (
	"github.com/hlf2016/snippetbox/internal/models"
	"github.com/hlf2016/snippetbox/ui"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"
)

// 定义 templateData 类型，作为我们要传递给 HTML 模板的任何动态数据的存储结构。目前，它只包含一个字段，但随着构建的进行，我们将添加更多的字段
type templateData struct {
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	CurrentYear     int
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	// 若传参值为零时刻 即 January 1, year 1, 00:00:00 UTC 返回空字符串
	if t.IsZero() {
		return ""
	}
	// 无论传参值是何种时区 先转成 UTC 时区 再进行转换 确保只输出 UTC 时区时间
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// 初始化 template.FuncMap 对象并将其存储在全局变量中。它本质上是一个字符串键值映射，在自定义模板函数名称和函数本身之间起查找作用。
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	//	使用 filepath.Glob() 函数获取与"./ui/html/pages/*.tmpl "模式匹配的所有文件路径的片段。
	//	这将基本上为我们提供应用程序 "页面 "模板的所有文件路径片段，例如[ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	// pages, err := filepath.Glob("./ui/html/pages/*.tmpl") // 替换成下面的 embedded filesystems

	// 使用 fs.Glob() 获取 ui.Files 内嵌文件系统中与 "html/pages/*.tmpl "模式匹配的所有文件路径的片段。就像之前一样，这基本上为我们提供了应用程序所有 "page"模板的片段。
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	// 逐个循环页面文件路径
	for _, page := range pages {
		// 从完整文件路径中提取文件名（如 "home.tmpl"）并将其赋值给名称变量
		name := filepath.Base(page)

		// 创建一个切片，其中包含我们要解析的模板的文件路径模式。
		patterns := []string{
			"html/base.tmpl",
			"html/pages/partials/*.tmpl",
			page,
		}
		// 将基本模板文件解析为模板集
		// 在调用 ParseFiles() 方法之前，template.FuncMap 必须与模板集一起注册。
		// 这意味着我们必须使用 template.New() 创建一个空模板集，使用 Funcs() 方法注册 template.FuncMap，然后按正常方法解析文件。
		//ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		//if err != nil {
		//	return nil, err
		//}
		// 调用此模板集 ts 上的 ParseGlob() 来添加任何 partials 页面。
		//ts, err = ts.ParseGlob("./ui/html/pages/partials/*.tmpl")
		//if err != nil {
		//	return nil, err
		//}
		// 在此模板集 ts 上调用 ParseFiles() 添加页面模板
		//ts, err = ts.ParseFiles(page)
		//if err != nil {
		//	return nil, err
		//}

		// 改用 embedded filesystems
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// 使用页面名称（如 "home.tmpl"）作为关键字，将模板集添加到map中。
		cache[name] = ts
	}
	return cache, nil
}
