package admin

import (
	"chigua-backend/internal/model"
	"chigua-backend/internal/service"
	"chigua-backend/utils/jwt"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AdminLogin(c *gin.Context) {
	var userLogin model.UserLogin
	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	user, err := service.LoginUser(userLogin)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.PasswordError))
		return
	}

	// 获取权限：admin 拥有全部权限
	roles, _ := service.GetUserRoles(user.ID)
	var permissions []string
	isAdmin := false
	for _, r := range roles {
		if r == "admin" {
			isAdmin = true
			break
		}
	}
	if isAdmin {
		permissions, _ = service.GetAdminPermissions()
	} else {
		permissions, _ = service.GetUserPermissions(user.ID)
	}

	token, err := jwt.GenerateTokenWithPermissions(user.ID, permissions)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	user.Roles = roles
	user.Permissions = permissions

	c.JSON(int(model.Success), model.SuccessResponse(gin.H{"token": token, "user": user}))
}

func CreateUser(c *gin.Context) {
	var userCreate struct {
		Username string  `json:"username" binding:"required"`
		Password string  `json:"password" binding:"required"`
		Nickname string  `json:"nickname"`
		Avatar   string  `json:"avatar"`
		RoleIDs  []int64 `json:"roleIds"`
	}
	if err := c.ShouldBindJSON(&userCreate); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	user, err := service.CreateUser(userCreate.Username, userCreate.Password, userCreate.Nickname, userCreate.Avatar)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			c.JSON(int(model.BadRequest), model.ErrorResponse(model.UserExists))
			return
		}
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	// 分配角色（默认 user）
	if len(userCreate.RoleIDs) > 0 {
		service.SetUserRoles(user.ID, userCreate.RoleIDs)
	} else {
		userRole, _ := service.GetRoleByCode("user")
		if userRole != nil {
			service.AssignUserRole(user.ID, userRole.ID)
		}
	}

	c.JSON(int(model.Success), model.SuccessResponse(model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		CreatedAt: user.CreatedAt,
		UpdateAt:  user.UpdateAt,
	}))
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	var userUpdate struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	err = service.AdminUpdateUser(id, userUpdate.Nickname, userUpdate.Avatar)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}

func GetUserList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	keyword := c.Query("keyword")

	users, err := service.GetUserList(page, pageSize, keyword)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	// 为每个用户填充角色和权限
	for i := range users.Items {
		roles, _ := service.GetUserRoles(users.Items[i].ID)
		users.Items[i].Roles = roles
	}

	c.JSON(int(model.Success), model.SuccessResponse(users))
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	err = service.DeleteUserByID(id)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}
