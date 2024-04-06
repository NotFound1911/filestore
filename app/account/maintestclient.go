package main

import (
	"context"
	"fmt"
	userv1 "github.com/NotFound1911/filestore/app/account/api/proto/gen/user/v1"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	resp, err := client.Signup(ctx, &userv1.SignupReq{User: &userv1.User{
		Email: "123@123.com",
	}})
	fmt.Println("resp:", resp)
	fmt.Println("err:", err)
	cancel()

}
