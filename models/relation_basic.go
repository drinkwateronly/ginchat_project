package models

import "gorm.io/gorm"

// 人员关系

type RelationBasic struct {
	gorm.Model
	OwnerId  string // 属于谁的关系
	TargetId string // 对应的谁
	Type     int    // 什么类型
}

func (rb *RelationBasic) TableName() string {
	return "relation_basic"
}
