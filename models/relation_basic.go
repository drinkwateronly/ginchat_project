package models

import (
	"ginchat/utils"
	"gorm.io/gorm"
)

// 人员关系

type RelationBasic struct {
	gorm.Model
	OwnerIdentity  string // 属于谁的关系
	TargetIdentity string // 对应的谁
	Type           int    // 什么类型
}

func (rb *RelationBasic) TableName() string {
	return "relation_basic"
}

func SearchFriends(userId string) []UserBasic {
	relations := make([]RelationBasic, 0)
	friendsInfoList := make([]UserBasic, 0)
	utils.DB.Where("owner_identity = ? and type = 1", userId).Find(&relations)
	for _, relation := range relations {
		ub := UserBasic{}
		utils.DB.Where("user_id = ?", relation.TargetIdentity).Find(&ub)
		friendsInfoList = append(friendsInfoList, ub)
	}
	return friendsInfoList
}

func IsRelationExist(userId, targetId string) bool {
	rb := RelationBasic{}
	row1 := utils.DB.Where("owner_identity = ? AND target_identity = ? AND type = 1", userId, targetId).Find(&rb).RowsAffected
	row2 := utils.DB.Where("owner_identity = ? AND target_identity = ? AND type = 1", targetId, userId).Find(&rb).RowsAffected
	if row1 == 0 && row2 == 0 {
		return false
	} else {
		return true
	}
}

// 好友添加
func AddFriend(userId, targetId string) (string, int) {
	// 用户是否存在？
	_, isExist := FindUserByUserId(targetId)
	if !isExist {
		return "用户不存在", -1
	}
	// 好友关系是否存在？
	if IsRelationExist(userId, targetId) {
		return "关系已存在", -1
	}
	// 添加关系
	rb1 := RelationBasic{
		OwnerIdentity:  userId,
		TargetIdentity: targetId,
		Type:           1,
	}
	rb2 := RelationBasic{
		OwnerIdentity:  targetId,
		TargetIdentity: userId,
		Type:           1,
	}
	// 开启事务
	tx := utils.DB.Begin()
	// 不论遇到什么异常都会rollback()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(&rb1).Error; err != nil {
		tx.Rollback()
		return "事务回滚", -1
	}
	if err := tx.Create(&rb2).Error; err != nil {
		tx.Rollback()
		return "事务回滚", -1
	}
	tx.Commit()
	return "添加成功", 0
}
