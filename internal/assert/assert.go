package assert

import "testing"

func Equal[T comparable](t *testing.T, actual, expected T) {
	// t.Helper() 函数向 Go 测试运行程序表明，我们的 Equal() 函数是一个测试辅助函数。这意味着，当从我们的 Equal() 函数调用 t.Errorf() 时，Go 测试运行程序将在输出中报告调用我们 Equal() 函数的代码的文件名和行号。
	t.Helper()

	if actual != expected {
		t.Errorf("got %v; want %v", actual, expected)
	}
}
