package api

import (
	"chigua-backend/internal/model"
	"chigua-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(int(model.Unauthorized), model.ErrorResponse(model.Unauthorized))
		return
	}

	var commentCreate model.CommentCreate
	if err := c.ShouldBindJSON(&commentCreate); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	comment, err := service.CreateComment(commentCreate, userID.(int64))
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(comment))
}

func GetCommentsByArticleID(c *gin.Context) {
	articleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	comments, err := service.GetCommentsByArticleID(articleID)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(comments))
}

func DeleteComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(int(model.Unauthorized), model.ErrorResponse(model.Unauthorized))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	err = service.DeleteComment(id, userID.(int64))
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.InvalidParams))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}
