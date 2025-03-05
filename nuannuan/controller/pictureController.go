package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"net/http"
	"nuannuan/io"
)

func PicturePush(ctx *gin.Context) {

	// 解析多部分表单数据
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无法解析表单数据",
		})
		return
	}

	// 从表单中获取文件
	fileHeaders := form.File["file"]
	if len(fileHeaders) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "没有上传文件",
		})
		return
	}

	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无法打开文件",
			})
			return
		}
		defer file.Close()
		name := fileHeader.Filename

		bucketName := "nuannuan" // TODO:替换为您的桶名
		objectName := name
		_, err3 := io.MinioClient.PutObject(
			ctx,             // 上下文
			bucketName,      // 存储桶名称
			objectName,      // 对象名称
			file,            // 文件数据（实现了 io.Reader 接口）
			fileHeader.Size, // 数据大小
			minio.PutObjectOptions{
				ContentType: fileHeader.Header.Get("Content-Type"), // 内容类型
			},
		)

		if err3 != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "文件上传失败",
			})
			return
		}

	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件上传成功",
	})

}
