package v1

import (
	"github.com/NotFound1911/filestore/api/proto/gen/account/v1"
	"github.com/NotFound1911/filestore/errs"
	"github.com/NotFound1911/filestore/internal/web/jwt"
	serv "github.com/NotFound1911/filestore/pkg/server"
	"github.com/NotFound1911/filestore/service/account/service"
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
	jwt.Handler
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	client         accountv1.AccountServiceClient
}

func NewUserHandler(client accountv1.AccountServiceClient, hdl jwt.Handler) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		client:         client,
		Handler:        hdl,
	}
}
func (u *UserHandler) Signup(ctx *gin.Context, req SignupReq) (serv.Result, error) {
	isEmail, err := u.emailRexExp.MatchString(req.Email)
	if err != nil {
		return serv.Result{
			Code: errs.UserInvalidInput,
			Msg:  "系统错误",
		}, err
	}
	if !isEmail {
		return serv.Result{
			Code: errs.UserInvalidInput,
			Msg:  "非法邮箱格式",
		}, nil
	}
	if req.Password != req.ConfirmPassword {
		return serv.Result{
			Code: errs.UserInvalidInput,
			Msg:  "两次输入的密码不相等",
		}, nil
	}
	isPassword, err := u.passwordRexExp.MatchString(req.Password)
	if err != nil {
		return serv.Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统错误",
		}, err
	}
	if !isPassword {
		return serv.Result{
			Code: errs.UserInvalidInput,
			Msg:  "密码必须包含字母、数字、特殊字符,并且不少于八位",
		}, nil
	}
	_, err = u.client.Signup(ctx.Request.Context(), &accountv1.SignupReq{
		User: &accountv1.User{
			Email:    req.Email,
			Password: req.Password,
		},
	})
	if err != nil {
		return serv.Result{
			Code: errs.UserInternalServerError,
			Msg:  "注册失败",
			Data: err.Error(),
		}, err
	}
	return serv.Result{
		Code: 2000,
		Msg:  "注册成功",
	}, nil
}
func (u *UserHandler) LoginJWT(ctx *gin.Context, req LoginJWTReq) (serv.Result, error) {
	user, err := u.client.Login(ctx, &accountv1.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		err = u.SetLoginToken(ctx, user.Id)
		if err != nil {
			return serv.Result{
				Code: -1,
				Msg:  "系统错误",
			}, err
		}
		return serv.Result{
			Msg: "OK",
		}, nil
	case service.ErrInvalidUserOrPassword:
		return serv.Result{Msg: "用户名或者密码错误"}, nil
	default:
		return serv.Result{Msg: "系统错误"}, err
	}
}
func (u *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := u.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, serv.Result{Code: -1, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, serv.Result{Msg: "退出登录成功"})
}
func (u *UserHandler) Profile(ctx *gin.Context, uc jwt.UserClaims) (serv.Result, error) {
	user, err := u.client.Profile(ctx, &accountv1.ProfileReq{
		Id: uc.UId,
	})
	if err != nil {
		return serv.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}
	return serv.Result{
		Code: 2000,
		Msg:  "ok",
		Data: User{
			Name:  user.User.Name,
			Email: user.User.Email,
			Phone: user.User.Phone,
		},
	}, nil
}
func (u *UserHandler) RegisterUserRoutes(core *gin.Engine) {
	ug := core.Group("/api/storage/v1/users")
	ug.POST("/signup", serv.WrapBody(u.Signup))
	ug.POST("/login", serv.WrapBody(u.LoginJWT))
	ug.POST("/logout", u.LogoutJWT)
	ug.GET("/profile", serv.WrapClaims(u.Profile))
}
