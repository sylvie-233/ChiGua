package sql

// CommentSQL 评论相关SQL常量

// CommentInsert 插入评论
const CommentInsert = `
INSERT INTO comment (
	parant_id, 
	article_id, 
	reply_user_id, 
	user_id, 
	content, 
	created_at
) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

// CommentSelectByArticle 按文章ID查询评论
const CommentSelectByArticle = `SELECT id, parant_id, article_id, reply_user_id, user_id, content, created_at FROM comment WHERE article_id = $1 ORDER BY created_at DESC`

// CommentSelectUserID 查询评论用户ID
const CommentSelectUserID = `SELECT id, user_id FROM comment WHERE id = $1`

// CommentDelete 删除评论
const CommentDelete = `DELETE FROM comment WHERE id = $1`

// CommentSelectFirstLevel 分页查询一级评论（ParentID = 0）
const CommentSelectFirstLevel = `
SELECT 
	id, 
	parant_id, 
	article_id, 
	reply_user_id, 
	user_id, 
	content, 
	created_at 
FROM comment 
WHERE article_id = $1 AND parant_id = 0 
ORDER BY created_at DESC 
LIMIT $2 OFFSET $3
`

// CommentSelectFirstLevelCount 统计一级评论总数
const CommentSelectFirstLevelCount = `
SELECT 
	COUNT(*) 
FROM comment 
WHERE article_id = $1 AND parant_id = 0
`

// CommentSelectSecondLevel 分页查询二级评论（指定ParentID）
const CommentSelectSecondLevel = `
SELECT 
	id, 
	parant_id, 
	article_id, 
	reply_user_id, 
	user_id, 
	content, 
	created_at 
FROM comment 
WHERE article_id = $1 AND parant_id = $2 
ORDER BY created_at DESC 
LIMIT $3 OFFSET $4
`

// CommentSelectSecondLevelCount 统计二级评论总数
const CommentSelectSecondLevelCount = `
SELECT COUNT(*) 
FROM comment 
WHERE article_id = $1 AND parant_id = $2
`

// CommentCountAll 统计所有评论数
const CommentCountAll = `
SELECT COUNT(*) 
FROM comment
`

// CommentSelectAll 查询所有评论列表
const CommentSelectAll = `
SELECT
	id,
	parant_id,
	article_id,
	reply_user_id,
	user_id,
	content,
	created_at
FROM comment
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

// CommentCountAllByContent 按内容关键词统计评论数
const CommentCountAllByContent = `
SELECT COUNT(*)
FROM comment
WHERE content ILIKE '%' || $1 || '%'
`

// CommentSelectAllByContent 按内容关键词查询评论列表
const CommentSelectAllByContent = `
SELECT
	id,
	parant_id,
	article_id,
	reply_user_id,
	user_id,
	content,
	created_at
FROM comment
WHERE content ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`
