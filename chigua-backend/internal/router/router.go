package router

import (
	api "chigua-backend/internal/api"
	admin "chigua-backend/internal/admin"
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
		users.PUT("/me", middleware.AuthMiddleware(), api.UpdateUser)
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
		articles.PUT("/:id/status", middleware.AuthMiddleware(), api.UpdateArticleStatus)
	}

	// 评论路由
	comments := apiGroup.Group("/comment")
	{
		comments.POST("", middleware.AuthMiddleware(), api.CreateComment)
		comments.GET("/article/:id", api.GetCommentsWithPagination)           // 两级评论分页列表
		comments.GET("/article/:id/children", api.GetMoreSecondLevelComments) // 获取更多二级评论
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

	// 文件上传路由
	upload := apiGroup.Group("/upload")
	{
		upload.POST("/file", middleware.AuthMiddleware(), api.UploadFile)
		upload.POST("/files", middleware.AuthMiddleware(), api.UploadFiles)
		upload.GET("/*filepath", api.GetFile)
		upload.DELETE("/file", middleware.AuthMiddleware(), api.DeleteFile)
	}

	// Admin 登录路由（不经过认证中间件）
	apiGroup.POST("/admin/login", admin.AdminLogin)

	// Admin 路由组
	adminGroup := apiGroup.Group("/admin", middleware.AuthMiddleware(), middleware.AdminRoleMiddleware())
	{
		// Admin 文章路由
		adminArticles := adminGroup.Group("/article")
		{
			adminArticles.POST("", admin.CreateArticle)
			adminArticles.GET("", admin.GetArticleList)
			adminArticles.GET("/:id", admin.GetArticleByID)
			adminArticles.PUT("/:id", admin.UpdateArticle)
			adminArticles.DELETE("/:id", admin.DeleteArticle)
			adminArticles.POST("/:id/publish", admin.PublishArticle)
			adminArticles.PUT("/:id/status", admin.UpdateArticleStatus)
		}

		// Admin 分类路由
		adminCategories := adminGroup.Group("/category")
		{
			adminCategories.POST("", admin.CreateCategory)
			adminCategories.GET("", admin.GetAllCategories)
			adminCategories.PUT("/:id", admin.UpdateCategory)
			adminCategories.DELETE("/:id", admin.DeleteCategory)
		}

		// Admin 标签路由
		adminTags := adminGroup.Group("/tag")
		{
			adminTags.POST("", admin.CreateTag)
			adminTags.GET("", admin.GetAllTags)
			adminTags.PUT("/:id", admin.UpdateTag)
			adminTags.DELETE("/:id", admin.DeleteTag)
		}

		// Admin 评论路由
		adminComments := adminGroup.Group("/comment")
		{
			adminComments.GET("", admin.GetCommentList)
			adminComments.DELETE("/:id", admin.DeleteComment)
		}

		// Admin 用户路由
		adminUsers := adminGroup.Group("/user")
		{
			adminUsers.POST("", admin.CreateUser)
			adminUsers.GET("", admin.GetUserList)
			adminUsers.PUT("/:id", admin.UpdateUser)
			adminUsers.DELETE("/:id", admin.DeleteUser)
			adminUsers.PUT("/:id/role", admin.UpdateUserRole)
		}

		// Admin 仪表盘统计路由
		adminGroup.GET("/stats", admin.GetDashboardStats)
	}
}
