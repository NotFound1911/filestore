package run

import (
	"context"
	"fmt"
	file_managerv1 "github.com/NotFound1911/filestore/api/proto/gen/file_manager/v1"
	v1 "github.com/NotFound1911/filestore/api/rest/download/v1"
	"github.com/NotFound1911/filestore/config"
	"github.com/NotFound1911/filestore/internal/logger"
	"github.com/NotFound1911/filestore/internal/web/jwt"
	"github.com/NotFound1911/filestore/internal/web/middleware"
	"github.com/NotFound1911/filestore/service/download/ioc"
	"github.com/NotFound1911/filestore/service/download/repository"
	"github.com/NotFound1911/filestore/service/download/repository/dao"
	"github.com/NotFound1911/filestore/service/download/service"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/redis/go-redis/v9"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

func Run() {
	conf := config.NewConfig("")
	gin.SetMode(conf.Service.Download.Http.Mode)
	server := gin.Default()
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       conf.Redis.Db,
	})
	hdl := jwt.NewRedisJWTHandler(rdb)
	server.Use(middleware.NewLoginJWTMiddlewareBuilder(hdl).CheckLogin())
	orm := ioc.InitDb(conf)
	downloadDao := dao.NewOrmDownload(orm)
	downRepo := repository.NewDownloadRepository(downloadDao)
	downService := service.NewDownloadService(downRepo)
	cli, err := etcdv3.New(etcdv3.Config{
		Endpoints: conf.Etcd.Endpoints,
	})
	if err != nil {
		panic(err)
	}
	r := etcd.New(cli)
	cc, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%s", conf.Service.FileManager.Name)),
		grpc.WithDiscovery(r),
	)
	defer cc.Close()
	client := file_managerv1.NewFileManagerServiceClient(cc)
	log := logger.New(conf, conf.Service.Download.Name)
	downloadHandler := v1.NewHandler(downService, hdl, client, log)
	downloadHandler.RegisterDownloadRoutes(server)
	server.Run(conf.Service.Download.Http.Addr...)
}
