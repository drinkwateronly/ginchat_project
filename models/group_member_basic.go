package models

import (
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

func LoadGroups(userId string) ([]GroupBasic, int) {
	gbl := make([]GroupBasic, 0)
	err := utils.DB.Table("group_basic AS gb").
		Where("member_identity = ?", userId).
		Joins("join group_member_basic AS gmb ON gmb.member_identity = ? AND gb.group_identity = gmb.group_identity", userId).
		Find(&gbl).Error
	//fmt.Println(userId, gmbList)
	if err != nil {
		return gbl, -1
	}
	return gbl, 0
}

func JoinGroup(userId string, groupId string) (code int, msg string) {
	// 先查找群聊是否存在？
	if !IsGroupExist(groupId) {
		return -1, "该群聊不存在"
	}
	if IsGroupJoined(userId, groupId) {
		return -1, "该群聊已加入"
	}
	gmb := GroupMemberBasic{
		GroupIdentity:  groupId,
		MemberIdentity: userId,
		Type:           0,
	}
	if err := utils.DB.Create(&gmb).Error; err != nil {
		return -1, "加入群聊失败"
	}
	return 0, "加入群聊成功"
}

func FindMembersByGroupId(groupId string) []GroupMemberBasic {
	var gmbList []GroupMemberBasic
	utils.DB.Where("group_identity = ?", groupId).Find(&gmbList)

	return gmbList
}
