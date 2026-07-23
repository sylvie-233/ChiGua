package admin

import (
	"chigua-backend/internal/model"
	"chigua-backend/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

func GetDashboardStats(c *gin.Context) {
	var articleTotal, articleDraft, articlePublished, articleUnpublished, articlePending, categoryTotal, tagTotal, commentTotal, userTotal int64
	var err error

	if articleDraft, err = service.CountArticleByStatus(int(model.ArticleStatusDraft)); err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}
	if articlePublished, err = service.CountArticleByStatus(int(model.ArticleStatusPublished)); err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}
	if articleUnpublished, err = service.CountArticleByStatus(int(model.ArticleStatusUnpublished)); err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}
	if articlePending, err = service.CountArticleByStatus(int(model.ArticleStatusPending)); err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}
	articleTotal = articleDraft + articlePublished + articleUnpublished + articlePending

	if categoryTotal, err = service.CountCategory(); err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}
	if tagTotal, err = service.CountTag(); err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}
	if commentTotal, err = service.CountComment(); err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}
	if userTotal, err = service.CountUser(); err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	recentArticles, err := service.GetRecentArticles(5)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	startDate := time.Now().AddDate(-1, 0, 0)
	articleDailyStats, err := service.GetArticleCountByDate(startDate)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(gin.H{
		"articles": gin.H{
			"total":       articleTotal,
			"draft":       articleDraft,
			"published":   articlePublished,
			"unpublished": articleUnpublished,
			"pending":     articlePending,
		},
		"categories":        categoryTotal,
		"tags":              tagTotal,
		"comments":          commentTotal,
		"users":             userTotal,
		"recentArticles":    recentArticles,
		"articleDailyStats": articleDailyStats,
	}))
}
