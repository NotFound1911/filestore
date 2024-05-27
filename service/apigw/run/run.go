package run

import (
	"context"
	"fmt"
	accountv1 "github.com/NotFound1911/filestore/api/proto/gen/account/v1"
	v1 "github.com/NotFound1911/filestore/api/rest/apigw/v1"
	"github.com/NotFound1911/filestore/config"
	"github.com/NotFound1911/filestore/internal/web/jwt"
	"github.com/NotFound1911/filestore/internal/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/redis/go-redis/v9"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

func Run() {
	conf := config.NewConfig("")
	gin.SetMode(conf.Service.Apigw.Http.Mode)
	server := gin.Default()
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       conf.Redis.Db,
	})
	hdl := jwt.NewRedisJWTHandler(rdb)
	server.Use(middleware.NewLoginJWTMiddlewareBuilder(hdl).CheckLogin())

	cli, err := etcdv3.New(etcdv3.Config{
		Endpoints: conf.Etcd.Endpoints,
	})
	if err != nil {
		panic(err)
	}
	// 默认是 WRR 负载均衡算法
	r := etcd.New(cli)
	cc, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%s", conf.Service.Account.Name)),
		grpc.WithDiscovery(r),
	)
	defer cc.Close()
	client := accountv1.NewAccountServiceClient(cc)

	userHandler := v1.NewUserHandler(client, hdl)
	userHandler.RegisterUserRoutes(server)
	server.Run(conf.Service.Apigw.Http.Addr...)
}
