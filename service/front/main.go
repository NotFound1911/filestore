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

func main() {
	server := gin.Default()
	server.HTMLRender = loadTemplates(staticPath)
	registerRoutes(server)
	if err := server.Run(":8080"); err != nil {
		fmt.Println("err:", err)
	}
}
