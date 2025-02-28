package main

import (
	"flag"
	"fmt"
	"log"
	"mxshop_sers/goods_srv/global"
	"mxshop_sers/goods_srv/utils"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"mxshop_sers/goods_srv/handler"
	"mxshop_sers/goods_srv/initialize"
	pb "mxshop_sers/goods_srv/proto"
)

func main() {
	//Ip := flag.String("ip", "127.0.0.1", "ip地址")
	//测试的时候把value改成了50051，正确的应该是0
	Port := flag.Int("port", 0, "端口")

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	zap.S().Info("config:", global.ServerConfig)

	flag.Parse()
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}
	zap.S().Info("ip: ", global.ServerConfig.Host, " port: ", *Port)

	server := grpc.NewServer()
	pb.RegisterGoodsServer(server, &handler.GoodsServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//注册健康服务检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	registration.ID = uuid.New().String()
	registration.Port = *Port
	registration.Check = check
	registration.Address = global.ServerConfig.Host
	registration.Tags = global.ServerConfig.Tags

	go func() {
		//注册服务
		if err := client.Agent().ServiceRegister(registration); err != nil {
			panic(err)
		}
		//监听grpc请求
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(registration.ID); err != nil {
		zap.S().Fatalf("注销失败: %v", err)
	}
	zap.S().Info("注销成功")
}
