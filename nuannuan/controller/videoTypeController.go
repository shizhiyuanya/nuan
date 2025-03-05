package controller

import (
	"net/http"
	"nuannuan/database"
	"nuannuan/model/dto"
	"nuannuan/model/vm"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var mutex = &sync.Mutex{}

func VideoTypePush(ctx *gin.Context) {

	db := database.GetDB()
	// 确保请求方法是 POST
	if ctx.Request.Method != http.MethodPost {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Invalid request method"})
		return
	}

	var videoTypesRequest vm.VideoTypes
	if err := ctx.ShouldBindJSON(&videoTypesRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse JSON"})
		return
	}

	var wg sync.WaitGroup
	var typeIds = make(chan string, 32) // 带缓冲区的通道
	for _, tmpType := range videoTypesRequest.Types {
		wg.Add(1)
		go func(tmpType string) {
			defer wg.Done()
			mutex.Lock()
			defer mutex.Unlock()

			// 在这里安全地访问数据库
			var typeFind dto.VideoType
			db.Where("type_name = ?", tmpType).First(&typeFind)
			var id string
			if typeFind.ID == 0 {
				// 创建新记录的逻辑
				newType := dto.VideoType{
					TypeName: tmpType,
				}
				db.Create(&newType)

				id = strconv.FormatUint(uint64(newType.ID), 10)
			} else {
				db.Where("type_name = ?", tmpType).First(&typeFind)
				id = strconv.FormatUint(uint64(typeFind.ID), 10)
			}
			// 发送 ID 到通道
			typeIds <- id + ","
		}(tmpType)
	}

	go func() {
		wg.Wait()
		close(typeIds) // 所有 goroutine 完成后关闭通道 这样就能拿数据
	}()
	// 收集所有 ID 并构建最终的字符串
	var finalTypeIds string
	for id := range typeIds {
		finalTypeIds += id
	}
	finalTypeIds = strings.TrimRight(finalTypeIds, ",")

	var videoFind dto.Video

	db.Where("Name = ?", videoTypesRequest.Name).First(&videoFind)
	if videoFind.ID == 0 {
		// 保存到数据库
		newVideo := dto.Video{
			Name:   videoTypesRequest.Name,
			TypeId: finalTypeIds,
		}
		db.Create(&newVideo)
	} else {
		// 如果视频已存在，您可能想更新它或返回错误，取决于您的业务逻辑
		ctx.JSON(http.StatusConflict, gin.H{"error": "Video already exists"})
		return
	}

	// 发送响应给客户端
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "video上传成功",
	})

}
