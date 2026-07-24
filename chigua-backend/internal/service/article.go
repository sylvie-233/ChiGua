package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
	"chigua-backend/internal/sql"
	"chigua-backend/utils/logger"
	"errors"
	"strings"
	"time"
)

var ErrArticleNoPermission = errors.New("无权限操作此文章")
var ErrArticleNotPending = errors.New("文章不在审核中状态")

var (
	ErrArticleNoTitle      = errors.New("标题不能为空")
	ErrArticleNoContent    = errors.New("内容不能为空")
	ErrArticleNoCategory   = errors.New("请选择分类")
	ErrArticleNoCoverImage = errors.New("请上传3张封面图片")
)

// validateArticleForPublish 发布前校验必填项
func validateArticleForPublish(article *model.Article) error {
	if strings.TrimSpace(article.Title) == "" || article.Title == "无标题" {
		return ErrArticleNoTitle
	}
	if strings.TrimSpace(article.Content) == "" {
		return ErrArticleNoContent
	}
	if article.CategoryID == 0 {
		return ErrArticleNoCategory
	}
	covers := strings.Split(article.CoverImage, ",")
	if len(covers) != 3 || strings.TrimSpace(article.CoverImage) == "" {
		return ErrArticleNoCoverImage
	}
	return nil
}

// IsAdminOrReviewer 判断用户是否为管理员或审核员
func IsAdminOrReviewer(userID int64) (bool, error) {
	user, err := GetUserByID(userID)
	if err != nil {
		return false, err
	}
	return user.Role == "admin" || user.Role == "reviewer", nil
}

// userHasReviewPermission 检查用户是否有权限操作他人文章
func userHasReviewPermission(userID int64) (bool, error) {
	return IsAdminOrReviewer(userID)
}

// buildArticleResponse 构建文章响应（含标签、分类、作者、审核人）
func buildArticleResponse(article model.Article) model.ArticleResponse {
	var tags []model.Tag
	database.DB.Select(&tags, sql.ArticleSelectTags, article.ID)

	var category model.Category
	database.DB.Get(&category, sql.ArticleSelectCategory, article.CategoryID)

	var author model.User
	database.DB.Get(&author, sql.UserSelectByID, article.AuthorID)

	resp := model.ArticleResponse{
		ID:            article.ID,
		AuthorID:      article.AuthorID,
		CategoryID:    article.CategoryID,
		Title:         article.Title,
		Content:       article.Content,
		CoverImage:    article.CoverImage,
		Status:        article.Status,
		ReviewerID:    article.ReviewerID,
		ReviewComment: article.ReviewComment,
		SubmittedAt:   article.SubmittedAt,
		PublishAt:     article.PublishAt,
		CreatedAt:     article.CreatedAt,
		UpdateAt:      article.UpdateAt,
		Tags:          tags,
		Category:      category,
		Author:        *author.ToResponse(),
	}

	// 如果有审核人，查询审核人信息
	if article.ReviewerID > 0 {
		var reviewer model.User
		if err := database.DB.Get(&reviewer, sql.UserSelectByID, article.ReviewerID); err == nil {
			resp.Reviewer = reviewer.ToResponse()
		}
	}

	return resp
}

func CreateArticle(article model.ArticleCreate, authorID int64) (*model.Article, error) {
	now := time.Now()
	title := article.Title
	if strings.TrimSpace(title) == "" {
		title = "无标题"
	}
	newArticle := model.Article{
		AuthorID:   authorID,
		CategoryID: article.CategoryID,
		Title:      title,
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

	err = tx.QueryRow(sql.ArticleInsert, newArticle.AuthorID, newArticle.CategoryID, newArticle.Title, newArticle.Content, newArticle.CoverImage, newArticle.Status, newArticle.ReviewerID, newArticle.ReviewComment, newArticle.SubmittedAt, newArticle.CreatedAt, newArticle.UpdateAt).Scan(&newArticle.ID)
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

func GetAllArticleList(page, pageSize int, keyword string) (*model.ArticleList, error) {
	var total int64
	var articles []model.Article
	var err error

	offset := (page - 1) * pageSize

	if keyword != "" {
		err = database.DB.Get(&total, sql.ArticleCountAllByTitle, keyword)
		if err != nil {
			return nil, err
		}
		err = database.DB.Select(&articles, sql.ArticleSelectAllByTitle, keyword, pageSize, offset)
	} else {
		err = database.DB.Get(&total, sql.ArticleCountAll)
		if err != nil {
			return nil, err
		}
		err = database.DB.Select(&articles, sql.ArticleSelectAll, pageSize, offset)
	}

	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page < totalPages
	hasPrev := page > 1

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
		response.Items = append(response.Items, buildArticleResponse(article))
	}

	return response, nil
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
		response.Items = append(response.Items, buildArticleResponse(article))
	}

	return response, nil
}

func GetArticleByID(id int64) (*model.ArticleResponse, error) {
	var article model.Article
	err := database.DB.Get(&article, sql.ArticleSelectByID, id)
	if err != nil {
		return nil, err
	}

	resp := buildArticleResponse(article)
	return &resp, nil
}

func UpdateArticle(id int64, update model.ArticleUpdate, authorID int64) (*model.Article, error) {
	var article model.Article
	err := database.DB.Get(&article, sql.ArticleSelectAuthorID, id)
	if err != nil {
		return nil, err
	}

	// 权限检查：作者本人 或 admin/reviewer 可以编辑
	if article.AuthorID != authorID {
		isAdminReviewer, err := userHasReviewPermission(authorID)
		if err != nil || !isAdminReviewer {
			return nil, ErrArticleNoPermission
		}
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

	// 已发布的文章编辑后重新进入审核中状态
	if article.Status == int(model.ArticleStatusPublished) {
		_, err = tx.Exec(sql.ArticleSubmitForReview, model.ArticleStatusPendingReview, now, id)
		if err != nil {
			return nil, err
		}
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

// SubmitForReview 提交文章审核
func SubmitForReview(id int64, authorID int64) error {
	var article model.Article
	err := database.DB.Get(&article, sql.ArticleSelectByID, id)
	if err != nil {
		return err
	}

	// 权限检查：作者本人 或 admin/reviewer 可以提交审核
	if article.AuthorID != authorID {
		isAdminReviewer, err := userHasReviewPermission(authorID)
		if err != nil || !isAdminReviewer {
			return ErrArticleNoPermission
		}
	}

	// 只有草稿或已下架状态可以提交审核
	if article.Status != int(model.ArticleStatusDraft) && article.Status != int(model.ArticleStatusUnpublished) {
		return ErrArticleNotPending
	}

	// 提交审核前校验必填项
	if err := validateArticleForPublish(&article); err != nil {
		return err
	}

	now := time.Now()
	_, err = database.DB.Exec(sql.ArticleSubmitForReview, model.ArticleStatusPendingReview, now, id)
	return err
}

// ApproveArticle 审核通过
func ApproveArticle(id int64, reviewerID int64) error {
	var article model.Article
	err := database.DB.Get(&article, sql.ArticleSelectByID, id)
	if err != nil {
		return err
	}

	// 只有审核中状态才能通过
	if article.Status != int(model.ArticleStatusPendingReview) {
		return ErrArticleNotPending
	}

	// 审核通过前再次校验必填项
	if err := validateArticleForPublish(&article); err != nil {
		return err
	}

	now := time.Now()
	// 更新文章状态、审核人、发布时间
	_, err = database.DB.Exec(sql.ArticleUpdateReview, model.ArticleStatusPublished, reviewerID, "", now, id)
	if err != nil {
		return err
	}

	// 记录审核历史
	return CreateReviewRecord(id, reviewerID, "approve", "")
}

// RejectArticle 驳回文章
func RejectArticle(id int64, reviewerID int64, comment string) error {
	var article model.Article
	err := database.DB.Get(&article, sql.ArticleSelectByID, id)
	if err != nil {
		return err
	}

	// 只有审核中状态才能驳回
	if article.Status != int(model.ArticleStatusPendingReview) {
		return ErrArticleNotPending
	}

	now := time.Now()
	// 驳回：状态回到草稿，记录审核人和审核意见
	_, err = database.DB.Exec(sql.ArticleUpdateReview, model.ArticleStatusDraft, reviewerID, comment, now, id)
	if err != nil {
		return err
	}

	// 记录审核历史
	return CreateReviewRecord(id, reviewerID, "reject", comment)
}

// UnpublishArticle 下架文章
func UnpublishArticle(id int64, reviewerID int64, comment string) error {
	var article model.Article
	err := database.DB.Get(&article, sql.ArticleSelectByID, id)
	if err != nil {
		return err
	}

	// 只有已发布状态才能下架
	if article.Status != int(model.ArticleStatusPublished) {
		return ErrArticleNotPending
	}

	now := time.Now()
	_, err = database.DB.Exec(sql.ArticleUpdateReview, model.ArticleStatusUnpublished, reviewerID, comment, now, id)
	if err != nil {
		return err
	}

	// 记录审核历史
	return CreateReviewRecord(id, reviewerID, "unpublish", comment)
}

// CreateReviewRecord 创建审核记录
func CreateReviewRecord(articleID, reviewerID int64, action, comment string) error {
	_, err := database.DB.Exec(sql.ReviewRecordInsert, articleID, reviewerID, action, comment, time.Now())
	return err
}

// GetReviewRecords 获取文章审核记录列表
func GetReviewRecords(articleID int64) ([]model.ArticleReviewRecordResponse, error) {
	type row struct {
		ID         int64     `db:"id"`
		ArticleID  int64     `db:"article_id"`
		ReviewerID int64     `db:"reviewer_id"`
		Action     string    `db:"action"`
		Comment    string    `db:"comment"`
		CreatedAt  time.Time `db:"created_at"`
		// reviewer fields
		ReviewerID2       int64     `db:"id"`
		ReviewerUsername  string    `db:"username"`
		ReviewerPassword  string    `db:"password"`
		ReviewerNickname  string    `db:"nickname"`
		ReviewerRole      string    `db:"role"`
		ReviewerCreatedAt time.Time `db:"created_at"`
		ReviewerUpdateAt  time.Time `db:"update_at"`
	}

	var rows []row
	err := database.DB.Select(&rows, sql.ReviewRecordSelectByArticle, articleID)
	if err != nil {
		return nil, err
	}

	records := make([]model.ArticleReviewRecordResponse, 0, len(rows))
	for _, r := range rows {
		records = append(records, model.ArticleReviewRecordResponse{
			ID:        r.ID,
			ArticleID: r.ArticleID,
			Action:    r.Action,
			Comment:   r.Comment,
			CreatedAt: r.CreatedAt,
			Reviewer: model.UserResponse{
				ID:        r.ReviewerID2,
				Username:  r.ReviewerUsername,
				Nickname:  r.ReviewerNickname,
				Role:      r.ReviewerRole,
				CreatedAt: r.ReviewerCreatedAt,
				UpdateAt:  r.ReviewerUpdateAt,
			},
		})
	}

	return records, nil
}

// GetAllReviewRecords 获取所有审核记录（分页）
func GetAllReviewRecords(page, pageSize int) (*model.ArticleReviewRecordList, error) {
	var total int64
	err := database.DB.Get(&total, sql.ReviewRecordCountAll)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize

	type row struct {
		ID           int64     `db:"id"`
		ArticleID    int64     `db:"article_id"`
		ArticleTitle string    `db:"article_title"`
		Action       string    `db:"action"`
		Comment      string    `db:"comment"`
		CreatedAt    time.Time `db:"created_at"`
		// reviewer fields
		ReviewerID       int64     `db:"id"`
		ReviewerUsername string    `db:"username"`
		ReviewerPassword string    `db:"password"`
		ReviewerNickname string    `db:"nickname"`
		ReviewerRole     string    `db:"role"`
		ReviewerCreatedAt time.Time `db:"created_at"`
		ReviewerUpdateAt  time.Time `db:"update_at"`
	}

	var rows []row
	err = database.DB.Select(&rows, sql.ReviewRecordSelectAll, pageSize, offset)
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page < totalPages
	hasPrev := page > 1

	records := make([]model.ArticleReviewRecordListResponse, 0, len(rows))
	for _, r := range rows {
		records = append(records, model.ArticleReviewRecordListResponse{
			ID:           r.ID,
			ArticleID:    r.ArticleID,
			ArticleTitle: r.ArticleTitle,
			Action:       r.Action,
			Comment:      r.Comment,
			CreatedAt:    r.CreatedAt,
			Reviewer: model.UserResponse{
				ID:        r.ReviewerID,
				Username:  r.ReviewerUsername,
				Nickname:  r.ReviewerNickname,
				Role:      r.ReviewerRole,
				CreatedAt: r.ReviewerCreatedAt,
				UpdateAt:  r.ReviewerUpdateAt,
			},
		})
	}

	return &model.ArticleReviewRecordList{
		PageResponse: model.PageResponse{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
		Items: records,
	}, nil
}

// GetPendingReviewArticles 获取待审核文章列表
func GetPendingReviewArticles(page, pageSize int) (*model.ArticleList, error) {
	var total int64
	err := database.DB.Get(&total, sql.ArticleCountByStatus, model.ArticleStatusPendingReview)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	var articles []model.Article
	err = database.DB.Select(&articles, sql.ArticleSelectPendingReview, model.ArticleStatusPendingReview, pageSize, offset)
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page < totalPages
	hasPrev := page > 1

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
		response.Items = append(response.Items, buildArticleResponse(article))
	}

	return response, nil
}

func CountArticleByStatus(status int) (int64, error) {
	var count int64
	err := database.DB.Get(&count, sql.ArticleCountByStatus, status)
	if err != nil {
		return 0, err
	}
	return count, nil
}

type RecentArticle struct {
	ID         int64     `json:"id" db:"id"`
	Title      string    `json:"title" db:"title"`
	AuthorName string    `json:"authorName" db:"author_name"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}

func GetRecentArticles(limit int) ([]RecentArticle, error) {
	articles := make([]RecentArticle, 0)
	err := database.DB.Select(&articles, sql.ArticleSelectRecent, limit)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

type ArticleDailyCount struct {
	Date  string `json:"date" db:"date"`
	Count int64  `json:"count" db:"count"`
}

func GetArticleCountByDate(startDate time.Time) ([]ArticleDailyCount, error) {
	counts := make([]ArticleDailyCount, 0)
	err := database.DB.Select(&counts, sql.ArticleCountByDate, startDate)
	if err != nil {
		return nil, err
	}
	return counts, nil
}

func DeleteArticle(id int64, authorID int64) error {
	var article model.Article
	err := database.DB.Get(&article, sql.ArticleSelectAuthorID, id)
	if err != nil {
		return err
	}

	// 权限检查：作者本人 或 admin/reviewer 可以删除
	if article.AuthorID != authorID {
		isAdminReviewer, err := userHasReviewPermission(authorID)
		if err != nil || !isAdminReviewer {
			return ErrArticleNoPermission
		}
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
	// 检查用户角色
	isAdminReviewer, err := userHasReviewPermission(authorID)
	if err != nil {
		return err
	}

	// 普通用户不能直接发布，必须走审核流程
	if !isAdminReviewer {
		return SubmitForReview(id, authorID)
	}

	// admin/reviewer 可以直接发布
	var article model.Article
	err = database.DB.Get(&article, sql.ArticleSelectByID, id)
	if err != nil {
		return err
	}
	if article.AuthorID != authorID {
		return ErrArticleNoPermission
	}

	// 发布前校验必填项
	if err := validateArticleForPublish(&article); err != nil {
		return err
	}

	now := time.Now()
	_, err = database.DB.Exec(sql.ArticleUpdateStatus, model.ArticleStatusPublished, now, id)
	return err
}

func UpdateArticleStatus(id int64, authorID int64, status int) error {
	var article model.Article
	err := database.DB.Get(&article, sql.ArticleSelectByID, id)
	if err != nil {
		return err
	}

	// 权限检查：作者本人 或 admin/reviewer 可以修改
	if article.AuthorID != authorID {
		isAdminReviewer, err := userHasReviewPermission(authorID)
		if err != nil || !isAdminReviewer {
			return ErrArticleNoPermission
		}
	}

	// 不允许通过此接口设置为审核中状态（必须通过 SubmitForReview）
	if status == int(model.ArticleStatusPendingReview) {
		return ErrArticleNotPending
	}

	// 发布时校验必填项
	if status == int(model.ArticleStatusPublished) {
		if err := validateArticleForPublish(&article); err != nil {
			return err
		}
	}

	now := time.Now()
	_, err = database.DB.Exec(sql.ArticleUpdateStatus, status, now, id)
	return err
}
