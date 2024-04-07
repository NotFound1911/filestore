package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type RedisJWTHandler struct {
	client        redis.Cmdable
	signingMethod jwt.SigningMethod
	rcExpiration  time.Duration
}

var _ Handler = &RedisJWTHandler{}

func (r *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	uc, ok := ctx.Value(userKey{}).(UserClaims)
	if !ok {
		return fmt.Errorf("信息丢失")
	}
	return r.client.Set(ctx,
		fmt.Sprintf("users:ssid:%s", uc.Ssid),
		"", r.rcExpiration).Err()
}

// ExtractToken 根据约定，token 在 Authorization 头部
// X-Refresh-Token+" "+X-Jwt-Token
func (r *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	authCode := ctx.Request.Header.Get("Authorization")
	if authCode == "" {
		return authCode
	}
	segs := strings.Split(authCode, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

func (r *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := r.setRefreshToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	return r.SetJWTToken(ctx, uid, ssid)
}

func (r *RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	uc := UserClaims{
		Uid:       uid,
		Ssid:      ssid,
		UserAgent: ctx.Request.Header.Get("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			// 30 分钟过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	token := jwt.NewWithClaims(r.signingMethod, uc)
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (r *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	fmt.Println("ssid:", ssid)
	cnt, err := r.client.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	fmt.Println("cnt:", cnt)
	fmt.Println("err:", err)
	if err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("token 无效")
	}
	return nil
}

func NewRedisJWTHandler(client redis.Cmdable) Handler {
	return &RedisJWTHandler{
		client:        client,
		signingMethod: jwt.SigningMethodHS512,
		rcExpiration:  time.Hour * 24 * 7,
	}
}
func (r *RedisJWTHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	rc := RefreshClaims{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(r.rcExpiration)),
		},
	}
	token := jwt.NewWithClaims(r.signingMethod, rc)
	tokenStr, err := token.SignedString(RCJWTKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")
var RCJWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgA")

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid  int64
	Ssid string
}
type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	Ssid      string
	UserAgent string
}

type userKey struct {
}