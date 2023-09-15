package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username      string
	Password      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Avatar        string //头像
	Identity      string //唯一标识
	ClientIP      string //设备
	ClientPort    string
	Salt          string
	LoginTime     time.Time //登入时间
	HeartbeatTime time.Time //心跳时间
	LoginOutTime  time.Time `gorm:"column:login_out_time" json:"login_out_time"` //下线时间
	IsLogout      bool
	DeviceInfo    string //设备信息
}

func (table *User) TableName() string {
	return "user_subject"
}
