package models

import (
	"ginchat/utils"
	"gorm.io/gorm"
)

type GroupBasic struct {
	gorm.Model
	GroupIdentity string
	GroupName     string
	OwnerIdentity string // 属于谁的关系
	Icon          string // 对应的谁
	Type          int    // 什么类型
	Desc          string //
}

func (rb *GroupBasic) TableName() string {
	return "group_basic"
}

func CreateGroup(gb GroupBasic) (int, string) {
	if len(gb.GroupName) == 0 {
		return -1, "请输入群昵称"
	}
	if gb.OwnerIdentity == "" {
		return -1, "用户未登录"
	}
	ub := UserBasic{}
	if row := utils.DB.Where("user_id = ?", gb.OwnerIdentity).Find(&ub).RowsAffected; row == 0 {
		return -1, "用户不存在"
	}
	gb.GroupIdentity = utils.MakeGroupId()
	for utils.DB.Where("group_identity = ?", gb.GroupIdentity).Find(&GroupBasic{}).RowsAffected != 0 {
		gb.GroupIdentity = utils.MakeGroupId()
	}
	// 事务开启
	tx := utils.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 建群
	if err := tx.Create(&gb).Error; err != nil {
		return -1, "建群出错"
		tx.Rollback()
	}
	// 群主入群
	gmb := GroupMemberBasic{
		GroupIdentity:  gb.GroupIdentity,
		MemberIdentity: gb.OwnerIdentity,
		Type:           2, // 群主 2
	}
	if err := tx.Create(&gmb).Error; err != nil {
		return -1, "建群出错"
		tx.Rollback()
	}
	tx.Commit()
	return 0, "建群成功"
}

// IsGroupExist 群聊是否存在
func IsGroupExist(groupIdentity string) bool {
	count := utils.DB.Where("group_identity = ?", groupIdentity).Find(&GroupBasic{}).RowsAffected
	if count == 0 {
		return false
	}
	return true
}

// IsGroupJoined 用户是否已加入某群聊，并不判断群聊是否存在：如果群聊不存在，用户自然不在该群聊
func IsGroupJoined(groupIdentity string, memberIdentity string) bool {
	count := utils.DB.Where("group_identity = ? AND member_identity = ?", groupIdentity, memberIdentity).Find(&GroupMemberBasic{}).RowsAffected
	if count == 0 {
		return false
	}
	return true
}
