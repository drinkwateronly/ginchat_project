package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	UserId        string `valid:"matches(^[a-zA-Z0-9]{6,}$)"` // 账号，数字或字母，6~20位
	Name          string
	Password      string `valid:"matches(^.{6,20}$)"` // 任意字符6~20位
	Salt          string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"` // 电话号码校验
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string
	ClientPort    string
	LoginTime     time.Time
	HeartbeatTime time.Time
	LoginOutTime  time.Time `gorm:"column:login_out_time" json:"login_out_time"`
	IsLoginOut    bool      `gorm:"column:is_login_out" json:"is_login_out"`
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, user := range data {
		fmt.Println(user)
	}
	return data
}

func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(
		UserBasic{
			Name:     user.Name,
			Password: user.Password,
			Email:    user.Email,
			Phone:    user.Phone,
		})
}

func FindUserByUserId(userId string) (*UserBasic, bool) {
	user := UserBasic{}
	rowsAffected := utils.DB.Where("user_id = ?", userId).Find(&user).RowsAffected
	if rowsAffected == 0 {
		return nil, false
	}
	return &user, true
}

func FindUserByEmail(email string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("email = ?", email).First(&user)
}

//func FindUserByNameAndPassword(phone string) *gorm.DB {
//	user := UserBasic{}
//	return utils.DB.Where("email = ", phone).First(&user)
//}
