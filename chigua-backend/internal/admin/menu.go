package admin

import (
	"chigua-backend/internal/model"
	"chigua-backend/internal/service"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPermissionTree 获取当前用户的菜单树
func GetPermissionTree(c *gin.Context) {
	permissions, _ := c.Get("permissions")
	perms, _ := permissions.([]string)

	tree, err := service.GetPermissionTree(perms)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(tree))
}

// GetAllPermissions 获取所有菜单（管理用）
func GetAllPermissions(c *gin.Context) {
	menus, err := service.GetAllPermissions()
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(menus))
}

// CreateMenu 创建菜单
func CreateMenu(c *gin.Context) {
	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	result, err := service.CreateMenu(menu)
	if err != nil {
		if errors.Is(err, service.ErrPermissionCodeRequired) || errors.Is(err, service.ErrPermissionTitleRequired) || errors.Is(err, service.ErrPermissionCodeExists) {
			c.JSON(int(model.BadRequest), model.ErrorResponseWithMsg(model.BadRequest, err.Error()))
		} else {
			c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		}
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(result))
}

// UpdateMenu 更新菜单
func UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	var menu model.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	err = service.UpdateMenu(id, menu)
	if err != nil {
		if errors.Is(err, service.ErrPermissionCodeRequired) || errors.Is(err, service.ErrPermissionTitleRequired) || errors.Is(err, service.ErrPermissionCodeExists) {
			c.JSON(int(model.BadRequest), model.ErrorResponseWithMsg(model.BadRequest, err.Error()))
		} else {
			c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		}
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}

// DeleteMenu 删除菜单
func DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	err = service.DeleteMenu(id)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}
