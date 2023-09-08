package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// 需要指出的是，您传递给template.ParseFiles()函数的文件路径必须是相对于当前工作目录的，也就是运行 go run 的目录，或者是绝对路径。在下面的代码中，我设置了相对于项目目录根目录的路径。
	files := []string{
		"./ui/html/pages/base.tmpl",
		"./ui/html/pages/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// 然后，我们在模板集上使用 Execute() 方法将模板内容写入响应体。Execute() 的最后一个参数代表我们要传入的任何动态数据，现在我们暂且将其保持为nil
	// err = ts.Execute(w, nil)
	// 使用 ExecuteTemplate() 方法将 "base "模板的内容写入响应体。
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}
func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
