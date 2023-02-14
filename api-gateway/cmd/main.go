package main

import (
	"api-gateway/config"
	"api-gateway/discovery"
	"api-gateway/internal/service"
	"api-gateway/middleware/wrapper"
	"api-gateway/pkg/util"
	"api-gateway/routes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.InitConfig()

	go startListen()
	{
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		s := <-osSignal
		fmt.Println("exit!", s)
	}
	fmt.Println("gateway listen on :4000")

}

func startListen() {
	// etcd注册
	etcdAddress := []string{viper.GetString("etcd.address")}
	etcdRegister := discovery.NewResolver(etcdAddress, logrus.New())
	resolver.Register(etcdRegister)
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	// 服务名
	userServiceName := viper.GetString("domain.user")
	videoServiceName := viper.GetString("domain.video")
	commentServiceName := viper.GetString("domain.comment")

	// RPC 连接
	connUser, err := RPCConnect(ctx, userServiceName, etcdRegister)
	if err != nil {
		return
	}
	userService := service.NewUserServiceClient(connUser)

	connVideo, err := RPCConnect(ctx, videoServiceName, etcdRegister)
	if err != nil {
		return
	}
	videoService := service.NewVideoServiceClient(connVideo)

	connComment, err := RPCConnect(ctx, commentServiceName, etcdRegister)
	if err != nil {
		return
	}
	commentService := service.NewCommentServiceClient(connComment)

	// 加入熔断 TODO main太臃肿了
	wrapper.NewServiceWrapper(userServiceName)
	wrapper.NewServiceWrapper(videoServiceName)
	wrapper.NewServiceWrapper(commentServiceName)

	ginRouter := routes.NewRouter(userService, videoService, commentService)
	server := &http.Server{
		Addr:           viper.GetString("server.port"),
		Handler:        ginRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("绑定HTTP到 %s 失败！可能是端口已经被占用，或用户权限不足")
		fmt.Println(err)
	}
	go func() {
		// 优雅关闭
		util.GracefullyShutdown(server)
	}()
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("gateway启动失败, err: ", err)
	}
}

func RPCConnect(ctx context.Context, serviceName string, etcdRegister *discovery.Resolver) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	addr := fmt.Sprintf("%s:///%s", etcdRegister.Scheme(), serviceName)
	conn, err = grpc.DialContext(ctx, addr, opts...)
	return
}
