package model

// DeleteFileRequest 删除文件请求
type DeleteFileRequest struct {
	FilePath string `json:"filePath" binding:"required"`
}

// UploadFileResponse 文件上传响应
type UploadFileResponse struct {
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
	FileURL  string `json:"fileUrl"`
}
