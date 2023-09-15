package service

import (
	g "InstantMessaging/app/global"
	"InstantMessaging/app/internal/model"
)

// LoadFriend 好友列表
func LoadFriend(userId uint) []*model.User {
	contacts := make([]model.Contact, 0)
	objIds := make([]uint64, 0)
	// 在数据库中查找所有属于该用户且类型为1（表示好友）的联系人
	g.MysqlDB.Where("owner_id = ? and type=1", userId).Find(&contacts)
	// 遍历联系人列表，获取每个联系人的目标ID（即好友的ID）
	for _, v := range contacts {
		objIds = append(objIds, uint64(v.TargetId))
	}
	users := make([]*model.User, 0)
	g.MysqlDB.Where("id in (?)", objIds).Find(&users)
	return users
}

// AddFriend 添加好友   自己的ID，好友名
func AddFriend(userId uint, targetName string) (bool, string) {

	if targetName != "" {
		targetUser := CheckUserIsExist(targetName)
		//fmt.Println(targetUser, " userId        ", )
		if targetUser.Salt != "" {
			if targetUser.ID == userId {
				return false, "不能加自己"
			}
			contact0 := model.Contact{}
			g.MysqlDB.Where("owner_id =? and target_id =? and type=1", userId, targetUser.ID).Find(&contact0)
			if contact0.ID != 0 {
				return false, "不能重复添加"
			}
			tx := g.MysqlDB.Begin()
			//事务一旦开始，不论什么异常最终都会 Rollback
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()
			contact := model.Contact{}
			contact.OwnerId = userId
			contact.TargetId = targetUser.ID
			contact.Type = 1
			if err := g.MysqlDB.Create(&contact).Error; err != nil {
				tx.Rollback() // 删除联系人失败，回滚事务
				return false, "添加好友失败"
			}
			contact1 := model.Contact{}
			contact1.OwnerId = targetUser.ID
			contact1.TargetId = userId
			contact1.Type = 1
			if err := g.MysqlDB.Create(&contact1).Error; err != nil {
				tx.Rollback()
				return false, "添加好友失败"
			}
			tx.Commit()
			return true, "添加好友成功"
		}
		return false, "没有找到此用户"
	}
	return false, "好友ID不能为空"
}

// DeleteFriend 删除好友
func DeleteFriend(userId uint, targetName string) (bool, string) {
	if targetName != "" {
		targetUser := CheckUserIsExist(targetName)
		if targetUser.Salt != "" {
			contact := model.Contact{}
			g.MysqlDB.Where("owner_id = ? and target_id = ? and type = 1", userId, targetUser.ID).Find(&contact)
			if contact.ID == 0 {
				return false, "你们不是好友关系"
			}

			tx := g.MysqlDB.Begin()
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()

			if err := g.MysqlDB.Delete(&contact).Error; err != nil {
				tx.Rollback() // 删除联系人失败，回滚事务
				return false, "删除好友失败"
			}

			// 同时删除对方好友关系
			contactReverse := model.Contact{}
			g.MysqlDB.Where("owner_id = ? and target_id = ? and type = 1", targetUser.ID, userId).Find(&contactReverse)
			if contactReverse.ID != 0 {
				if err := g.MysqlDB.Delete(&contactReverse).Error; err != nil {
					tx.Rollback()
					return false, "删除好友失败"
				}
			}

			tx.Commit()
			return true, "删除好友成功"
		}
		return false, "没有找到此用户"
	}
	return false, "好友ID不能为空"
}

func SearchUserByGroupId(communityId uint) []uint {
	contacts := make([]model.Contact, 0)
	objIds := make([]uint, 0)
	g.MysqlDB.Where("target_id = ? and type=2", communityId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, v.OwnerId)
	}
	return objIds
}
