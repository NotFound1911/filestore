package run

import (
	"context"
	"fmt"
	file_managerv1 "github.com/NotFound1911/filestore/api/proto/gen/file_manager/v1"
	"github.com/NotFound1911/filestore/api/rest/upload/v1"
	"github.com/NotFound1911/filestore/config"
	"github.com/NotFound1911/filestore/internal/logger"
	"github.com/NotFound1911/filestore/internal/mq"
	"github.com/NotFound1911/filestore/internal/storage"
	"github.com/NotFound1911/filestore/internal/web/jwt"
	"github.com/NotFound1911/filestore/internal/web/middleware"
	"github.com/NotFound1911/filestore/service/upload/ioc"
	"github.com/NotFound1911/filestore/service/upload/repository"
	"github.com/NotFound1911/filestore/service/upload/repository/cache"
	"github.com/NotFound1911/filestore/service/upload/repository/dao"
	"github.com/NotFound1911/filestore/service/upload/service"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/redis/go-redis/v9"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

func Run() {
	conf := config.NewConfig("")
	gin.SetMode(conf.Service.Upload.Http.Mode)
	server := gin.Default()
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password, // no password set
		DB:       conf.Redis.Db,       // use default DB
	})
	hdl := jwt.NewRedisJWTHandler(rdb)
	server.Use(middleware.NewLoginJWTMiddlewareBuilder(hdl).CheckLogin())
	// **************
	// gorm
	orm := ioc.InitDb()
	// dao
	uploadDao := dao.NewOrmUpload(orm)
	// cache
	uploadCache := cache.NewChunkCache(rdb)
	// repository
	uploadRepo := repository.NewUploadRepository(uploadDao, uploadCache)
	// service
	uploadService := service.NewUploadService(uploadRepo)
	// 服务发现
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
	log := logger.New(conf, conf.Service.Upload.Name)
	uploadHandler := v1.NewHandler(uploadService,
		hdl, client, v1.DiHandler{
			Storage:      storage.New(conf, log),
			MessageQueue: mq.New(conf, log),
			Logger:       log,
		})
	uploadHandler.RegisterUploadRoutes(server)

	server.Run(conf.Service.Upload.Http.Addr...)
}
