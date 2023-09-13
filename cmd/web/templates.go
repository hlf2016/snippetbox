package main

import (
	"github.com/hlf2016/snippetbox/internal/models"
	"html/template"
	"path/filepath"
)

// 定义 templateData 类型，作为我们要传递给 HTML 模板的任何动态数据的存储结构。目前，它只包含一个字段，但随着构建的进行，我们将添加更多的字段
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	//	使用 filepath.Glob() 函数获取与"./ui/html/pages/*.tmpl "模式匹配的所有文件路径的片段。
	//	这将基本上为我们提供应用程序 "页面 "模板的所有文件路径片段，例如[ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	// 逐个循环页面文件路径
	for _, page := range pages {
		// 从完整文件路径中提取文件名（如 "home.tmpl"）并将其赋值给名称变量
		name := filepath.Base(page)
		// 创建一个片段，其中包含基本模板、任何部分模板和页面的文件路径。
		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/pages/partials/nav.tmpl",
			page,
		}
		// 将文件解析为模板集
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		// 使用页面名称（如 "home.tmpl"）作为关键字，将模板集添加到map中。
		cache[name] = ts
	}
	return cache, nil
}
