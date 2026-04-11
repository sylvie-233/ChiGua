package model

import "time"

type Article struct {
	ID         int64     `json:"id" db:"id"`
	AuthorID   int64     `json:"authorId" db:"author_id"`
	CategoryID int64     `json:"categoryId" db:"category_id"`
	Title      string    `json:"title" db:"title" binding:"required"`
	Content    string    `json:"content" db:"content"`
	CoverImage string    `json:"coverImage" db:"cover_image"`
	Status     int       `json:"status" db:"status"`
	PublishAt  time.Time `json:"publishAt" db:"publish_at"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdateAt   time.Time `json:"updateAt" db:"update_at"`
}

type ArticleCreate struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content"`
	CoverImage string `json:"coverImage"`
	CategoryID int64  `json:"categoryId" binding:"required"`
	TagIDs     []int64 `json:"tagIds"`
}

type ArticleUpdate struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	CoverImage string `json:"coverImage"`
	CategoryID int64  `json:"categoryId"`
	TagIDs     []int64 `json:"tagIds"`
}

type ArticleResponse struct {
	ID         int64     `json:"id"`
	AuthorID   int64     `json:"authorId"`
	CategoryID int64     `json:"categoryId"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CoverImage string    `json:"coverImage"`
	Status     int       `json:"status"`
	PublishAt  time.Time `json:"publishAt"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdateAt   time.Time `json:"updateAt"`
	Tags       []Tag     `json:"tags"`
	Category   Category  `json:"category"`
	Author     User      `json:"author"`
}

type ArticleList struct {
	Total int64            `json:"total"`
	Items []ArticleResponse `json:"items"`
}
