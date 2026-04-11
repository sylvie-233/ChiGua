package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
	"errors"
	"time"
)

func CreateComment(comment model.CommentCreate, userID int64) (*model.Comment, error) {
	now := time.Now()
	newComment := model.Comment{
		ParentID:    comment.ParentID,
		ArticleID:   comment.ArticleID,
		ReplyUserID: comment.ReplyUserID,
		UserID:      userID,
		Content:     comment.Content,
		CreatedAt:   now,
	}

	// 插入评论
	query := `INSERT INTO comment (parant_id, article_id, reply_user_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := database.DB.QueryRow(query, newComment.ParentID, newComment.ArticleID, newComment.ReplyUserID, newComment.UserID, newComment.Content, newComment.CreatedAt).Scan(&newComment.ID)
	if err != nil {
		return nil, err
	}

	return &newComment, nil
}

func GetCommentsByArticleID(articleID int64) ([]model.CommentResponse, error) {
	var comments []model.Comment
	err := database.DB.Select(&comments, "SELECT id, parant_id, article_id, reply_user_id, user_id, content, created_at FROM comment WHERE article_id = $1 ORDER BY created_at DESC", articleID)
	if err != nil {
		return nil, err
	}

	// 构建响应
	response := make([]model.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		// 获取用户信息
		var user model.User
		database.DB.Get(&user, "SELECT id, username, nickname, role, created_at, update_at FROM users WHERE id = $1", comment.UserID)

		// 获取回复用户信息
		var replyUser *model.User
		if comment.ReplyUserID != 0 {
			replyUser = &model.User{}
			database.DB.Get(replyUser, "SELECT id, username, nickname, role, created_at, update_at FROM users WHERE id = $1", comment.ReplyUserID)
		}

		response = append(response, model.CommentResponse{
			ID:          comment.ID,
			ParentID:    comment.ParentID,
			ArticleID:   comment.ArticleID,
			ReplyUserID: comment.ReplyUserID,
			UserID:      comment.UserID,
			Content:     comment.Content,
			CreatedAt:   comment.CreatedAt,
			User:        user,
			ReplyUser:   replyUser,
		})
	}

	return response, nil
}

func DeleteComment(id int64, userID int64) error {
	// 检查评论是否存在且属于当前用户
	var comment model.Comment
	err := database.DB.Get(&comment, "SELECT id, user_id FROM comment WHERE id = $1", id)
	if err != nil {
		return err
	}
	if comment.UserID != userID {
		return errors.New("无权限删除此评论")
	}

	// 删除评论
	_, err = database.DB.Exec(`DELETE FROM comment WHERE id = $1`, id)
	return err
}
