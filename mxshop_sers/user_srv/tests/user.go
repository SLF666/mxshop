package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "mxshop_sers/user_srv/proto"
)

var conn *grpc.ClientConn
var userClient pb.UserClient

func Init() {
	var err error
	conn, err = grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	userClient = pb.NewUserClient(conn)
}

func TestGetUserList() {
	rep, err := userClient.GetUserList(context.Background(), &pb.PageInfo{
		Pn:    1,
		PSize: 3,
	})
	if err != nil {
		panic(err)
	}
	for _, v := range rep.Data {
		fmt.Println(v)
		check, err := userClient.CheckPassWord(context.Background(), &pb.PasswordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: v.Password,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(check.Success)
	}
}

func main() {
	Init()
	TestGetUserList()
	conn.Close()
}
