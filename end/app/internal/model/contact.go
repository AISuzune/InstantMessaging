package model

import "github.com/jinzhu/gorm"

// Contact 人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint   //谁的关系信息
	TargetId uint   //对应的谁 /群 ID
	Type     int    //对应的类型  1好友  2群
	Desc     string //描述信息
}

func (table *Contact) TableName() string {
	return "contact"
}
