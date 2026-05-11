package sql

// ArticleSQL 文章相关SQL常量

// ArticleInsert 插入文章
const ArticleInsert = `INSERT INTO article (author_id, category_id, title, content, cover_image, status, created_at, update_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

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
	publish_at, 
	created_at, 
	update_at 
FROM article 
WHERE status = $1 
ORDER BY publish_at DESC 
LIMIT $2 OFFSET $3
`

// ArticleSelectByID 按ID查询文章
const ArticleSelectByID = `SELECT id, author_id, category_id, title, content, cover_image, status, publish_at, created_at, update_at FROM article WHERE id = $1`

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

// ArticleDelete 删除文章
const ArticleDelete = `DELETE FROM article WHERE id = $1`

// ArticleDeleteTags 删除文章标签关联
const ArticleDeleteTags = `DELETE FROM article_tag WHERE article_id = $1`

// ArticleSelectAuthorID 查询文章作者ID
const ArticleSelectAuthorID = `
SELECT 
	id, 
	author_id 
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
