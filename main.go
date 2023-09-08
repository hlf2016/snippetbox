package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from snippetBox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	// 获取请求中的 id 参数 并转 int
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	//使用 fmt.Fprintf() 函数将 id 值与我们的响应进行插值，并将其写入 http.ResponseWriter 中。
	fmt.Fprintf(w, "Display specific snippet with ID %d", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// 在调用 w.WriteHeader() 或 w.Write() 后更改响应标头映射不会影响用户收到的标头。
		// 在调用这些方法之前，您需要确保响应标头映射包含您想要的所有标头。
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create snippet"))
}

func main() {
	/*
		不建议直接使用 http 内建的 默认 mux serve 因为比较危险
		由于 DefaultServeMux 是一个全局变量，因此任何软件包都可以访问它并注册路由，包括应用程序导入的任何第三方软件包。
		如果这些第三方软件包中的一个受到攻击，它们就会使用 DefaultServeMux 向网络暴露恶意处理程序。
		因此，为了安全起见，避免使用 DefaultServeMux 和相应的辅助函数通常是个好主意。
		请使用你自己的本地范围 servicemux，就像我们迄今为止在本项目中所做的那样。
	*/
	mux := http.NewServeMux()
	// 兜底路由 其他路由都不匹配的 走这里
	// tree path 相当于 /**
	mux.HandleFunc("/", home)

	// fixed path 只有精确匹配才会调用响应的方法 跟上面的区别就是 是否以 / 结尾
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
