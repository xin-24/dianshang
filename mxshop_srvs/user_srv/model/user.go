package model

import(
	"time"

	"gorm.io/gorm"
)

type BaseModel struct{
ID int32 `gorm:"primarykey"`//设置ID为主键
CreatedAt time.Time `gorm:"column:add_time"`//表示时间点指定对应列名 column自定义在数据库中列的名称
UpdatedAt time.Time `gorm:"column:update_time"`
DeletedAt gorm.DeletedAt
IsDeleted bool `gorm:"column:is_deleted"`
}


/*
1.密文 2.不可反解
	1.对称加密（加密解密同一吧钥匙）
	2.非对称加密（加密解密不是同一把钥匙）
	3.md5 信息摘要法
*/
type User struct{
BaseModel
PassWord string `gorm:"type:varchar(100);not null"`
Mobile string `gorm:"index:idx_mobile;unique;type:varchar(11);not null "`
NickName string `gorm:"type:varchar(20)"`
Birthday *time.Time `gorm:"type:datetime"`
Gender string `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女,male表示男'"`
Role int `gorm:"colum:role;default:1;type:int comment '1表示普通用户,2表示管理员'"`
}