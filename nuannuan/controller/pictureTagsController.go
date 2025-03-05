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

func TagsPush(ctx *gin.Context) {

	db := database.GetDB()
	// 确保请求方法是 POST
	if ctx.Request.Method != http.MethodPost {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Invalid request method"})
		return
	}

	// 解析 JSON 请求体
	var tagsRequest vm.PictureTags
	if err := ctx.ShouldBindJSON(&tagsRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse JSON"})
		return
	}

	var wg sync.WaitGroup
	var tagIds = make(chan string, 32) // 带缓冲区的通道

	for _, tmpTag := range tagsRequest.Tags {
		wg.Add(1)
		go func(tmpTag string) {
			defer wg.Done()
			mutex.Lock()
			defer mutex.Unlock()

			// 在这里安全地访问数据库

			var tagFind dto.Tag
			db.Where("tag_name = ?", tmpTag).First(&tagFind)
			var id string
			if tagFind.ID == 0 {
				// 创建新记录的逻辑
				newTag := dto.Tag{
					TagName: tmpTag,
				}
				db.Create(&newTag)
				id = strconv.FormatUint(uint64(newTag.ID), 10)
			} else {
				db.Where("type_name = ?", tmpTag).First(&tmpTag)
				id = strconv.FormatUint(uint64(tagFind.ID), 10)
			}
			// 发送 ID 到通道
			tagIds <- id + ","
		}(tmpTag)
	}

	go func() {
		wg.Wait()
		close(tagIds) // 所有 goroutine 完成后关闭通道 这样就能拿数据
	}()

	// 收集所有 ID 并构建最终的字符串
	var finalTagIds string
	for id := range tagIds {
		finalTagIds += id
	}
	// 去除最后一个逗号（如果需要）
	finalTagIds = strings.TrimRight(finalTagIds, ",")

	var pictureFind dto.Picture

	db.Where("Name = ?", tagsRequest.Name).First(&pictureFind)
	println("看看最后的type_id： ", finalTagIds)
	if pictureFind.ID == 0 {
		newVideo := dto.Video{
			Name:   tagsRequest.Name,
			TypeId: finalTagIds,
		}
		db.Create(&newVideo)
	} else {
		// 如果图片已存在，可能想更新它或返回错误，取决于业务逻辑
		ctx.JSON(http.StatusConflict, gin.H{"error": "Picture already exists"})
		return
	}

	// 发送响应给客户端
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "video上传成功",
	})

}
