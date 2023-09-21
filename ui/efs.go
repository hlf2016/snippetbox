package ui

import "embed"

// Files 该注释指令指示 Go 将 ui/html 和 ui/static 文件夹中的文件存储到由全局变量 Files 引用的 embed.FS 嵌入式文件系统中。
//
//go:embed "html" "static"
var Files embed.FS
