package router

import (
	"ginchat/docs" // 通过swag init生成的
	"ginchat/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "" //  Swagger 文档将位于根路径下
	/*
		ginSwagger.WrapHandler 用于将 Swagger 文档的请求处理程序包装成一个 Gin 处理程序，以便将其与 Gin 框架集成。
		swaggerfiles.Handler 是 Swagger 自动生成的文档处理程序，它用于提供 Swagger 文档的静态文件（HTML、CSS、JavaScript等）和 API 规范。
	*/
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/index", service.GetIndex)
	r.GET("/user/getUserList", service.GetUserList)
	r.POST("/user/register", service.UserRegister)
	r.GET("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("user/userLogin", service.UserLogin)
	return r
}
