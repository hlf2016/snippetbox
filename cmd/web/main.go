package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
)

// 聚合 config 设置 然后使用 flag.StringVar 读取环境变量赋值
type config struct {
	addr      string
	staticDir string
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	// 重要的是，我们使用 flag.Parse() 函数来解析命令行标志。它会读入命令行标志值并将其赋值给 addr 变量。
	// 您需要在使用 addr 变量之前调用该函数，否则它将始终包含默认值":4000"。如果在解析过程中遇到任何错误，应用程序将被终止。
	flag.Parse()
	mux := http.NewServeMux()
	// 当该处理程序接收到一个请求时，它会删除 URL 路径中的前导斜线，然后在 ./ui/static 目录中搜索相应的文件发送给用户。
	// 因此，为了使该处理程序正常工作，我们必须在将 URL 路径传递给 http.FileServer 之前，去掉 URL 路径中以"/static "开头的斜线。
	// 否则，它将寻找一个不存在的文件，用户将收到未找到的 404 页面响应。幸运的是，Go 包含了一个 http.StripPrefix() 助手，专门用于完成这项任务。
	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	// log.Print(cfg)
	mux.Handle("/static/", http.StripPrefix("/static", neuter(fileServer)))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	log.Printf("Starting server on %s", cfg.addr)
	err := http.ListenAndServe(cfg.addr, mux)
	log.Fatal(err)
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
