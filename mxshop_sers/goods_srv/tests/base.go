package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "mxshop_sers/goods_srv/proto"
)

var conn *grpc.ClientConn
var GoodClient pb.GoodsClient

func Init() {
	var err error
	conn, err = grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	GoodClient = pb.NewGoodsClient(conn)
}
