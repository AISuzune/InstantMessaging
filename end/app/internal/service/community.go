package service

import (
	g "InstantMessaging/app/global"
	"InstantMessaging/app/internal/model"
	"fmt"
)

// CreateCommunity 建群
func CreateCommunity(community model.Community) (bool, string) {
	tx := g.MysqlDB.Begin()
	//事务一旦开始，不论什么异常最终都会 Rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if len(community.Name) == 0 {
		return false, "群名称不能为空"
	}
	if community.OwnerId == 0 {
		return false, "请先登录"
	}
	if err := g.MysqlDB.Create(&community).Error; err != nil {
		fmt.Println(err)
		tx.Rollback()
		return false, "建群失败"
	}
	contact := model.Contact{}
	contact.OwnerId = community.OwnerId
	contact.TargetId = community.ID
	contact.Type = 2 //群关系
	if err := g.MysqlDB.Create(&contact).Error; err != nil {
		tx.Rollback()
		return false, "添加群关系失败"
	}

	tx.Commit()
	return true, "建群成功"

}

// UpdateCommunity 更新群信息
func UpdateCommunity(community model.Community, comId string) (bool, string) {
	tx := g.MysqlDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查用户是否具有更新群的权限
	ok, _ := g.Enforcer.Enforce(fmt.Sprint(community.OwnerId), comId, "update")

	if !ok {
		return false, "您没有权限更新此群"
	}

	if len(community.Name) == 0 {
		return false, "群名称不能为空"
	}
	if err := tx.Model(&community).Updates(community).Error; err != nil {
		fmt.Println(err)
		tx.Rollback()
		return false, "更新群信息失败"
	}

	tx.Commit()
	return true, "更新群信息成功"
}

// DeleteCommunity 删除群
func DeleteCommunity(userId uint, comId string) (bool, string) {
	community := model.Community{}
	g.MysqlDB.Where("id=? or name=?", comId, comId).Find(&community)

	// 检查用户是否具有删除群的权限
	ok, _ := g.Enforcer.Enforce(fmt.Sprint(userId), comId, "delete")
	if !ok {
		return false, "您没有权限删除此群"
	}
	if community.Name == "" {
		return false, "没有找到群"
	}

	tx := g.MysqlDB.Begin()
	contact := model.Contact{}
	if err := tx.Delete(&contact, "target_id=? and type=2", comId).Error; err != nil {
		tx.Rollback()
		return false, "删除联系人失败"
	}
	if err := tx.Delete(&community).Error; err != nil {
		tx.Rollback()
		return false, "删除群失败"
	}
	tx.Commit()
	return true, "解散群成功"
}

// LoadCommunity 群列表
func LoadCommunity(ownerId uint) []*model.Community {
	contacts := make([]model.Contact, 0)
	objIds := make([]uint64, 0)
	g.MysqlDB.Where("owner_id = ? and type=2", ownerId).Find(&contacts) // 查找所有属于该用户的群列表
	for _, v := range contacts {
		objIds = append(objIds, uint64(v.TargetId)) // 将所有群的 ID 添加到列表中
	}
	communities := make([]*model.Community, 0)
	g.MysqlDB.Where("id in (?)", objIds).Find(&communities)
	return communities
}

// JoinCommunity 加群
func JoinCommunity(userId uint, comId string) (bool, string) {
	contact := model.Contact{}
	contact.OwnerId = userId
	contact.Type = 2
	community := model.Community{}

	g.MysqlDB.Where("id=? or name=?", comId, comId).Find(&community)
	if community.Name == "" {
		return false, "没有找到群"
	}
	g.MysqlDB.Where("owner_id=? and target_id=? and type =2 ", userId, comId).Find(&contact)
	if !contact.CreatedAt.IsZero() {
		return false, "已加过此群"
	} else {
		contact.TargetId = community.ID
		g.MysqlDB.Create(&contact)
		return true, "加群成功"
	}
}

// LeaveCommunity 退群
func LeaveCommunity(userId uint, comId string) (bool, string) {
	contact := model.Contact{}
	community := model.Community{}

	g.MysqlDB.Where("id=? or name=?", comId, comId).Find(&community)
	if community.Name == "" {
		return false, "没有找到群"
	}
	g.MysqlDB.Where("owner_id=? and target_id=? and type =2 ", userId, comId).Find(&contact)
	if contact.CreatedAt.IsZero() {
		return false, "您不在此群中"
	} else {
		g.MysqlDB.Delete(&contact)
		return true, "退出群成功"
	}

}
