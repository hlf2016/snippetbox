package main

import (
	"bytes"
	"github.com/hlf2016/snippetbox/internal/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	// 它本质上是 http.ResponseWriter 的一种实现，用于记录响应状态代码、标题和正文，而不是将其实际写入 HTTP 连接。
	rr := httptest.NewRecorder()

	// 初始化一个新的虚拟http.Request
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 调用ping处理程序函数，传入httptest.ResponseRecorder和Http.Request。
	ping(rr, r)

	// 调用 http.ResponseRecorder 的 Result() 方法，获取 ping 处理程序生成的 http.Response。
	rs := rr.Result()

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// 我们可以检查 ping 处理程序写入的响应体是否等于 "OK"。
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		// 我们在几个地方使用了 t.Fatal() 函数来处理测试代码中出现意外错误的情况。
		// 调用时，t.Fatal() 会将测试标记为失败，记录错误，然后完全停止当前测试（或子测试）的执行。
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
