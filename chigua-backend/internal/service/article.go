package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
	"chigua-backend/internal/sql"
	"chigua-backend/utils/logger"
	"errors"
	"time"
)

var ErrArticleNoPermission = errors.New("无权限操作此文章")

func CreateArticle(article model.ArticleCreate, authorID int64) (*model.Article, error) {
	now := time.Now()
	newArticle := model.Article{
		AuthorID:   authorID,
		CategoryID: article.CategoryID,
		Title:      article.Title,
		Content:    article.Content,
		CoverImage: article.CoverImage,
		Status:     int(model.ArticleStatusDraft),
		CreatedAt:  now,
		UpdateAt:   now,
	}

	// 开始事务
	tx, err := database.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = tx.QueryRow(sql.ArticleInsert, newArticle.AuthorID, newArticle.CategoryID, newArticle.Title, newArticle.Content, newArticle.CoverImage, newArticle.Status, newArticle.CreatedAt, newArticle.UpdateAt).Scan(&newArticle.ID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if len(article.TagIDs) > 0 {
		for _, tagID := range article.TagIDs {
			_, err = tx.Exec(sql.ArticleInsertTag, newArticle.ID, tagID)
			if err != nil {
				logger.Error(err)
				return nil, err
			}
		}
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &newArticle, nil
}

func GetArticleList(page, pageSize int) (*model.ArticleList, error) {
	var total int64
	err := database.DB.Get(&total, sql.ArticleCountByStatus, model.ArticleStatusPublished)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	var articles []model.Article
	err = database.DB.Select(&articles, sql.ArticleSelectByStatus, model.ArticleStatusPublished, pageSize, offset)
	if err != nil {
		return nil, err
	}

	// 计算分页参数
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page < totalPages
	hasPrev := page > 1

	// 构建响应
	response := &model.ArticleList{
		PageResponse: model.PageResponse{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
		Items: make([]model.ArticleResponse, 0, len(articles)),
	}

	for _, article := range articles {
		// 查询文章标签
		var tags []model.Tag
		database.DB.Select(&tags, sql.ArticleSelectTags, article.ID)

		// 查询文章分类
		var category model.Category
		database.DB.Get(&category, sql.ArticleSelectCategory, article.CategoryID)

		// 查询文章作者
		var author model.User
		database.DB.Get(&author, sql.UserSelectByID, article.AuthorID)

		response.Items = append(response.Items, model.ArticleResponse{
			ID:         article.ID,
			AuthorID:   article.AuthorID,
			CategoryID: article.CategoryID,
			Title:      article.Title,
			Content:    article.Content,
			CoverImage: article.CoverImage,
			Status:     article.Status,
			PublishAt:  article.PublishAt,
			CreatedAt:  article.CreatedAt,
			UpdateAt:   article.UpdateAt,
			Tags:       tags,
			Category:   category,
			Author:     *author.ToResponse(),
		})
	}

	return response, nil
}

func GetArticleByID(id int64) (*model.ArticleResponse, error) {
	var article model.Article
	err := database.DB.Get(&article, sql.ArticleSelectByID, id)
	if err != nil {
		return nil, err
	}

	var tags []model.Tag
	database.DB.Select(&tags, sql.ArticleSelectTags, article.ID)

	var category model.Category
	database.DB.Get(&category, sql.ArticleSelectCategory, article.CategoryID)

	var author model.User
	database.DB.Get(&author, sql.UserSelectByID, article.AuthorID)

	response := &model.ArticleResponse{
		ID:         article.ID,
		AuthorID:   article.AuthorID,
		CategoryID: article.CategoryID,
		Title:      article.Title,
		Content:    article.Content,
		CoverImage: article.CoverImage,
		Status:     article.Status,
		PublishAt:  article.PublishAt,
		CreatedAt:  article.CreatedAt,
		UpdateAt:   article.UpdateAt,
		Tags:       tags,
		Category:   category,
		Author:     *author.ToResponse(),
	}

	return response, nil
}

func UpdateArticle(id int64, update model.ArticleUpdate, authorID int64) (*model.Article, error) {
	var article model.Article
	err := database.DB.Get(&article, sql.ArticleSelectAuthorID, id)
	if err != nil {
		return nil, err
	}
	if article.AuthorID != authorID {
		return nil, ErrArticleNoPermission
	}

	tx, err := database.DB.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 更新文章主体
	now := time.Now()
	_, err = tx.Exec(sql.ArticleUpdate, update.Title, update.Content, update.CoverImage, update.CategoryID, now, id)
	if err != nil {
		return nil, err
	}

	// 更新文章标签
	if update.TagIDs != nil {
		// 删除旧标签
		_, err = tx.Exec(sql.ArticleDeleteTags, id)
		if err != nil {
			return nil, err
		}

		// 插入新标签
		for _, tagID := range update.TagIDs {
			_, err = tx.Exec(sql.ArticleInsertTag, id, tagID)
			if err != nil {
				return nil, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	var updatedArticle model.Article
	database.DB.Get(&updatedArticle, sql.ArticleSelectByID, id)

	return &updatedArticle, nil
}

func DeleteArticle(id int64, authorID int64) error {
	var article model.Article
	err := database.DB.Get(&article, sql.ArticleSelectAuthorID, id)
	if err != nil {
		return err
	}
	if article.AuthorID != authorID {
		return ErrArticleNoPermission
	}

	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 删除文章关联标签（中间表）
	_, err = tx.Exec(sql.ArticleDeleteTags, id)
	if err != nil {
		return err
	}

	// 删除文章
	_, err = tx.Exec(sql.ArticleDelete, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func PublishArticle(id int64, authorID int64) error {
	var article model.Article
	err := database.DB.Get(&article, sql.ArticleSelectAuthorID, id)
	if err != nil {
		return err
	}
	if article.AuthorID != authorID {
		return ErrArticleNoPermission
	}

	now := time.Now()
	_, err = database.DB.Exec(sql.ArticleUpdateStatus, model.ArticleStatusPublished, now, id)
	return err
}

func UpdateArticleStatus(id int64, authorID int64, status int) error {
	var article model.Article
	err := database.DB.Get(&article, sql.ArticleSelectAuthorID, id)
	if err != nil {
		return err
	}
	if article.AuthorID != authorID {
		return ErrArticleNoPermission
	}

	now := time.Now()
	_, err = database.DB.Exec(sql.ArticleUpdateStatus, status, now, id)
	return err
}
