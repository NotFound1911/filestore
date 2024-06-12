package run

import (
	file_managerv1 "github.com/NotFound1911/filestore/api/proto/gen/file_manager/v1"
	"github.com/NotFound1911/filestore/config"
	gprcserv "github.com/NotFound1911/filestore/service/file_manager/grpc"
	"github.com/NotFound1911/filestore/service/file_manager/ioc"
	"github.com/NotFound1911/filestore/service/file_manager/repository"
	"github.com/NotFound1911/filestore/service/file_manager/repository/dao"
	"github.com/NotFound1911/filestore/service/file_manager/service"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

func Run() {
	conf := config.NewConfig("")
	cli, err := etcdv3.New(etcdv3.Config{
		Endpoints: conf.Etcd.Endpoints,
	})
	if err != nil {
		panic(err)
	}
	grpcSrv := grpc.NewServer(
		grpc.Address(conf.Service.FileManager.Grpc.Addr),
		grpc.Middleware(recovery.Recovery()),
	)
	// **************
	// gorm
	orm := ioc.InitDb(conf)
	// dao
	fileMangerDao := dao.NewOrmFileManager(orm)
	// repository
	fileManagerRepo := repository.NewFileManagerRepository(fileMangerDao)
	// service
	fileManagerService := service.NewFileManagerService(fileManagerRepo)
	// grpc
	fileManagerGrpc := gprcserv.NewFileManagerServiceServer(fileManagerService)
	// **************
	file_managerv1.RegisterFileManagerServiceServer(grpcSrv, fileManagerGrpc)
	// etcd 注册中心
	r := etcd.New(cli)
	app := kratos.New(
		kratos.Name(conf.Service.FileManager.Name),
		kratos.Server(
			grpcSrv,
		),
		kratos.Registrar(r),
	)
	err = app.Run()
}
