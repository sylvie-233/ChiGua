package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
	"chigua-backend/internal/sql"
	"chigua-backend/utils/logger"
	"errors"
	"time"
)

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

	// 构建响应
	response := &model.ArticleList{
		Total: total,
		Items: make([]model.ArticleResponse, 0, len(articles)),
	}

	for _, article := range articles {
		var tags []model.Tag
		database.DB.Select(&tags, sql.ArticleSelectTags, article.ID)

		var category model.Category
		database.DB.Get(&category, sql.ArticleSelectCategory, article.CategoryID)

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
			Author:     author,
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
		Author:     author,
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
		return nil, errors.New("无权限修改此文章")
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

	now := time.Now()
	_, err = tx.Exec(sql.ArticleUpdate, update.Title, update.Content, update.CoverImage, update.CategoryID, now, id)
	if err != nil {
		return nil, err
	}

	if update.TagIDs != nil {
		_, err = tx.Exec(sql.ArticleDeleteTags, id)
		if err != nil {
			return nil, err
		}

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
		return errors.New("无权限删除此文章")
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

	_, err = tx.Exec(sql.ArticleDeleteTags, id)
	if err != nil {
		return err
	}

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
		return errors.New("无权限发布此文章")
	}

	now := time.Now()
	_, err = database.DB.Exec(sql.ArticleUpdateStatus, model.ArticleStatusPublished, now, id)
	return err
}
