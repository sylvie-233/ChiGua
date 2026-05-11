package service

import (
	"context"
	"mime/multipart"

	"chigua-backend/config"
	"chigua-backend/internal/model"
	minioUtil "chigua-backend/utils/minio"

	"github.com/minio/minio-go/v7"
)

// UploadFile 上传文件到MinIO（根据文件类型存储到不同目录）
func UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (*model.UploadFileResponse, error) {
	contentType := fileHeader.Header.Get("Content-Type")
	objectPath := minioUtil.GetObjectPath(fileHeader.Filename, contentType)

	cfg := config.AppConfig.MinIO

	_, err := minioUtil.MinIOClient.PutObject(
		context.Background(),
		cfg.BucketName,
		objectPath,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return nil, err
	}

	return &model.UploadFileResponse{
		FileName: fileHeader.Filename,
		FilePath: objectPath,
		FileURL:  minioUtil.BuildFileURL(objectPath),
	}, nil
}

// UploadFiles 批量上传文件
func UploadFiles(files []*multipart.FileHeader) ([]*model.UploadFileResponse, error) {
	var results []*model.UploadFileResponse

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}

		result, err := UploadFile(file, fileHeader)
		if err != nil {
			file.Close()
			return nil, err
		}

		results = append(results, result)
		file.Close()
	}

	return results, nil
}

// GetFileURL 获取文件访问URL
func GetFileURL(filePath string) (string, error) {
	return minioUtil.BuildFileURL(filePath), nil
}

// DeleteFile 删除文件
func DeleteFile(filePath string) error {
	cfg := config.AppConfig.MinIO

	err := minioUtil.MinIOClient.RemoveObject(
		context.Background(),
		cfg.BucketName,
		filePath,
		minio.RemoveObjectOptions{},
	)
	return err
}
