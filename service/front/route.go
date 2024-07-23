package main

import (
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	X_Jwt_Token         = ""
	X_Refresh_Token     = ""
	Key_X_Jwt_Token     = "X-Jwt-Token"
	Key_X_Refresh_Token = "X-Refresh-Token"
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
	fs.GET("/profile", profileHandler())
	fs.GET("/list", listHandler())
	fs.POST("/download", downloadHandler())
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
		res, err := loginRequest(req)
		if err != nil || res.RespStatusCode != 200 {
			c.JSON(http.StatusBadRequest, "请求错误")
			return
		}
		X_Jwt_Token = res.RespHeader.Get(Key_X_Jwt_Token)
		X_Refresh_Token = res.RespHeader.Get(Key_X_Refresh_Token)
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
		res, err := signupRequest(req)
		if err != nil || res.RespStatusCode != 200 {
			c.JSON(http.StatusBadRequest, "请求错误")
			return
		}
		// 登录成功页面
		c.JSON(http.StatusOK, res.Result)
	}
}

// registerSuccessfulHandler 注册成功页面
func registerSuccessfulHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "register_successful.html", gin.H{"title": "filestore"})
	}
}

// profileHandler 基本信息
func profileHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := profileRequest()
		if err != nil || res.RespStatusCode != 200 {
			c.JSON(http.StatusBadRequest, "请求错误")
			return
		}
		// 登录成功页面
		c.JSON(http.StatusOK, res.Result)
	}
}

func listHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := listRequest()
		if err != nil || res.RespStatusCode != 200 {
			c.JSON(http.StatusBadRequest, "请求错误")
			return
		}
		// 登录成功页面
		c.JSON(http.StatusOK, res.Result)
	}
}

func downloadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &DownloadReq{}
		if err := c.ShouldBind(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := downloadUrlRequest(req)
		if err != nil || res.RespStatusCode != 200 {
			c.JSON(http.StatusBadRequest, "请求错误")
			return
		}
		fmt.Println("res:", res)
	}
}
