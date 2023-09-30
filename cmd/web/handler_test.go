package main

import (
	"github.com/hlf2016/snippetbox/internal/assert"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	//创建应用程序结构体的新实例。目前，该结构只包含几个模拟日志记录器（它们会丢弃写入其中的任何内容）。
	app := newTestApplication(t)
	// 然后，我们使用 httptest.NewTLSServer() 函数创建一个新的测试服务器，并将 app.routes() 方法返回的值作为服务器的处理程序。
	// 这样就启动了一个 HTTPS 服务器，在测试期间监听本地计算机随机选择的端口。请注意，我们推迟了对 ts.Close() 的调用，以便在测试结束时关闭服务器
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	// 测试服务器正在侦听的网络地址包含在ts.URL字段中。
	// 我们可以将其与ts.Client().Get()方法一起使用，以对测试服务器发出GET ping请求。这将返回一个包含响应的HTTP.Response结构。
	statusCode, _, body := ts.get(t, "/ping")
	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)
			assert.Equal(t, code, tt.wantCode)
			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}

}
