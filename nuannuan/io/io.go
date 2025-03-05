package io

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

// 本地连接
var MinioClient *minio.Client

func InitMinIO() {

	choose := 2

	if choose == 1 {
		LocalInitMinIO()

	} else if choose == 2 {
		RemoteInitMinIO()
	}

}

func LocalInitMinIO() {

	// TODO:填写MinIO内容
	accessKeyID := ""
	secretAccessKey := ""

	minioClient, err := minio.New("127.0.0.1:9000", &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		// 是否启用https 帮加头
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}
	MinioClient = minioClient
	log.Println(minioClient)
	// 使用客户端列出所有存储桶的名称。
	buckets, err := minioClient.ListBuckets(context.Background())
	log.Println("show buckets")
	for _, bucket := range buckets {
		fmt.Println(bucket.Name)
	}
	isExist, err := minioClient.BucketExists(context.Background(), "try")
	if isExist {
		fmt.Printf("%s exists!\n", "try")
	} else {
		fmt.Printf("%s not exists!\n", "try")
	}
}

func RemoteInitMinIO() {

	// TODO:填写MinIO内容
	accessKeyID := ""
	secretAccessKey := ""
	minioClient, err := minio.New("IP:端口", &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		// 是否启用https 帮加头
		Secure: false,
	})

	if err != nil {
		log.Fatalln(err)
	}
	MinioClient = minioClient
	log.Println(minioClient)
	// 使用客户端列出所有存储桶的名称。
	buckets, err := minioClient.ListBuckets(context.Background())
	log.Println("show buckets")
	for _, bucket := range buckets {
		fmt.Println(bucket.Name)
	}
}
