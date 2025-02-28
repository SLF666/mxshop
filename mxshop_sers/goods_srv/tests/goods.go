package main

import (
	"context"
	"fmt"
	pb "mxshop_sers/goods_srv/proto"
)

func TestGoodsList() {
	rsp, err := GoodClient.GetGoodsDetail(context.Background(), &pb.GoodInfoRequest{
		Id: 421,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("tot = ", rsp.Name)
}

func main() {
	Init()
	TestGoodsList()
	conn.Close()
}
