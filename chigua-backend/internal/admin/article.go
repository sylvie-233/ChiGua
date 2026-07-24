package admin

import (
	"chigua-backend/internal/model"
	"chigua-backend/internal/service"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateArticle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(int(model.Unauthorized), model.ErrorResponse(model.Unauthorized))
		return
	}

	var articleCreate model.ArticleCreate
	if err := c.ShouldBindJSON(&articleCreate); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	article, err := service.CreateArticle(articleCreate, userID.(int64))
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(article))
}

func GetArticleList(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	keyword := c.Query("keyword")

	articles, err := service.GetAllArticleList(page, pageSize, keyword)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(articles))
}

func GetArticleByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	article, err := service.GetArticleByID(id)
	if err != nil {
		c.JSON(int(model.NotFound), model.ErrorResponse(model.ArticleNotFound))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(article))
}

func UpdateArticle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(int(model.Unauthorized), model.ErrorResponse(model.Unauthorized))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	var articleUpdate model.ArticleUpdate
	if err := c.ShouldBindJSON(&articleUpdate); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	if articleUpdate.Title == nil && articleUpdate.Content == nil && articleUpdate.CoverImage == nil && articleUpdate.CategoryID == nil && articleUpdate.TagIDs == nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.InvalidParams))
		return
	}

	article, err := service.UpdateArticle(id, articleUpdate, userID.(int64))
	if err != nil {
		if errors.Is(err, service.ErrArticleNoPermission) {
			c.JSON(int(model.Forbidden), model.ErrorResponse(model.Forbidden))
		} else {
			c.JSON(int(model.NotFound), model.ErrorResponse(model.ArticleNotFound))
		}
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(article))
}

func DeleteArticle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(int(model.Unauthorized), model.ErrorResponse(model.Unauthorized))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	err = service.DeleteArticle(id, userID.(int64))
	if err != nil {
		if errors.Is(err, service.ErrArticleNoPermission) {
			c.JSON(int(model.Forbidden), model.ErrorResponse(model.Forbidden))
		} else {
			c.JSON(int(model.NotFound), model.ErrorResponse(model.ArticleNotFound))
		}
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}

func PublishArticle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(int(model.Unauthorized), model.ErrorResponse(model.Unauthorized))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	err = service.PublishArticle(id, userID.(int64))
	if err != nil {
		if errors.Is(err, service.ErrArticleNoPermission) {
			c.JSON(int(model.Forbidden), model.ErrorResponse(model.Forbidden))
		} else if errors.Is(err, service.ErrArticleNoTitle) ||
			errors.Is(err, service.ErrArticleNoContent) ||
			errors.Is(err, service.ErrArticleNoCategory) ||
			errors.Is(err, service.ErrArticleNoCoverImage) {
			c.JSON(int(model.BadRequest), model.ErrorResponseWithMsg(model.BadRequest, err.Error()))
		} else {
			c.JSON(int(model.NotFound), model.ErrorResponse(model.ArticleNotFound))
		}
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}

func UpdateArticleStatus(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(int(model.Unauthorized), model.ErrorResponse(model.Unauthorized))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	var statusUpdate model.ArticleStatusUpdate
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	if statusUpdate.Status < 0 || statusUpdate.Status > 3 {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	err = service.UpdateArticleStatus(id, userID.(int64), statusUpdate.Status)
	if err != nil {
		if errors.Is(err, service.ErrArticleNoPermission) {
			c.JSON(int(model.Forbidden), model.ErrorResponse(model.Forbidden))
		} else if errors.Is(err, service.ErrArticleNoTitle) ||
			errors.Is(err, service.ErrArticleNoContent) ||
			errors.Is(err, service.ErrArticleNoCategory) ||
			errors.Is(err, service.ErrArticleNoCoverImage) {
			c.JSON(int(model.BadRequest), model.ErrorResponseWithMsg(model.BadRequest, err.Error()))
		} else if errors.Is(err, service.ErrArticleNotPending) {
			c.JSON(int(model.BadRequest), model.ErrorResponseWithMsg(model.BadRequest, err.Error()))
		} else {
			c.JSON(int(model.NotFound), model.ErrorResponse(model.ArticleNotFound))
		}
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}

// ApproveArticle 审核通过文章
func ApproveArticle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(int(model.Unauthorized), model.ErrorResponse(model.Unauthorized))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	err = service.ApproveArticle(id, userID.(int64))
	if err != nil {
		if errors.Is(err, service.ErrArticleNotPending) {
			c.JSON(int(model.BadRequest), model.ErrorResponseWithMsg(model.BadRequest, err.Error()))
		} else if errors.Is(err, service.ErrArticleNoTitle) ||
			errors.Is(err, service.ErrArticleNoContent) ||
			errors.Is(err, service.ErrArticleNoCategory) ||
			errors.Is(err, service.ErrArticleNoCoverImage) {
			c.JSON(int(model.BadRequest), model.ErrorResponseWithMsg(model.BadRequest, err.Error()))
		} else {
			c.JSON(int(model.NotFound), model.ErrorResponse(model.ArticleNotFound))
		}
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}

// RejectArticle 驳回文章
func RejectArticle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(int(model.Unauthorized), model.ErrorResponse(model.Unauthorized))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	var reviewAction model.ReviewAction
	if err := c.ShouldBindJSON(&reviewAction); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	err = service.RejectArticle(id, userID.(int64), reviewAction.Comment)
	if err != nil {
		if errors.Is(err, service.ErrArticleNotPending) {
			c.JSON(int(model.BadRequest), model.ErrorResponseWithMsg(model.BadRequest, err.Error()))
		} else {
			c.JSON(int(model.NotFound), model.ErrorResponse(model.ArticleNotFound))
		}
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}

// UnpublishArticle 下架文章
func UnpublishArticle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(int(model.Unauthorized), model.ErrorResponse(model.Unauthorized))
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	var body struct {
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.BadRequest))
		return
	}

	err = service.UnpublishArticle(id, userID.(int64), body.Comment)
	if err != nil {
		if errors.Is(err, service.ErrArticleNotPending) {
			c.JSON(int(model.BadRequest), model.ErrorResponseWithMsg(model.BadRequest, err.Error()))
		} else {
			c.JSON(int(model.NotFound), model.ErrorResponse(model.ArticleNotFound))
		}
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}

// GetReviewRecords 获取所有审核记录（分页）
func GetReviewRecords(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	records, err := service.GetAllReviewRecords(page, pageSize)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(records))
}

// GetPendingReviewArticles 获取待审核文章列表
func GetPendingReviewArticles(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	articles, err := service.GetPendingReviewArticles(page, pageSize)
	if err != nil {
		c.JSON(int(model.InternalServerError), model.ErrorResponse(model.InternalServerError))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(articles))
}