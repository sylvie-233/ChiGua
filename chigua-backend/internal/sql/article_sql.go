package sql

// ArticleSQL 文章相关SQL常量

// ArticleInsert 插入文章
const ArticleInsert = `INSERT INTO article (author_id, category_id, title, content, cover_image, status, reviewer_id, review_comment, submitted_at, created_at, update_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`

// ArticleInsertTag 插入文章标签关联
const ArticleInsertTag = `INSERT INTO article_tag (article_id, tag_id) VALUES ($1, $2)`

// ArticleCountByStatus 按状态统计文章数
const ArticleCountByStatus = `
SELECT
	COUNT(*)
FROM article
WHERE status = $1
`

// ArticleSelectByStatus 按状态查询文章列表
const ArticleSelectByStatus = `
SELECT
	id,
	author_id,
	category_id,
	title,
	content,
	cover_image,
	status,
	reviewer_id,
	review_comment,
	submitted_at,
	publish_at,
	created_at,
	update_at
FROM article
WHERE status = $1
ORDER BY publish_at DESC
LIMIT $2 OFFSET $3
`

// ArticleSelectByID 按ID查询文章
const ArticleSelectByID = `SELECT id, author_id, category_id, title, content, cover_image, status, reviewer_id, review_comment, submitted_at, publish_at, created_at, update_at FROM article WHERE id = $1`

// ArticleUpdate 更新文章
const ArticleUpdate = `
UPDATE article
SET
	title = COALESCE($1, title),
	content = COALESCE($2, content),
	cover_image = COALESCE($3, cover_image),
	category_id = COALESCE($4, category_id),
	update_at = $5
WHERE id = $6
`

// ArticleUpdateStatus 更新文章状态
const ArticleUpdateStatus = `UPDATE article SET status = $1, publish_at = $2, update_at = $2 WHERE id = $3`

// ArticleSubmitForReview 提交审核（更新状态和提交时间）
const ArticleSubmitForReview = `UPDATE article SET status = $1, submitted_at = $2, update_at = $2 WHERE id = $3`

// ArticleUpdateReview 审核操作（更新状态、审核人、审核意见、发布时间）
const ArticleUpdateReview = `UPDATE article SET status = $1, reviewer_id = $2, review_comment = $3, publish_at = $4, update_at = $4 WHERE id = $5`

// ArticleDelete 删除文章
const ArticleDelete = `DELETE FROM article WHERE id = $1`

// ArticleDeleteTags 删除文章标签关联
const ArticleDeleteTags = `DELETE FROM article_tag WHERE article_id = $1`

// ArticleSelectAuthorID 查询文章作者ID和状态
const ArticleSelectAuthorID = `
SELECT
	id,
	author_id,
	status
FROM article
WHERE id = $1
`

// ArticleSelectTags 查询文章标签
const ArticleSelectTags = `
SELECT
	t.id,
	t.name,
	t.created_at,
	t.update_at
FROM tag t
JOIN article_tag at ON t.id = at.tag_id
WHERE at.article_id = $1
`

// ArticleSelectCategory 查询文章分类
const ArticleSelectCategory = `SELECT id, name, created_at, update_at FROM category WHERE id = $1`

// ArticleCountAll 统计所有文章数
const ArticleCountAll = `
SELECT
	COUNT(*)
FROM article
`

// ArticleCountByStatusAndCategory 按状态+分类统计文章数
const ArticleCountByStatusAndCategory = `
SELECT
	COUNT(*)
FROM article
WHERE status = $1 AND category_id = $2
`

// ArticleSelectByStatusAndCategory 按状态+分类查询文章列表
const ArticleSelectByStatusAndCategory = `
SELECT
	id,
	author_id,
	category_id,
	title,
	content,
	cover_image,
	status,
	reviewer_id,
	review_comment,
	submitted_at,
	publish_at,
	created_at,
	update_at
FROM article
WHERE status = $1 AND category_id = $2
ORDER BY publish_at DESC
LIMIT $3 OFFSET $4
`

// ArticleSelectRecent 查询最近发布的文章
const ArticleSelectRecent = `
SELECT a.id, a.title, u.username AS author_name, a.created_at
FROM article a
LEFT JOIN users u ON u.id = a.author_id
WHERE a.status = 1
ORDER BY a.created_at DESC
LIMIT $1
`

// ArticleCountByDate 按日期统计文章数量（最近一年）
const ArticleCountByDate = `
SELECT
	DATE(created_at) AS date,
	COUNT(*) AS count
FROM article
WHERE created_at >= $1
GROUP BY DATE(created_at)
ORDER BY date ASC
`

// ArticleSelectAll 查询所有文章列表
const ArticleSelectAll = `
SELECT
	id,
	author_id,
	category_id,
	title,
	content,
	cover_image,
	status,
	reviewer_id,
	review_comment,
	submitted_at,
	publish_at,
	created_at,
	update_at
FROM article
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

// ArticleCountAllByTitle 按标题关键词统计文章数
const ArticleCountAllByTitle = `
SELECT COUNT(*)
FROM article
WHERE title ILIKE '%' || $1 || '%'
`

// ArticleSelectAllByTitle 按标题关键词查询文章列表
const ArticleSelectAllByTitle = `
SELECT
	id,
	author_id,
	category_id,
	title,
	content,
	cover_image,
	status,
	reviewer_id,
	review_comment,
	submitted_at,
	publish_at,
	created_at,
	update_at
FROM article
WHERE title ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

// ArticleSelectPendingReview 查询审核中的文章（按提交时间倒序）
const ArticleSelectPendingReview = `
SELECT
	id,
	author_id,
	category_id,
	title,
	content,
	cover_image,
	status,
	reviewer_id,
	review_comment,
	submitted_at,
	publish_at,
	created_at,
	update_at
FROM article
WHERE status = $1
ORDER BY submitted_at DESC
LIMIT $2 OFFSET $3
`

// --- 审核记录表 ---

// ReviewRecordInsert 插入审核记录
const ReviewRecordInsert = `INSERT INTO article_review_record (article_id, reviewer_id, action, comment, created_at) VALUES ($1, $2, $3, $4, $5)`

// ReviewRecordSelectByArticle 查询某篇文章的审核记录（按时间倒序，带审核人信息）
const ReviewRecordSelectByArticle = `
SELECT
	r.id, r.article_id, r.reviewer_id, r.action, r.comment, r.created_at,
	u.id, u.username, u.password, u.nickname, u.created_at, u.update_at
FROM article_review_record r
JOIN users u ON r.reviewer_id = u.id
WHERE r.article_id = $1
ORDER BY r.created_at DESC
`

// ReviewRecordCountAll 统计所有审核记录数
const ReviewRecordCountAll = `SELECT COUNT(*) FROM article_review_record`

// ReviewRecordSelectAll 查询所有审核记录（分页，带文章标题和审核人信息）
const ReviewRecordSelectAll = `
SELECT
	r.id, r.article_id, a.title AS article_title, r.action, r.comment, r.created_at,
	u.id, u.username, u.password, u.nickname, u.created_at, u.update_at
FROM article_review_record r
JOIN article a ON r.article_id = a.id
JOIN users u ON r.reviewer_id = u.id
ORDER BY r.created_at DESC
LIMIT $1 OFFSET $2
`
