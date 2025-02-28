package main

import (
	"context"
	"fmt"
	pb "mxshop_sers/goods_srv/proto"
)

func TestGetCategoryBrandList() {
	rsp, err := GoodClient.GetCategoryBrandList(context.Background(), &pb.CategoryInfoRequest{
		Id: 130366,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("tot = ", rsp.Total)
	fmt.Println("data = ", rsp.Data)
}

//func main() {
//	Init()
//	TestGetCategoryBrandList()
//	//TestGetBrand()
//	conn.Close()
//}
