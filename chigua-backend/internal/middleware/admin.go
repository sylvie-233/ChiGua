package middleware

import (
	"chigua-backend/internal/model"

	"github.com/gin-gonic/gin"
)

// RequirePermission 检查 JWT claims 中是否包含指定权限（无需查数据库）
func RequirePermission(permCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, exists := c.Get("permissions")
		if !exists {
			c.JSON(int(model.Forbidden), model.ErrorResponse(model.Forbidden))
			c.Abort()
			return
		}

		perms, ok := permissions.([]string)
		if !ok {
			c.JSON(int(model.Forbidden), model.ErrorResponse(model.Forbidden))
			c.Abort()
			return
		}

		for _, p := range perms {
			if p == permCode {
				c.Next()
				return
			}
		}

		c.JSON(int(model.Forbidden), model.ErrorResponse(model.Forbidden))
		c.Abort()
	}
}

// RequireAnyPermission 检查是否包含任一权限
func RequireAnyPermission(permCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, exists := c.Get("permissions")
		if !exists {
			c.JSON(int(model.Forbidden), model.ErrorResponse(model.Forbidden))
			c.Abort()
			return
		}

		perms, ok := permissions.([]string)
		if !ok {
			c.JSON(int(model.Forbidden), model.ErrorResponse(model.Forbidden))
			c.Abort()
			return
		}

		for _, p := range perms {
			for _, required := range permCodes {
				if p == required {
					c.Next()
					return
				}
			}
		}

		c.JSON(int(model.Forbidden), model.ErrorResponse(model.Forbidden))
		c.Abort()
	}
}
