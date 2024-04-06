package main

import (
	"fmt"
	userv1 "github.com/NotFound1911/filestore/app/account/api/proto/gen/user/v1"
	gprcserv "github.com/NotFound1911/filestore/app/account/grpc"
	"github.com/NotFound1911/filestore/app/account/ioc"
	"github.com/NotFound1911/filestore/app/account/repository"
	"github.com/NotFound1911/filestore/app/account/repository/dao"
	"github.com/NotFound1911/filestore/app/account/service"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	cli, err := etcdv3.New(etcdv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		panic(err)
	}
	grpcSrv := grpc.NewServer(
		grpc.Address(":8090"),
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
	userv1.RegisterUserServiceServer(grpcSrv, account)
	// etcd 注册中心
	r := etcd.New(cli)
	app := kratos.New(
		kratos.Name("user"),
		kratos.Server(
			grpcSrv,
		),
		kratos.Registrar(r),
	)
	err = app.Run()
	fmt.Println("err:", err)
}
