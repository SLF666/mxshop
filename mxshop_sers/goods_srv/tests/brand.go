package main

import (
	"context"
	"fmt"
	pb "mxshop_sers/goods_srv/proto"
)

func TestGetBrand() {
	rsp, err := GoodClient.BrandList(context.Background(), &pb.BrandFilterRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, v := range rsp.Data {
		fmt.Println(v.Name)
	}
}
