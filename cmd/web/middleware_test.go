package main

import (
	"github.com/hlf2016/snippetbox/internal/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 创建一个模拟 HTTP 处理程序，我们可以将其传递给 secureHeaders 中间件，由它写入 200 状态代码和 "OK "响应体。
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// 将模拟的HTTP处理程序传递给我们的secureHeaders中间件。
	// 因为secureHeaders返回一个HTTP.Handler，所以我们可以调用它的ServeHTTP()方法，传入HTTP.ResponseRecorder和虚拟的Http.Request来执行它。
	secureHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	// 检查中间件是否正确设置了响应的内容-安全-策略(Content-Security-Policy)标头。
	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	// 检查中间件是否在响应中正确设置了 Referrer-Policy 标头。
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)

	// 检查中间件是否正确设置了响应的 X-Content-Type-Options 头信息。
	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)
	// 检查中间件是否正确设置了响应的 X-Frame-Options 标头.
	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)
	// 检查中间件是否在响应中正确设置了 X-XSS-Protection 标头
	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)
	// 检查中间件是否正确调用了下一个处理程序，响应状态代码和正文是否符合预期。
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(body), "OK")

}
