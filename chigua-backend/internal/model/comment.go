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
	ID          int64         `json:"id"`
	ParentID    int64         `json:"parentId"`
	ArticleID   int64         `json:"articleId"`
	ReplyUserID int64         `json:"replyUserId"`
	UserID      int64         `json:"userId"`
	Content     string        `json:"content"`
	CreatedAt   time.Time     `json:"createdAt"`
	User        *UserResponse `json:"user"`
	ReplyUser   *UserResponse `json:"replyUser"`
}

// SecondLevelCommentList 二级评论分页列表
type SecondLevelCommentList struct {
	PageResponse
	Items []CommentResponse `json:"items"` // 当前页二级评论数据
}

// CommentWithChildren 带二级评论的响应结构
type CommentWithChildren struct {
	ID          int64                  `json:"id"`
	ParentID    int64                  `json:"parentId"`
	ArticleID   int64                  `json:"articleId"`
	ReplyUserID int64                  `json:"replyUserId"`
	UserID      int64                  `json:"userId"`
	Content     string                 `json:"content"`
	CreatedAt   time.Time              `json:"createdAt"`
	User        *UserResponse          `json:"user"`
	ReplyUser   *UserResponse          `json:"replyUser"`
	Children    SecondLevelCommentList `json:"children"` // 二级评论分页列表
}

// CommentListResponse 评论列表响应（与文章分页结构一致）
type CommentListResponse struct {
	PageResponse
	Items []CommentWithChildren `json:"items"` // 当前页数据
}
