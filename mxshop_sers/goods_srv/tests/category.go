package main

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "mxshop_sers/goods_srv/proto"
)

func TestGetCategotryList() {
	rsp, err := GoodClient.GetAllCategorysList(context.Background(), &emptypb.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println("tot = ", rsp.Total)
	fmt.Println("data = ", rsp.Data)
	fmt.Println(rsp.JsonData)
}

func TestGetSubCategotryList() {
	rsp, err := GoodClient.GetSubCategory(context.Background(), &pb.CategoryListRequest{
		Id: 130358,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("tot = ", rsp.Total)
	fmt.Println("data = ", rsp.SubCategorys)
}

//func main() {
//	Init()
//	TestGetSubCategotryList()
//	//TestGetBrand()
//	conn.Close()
//}
