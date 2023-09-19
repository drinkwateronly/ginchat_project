package models

import "gorm.io/gorm"

type GroupBasic struct {
	gorm.Model
	GroupId   string
	GroupName string
	OwnerId   string // 属于谁的关系
	Icon      string // 对应的谁
	Type      int    // 什么类型
}

func (rb *GroupBasic) TableName() string {
	return "group_basic"
}
