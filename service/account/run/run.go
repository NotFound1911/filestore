package run

import (
	"github.com/NotFound1911/filestore/api/proto/gen/account/v1"
	"github.com/NotFound1911/filestore/config"
	gprcserv "github.com/NotFound1911/filestore/service/account/grpc"
	"github.com/NotFound1911/filestore/service/account/ioc"
	"github.com/NotFound1911/filestore/service/account/repository"
	"github.com/NotFound1911/filestore/service/account/repository/dao"
	"github.com/NotFound1911/filestore/service/account/service"
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
		grpc.Address(conf.Service.Account.Grpc.Addr),
		grpc.Middleware(recovery.Recovery()),
	)
	// **************
	// gorm
	orm := ioc.InitDb()
	// dao
	userDao := dao.NewOrmUser(orm)
	// repository
	userRepo := repository.NewCachedUserRepository(userDao)
	// service
	userService := service.NewUserService(userRepo)
	// grpc
	account := gprcserv.NewAccountServiceServer(userService)
	// **************
	accountv1.RegisterAccountServiceServer(grpcSrv, account)
	// etcd 注册中心
	r := etcd.New(cli)
	app := kratos.New(
		kratos.Name(conf.Service.Account.Name),
		kratos.Server(
			grpcSrv,
		),
		kratos.Registrar(r),
	)
	err = app.Run()
}
