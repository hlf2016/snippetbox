package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hlf2016/snippetbox/internal/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	infoLogger     *log.Logger
	errorLogger    *log.Logger
	fileInfoLogger *log.Logger
	cfg            config
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	// 存储在 session 中的用于判断用户是否已经登录的key
	authId string
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

	// 生成 页面缓存 注入 application 中 方便各处使用
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLogger.Fatal(err)
	}

	// 初始化decoder实例
	formDecoder := form.NewDecoder()

	// 初始化 session manager
	sessionManager := scs.New()
	// 可以更改会话 cookie，使用 SameSite=Strict 设置，而不是（默认设置）SameSite=Lax
	// 但需要注意的是，使用 SameSite=Strict 会阻止用户浏览器在所有跨站使用中发送会话 cookie，包括使用 GET 和 HEAD 等 HTTP 方法的安全请求。
	// 虽然这听起来更安全（确实如此！），但缺点是当用户从其他网站点击链接到您的应用程序时，不会发送会话 cookie。反过来，这意味着￼￼您的应用程序最初会将用户视为 "未登录"，即使他们有一个包含其 "authenticatedUserID "值的活动会话。
	// sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	// 确保在会话 cookie 上设置 Secure 属性。设置该属性意味着用户的网络浏览器只有在使用 HTTPS 连接时才会发送 cookie（而不会通过不安全的 HTTP 连接发送）。
	sessionManager.Cookie.Secure = true
	app := &application{
		infoLogger:     infoLogger,
		errorLogger:    errorLogger,
		fileInfoLogger: fileInfoLogger,
		cfg:            cfg,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		authId:         "authenticatedUserID",
	}

	infoLogger.Printf("Starting server on %s", cfg.addr)
	fileInfoLogger.Printf("Starting server on %s", cfg.addr)

	// 初始化一个 tls.Config 结构，用于保存我们希望服务器使用的非默认 TLS 设置。在本例中，我们唯一要更改的是曲线优选值，以便只使用具有汇编实现的椭圆曲线。
	// 基本上，使用 tls.Config 设置受支持密码套件的自定义列表只会影响 TLS 1.0-1.2 连接 TLS 1.3 则与之无关 被普遍认为是安全的
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// 自定义 http Server 错误日志输出器
	srv := &http.Server{
		Addr: cfg.addr,
		// 设置 错误日志输出为 自定义格式
		ErrorLog:  errorLogger,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// 限制 连接闲置 1 分钟后 自动断开
		IdleTimeout: time.Minute,
		// 如果在请求被接受 5 秒后仍在读取请求头或正文，Go 就会关闭底层连接。由于这是对连接的 "硬 "关闭，用户不会收到任何 HTTP(S) 响应。
		// 降低慢客户端攻击的风险
		// 如果设置了 ReadTimeout 但没有设置 IdleTimeout，那么 IdleTimeout 将默认使用与 ReadTimeout 相同的设置。
		// 例如，如果将 ReadTimeout 设置为 3 秒，那么就会产生一个副作用，即所有保持连接也会在 3 秒未活动后关闭。一般来说，我的建议是避免任何歧义，始终为服务器设置明确的 IdleTimeout 值。
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// err = srv.ListenAndServe() // 改用 https
	// 使用 ListenAndServeTLS() 方法启动 HTTPS 服务器。我们将 TLS 证书的路径和相应的私钥作为两个参数传递进去。
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLogger.Fatal(err)
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
