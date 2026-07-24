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

	// 用户路由（公共）
	users := apiGroup.Group("/user")
	{
		users.POST("/register", api.Register)
		users.POST("/login", api.Login)
		users.GET("/me", middleware.AuthMiddleware(), api.GetCurrentUser)
		users.PUT("/me", middleware.AuthMiddleware(), api.UpdateUser)
		users.PUT("/avatar", middleware.AuthMiddleware(), api.UpdateAvatar)
	}

	// 文章路由（公共）
	articles := apiGroup.Group("/article")
	{
		articles.POST("", middleware.AuthMiddleware(), api.CreateArticle)
		articles.GET("", api.GetArticleList)
		articles.GET("/:id", api.GetArticleByID)
		articles.PUT("/:id", middleware.AuthMiddleware(), api.UpdateArticle)
		articles.DELETE("/:id", middleware.AuthMiddleware(), api.DeleteArticle)
		articles.POST("/:id/publish", middleware.AuthMiddleware(), api.PublishArticle)
		articles.PUT("/:id/status", middleware.AuthMiddleware(), api.UpdateArticleStatus)
		articles.POST("/:id/submit", middleware.AuthMiddleware(), api.SubmitForReview)
		articles.GET("/:id/reviews", api.GetArticleReviewRecords)
	}

	// 评论路由
	comments := apiGroup.Group("/comment")
	{
		comments.POST("", middleware.AuthMiddleware(), api.CreateComment)
		comments.GET("/article/:id", api.GetCommentsWithPagination)
		comments.GET("/article/:id/children", api.GetMoreSecondLevelComments)
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

	// Admin 登录路由
	apiGroup.POST("/admin/login", admin.AdminLogin)

	// ====== Admin 路由（权限控制） ======

	// 文章管理路由
	adminArticles := apiGroup.Group("/admin/article", middleware.AuthMiddleware(), middleware.RequireAnyPermission("article:list", "article:create"))
	{
		adminArticles.POST("", middleware.RequirePermission("article:create"), admin.CreateArticle)
		adminArticles.GET("", middleware.RequirePermission("article:list"), admin.GetArticleList)
		adminArticles.GET("/pending", middleware.RequirePermission("article:review"), admin.GetPendingReviewArticles)
		adminArticles.GET("/reviews", middleware.RequirePermission("article:records"), admin.GetReviewRecords)
		adminArticles.GET("/:id", middleware.RequirePermission("article:list"), admin.GetArticleByID)
		adminArticles.PUT("/:id", middleware.RequirePermission("article:edit"), admin.UpdateArticle)
		adminArticles.DELETE("/:id", middleware.RequirePermission("article:delete"), admin.DeleteArticle)
		adminArticles.POST("/:id/publish", middleware.RequirePermission("article:publish"), admin.PublishArticle)
		adminArticles.PUT("/:id/status", middleware.RequirePermission("article:publish"), admin.UpdateArticleStatus)
		adminArticles.POST("/:id/approve", middleware.RequirePermission("article:approve"), admin.ApproveArticle)
		adminArticles.POST("/:id/reject", middleware.RequirePermission("article:reject"), admin.RejectArticle)
		adminArticles.POST("/:id/unpublish", middleware.RequirePermission("article:unpublish"), admin.UnpublishArticle)
	}

	// 分类管理路由
	adminCategories := apiGroup.Group("/admin/category", middleware.AuthMiddleware(), middleware.RequirePermission("category:list"))
	{
		adminCategories.POST("", middleware.RequirePermission("category:create"), admin.CreateCategory)
		adminCategories.GET("", admin.GetAllCategories)
		adminCategories.PUT("/:id", middleware.RequirePermission("category:update"), admin.UpdateCategory)
		adminCategories.DELETE("/:id", middleware.RequirePermission("category:delete"), admin.DeleteCategory)
	}

	// 标签管理路由
	adminTags := apiGroup.Group("/admin/tag", middleware.AuthMiddleware(), middleware.RequirePermission("tag:list"))
	{
		adminTags.POST("", middleware.RequirePermission("tag:create"), admin.CreateTag)
		adminTags.GET("", admin.GetAllTags)
		adminTags.PUT("/:id", middleware.RequirePermission("tag:update"), admin.UpdateTag)
		adminTags.DELETE("/:id", middleware.RequirePermission("tag:delete"), admin.DeleteTag)
	}

	// 评论管理路由
	adminComments := apiGroup.Group("/admin/comment", middleware.AuthMiddleware(), middleware.RequirePermission("comment:list"))
	{
		adminComments.GET("", admin.GetCommentList)
		adminComments.DELETE("/:id", middleware.RequirePermission("comment:delete"), admin.DeleteComment)
	}

	// 用户管理路由
	adminUsers := apiGroup.Group("/admin/user", middleware.AuthMiddleware(), middleware.RequirePermission("user:list"))
	{
		adminUsers.POST("", middleware.RequirePermission("user:create"), admin.CreateUser)
		adminUsers.GET("", admin.GetUserList)
		adminUsers.PUT("/:id", middleware.RequirePermission("user:update"), admin.UpdateUser)
		adminUsers.DELETE("/:id", middleware.RequirePermission("user:delete"), admin.DeleteUser)
		adminUsers.PUT("/:id/roles", middleware.RequirePermission("user:update"), admin.UpdateUserRoles)
	}

	// 角色管理路由
	adminRoles := apiGroup.Group("/admin/role", middleware.AuthMiddleware(), middleware.RequirePermission("role:list"))
	{
		adminRoles.GET("", admin.GetRoles)
		adminRoles.POST("", middleware.RequirePermission("role:manage"), admin.CreateRole)
		adminRoles.PUT("/:id", middleware.RequirePermission("role:manage"), admin.UpdateRole)
		adminRoles.DELETE("/:id", middleware.RequirePermission("role:manage"), admin.DeleteRoleHandler)
		adminRoles.GET("/:id/menus", admin.GetRolePermissionIDsHandler)
		adminRoles.PUT("/:id/menus", middleware.RequirePermission("role:manage"), admin.UpdateRoleMenusHandler)
	}

	// 菜单管理路由
	adminMenus := apiGroup.Group("/admin/menu", middleware.AuthMiddleware())
	{
		adminMenus.GET("/tree", admin.GetPermissionTree)
		adminMenus.GET("", admin.GetAllPermissions)
		adminMenus.POST("", middleware.RequirePermission("role:manage"), admin.CreateMenu)
		adminMenus.PUT("/:id", middleware.RequirePermission("role:manage"), admin.UpdateMenu)
		adminMenus.DELETE("/:id", middleware.RequirePermission("role:manage"), admin.DeleteMenu)
	}

	// 仪表盘统计
	apiGroup.GET("/admin/stats", middleware.AuthMiddleware(), middleware.RequirePermission("dashboard:view"), admin.GetDashboardStats)
}
