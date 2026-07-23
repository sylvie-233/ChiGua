package admin

import (
	"chigua-backend/internal/model"
	"chigua-backend/internal/service"
	"errors"
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
		if errors.Is(err, service.ErrTagExists) {
			c.JSON(int(model.BadRequest), model.ErrorResponse(model.TagExists))
			return
		}
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.InvalidParams))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(tag))
}

func GetAllTags(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	keyword := c.Query("keyword")

	tags, err := service.GetTagList(page, pageSize, keyword)
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

func UpdateTag(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	var tagCreate model.TagCreate
	if err := c.ShouldBindJSON(&tagCreate); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	tag, err := service.UpdateTag(id, tagCreate.Name)
	if err != nil {
		if errors.Is(err, service.ErrTagExists) {
			c.JSON(int(model.BadRequest), model.ErrorResponse(model.TagExists))
			return
		}
		if errors.Is(err, service.ErrTagNotFound) {
			c.JSON(int(model.NotFound), model.ErrorResponse(model.TagNotFound))
			return
		}
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(tag))
}
