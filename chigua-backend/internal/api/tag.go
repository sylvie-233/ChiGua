package api

import (
	"chigua-backend/internal/model"
	"chigua-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateTag(c *gin.Context) {
	var tagCreate model.TagCreate
	if err := c.ShouldBindJSON(&tagCreate); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	tag, err := service.CreateTag(tagCreate)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(tag))
}

func GetAllTags(c *gin.Context) {
	tags, err := service.GetAllTags()
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(tags))
}

func DeleteTag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	err = service.DeleteTag(id)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}
