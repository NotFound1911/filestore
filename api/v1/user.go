package v1

import (
	ijwt "github.com/NotFound1911/filestore/api/v1/jwt"
	"github.com/NotFound1911/filestore/domain"
	"github.com/NotFound1911/filestore/errs"
	"github.com/NotFound1911/filestore/internal/server"
	"github.com/NotFound1911/filestore/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	bizLogin             = "login"
)

type UserHandler struct {
	ijwt.Handler
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            service.UserService
}

func NewUserHandler(svc service.UserService,
	hdl ijwt.Handler) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
		Handler:        hdl,
	}
}
func (u *UserHandler) SignUp(ctx *gin.Context, req SignUpReq) (server.Result, error) {
	isEmail, err := u.emailRexExp.MatchString(req.Email)
	if err != nil {
		return server.Result{
			Code: errs.UserInvalidInput,
			Msg:  "系统错误",
		}, err
	}
	if !isEmail {
		return server.Result{
			Code: errs.UserInvalidInput,
			Msg:  "非法邮箱格式",
		}, nil
	}
	if req.Password != req.ConfirmPassword {
		return server.Result{
			Code: errs.UserInvalidInput,
			Msg:  "两次输入的密码不相等",
		}, nil
	}
	isPassword, err := u.passwordRexExp.MatchString(req.Password)
	if err != nil {
		return server.Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统错误",
		}, err
	}
	if !isPassword {
		return server.Result{
			Code: errs.UserInvalidInput,
			Msg:  "密码必须包含字母、数字、特殊字符,并且不少于八位",
		}, nil
	}
	err = u.svc.Signup(ctx.Request.Context(), domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		return server.Result{
			Msg: "OK",
		}, nil
	case service.ErrDuplicateEmail:
		return server.Result{
			Code: errs.UserDuplicateEmail,
			Msg:  "邮箱冲突",
		}, nil
	default:
		return server.Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统错误",
		}, err
	}
}
func (u *UserHandler) LoginJWT(ctx *gin.Context, req LoginJWTReq) (server.Result, error) {
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		err = u.SetLoginToken(ctx, user.Id)
		if err != nil {
			return server.Result{
				Code: -1,
				Msg:  "系统错误",
			}, err
		}
		return server.Result{
			Msg: "OK",
		}, nil
	case service.ErrInvalidUserOrPassword:
		return server.Result{Msg: "用户名或者密码错误"}, nil
	default:
		return server.Result{Msg: "系统错误"}, err
	}
}
func (u *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := u.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, server.Result{Code: -1, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, server.Result{Msg: "退出登录成功"})
}
func (u *UserHandler) Profile(ctx *gin.Context,
	uc ijwt.UserClaims) (server.Result, error) {
	//us := ctx.MustGet("user").(UserClaims)
	//ctx.String(http.StatusOK, "这是 profile")
	user, err := u.svc.FindById(ctx, uc.Uid)
	if err != nil {
		return server.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	return server.Result{
		Data: User{
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
func (u *UserHandler) RegisterUserRoutes(core *gin.Engine) {

	ug := core.Group("/api/storage/v1/users")
	ug.POST("/signup", server.WrapBody(u.SignUp))
	ug.POST("/login", server.WrapBody(u.LoginJWT))
	ug.POST("/logout", u.LogoutJWT)
	ug.GET("/profile", server.WrapClaims(u.Profile))
}
