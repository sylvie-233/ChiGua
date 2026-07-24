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

	if user.Role != "admin" && user.Role != "reviewer" {
		c.JSON(int(model.Forbidden), model.ErrorResponse(model.Forbidden))
		return
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(gin.H{"token": token, "user": user}))
}

func CreateUser(c *gin.Context) {
	var userCreate struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&userCreate); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	if userCreate.Role != "admin" && userCreate.Role != "user" && userCreate.Role != "reviewer" {
		userCreate.Role = "user"
	}

	user, err := service.CreateUser(userCreate.Username, userCreate.Password, userCreate.Nickname, userCreate.Role)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			c.JSON(int(model.BadRequest), model.ErrorResponse(model.UserExists))
			return
		}
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(user.ToResponse()))
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	var userUpdate struct {
		Nickname string `json:"nickname"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	if userUpdate.Role != "" && userUpdate.Role != "admin" && userUpdate.Role != "user" && userUpdate.Role != "reviewer" {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	err = service.AdminUpdateUser(id, userUpdate.Nickname, userUpdate.Role)
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

func UpdateUserRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	var roleUpdate struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&roleUpdate); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	if roleUpdate.Role != "admin" && roleUpdate.Role != "user" && roleUpdate.Role != "reviewer" {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	err = service.UpdateUserRole(id, roleUpdate.Role)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}