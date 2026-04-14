package api

import (
	"chigua-backend/internal/model"
	"chigua-backend/internal/service"
	"chigua-backend/utils/jwt"

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
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.InvalidParams))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(user))
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

	c.JSON(int(model.Success), model.SuccessResponse(user))
}
