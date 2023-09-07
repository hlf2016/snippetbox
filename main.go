package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from snippetBox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from snippet view"))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create snippet"))
}

func main() {
	mux := http.NewServeMux()
	//  兜底路由 其他路由都不匹配的 走这里
	// tree path 相当于 /**
	mux.HandleFunc("/", home)

	// fixed path 只有精确匹配才会调用响应的方法 跟上面的区别就是 是否以 / 结尾
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
