package minio

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"chigua-backend/config"
	"chigua-backend/utils/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOClient MinIO客户端实例
var MinIOClient *minio.Client

// InitMinIO 初始化MinIO客户端
func InitMinIO() error {
	cfg := config.AppConfig.MinIO

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return err
	}

	MinIOClient = client

	// 检查并创建bucket
	err = createBucketIfNotExists(cfg.BucketName)
	if err != nil {
		return err
	}

	logger.Info("MinIO客户端初始化成功")
	return nil
}

// createBucketIfNotExists 如果bucket不存在则创建
func createBucketIfNotExists(bucketName string) error {
	exists, err := MinIOClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = MinIOClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		logger.Infof("创建Bucket: %s", bucketName)

		// 设置bucket为公开可读
		policyStr := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": "*",
					"Action": [
						"s3:GetObject",
						"s3:ListBucket"
					],
					"Resource": [
						"arn:aws:s3:::` + bucketName + `",
						"arn:aws:s3:::` + bucketName + `/*"
					]
				}
			]
		}`
		err = MinIOClient.SetBucketPolicy(context.Background(), bucketName, policyStr)
		if err != nil {
			return err
		}
		logger.Infof("Bucket %s 已设置为公开权限", bucketName)
	}

	return nil
}

// GetFileType 获取文件类型（image, video, audio, default）
func GetFileType(contentType string) string {
	if strings.HasPrefix(contentType, "image/") {
		return "image"
	}
	if strings.HasPrefix(contentType, "video/") {
		return "video"
	}
	if strings.HasPrefix(contentType, "audio/") {
		return "audio"
	}
	return "default"
}

// GetFileTypePath 根据文件类型获取存储目录
func GetFileTypePath(contentType string) string {
	fileType := GetFileType(contentType)
	if path, ok := config.AppConfig.MinIO.FileTypePaths[fileType]; ok {
		return path
	}
	if path, ok := config.AppConfig.MinIO.FileTypePaths["default"]; ok {
		return path
	}
	return "others"
}

// BuildFileURL 构建文件的公开访问URL
func BuildFileURL(objectPath string) string {
	cfg := config.AppConfig.MinIO
	protocol := "http"
	if cfg.UseSSL {
		protocol = "https"
	}
	endpoint := cfg.PublicEndpoint
	if endpoint == "" {
		endpoint = cfg.Endpoint
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, endpoint, cfg.BucketName, objectPath)
}

// GetObjectPath 生成对象存储路径（格式：目录/原文件名_时间戳.扩展名）
func GetObjectPath(fileName, contentType string) string {
	dir := GetFileTypePath(contentType)
	ext := filepath.Ext(fileName)
	baseName := strings.TrimSuffix(fileName, ext)
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s/%s_%s%s", dir, sanitizeFileName(baseName), timestamp, ext)
}

// sanitizeFileName 清理文件名，移除特殊字符
func sanitizeFileName(name string) string {
	// 移除文件名中的特殊字符，只保留字母、数字、中文和下划线
	result := strings.Builder{}
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') ||
			(r >= '\u4e00' && r <= '\u9fa5') || r == '_' || r == '-' {
			result.WriteRune(r)
		} else {
			result.WriteRune('_')
		}
	}
	return result.String()
}
