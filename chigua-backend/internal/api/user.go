package api

import (
	"chigua-backend/internal/model"
	"chigua-backend/internal/service"
	"chigua-backend/utils/jwt"
	"errors"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var userRegister model.UserRegister
	if err := c.ShouldBindJSON(&userRegister); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	user, err := service.RegisterUser(userRegister)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			c.JSON(int(model.BadRequest), model.ErrorResponse(model.UserExists))
			return
		}
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.InvalidParams))
		return
	}

	// 新用户默认分配 user 角色
	userRole, _ := service.GetRoleByCode("user")
	if userRole != nil {
		service.AssignUserRole(user.ID, userRole.ID)
	}

	c.JSON(int(model.Success), model.SuccessResponse(user.ToResponse()))
}

func Login(c *gin.Context) {
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

	// 生成token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(gin.H{"token": token, "user": user}))
}

func GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(int(model.Unauthorized), model.ErrorResponse(model.Unauthorized))
		return
	}

	user, err := service.GetUserByID(userID.(int64))
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	resp := user.ToResponse()
	// 填充角色和权限
	roles, _ := service.GetUserRoles(user.ID)
	resp.Roles = roles

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
	resp.Permissions = permissions

	c.JSON(int(model.Success), model.SuccessResponse(resp))
}

func UpdateUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(int(model.Unauthorized), model.ErrorResponse(model.Unauthorized))
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.InvalidParams))
		return
	}

	if err := service.UpdateUser(userID.(int64), updateData); err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse("更新成功"))
}

// UpdateAvatar 更新用户头像
func UpdateAvatar(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(int(model.Unauthorized), model.ErrorResponse(model.Unauthorized))
		return
	}

	var body struct {
		Avatar string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.InvalidParams))
		return
	}

	if err := service.UpdateUserAvatar(userID.(int64), body.Avatar); err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}
