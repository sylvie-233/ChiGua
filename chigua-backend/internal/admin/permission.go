package admin

import (
	"chigua-backend/internal/model"
	"chigua-backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetRoles 角色列表（分页+搜索）
func GetRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 10 }

	roles, total, err := service.GetRoleList(page, pageSize, keyword)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	c.JSON(int(model.Success), model.SuccessResponse(gin.H{
		"items":      roles,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": totalPages,
		"hasNext":    page < totalPages,
		"hasPrev":    page > 1,
	}))
}

// CreateRole 新增角色
func CreateRole(c *gin.Context) {
	var body struct {
		Code        string `json:"code" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	role, err := service.CreateRole(body.Code, body.Name, body.Description)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}
	c.JSON(int(model.Success), model.SuccessResponse(role))
}

// UpdateRole 更新角色信息
func UpdateRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	var body struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	if err := service.UpdateRoleInfo(id, body.Name, body.Description); err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}
	c.JSON(int(model.Success), model.SuccessResponse(nil))
}

// DeleteRoleHandler 删除角色
func DeleteRoleHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	if err := service.DeleteRole(id); err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}
	c.JSON(int(model.Success), model.SuccessResponse(nil))
}

// GetRolePermissionIDsHandler 查询角色拥有的菜单ID
func GetRolePermissionIDsHandler(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	ids, err := service.GetRoleMenuIDs(roleID)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(ids))
}

// UpdateRoleMenusHandler 更新角色菜单权限
func UpdateRoleMenusHandler(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	var body struct {
		MenuIDs []int64 `json:"menuIds"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	err = service.UpdateRoleMenus(roleID, body.MenuIDs)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}

// UpdateUserRoles 更新用户角色
func UpdateUserRoles(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	var body struct {
		RoleIDs []int64 `json:"roleIds"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	err = service.SetUserRoles(userID, body.RoleIDs)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}
