package middleware

import (
	jwt2 "github.com/NotFound1911/filestore/internal/web/jwt"
	"github.com/NotFound1911/filestore/pkg/server"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type LoginJWTMiddlewareBuilder struct {
	jwt2.Handler
}

func NewLoginJWTMiddlewareBuilder(hdl jwt2.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: hdl,
	}
}

func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/api/storage/v1/users/signup" ||
			path == "/api/storage/v1/users/login" {
			// 不需要登录校验
			return
		}
		tokenStr := m.ExtractToken(ctx)
		var uc jwt2.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return jwt2.JWTKey, nil
		})
		if err != nil {
			// token 不对，token 是伪造的
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Result{
				Code: -1,
				Msg:  "认证失败",
			})
			return
		}
		if token == nil || !token.Valid {
			// 在这里发现 access_token 过期了，生成一个新的 access_token
			// token 解析出来了，但是 token 可能是非法的，或者过期了的
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Result{
				Code: -1,
				Msg:  "认证失败",
			})
			return
		}
		// 这里看
		err = m.CheckSession(ctx, uc.Ssid)
		if err != nil {
			// token 无效或者 redis 有问题
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Result{
				Code: -1,
				Msg:  "认证失败",
			})
			return
		}
		ctx.Set("user", uc)
	}
}
