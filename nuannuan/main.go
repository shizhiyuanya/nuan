package main

import (
	"nuannuan/controller"
	"nuannuan/database"
	"nuannuan/io"

	"github.com/gin-gonic/gin"
)

func main() {

	//初始化 MinIO连接
	io.InitMinIO()

	//获取初始化的数据库
	db := database.InitDB()
	//延迟关闭数据库
	defer db.Close()
	//创建一个默认的路由引擎
	r := gin.Default()

	//启动路由
	CollectRoutes(r)

	//9090端口
	panic(r.Run(":9090"))

}

func CollectRoutes(r *gin.Engine) *gin.Engine {

	//注册
	r.POST("/register", controller.Register)
	//登录
	r.POST("/login", controller.Login)
	// 图片上传
	r.POST("/picture/push", controller.PicturePush)
	// tag上传
	r.POST("/tags/push", controller.TagsPush)
	// english上传
	r.POST("/english/push", controller.PushEnglish)
	//video上传
	r.POST("/video/push", controller.VideoPush)
	r.POST("/video/type/push", controller.VideoTypePush)
	// 富文本图片上传
	//r.POST("/upload/image")

	// english获取
	r.GET("/english/get", controller.GetEnglish)
	r.GET("/english/page", controller.GetPage)

	return r
}
