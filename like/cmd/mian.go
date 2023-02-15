package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"like/config"
	"like/discovery"
	"like/internal/handler"
	"like/internal/repository"
	"like/internal/service"
	"net"
)

func main() {
	config.InitConfig()
	repository.InitDB()
	// etcd地址
	etcdAddress := []string{viper.GetString("etcd.address")}
	// 服务的注册
	etcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	grpcAddress := viper.GetString("server.grpcAddress")
	userNode := discovery.Server{
		Name: viper.GetString("server.domain"),
		Addr: grpcAddress,
	}
	server := grpc.NewServer()
	defer server.Stop()
	// 绑定服务
	service.RegisterLikeServiceServer(server, handler.NewLikeService())
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if _, err = etcdRegister.Register(userNode, 10); err != nil {
		panic(err)
	}
	if err = server.Serve(lis); err != nil {
		panic(err)
	}
}
