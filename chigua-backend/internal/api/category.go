package api

import (
	"chigua-backend/internal/model"
	"chigua-backend/internal/service"
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
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(category))
}

func GetAllCategories(c *gin.Context) {
	categories, err := service.GetAllCategories()
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
