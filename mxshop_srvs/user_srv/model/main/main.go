// 连接数库和数据库同步好
package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/xin-24/go/mxshop_srvs/user_srv/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}
func main() {
	//数据表连接信息

	// 定义数据库连接字符串（DSN，Data Source Name）
	// 包含了数据库的用户名（root）、密码（root）、主机地址（192.168.1.10）、端口（3306）、数据库名（mxshop_user_srv）
	// 以及一些参数，如字符集（utf8mb4）、是否解析时间（parseTime=True）、时区（loc=Local）
	dsn := "root:root@tcp(192.168.129.149:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

	//配置GORM日志
	newLogger := logger.New( //创建日记记录器实例
		//标准输出格式打印到控制台，确保每一行日志都在新的一行前面，日志的标识格式
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{ //用于定义日志格式
			SlowThreshold: time.Second, //定义慢查询阈值
			LogLevel:      logger.Info, //定义日志级别：显示所有日志
			Colorful:      true,        //日志会以彩色颜色显示（如慢查询会显示）
		},
	)
	//打开数据库连接
	//全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{ //gorm的配置结构体，用于定义全局行为
		//这是命名策略
		NamingStrategy: schema.NamingStrategy{ //默认蛇形命名

			SingularTable: true, //是否会使用单数表名
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode("admine123", options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	fmt.Println(newPassword)

	//创建测试用户
	for i := 0; i < 10; i++ {

		user := model.User{
			NickName: fmt.Sprintf("bobby%d", i),
			Mobile:   fmt.Sprintf("1461874881%d", i),
			PassWord: newPassword,
		}
		db.Save(&user)
	}

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
	// options := &password.Options{16, 100, 32, sha512.New}
	// salt, encodedPwd := password.Encode("generic password", options)
	// newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	// fmt.Println(len(newPassword))
	// fmt.Println(newPassword)

	// passwordInfo := strings.Split(newPassword, "$")
	// fmt.Println(passwordInfo)
	// check := password.Verify("generic password", passwordInfo[1], passwordInfo[2], options)
	// fmt.Println(check) // true

}
