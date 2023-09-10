package main

import (
	"flag"
	"log"
	"net/http"
	"os"
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

	// 使用 log.New() 创建一个日志记录器，用于写入信息消息。它需要三个参数：写入日志的目的地（os.Stdout）、
	// 信息的字符串前缀（INFO，后跟一个制表符），以及指示包含哪些附加信息（本地日期和时间）的标志。请注意，这些标志是用位运算符 | 连接起来的。
	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	//:如果想在日志输出中包含完整的文件路径，而不仅仅是文件名，可以在创建自定义日志记录器时使用 log.Llongfile 标志，而不是 log.Lshortfile。
	// 还可以通过添加 log.LUTC 标志，强制日志记录器使用 UTC 日期（而不是本地日期）。
	errorLogger := log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

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
	infoLogger.Printf("Starting server on %s", cfg.addr)

	// 自定义 http Server 错误日志输出器
	srv := &http.Server{
		Addr: cfg.addr,
		// 设置 错误日志输出为 自定义格式
		ErrorLog: errorLogger,
		Handler:  mux,
	}

	err := srv.ListenAndServe()
	errorLogger.Fatal(err)
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
