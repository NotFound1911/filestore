package main

import (
	"fmt"
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
	r.AddFromFiles("home.html", templatesDir+"/home.html")
	r.AddFromFiles("register_successful.html", templatesDir+"/register_successful.html")
	return r
}
func registerRoutes(core *gin.Engine) {
	fs := core.Group("/")
	fs.GET("/", indexHandler())     // 首页
	fs.GET("/register", register()) // 注册页
	fs.POST("/login", loginHandler())
	fs.POST("/signup", signupHandler())
	fs.GET("/register-successful", registerSuccessfulHandler())
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

// loginHandler 处理登录请求
func loginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 绑定表单数据到结构体
		req := &LoginReq{}
		if err := c.ShouldBind(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// todo
		fmt.Println("req:", req)
		// 登录成功页面
		c.HTML(http.StatusOK, "home.html", gin.H{"title": "filestore"})
	}
}

// signupHandler 处理注册请求
func signupHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &SignupReq{}
		if err := c.ShouldBind(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// todo
		fmt.Println("req:", req)
		// 登录成功页面
		c.JSON(http.StatusOK, "注册成功")
	}
}

// registerSuccessfulHandler 注册成功页面
func registerSuccessfulHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "register_successful.html", gin.H{"title": "filestore"})
	}
}
