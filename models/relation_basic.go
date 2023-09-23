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
