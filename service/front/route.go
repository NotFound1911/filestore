package main

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"net/http"
)

// loadTemplates 加载模板文件
func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	// 加载模板文件
	r.AddFromFiles("index.html", templatesDir+"/index.html")
	r.AddFromFiles("register.html", templatesDir+"/register.html")
	return r
}
func registerRoutes(core *gin.Engine) {
	fs := core.Group("/")
	fs.GET("/", indexHandler())     // 首页
	fs.GET("/register", register()) // 注册页
}

// indexHandler 首页
func indexHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "filestore"})
	}
}

// register 注册页
func register() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{"title": "filestore"})
	}
}
