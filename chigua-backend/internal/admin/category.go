package admin

import (
	"chigua-backend/internal/model"
	"chigua-backend/internal/service"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	var categoryCreate model.CategoryCreate
	if err := c.ShouldBindJSON(&categoryCreate); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	category, err := service.CreateCategory(categoryCreate)
	if err != nil {
		if errors.Is(err, service.ErrCategoryExists) {
			c.JSON(int(model.BadRequest), model.ErrorResponse(model.CategoryExists))
			return
		}
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.InvalidParams))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(category))
}

func GetAllCategories(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	keyword := c.Query("keyword")

	categories, err := service.GetCategoryList(page, pageSize, keyword)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(categories))
}

func DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	err = service.DeleteCategory(id)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}

func UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	var categoryCreate model.CategoryCreate
	if err := c.ShouldBindJSON(&categoryCreate); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	category, err := service.UpdateCategory(id, categoryCreate.Name)
	if err != nil {
		if errors.Is(err, service.ErrCategoryExists) {
			c.JSON(int(model.BadRequest), model.ErrorResponse(model.CategoryExists))
			return
		}
		if errors.Is(err, service.ErrCategoryNotFound) {
			c.JSON(int(model.NotFound), model.ErrorResponse(model.CategoryNotFound))
			return
		}
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(category))
}
