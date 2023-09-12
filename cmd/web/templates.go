package main

import "github.com/hlf2016/snippetbox/internal/models"

// 定义 templateData 类型，作为我们要传递给 HTML 模板的任何动态数据的存储结构。目前，它只包含一个字段，但随着构建的进行，我们将添加更多的字段
type templateData struct {
	Snippet *models.Snippet
}
