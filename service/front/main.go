package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"path/filepath"
)

const staticPath = "./service/front/static/view"

// 处理注册页面请求，返回HTML页面
func registerHandler(w http.ResponseWriter, r *http.Request) {
	// 读取并解析HTML文件
	tmpl := template.Must(template.ParseFiles(filepath.Join(staticPath, "register.html")))
	// 将模板渲染并发送给客户端
	tmpl.Execute(w, nil)
}

// 处理登录请求
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// 解析表单数据
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// 获取用户名和密码
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	fmt.Println("username:", username)
	fmt.Println("password:", password)
	// 假设这里进行简单的验证，实际情况需要更复杂的逻辑
	if username == "admin" && password == "admin123" {
		fmt.Fprintf(w, "Welcome, %s!", username)
	} else {
		fmt.Fprintf(w, "Invalid username or password.")
	}
}

func main() {
	server := gin.Default()
	server.HTMLRender = loadTemplates(staticPath)
	registerRoutes(server)
	if err := server.Run(":8080"); err != nil {
		fmt.Println("err:", err)
	}
}
