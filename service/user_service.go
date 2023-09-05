package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/gin-gonic/gin"
	"math/rand"
	"regexp"
	"strconv"
)

// GetUserList
// @Tags 首页
// @Success 200 {string} json{"code",message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"message": data,
	})
}

// UserRegister
// @Summary 用户注册
// @Tags 用户模块
// @param account formData string false "账号"
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param rePassword formData string false "确认密码"
// @param email formData string false "邮箱"
// @Success 200 {string} json{"code",message"}
// @Router /user/register [post]
func UserRegister(c *gin.Context) {
	// todo: 校验邮箱验证码，
	user := models.UserBasic{}
	// 校验账号密码格式
	account := c.PostForm("account")
	isMatch, _ := regexp.MatchString("^[a-zA-Z0-9]{6,}$", account)
	if isMatch == false {
		c.JSON(400, gin.H{
			"message": "账号格式不符",
		})
		return
	}
	password := c.PostForm("password")
	rePassword := c.PostForm("rePassword")
	if password != rePassword {
		c.JSON(-1, gin.H{
			"message": "两次密码不一致",
		})
		return
	}
	isMatch, _ = regexp.MatchString("^.{6,20}$", password)
	if isMatch == false {
		c.JSON(400, gin.H{
			"message": "密码格式不符",
		})
		return
	}

	// 检查账号是否存在
	_, isExist := models.FindUserByAccount(account)
	if isExist {
		c.JSON(400, gin.H{
			"message": "账号已存在",
		})
		return
	}
	// 通过校验，开始新建账号
	user.Account = account
	user.Name = c.PostForm("name")
	user.Email = c.PostForm("email")
	// Go 的 math/rand 包默认使用的伪随机数生成器是线性同余生成器（Linear Congruential Generator，LCG），它的随机性可能不足够强大
	salt := fmt.Sprintf("%06d", rand.Int31()) // todo 更换其他
	user.Salt = salt
	user.Password = utils.MakePassword(password, salt) // 密码不要明文存储

	models.CreateUser(user) // TODO: 错误处理
	c.JSON(200, gin.H{
		"message": "注册成功",
	})
}

// UserLogin
// @Summary 用户登录
// @Tags 用户模块
// @param account formData string id "账号"
// @param password formData string id "密码"
// @Success 200 {string} json{"code",message"}
// @Router /user/userLogin [post]
func UserLogin(c *gin.Context) {
	account := c.PostForm("account")
	password := c.PostForm("password")
	userBasic, isExist := models.FindUserByAccount(account)
	if !isExist {
		c.JSON(400, gin.H{
			"message": "用户不存在",
		})
	}
	isPass := utils.ValidatePassword(password, userBasic.Salt, userBasic.Password)
	if !isPass {
		c.JSON(400, gin.H{
			"message": "密码错误", // 应该不提示什么不正确
		})
		return
	}
	jwtoken, err := utils.GenerateJWT(userBasic.Account, userBasic.Name)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    -1,
			"message": "登陆失败", // 应该不提示什么不正确
		})
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "登录成功",
		"token":   jwtoken,
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string id "用户名id"
// @Success 200 {string} json{"code",message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user) // TODO: 错误处理
	c.JSON(200, gin.H{
		"message": "删除成功",
	})
}

// UpdateUser
// @Summary 用户信息修改
// @Tags 用户模块
// @param id formData string false "用户id"
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param email formData string false "邮箱"
// @param phone formData string false "电话"
// @Success 200 {string} json{"code",message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Email = c.PostForm("email")
	user.Phone = c.PostForm("phone")

	//_, err := govalidator.ValidateStruct(user)
	//if err != nil {
	//	fmt.Println(err)
	//	c.JSON(400, gin.H{
	//		"message": "格式出错",
	//	})
	//	return
	//}
	models.UpdateUser(user) // TODO: 错误处理
	c.JSON(200, gin.H{
		"message": "修改成功",
	})
}
