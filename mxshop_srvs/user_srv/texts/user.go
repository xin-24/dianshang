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

//测试创建用户
func TestCreakUser(){
	for i:=0;i<10;i++{
		rsp,err:=userClient.CreateUser(context.Background(),&proto.CreateUserInfo{
			NickName: fmt.Sprintf("bobby%d", i),
			Mobile:   fmt.Sprintf("1461874881%d", i),
			PassWord: "admin123",
		})
		if err!=nil{
			panic(err)
		}
		fmt.Println(rsp.Id)
	}
}
//测试手机号查询
func TestGetUserByMobile1(){
	for i:=0;i<30;i++{
	rsp,err:=userClient.GetUserMobile(context.Background(),&proto. MobileRequest{
		Mobile:"14618748810",
	})
	if err!=nil{
		panic(err)
	}
	if rsp!=nil{
		fmt.Println(rsp.Id)
	}else{
		fmt.Println("该号码不存在")
	}
	}
}
func TestGetUserByMobile() {
    var rsp *proto.UserInfoResponse
    var err error
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        rsp, err = userClient.GetUserMobile(context.Background(), &proto.MobileRequest{
            Mobile: "14618748810",
        })
        if err == nil {
            break
        }
        fmt.Printf("Retry %d: %v\n", i+1, err)
        time.Sleep(2 * time.Second)
    }
    if err != nil {
        panic(err)
    }
    if rsp != nil {
        fmt.Println(rsp.Id)
    } else {
        fmt.Println("该号码不存在")
    }
}
func main() {
	Init()

	// TestGetUserList()
    // TestCreakUser()
	// TestGetUserByMobile()//可能服务器资源耗尽运行之后服务器关闭了？？？
	conn.Close()
}
