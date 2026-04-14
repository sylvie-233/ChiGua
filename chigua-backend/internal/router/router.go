package router

import (
	api "chigua-backend/internal/api"
	"chigua-backend/internal/middleware"
	"chigua-backend/internal/model"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	apiGroup := r.Group("/api")

	apiGroup.GET("/ping", func(c *gin.Context) {
		c.JSON(int(model.Success), model.SuccessResponse(map[string]string{
			"message": "Pong",
		}))
	})

	// 用户路由
	users := apiGroup.Group("/user")
	{
		users.POST("/register", api.Register)
		users.POST("/login", api.Login)
		users.GET("/me", middleware.AuthMiddleware(), api.GetCurrentUser)
	}

	// 文章路由
	articles := apiGroup.Group("/article")
	{
		articles.POST("", middleware.AuthMiddleware(), api.CreateArticle)
		articles.GET("", api.GetArticleList)
		articles.GET("/:id", api.GetArticleByID)
		articles.PUT("/:id", middleware.AuthMiddleware(), api.UpdateArticle)
		articles.DELETE("/:id", middleware.AuthMiddleware(), api.DeleteArticle)
		articles.POST("/:id/publish", middleware.AuthMiddleware(), api.PublishArticle)
	}

	// 评论路由
	comments := apiGroup.Group("/comment")
	{
		comments.POST("", middleware.AuthMiddleware(), api.CreateComment)
		comments.GET("/article/:id", api.GetCommentsByArticleID)
		comments.DELETE("/:id", middleware.AuthMiddleware(), api.DeleteComment)
	}

	// 分类路由
	categories := apiGroup.Group("/categorie")
	{
		categories.POST("", middleware.AuthMiddleware(), api.CreateCategory)
		categories.GET("", api.GetAllCategories)
		categories.DELETE("/:id", middleware.AuthMiddleware(), api.DeleteCategory)
	}

	// 标签路由
	tags := apiGroup.Group("/tag")
	{
		tags.POST("", middleware.AuthMiddleware(), api.CreateTag)
		tags.GET("", api.GetAllTags)
		tags.DELETE("/:id", middleware.AuthMiddleware(), api.DeleteTag)
	}
}
