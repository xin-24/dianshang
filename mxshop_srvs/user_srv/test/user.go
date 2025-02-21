package main

import (
	"context"
	"fmt"
	"time"

	"github.com/xin-24/go/user_srv/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	ka := keepalive.ClientParameters{
		Time:                10 * time.Second, // 每10秒发送一次ping
		Timeout:             time.Second,      // 等待1秒钟以确认ping的响应
		PermitWithoutStream: true,             // 即使没有活动的流也发送ping
	}
	conn, err = grpc.Dial("0.0.0.0:50051", grpc.WithInsecure(), grpc.WithKeepaliveParams(ka))
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

// 测试列表
func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 10,
	})
	if err != nil {
		fmt.Printf("failed to get user list: %v\n", err)
		return
	}
	for _, user := range rsp.Data {
		fmt.Println(user.Mobile, user.NickName, user.PassWord, user.BirthDay)
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

// 测试创建用户
func TestCreakUser() {
	for i := 0; i < 10; i++ {
		rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			NickName: fmt.Sprintf("bobby%d", i),
			Mobile:   fmt.Sprintf("1461874881%d", i),
			PassWord: "admin123",
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(rsp.Id)
	}
}

// 测试手机号查询
func TestGetUserByMobile() {
	rsp, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: "14618748819",
	})
	if err != nil {
		fmt.Println("查询用户时发生错误:", err)

	}
	if rsp == nil {
		fmt.Println("未找到该手机号对应的用户")
	}
	fmt.Println(rsp.Id, rsp.NickName)
}

// 通过id查询
func TestGetUserById() {
	rsp, err := userClient.GetUserById(context.Background(), &proto.IdRequest{
		Id: 30,
	})
	if err != nil {
		fmt.Println("查询用户时发生错误:", err)
		return
	}
	if rsp == nil {
		fmt.Println("未找到该Id对应的用户")
		return
	}
	fmt.Println(rsp.Id, rsp.NickName, rsp.Mobile, rsp.PassWord)
}

// 更新用户
func TestUpdateUser() {
	rsp, err := userClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       30,
		NickName: "xxx",
		Gender:   "famle",
		BirthDay: uint64(time.Now().Unix()),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("更新成功", rsp)

}

// func TestGetUserByMobile(){

//		rsp,err:=userClient.GetUserMobile(context.Background(),&proto. MobileRequest{
//			Mobile:"14618748810",
//		})
//		if err!=nil{
//			fmt.Println("该号码不存在")
//		}
//			fmt.Println(rsp.Id,rsp.NickName,rsp.Mobile)
//	}
//
//	func TestGetUserByMobile() {
//	    var rsp *proto.UserInfoResponse
//	    var err error
//	    maxRetries := 3
//	    for i := 0; i < maxRetries; i++ {
//	        rsp, err = userClient.GetUserMobile(context.Background(), &proto.MobileRequest{
//	            Mobile: "14618748810",
//	        })
//	        if err == nil {
//	            break
//	        }
//	        fmt.Printf("Retry %d: %v\n", i+1, err)
//	        time.Sleep(2 * time.Second)
//	    }
//	    if err != nil {
//	        panic(err)
//	    }
//	    if rsp != nil {
//	        fmt.Println(rsp.Id)
//	    } else {
//	        fmt.Println("该号码不存在")
//	    }
//	}
func main() {
	Init()

	//  TestGetUserList()
	// TestCreakUser()
	TestGetUserByMobile()
	// TestGetUserById()
	// TestUpdateUser()
	conn.Close()
}
