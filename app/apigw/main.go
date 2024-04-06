package main

import (
	"context"
	userv1 "github.com/NotFound1911/filestore/app/account/api/proto/gen/user/v1"
	v1 "github.com/NotFound1911/filestore/app/apigw/api/v1"
	"github.com/NotFound1911/filestore/internal/web/jwt"
	"github.com/NotFound1911/filestore/internal/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/redis/go-redis/v9"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	gin.SetMode(gin.DebugMode)
	server := gin.Default()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
	})
	hdl := jwt.NewRedisJWTHandler(rdb)
	server.Use(middleware.NewLoginJWTMiddlewareBuilder(hdl).CheckLogin())

	cli, err := etcdv3.New(etcdv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		panic(err)
	}
	// 默认是 WRR 负载均衡算法
	r := etcd.New(cli)
	cc, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint("discovery:///user"),
		grpc.WithDiscovery(r),
	)
	defer cc.Close()
	client := userv1.NewUserServiceClient(cc)

	userHandler := v1.NewUserHandler(client, hdl)
	userHandler.RegisterUserRoutes(server)
	server.Run(":8888")
}
