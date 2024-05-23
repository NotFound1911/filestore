package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"log"
	"net/http"
)

func WrapBody[Req any](
	bizFn func(ctx *gin.Context, req Req) (Result, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		r := ctx.Request
		body, err := io.ReadAll(r.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("请求错误:%v", err))
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		ctx.Request = r
		if err := json.Unmarshal(body, &req); err != nil {
			ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("请求错误:%v", err))
			return
		}
		res, err := bizFn(ctx, req)
		if err != nil {
			fmt.Printf("执行失败:%v", err)
		}
		ctx.JSON(http.StatusOK, res)
	}
}
func WrapClaims[Claims any](
	bizFn func(ctx *gin.Context, uc Claims) (Result, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		val, ok := ctx.Get("user")
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, Result{
				Code: -1,
				Msg:  "认证失败,无效用户",
			})
			return
		}
		uc, ok := val.(Claims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, Result{
				Code: -1,
				Msg:  "认证失败,无效用户",
			})
			return
		}
		res, err := bizFn(ctx, uc)
		if err != nil {
			log.Printf("执行业务逻辑失败:%v", err)
		}
		ctx.JSON(http.StatusOK, res)
	}
}

// WrapBodyAndClaims bizFn 就是你的业务逻辑
func WrapBodyAndClaims[Req any, Claims jwt.Claims](
	bizFn func(ctx *gin.Context, req Req, uc Claims) (Result, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			log.Println("输入错误:", err)
			return
		}
		val, ok := ctx.Get("user")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		uc, ok := val.(Claims)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		res, err := bizFn(ctx, req, uc)
		if err != nil {
			log.Printf("执行业务逻辑失败:%v", err)
		}
		ctx.JSON(http.StatusOK, res)
	}
}
