package api

import (
	"chigua-backend/internal/model"
	"chigua-backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UploadFile 上传单个文件
func UploadFile(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.InvalidParams))
		return
	}

	defer file.Close()

	result, err := service.UploadFile(file, fileHeader)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(result))
}

// UploadFiles 批量上传文件
func UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.InvalidParams))
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.InvalidParams))
		return
	}

	results, err := service.UploadFiles(files)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(results))
}

// GetFile 获取文件（重定向到MinIO）
func GetFile(c *gin.Context) {
	filePath := c.Param("filepath")
	if filePath == "" {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.InvalidParams))
		return
	}

	url, err := service.GetFileURL(filePath)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.Redirect(http.StatusFound, url)
}

// DeleteFile 删除文件
func DeleteFile(c *gin.Context) {
	var req model.DeleteFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.InvalidParams))
		return
	}

	if req.FilePath == "" {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.InvalidParams))
		return
	}

	err := service.DeleteFile(req.FilePath)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}
