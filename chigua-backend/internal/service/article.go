package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
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
		Status:     0, // 未发布
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

	// 插入文章
	query := `INSERT INTO article (author_id, category_id, title, content, cover_image, status, created_at, update_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	err = tx.QueryRow(query, newArticle.AuthorID, newArticle.CategoryID, newArticle.Title, newArticle.Content, newArticle.CoverImage, newArticle.Status, newArticle.CreatedAt, newArticle.UpdateAt).Scan(&newArticle.ID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// 插入标签关联
	if len(article.TagIDs) > 0 {
		for _, tagID := range article.TagIDs {
			_, err = tx.Exec(`INSERT INTO article_tag (article_id, tag_id) VALUES ($1, $2)`, newArticle.ID, tagID)
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
	err := database.DB.Get(&total, "SELECT COUNT(*) FROM article WHERE status = 1")
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	var articles []model.Article
	err = database.DB.Select(&articles, "SELECT id, author_id, category_id, title, content, cover_image, status, publish_at, created_at, update_at FROM article WHERE status = 1 ORDER BY publish_at DESC LIMIT $1 OFFSET $2", pageSize, offset)
	if err != nil {
		return nil, err
	}

	// 构建响应
	response := &model.ArticleList{
		Total: total,
		Items: make([]model.ArticleResponse, 0, len(articles)),
	}

	for _, article := range articles {
		// 获取标签
		var tags []model.Tag
		database.DB.Select(&tags, "SELECT t.id, t.name, t.created_at, t.update_at FROM tag t JOIN article_tag at ON t.id = at.tag_id WHERE at.article_id = $1", article.ID)

		// 获取分类
		var category model.Category
		database.DB.Get(&category, "SELECT id, name, created_at, update_at FROM category WHERE id = $1", article.CategoryID)

		// 获取作者
		var author model.User
		database.DB.Get(&author, "SELECT id, username, nickname, role, created_at, update_at FROM users WHERE id = $1", article.AuthorID)

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
	err := database.DB.Get(&article, "SELECT id, author_id, category_id, title, content, cover_image, status, publish_at, created_at, update_at FROM article WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	// 获取标签
	var tags []model.Tag
	database.DB.Select(&tags, "SELECT t.id, t.name, t.created_at, t.update_at FROM tag t JOIN article_tag at ON t.id = at.tag_id WHERE at.article_id = $1", article.ID)

	// 获取分类
	var category model.Category
	database.DB.Get(&category, "SELECT id, name, created_at, update_at FROM category WHERE id = $1", article.CategoryID)

	// 获取作者
	var author model.User
	database.DB.Get(&author, "SELECT id, username, nickname, role, created_at, update_at FROM users WHERE id = $1", article.AuthorID)

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
	// 检查文章是否存在且属于当前用户
	var article model.Article
	err := database.DB.Get(&article, "SELECT id, author_id FROM article WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	if article.AuthorID != authorID {
		return nil, errors.New("无权限修改此文章")
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

	// 更新文章
	now := time.Now()
	query := `UPDATE article SET title = COALESCE($1, title), content = COALESCE($2, content), cover_image = COALESCE($3, cover_image), category_id = COALESCE($4, category_id), update_at = $5 WHERE id = $6`
	_, err = tx.Exec(query, update.Title, update.Content, update.CoverImage, update.CategoryID, now, id)
	if err != nil {
		return nil, err
	}

	// 更新标签关联
	if update.TagIDs != nil {
		// 删除旧标签关联
		_, err = tx.Exec(`DELETE FROM article_tag WHERE article_id = $1`, id)
		if err != nil {
			return nil, err
		}

		// 添加新标签关联
		for _, tagID := range update.TagIDs {
			_, err = tx.Exec(`INSERT INTO article_tag (article_id, tag_id) VALUES ($1, $2)`, id, tagID)
			if err != nil {
				return nil, err
			}
		}
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// 重新获取更新后的文章
	var updatedArticle model.Article
	database.DB.Get(&updatedArticle, "SELECT id, author_id, category_id, title, content, cover_image, status, publish_at, created_at, update_at FROM article WHERE id = $1", id)

	return &updatedArticle, nil
}

func DeleteArticle(id int64, authorID int64) error {
	// 检查文章是否存在且属于当前用户
	var article model.Article
	err := database.DB.Get(&article, "SELECT id, author_id FROM article WHERE id = $1", id)
	if err != nil {
		return err
	}
	if article.AuthorID != authorID {
		return errors.New("无权限删除此文章")
	}

	// 开始事务
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 删除标签关联
	_, err = tx.Exec(`DELETE FROM article_tag WHERE article_id = $1`, id)
	if err != nil {
		return err
	}

	// 删除文章
	_, err = tx.Exec(`DELETE FROM article WHERE id = $1`, id)
	if err != nil {
		return err
	}

	// 提交事务
	return tx.Commit()
}

func PublishArticle(id int64, authorID int64) error {
	// 检查文章是否存在且属于当前用户
	var article model.Article
	err := database.DB.Get(&article, "SELECT id, author_id FROM article WHERE id = $1", id)
	if err != nil {
		return err
	}
	if article.AuthorID != authorID {
		return errors.New("无权限发布此文章")
	}

	// 更新文章状态
	now := time.Now()
	_, err = database.DB.Exec(`UPDATE article SET status = 1, publish_at = $1, update_at = $1 WHERE id = $2`, now, id)
	return err
}
