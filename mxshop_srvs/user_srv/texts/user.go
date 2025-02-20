package main

import (
	"context"
	"fmt"
	"time"

	"github.com/xin-24/go/mxshop_srvs/user_srv/proto"
	"google.golang.org/grpc"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("0.0.0.0:50051", grpc.WithInsecure(), grpc.WithTimeout(10*time.Second))
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

// 测试列表
func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 5,
	})
	if err != nil {
		fmt.Printf("failed to get user list: %v\n", err)
		return
	}
	for _, user := range rsp.Data {
		fmt.Println(user.Mobile, user.NickName, user.PassWord)
		checkRsp, err := userClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			PassWord:          "admin123",
			EncryptedPassword: user.PassWord,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checkRsp.Success)
	}

}
func main() {
	Init()
	TestGetUserList()
	conn.Close()
}
