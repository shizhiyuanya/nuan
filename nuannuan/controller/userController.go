package controller

import (
	"net/http"
	"nuannuan/database"
	"nuannuan/model/dto"
	"sync"

	"github.com/gin-gonic/gin"
)

var usersMutex sync.Mutex

func Register(ctx *gin.Context) {
	db := database.GetDB()
	// 直接用dto会更安全一点
	var newUser dto.User
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	var existUser dto.User
	usersMutex.Lock()
	defer usersMutex.Unlock()
	db.Where("username = ?", newUser.Username).First(&existUser)
	if existUser.ID != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}
	db.Create(&newUser)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
	})
}

func Login(ctx *gin.Context) {

	db := database.GetDB()
	var loginUser dto.User
	if err := ctx.ShouldBindJSON(&loginUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 开始判断name和password
	var user dto.User
	db.Where("username = ?", loginUser.Username).First(&user)

	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "用户不存在",
		})
		return
	}

	if user.Password != loginUser.Password {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "密码错误",
		})
		return
	}
	//返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"id":      user.ID,
		"success": true,
	})
}
