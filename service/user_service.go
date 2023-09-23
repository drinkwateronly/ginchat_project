package service

import (
	"encoding/json"
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"
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
// @param userId formData string false "账号"
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
	userId := c.PostForm("userId")
	isMatch, _ := regexp.MatchString("^[a-zA-Z0-9]{6,}$", userId)
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
			"code":    -1,
			"message": "两次密码不一致",
		})
		return
	}
	isMatch, _ = regexp.MatchString("^.{6,20}$", password)
	if isMatch == false {
		c.JSON(400, gin.H{
			"code":    -1,
			"message": "密码格式不符",
		})
		return
	}

	// 检查账号是否存在
	_, isExist := models.FindUserByUserId(userId)
	if isExist {
		c.JSON(400, gin.H{
			"code":    -1,
			"message": "账号已存在",
		})
		return
	}
	// 通过校验，开始新建账号
	user.UserId = userId
	user.Name = c.PostForm("name")
	user.Email = c.PostForm("email")
	marshal, err := json.Marshal(user)
	if err != nil {
		return
	}
	user.Identity = utils.Md5Encode(string(marshal))
	// Go 的 math/rand 包默认使用的伪随机数生成器是线性同余生成器（Linear Congruential Generator，LCG），它的随机性可能不足够强大
	salt := fmt.Sprintf("%06d", rand.Int31()) // todo 更换其他
	user.Salt = salt
	user.Password = utils.MakePassword(password, salt) // 密码不要明文存储

	ret := models.CreateUser(user) // TODO: 错误处理
	if ret.Error != nil {
		fmt.Println(ret.Error)
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "注册失败",
		})
	}
	c.JSON(200, gin.H{
		"code":    0,
		"message": "注册成功",
	})
}

// UserLogin
// @Summary 用户登录
// @Tags 用户模块
// @param userId formData string id "账号"
// @param password formData string id "密码"
// @Success 200 {string} json{"code",message"}
// @Router /user/userLogin [post]
func UserLogin(c *gin.Context) {
	userId := c.PostForm("userId")
	password := c.PostForm("password")
	if userId == "" || password == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "用户名或密码不能为空",
		})
	}
	userBasic, isExist := models.FindUserByUserId(userId)
	fmt.Println(userId, password, isExist)
	if !isExist {
		c.JSON(400, gin.H{
			"code":    -1,
			"message": "用户不存在", // 实际上是用户不存在
		})
		return
	}
	isPass := utils.ValidatePassword(password, userBasic.Salt, userBasic.Password)
	if !isPass {
		c.JSON(400, gin.H{
			"code":    -1,
			"message": "用户名或密码错误", // 实际上是密码错误，但不应该提示是什么不正确
		})
		return
	}
	jwtoken, err := utils.GenerateJWT(userBasic.UserId, userBasic.Name)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    -1,
			"message": "登陆失败", // 应该不提示什么不正确
		})
	}

	c.JSON(200, gin.H{
		"code":    0,
		"message": "登录成功",
		"data":    userBasic,
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

// 用于升级HTTP连接到WebSocket连接的实例
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 它总是返回 true，表示接受任何来源的连接
	},
}

// SendMsg 仅发送一次数据，就断开websocket连接
func SendMsg(c *gin.Context) {
	// 升级，返回一个websocket连接
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// 升级失败
		fmt.Println(err)
		return
	}
	// websocket连接需要被释放
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	// 处理数据向websocket发送
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	// 订阅一个名为"123"的redis频道,并从中获取消息
	message, err := utils.Subscribe(c, "123")
	if err != nil {
		fmt.Println(err)
	}
	t := time.Now().Format("2006-01-02  15:04:05")
	m := fmt.Sprintf("[ginchat][%s]:[%s]", t, message)
	// 向websocket发送处理好的数据
	err = ws.WriteMessage(1, []byte(m))
	if err != nil {
		fmt.Println(err)
	}
}

func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func SearchFriends(c *gin.Context) {
	// 查找用户列表
	friendsInfoList := models.SearchFriends(c.PostForm("userId"))
	utils.RespOKList(c.Writer, friendsInfoList, len(friendsInfoList))
}

func UserSearch(c *gin.Context) {
	ub, isExist := models.FindUserByUserId(c.PostForm("userId"))
	if isExist {
		utils.RespFail(c.Writer, "用户不存在")
	}
	utils.RespOK(c.Writer, ub, "")
}
