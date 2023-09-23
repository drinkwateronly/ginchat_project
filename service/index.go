package service

import (
	"ginchat/models"
	"github.com/gin-gonic/gin"
	"text/template"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "index")
	// c.JSON(200, gin.H{
	// 	"message": "welcome !!  ",
	// })
}

func ToRegister(c *gin.Context) {
	ind, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	ind.Execute(c.Writer, "index")
	// c.JSON(200, gin.H{
	// 	"message": "welcome !!  ",
	// })
}

func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles(
		"views/chat/index.html",
		"views/chat/head.html",
		"views/chat/tabmenu.html",
		"views/chat/foot.html",
		"views/chat/group.html",
		"views/chat/userinfo.html",
		"views/chat/concat.html",
		"views/chat/profile.html",
		"views/chat/createcom.html",
		"views/chat/main.html")
	if err != nil {
		panic(err)
	}
	//c.Query("account")
	//c.Query("token")
	ub := models.UserBasic{
		UserId:   c.Query("userId"),
		Identity: c.Query("token"),
	}
	ind.Execute(c.Writer, ub)
}
