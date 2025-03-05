package controller

import (
	"math"
	"net/http"
	"nuannuan/database"
	"nuannuan/model/dto"
	"nuannuan/model/vm"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PushEnglish(ctx *gin.Context) {
	db := database.GetDB()

	var englishRequest vm.Word
	if err := ctx.ShouldBindJSON(&englishRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse JSON"})
		return
	}
	newEnglish := dto.Word{
		English: englishRequest.English,
		Chinese: englishRequest.Chinese,
	}
	db.Create(&newEnglish)
	// 发送响应给客户端
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "English上传成功",
	})
}

func GetEnglish(ctx *gin.Context) {

	db := database.GetDB()
	// 从 URL 查询参数中获取分页参数
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("pageSize")

	// 将字符串转换为整数
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageSize parameter"})
		return
	}

	// 计算偏移量和限制
	offset := (page - 1) * pageSize
	limit := pageSize

	// 定义结果切片
	var results []dto.Word // 假设你有一个名为 English 的结构体
	// var returnResults []vm.Word // 假设你有一个名为 English 的结构体

	// 使用 GORM 的 API 执行分页查询
	err = db.Offset(offset).Limit(limit).Find(&results).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		return
	}

	// 返回查询结果
	ctx.JSON(http.StatusOK, gin.H{"data": results})
}

func GetPage(ctx *gin.Context) {
	db := database.GetDB()

	// 前后联调 要改
	pageSize := 20

	var totalCount int64
	err := db.Model(&dto.Word{}).Count(&totalCount).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count records"})
		return
	}

	// 计算总页数（注意：如果总数不能被页面大小整除，则总页数需要加1）math.Ceil返回大于等于的数
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}
	// 构建响应数据
	response := gin.H{
		"totalPages": totalPages,
	}

	ctx.JSON(http.StatusOK, response)
}
