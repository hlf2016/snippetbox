package main

import (
	"bytes"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/hlf2016/snippetbox/internal/models/mocks"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}
	formDecoder := form.NewDecoder()
	// 和一个会话管理器实例。请注意，除了没有为会话管理器设置 "存储 "外，我们使用了与生产版相同的设置。如果不设置存储空间，SCS 软件包将默认使用瞬时内存存储空间，这非常适合测试目的。
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		errorLogger:    log.New(io.Discard, "", 0),
		infoLogger:     log.New(io.Discard, "", 0),
		snippets:       &mocks.SnippetModel{},
		users:          &mocks.UserModel{},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
}

// 定义嵌入 httptest.Server 实例的自定义 testServer 类型。
type testServer struct {
	*httptest.Server
}

// 创建 newTestServer 辅助程序，该程序将初始化并返回一个自定义 testServer 类型的新实例。
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	// 将 cookie jar 添加到测试服务器客户端。现在，使用该客户端时，任何响应 cookie 都将被存储并随后续请求一起发送。
	ts.Client().Jar = jar
	// 通过设置自定义 CheckRedirect 函数，禁用测试服务器客户端的重定向跟踪功能。
	// 每当客户端收到 3xx 响应时，该函数就会被调用，并通过始终返回 http.ErrUseLastResponse 错误，强制客户端立即返回收到的响应。
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{ts}
}

// 在自定义 testServer 类型上实现 get() 方法。该方法使用测试服务器客户端向给定的 url 路径发出 GET 请求，并返回响应状态代码、标题和正文。
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}
