package api

import (
	"chigua-backend/internal/model"
	"chigua-backend/internal/service"
	"chigua-backend/utils/logger"
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

	articles, err := service.GetArticleList(page, pageSize)
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
	logger.Error(err)
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

	article, err := service.UpdateArticle(id, articleUpdate, userID.(int64))
	if err != nil {
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.InvalidParams))
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
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.InvalidParams))
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
		c.JSON(int(model.BadRequest), model.ErrorResponse(model.InvalidParams))
		return
	}

	c.JSON(int(model.Success), model.SuccessResponse(nil))
}
