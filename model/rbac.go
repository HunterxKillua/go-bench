package models

import (
	"gorm.io/gorm"
)

type Tenant struct {
	ID   int    `json:"id" gorm:"column:id;primary_key;auto_increment;size:11;"`
	Name string `json:"name" gorm:"column:name;varchar(255);" form:"name"`
	Logo string `json:"logo" gorm:"column:logo;type:varchar(255);" form:"logo"`
}

type Admin struct {
	gorm.Model
	Password string `json:"password" gorm:"column:password;type:varchar(32);"`
	Phone    string `json:"phone" gorm:"column:phone;type:varchar(11);"`
	Name     string `json:"name" gorm:"column:name;type:varchar(64);"`
	Nick     string `json:"nick" gorm:"column:nick;type:varchar(64);"`
	Email    string `json:"email" gorm:"column:email;type:varchar(255);"`
	Roll     string `json:"roll" gorm:"column:roll;type:varchar(16);"`
	Sex      bool   `json:"sex" gorm:"column:sex;type:tinyint(1);"`
	Avatar   string `json:"avatar" gorm:"column:avatar;type:varchar(255);"`
}

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"column:name;type:varchar(64);"`
	Phone    string `json:"phone" gorm:"column:phone;type:varchar(11);"`
	NickName string `json:"nickname" gorm:"column:nickname;type:varchar(64);"`
	Sex      bool   `json:"sex" gorm:"column:sex;type:tinyint(1);"`
	Gh       string `json:"ghid" gorm:"column:sex;type:varchar(255);"`
	App      string `json:"appid" gorm:"column:appid;type:varchar(255);"`
	Union    string `json:"unionid" gorm:"column:unionid;type:varchar(255);"`
}

type Relations struct {
	ID       int `json:"id"`
	TenantID int `json:"tenant_id"`
	UserID   int
	AdminID  int
}

func init() {
	Models["tenant"] = &Tenant{}
	Models["admin"] = &Admin{}
	Models["user"] = &User{}
	Models["t-user-admin"] = &Relations{}
}

// gorm自定义表名
func (Tenant) TableName() string {
	return "tenant"
}

func (Admin) TableName() string {
	return "admin"
}

func (User) TableName() string {
	return "user"
}

func (Relations) TableName() string {
	return "t-user-admin"
}
 