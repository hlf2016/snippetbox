package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hlf2016/snippetbox/internal/models"
)

type application struct {
	infoLogger     *log.Logger
	errorLogger    *log.Logger
	fileInfoLogger *log.Logger
	cfg            config
	snippets       *models.SnippetModel
}

// 聚合 config 设置 然后使用 flag.StringVar 读取环境变量赋值
type config struct {
	addr      string
	staticDir string
	dsn       string
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")

	// DSN 中的 parseTime=true 部分是一个特定于驱动程序的参数，它指示我们的驱动程序将 SQL TIME 和 DATE 字段转换为 Go time.Time 对象。
	flag.StringVar(&cfg.dsn, "dsn", "goweb:25804769@/snippetbox?parseTime=true", "MySQL data source name")
	// 重要的是，我们使用 flag.Parse() 函数来解析命令行标志。它会读入命令行标志值并将其赋值给 addr 变量。
	// 您需要在使用 addr 变量之前调用该函数，否则它将始终包含默认值":4000"。如果在解析过程中遇到任何错误，应用程序将被终止。
	flag.Parse()

	// 使用 log.New() 创建一个日志记录器，用于写入信息消息。它需要三个参数：写入日志的目的地（os.Stdout）、
	// 信息的字符串前缀（INFO，后跟一个制表符），以及指示包含哪些附加信息（本地日期和时间）的标志。请注意，这些标志是用位运算符 | 连接起来的。
	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	//:如果想在日志输出中包含完整的文件路径，而不仅仅是文件名，可以在创建自定义日志记录器时使用 log.Llongfile 标志，而不是 log.Lshortfile。
	// 还可以通过添加 log.LUTC 标志，强制日志记录器使用 UTC 日期（而不是本地日期）。
	errorLogger := log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	// 启用文件日志
	f, err := os.OpenFile("./log/info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		errorLogger.Fatal(err)
	}
	defer f.Close()
	fileInfoLogger := log.New(f, "INFO\t", log.Ldate|log.Ltime)

	db, err := openDB(cfg.dsn)
	if err != nil {
		errorLogger.Fatal(err.Error())
	}
	defer db.Close()

	app := &application{
		infoLogger:     infoLogger,
		errorLogger:    errorLogger,
		fileInfoLogger: fileInfoLogger,
		cfg:            cfg,
		snippets:       &models.SnippetModel{DB: db},
	}

	infoLogger.Printf("Starting server on %s", cfg.addr)
	fileInfoLogger.Printf("Starting server on %s", cfg.addr)

	// 自定义 http Server 错误日志输出器
	srv := &http.Server{
		Addr: cfg.addr,
		// 设置 错误日志输出为 自定义格式
		ErrorLog: errorLogger,
		Handler:  app.routes(),
	}

	err = srv.ListenAndServe()
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

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
