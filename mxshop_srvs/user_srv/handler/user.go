package main

import (
	"context"

	"gorm.io/gorm"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"

	"github.com/xin-24/go/mxshop_srvs/user_srv/global"
	"github.com/xin-24/go/mxshop_srvs/user_srv/model"
	"github.com/xin-24/go/mxshop_srvs/user_srv/proto"
	
)

type UserServer struct{}

// 分页
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// gorm作用使用go代码方式操作数据库，不用手写复杂的sql语句类似智能助手
func ModelToRsponse(user model.User) proto.UserInfoResponse {
	//在grpc的message中字段有默认值，不能随便赋值nil进去，容易出错。
	//这里面需要明白，哪些字段是有默认值的。
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		PassWord: user.PassWord,//改了proto还是错了？？？//成功改了model/user.go中的PassWord
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday!=nil{
		userInfoRsp.BirthDay=uint64(user.Birthday.Unix())//Unix 时间戳 是从 1970年1月1日00:00:00 UTC 开始的秒数。	
	}
	return userInfoRsp
}

//用户列表
func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	//获取用户列表
	//gorm作用使用go代码方式操作数据库，不用手写复杂的sql语句类似智能助手
	//可以取数据库的数据了
	var users []model.User //切片表结构
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)
	for _, user := range users {
		userInfoRsp:=ModelToRsponse(user)
		rsp.Data=append(rsp.Data, &userInfoRsp)//go语法限制不允许改成一行，更安全。
	}
	return rsp,nil
}

//通过手机号码查询用户
func (s *UserServer)GetUserByMobile(ctx context.Context,req *proto.MobileRequest) (*proto.UserInfoResponse, error){
	//通过手机号码查询用户
	var user model.User
	result:=global.DB.Where(&model.User{Mobile:req.Mobile}).First(&user)
	if result.RowsAffected==0{
		return nil,status.Errorf(codes.NotFound,"用户不存在")//返回grpc状态码
	}
	if result.Error!=nil{
		return nil,result.Error
	}
	userInfoRsp:=ModelToRsponse(user)
	return &userInfoRsp,nil

}

//通过Id查询用户
func (s *UserServer)GetUserById(ctx context.Context,req *proto.IdRequest) (*proto.UserInfoResponse,error){//通过id查询用户{
//通过Id查询用户，Id为主键	
    var user model.User
	result:=global.DB.First(&user,req.Id)
	if result.RowsAffected==0{
		return nil,status.Errorf(codes.NotFound,"用户不存在")//返回grpc状态码
	}
	if result.Error!=nil{
		return nil,result.Error
	}
	userInfoRsp:=ModelToRsponse(user)
	return &userInfoRsp,nil

}

//创建用户
func (s *UserServer)CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error){
//新建用户
//新建前需要查询
var user model.User//如果不想加上前缀可在import中调用包的前面加上"."。
result:=global.DB.Where(&model.User{Mobile:req.Mobile}).First(&user)
if result.RowsAffected==1{
	return nil,status.Error(codes.AlreadyExists,"用户已存在")
}
user.Mobile=req.Mobile
user.NickName=req.NickName//初始化昵称为前端的

//密码加密
//测试
}