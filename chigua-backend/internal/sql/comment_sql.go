package sql

// CommentSQL 评论相关SQL常量

// CommentInsert 插入评论
const CommentInsert = `INSERT INTO comment (parant_id, article_id, reply_user_id, user_id, content, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

// CommentSelectByArticle 按文章ID查询评论
const CommentSelectByArticle = `SELECT id, parant_id, article_id, reply_user_id, user_id, content, created_at FROM comment WHERE article_id = $1 ORDER BY created_at DESC`

// CommentSelectUserID 查询评论用户ID
const CommentSelectUserID = `SELECT id, user_id FROM comment WHERE id = $1`

// CommentDelete 删除评论
const CommentDelete = `DELETE FROM comment WHERE id = $1`
