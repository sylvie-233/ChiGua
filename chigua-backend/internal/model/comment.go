package model

import "time"

type Comment struct {
	ID          int64     `json:"id" db:"id"`
	ParentID    int64     `json:"parentId" db:"parant_id"`
	ArticleID   int64     `json:"articleId" db:"article_id" binding:"required"`
	ReplyUserID int64     `json:"replyUserId" db:"reply_user_id"`
	UserID      int64     `json:"userId" db:"user_id"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type CommentCreate struct {
	ArticleID   int64  `json:"articleId" binding:"required"`
	ParentID    int64  `json:"parentId"`
	ReplyUserID int64  `json:"replyUserId"`
	Content     string `json:"content" binding:"required"`
}

type CommentResponse struct {
	ID          int64     `json:"id"`
	ParentID    int64     `json:"parentId"`
	ArticleID   int64     `json:"articleId"`
	ReplyUserID int64     `json:"replyUserId"`
	UserID      int64     `json:"userId"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
	User        User      `json:"user"`
	ReplyUser   *User     `json:"replyUser"`
}
