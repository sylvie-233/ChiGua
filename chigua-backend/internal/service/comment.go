package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
	"chigua-backend/internal/sql"
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

	err := database.DB.QueryRow(sql.CommentInsert, newComment.ParentID, newComment.ArticleID, newComment.ReplyUserID, newComment.UserID, newComment.Content, newComment.CreatedAt).Scan(&newComment.ID)
	if err != nil {
		return nil, err
	}

	return &newComment, nil
}

func GetCommentsByArticleID(articleID int64) ([]model.CommentResponse, error) {
	var comments []model.Comment
	err := database.DB.Select(&comments, sql.CommentSelectByArticle, articleID)
	if err != nil {
		return nil, err
	}

	response := make([]model.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		var user model.User
		database.DB.Get(&user, sql.UserSelectByID, comment.UserID)

		var replyUser *model.User
		if comment.ReplyUserID != 0 {
			replyUser = &model.User{}
			database.DB.Get(replyUser, sql.UserSelectByID, comment.ReplyUserID)
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
	var comment model.Comment
	err := database.DB.Get(&comment, sql.CommentSelectUserID, id)
	if err != nil {
		return err
	}
	if comment.UserID != userID {
		return errors.New("无权限删除此评论")
	}

	_, err = database.DB.Exec(sql.CommentDelete, id)
	return err
}
