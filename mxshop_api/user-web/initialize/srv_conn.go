package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver" //很重要
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mxshop_api/user-web/global"
	pb "mxshop_api/user-web/proto"
)

func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulConfig

	userConn, err := grpc.NewClient(
		//fmt.Sprintf("dns:///%s:%d", userSrvHost, userSrvPort),
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvConfig.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig":[{"round_robin":{}}]}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 用户服务失败")
		return
	}
	zap.S().Infof("连接用户服务成功")
	userSrvClient := pb.NewUserClient(userConn)
	global.UserSrvClien = userSrvClient
}
