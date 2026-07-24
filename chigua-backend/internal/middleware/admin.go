package middleware

import (
	"chigua-backend/internal/model"
	"chigua-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func AdminRoleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(int(model.Forbidden), model.ErrorResponse(model.Forbidden))
			c.Abort()
			return
		}

		user, err := service.GetUserByID(userID.(int64))
		if err != nil {
			c.JSON(int(model.Forbidden), model.ErrorResponse(model.Forbidden))
			c.Abort()
			return
		}

		if user.Role != "admin" {
			c.JSON(int(model.Forbidden), model.ErrorResponse(model.Forbidden))
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminOrReviewerMiddleware 允许 admin 和 reviewer 角色通过
func AdminOrReviewerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(int(model.Forbidden), model.ErrorResponse(model.Forbidden))
			c.Abort()
			return
		}

		user, err := service.GetUserByID(userID.(int64))
		if err != nil {
			c.JSON(int(model.Forbidden), model.ErrorResponse(model.Forbidden))
			c.Abort()
			return
		}

		if user.Role != "admin" && user.Role != "reviewer" {
			c.JSON(int(model.Forbidden), model.ErrorResponse(model.Forbidden))
			c.Abort()
			return
		}

		c.Next()
	}
}