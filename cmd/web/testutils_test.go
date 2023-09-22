package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestApplication(t *testing.T) *application {
	return &application{
		errorLogger: log.New(io.Discard, "", 0),
		infoLogger:  log.New(io.Discard, "", 0),
	}
}

// 定义嵌入 httptest.Server 实例的自定义 testServer 类型。
type testServer struct {
	*httptest.Server
}

// 创建 newTestServer 辅助程序，该程序将初始化并返回一个自定义 testServer 类型的新实例。
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
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
