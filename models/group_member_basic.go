package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
)

type GroupMemberBasic struct {
	gorm.Model
	GroupIdentity  string
	MemberIdentity string // 属于谁的关系
	Type           int    // 群主2，管理员1，群员0
}

func (gmb *GroupMemberBasic) TableName() string {
	return "group_member_basic"
}

func LoadGroups(userId string) ([]GroupMemberBasic, int) {
	gmbList := make([]GroupMemberBasic, 0)
	err := utils.DB.Table("group_member_basic AS gmb").Where("member_identity = ?", userId).
		Joins("join group_basic AS gb on gmb.member_identity = 'register' AND gb.group_identity = gmb.group_identity").Find(&gmbList).Error
	fmt.Println(gmbList)
	if err != nil {
		return gmbList, -1
	}
	return gmbList, 0
}
