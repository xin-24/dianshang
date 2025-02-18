// 连接数库
package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/anaskhan96/go-password-encoder"
)

func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}
func main() {
	// //数据表连接信息

	// // 定义数据库连接字符串（DSN，Data Source Name）
	// // 包含了数据库的用户名（root）、密码（root）、主机地址（192.168.1.10）、端口（3306）、数据库名（mxshop_user_srv）
	// // 以及一些参数，如字符集（utf8mb4）、是否解析时间（parseTime=True）、时区（loc=Local）
	// dsn := "root:root@tcp(172.23.210.216:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

	// //配置GORM日志
	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags),
	// 	logger.Config{
	// 		SlowThreshold: time.Second,
	// 		LogLevel:      logger.Info,
	// 		Colorful:      true,
	// 	},
	// )
	// //打开数据库连接
	// //全局模式
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	// 	NamingStrategy: schema.NamingStrategy{

	// 		SingularTable: true,
	// 	},
	// 	Logger: newLogger,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// //自动迁移数据表
	// _ = db.AutoMigrate(&model.User{})
	// fmt.Println(genMd5("123406"))
	// e10adc3949ba59abbe56e057f20f883e

	//加盐md5 推荐使用ssh
	
	//不需要了，可以以后理解

	// salt, encodedPwd := password.Encode("generic password", nil)
	// fmt.Println(salt)
	// fmt.Println(encodedPwd)
	// check := password.Verify("generic password", salt, encodedPwd, nil)
	// fmt.Println(check) // true

	// Using custom options
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode("generic password", options)
	newPassword:=fmt.Sprintf("$pbkdf2-sha512$%s$%s",salt,encodedPwd)
	fmt.Println(len(newPassword))
	fmt.Println(newPassword)

	passwordInfo:=strings.Split(newPassword,"$")
	fmt.Println(passwordInfo)
	check := password.Verify("generic password", passwordInfo[1], passwordInfo[2],options)
	fmt.Println(check) // true

}
