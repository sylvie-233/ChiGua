package sql

// TagSQL 标签相关SQL常量

// TagCheckExists 检查标签名称是否存在
const TagCheckExists = `SELECT EXISTS(SELECT 1 FROM tag WHERE name = $1)`

// TagCheckExistsExcludeID 检查标签名称是否存在（排除指定ID）
const TagCheckExistsExcludeID = `SELECT EXISTS(SELECT 1 FROM tag WHERE name = $1 AND id != $2)`

// TagSelectByID 按ID查询标签
const TagSelectByID = `SELECT id, name, created_at, update_at FROM tag WHERE id = $1`

// TagInsert 插入标签
const TagInsert = `INSERT INTO tag (name, created_at, update_at) VALUES ($1, $2, $3) RETURNING id`

// TagUpdate 更新标签
const TagUpdate = `UPDATE tag SET name = $1, update_at = $2 WHERE id = $3`

// TagSelectAll 查询所有标签
const TagSelectAll = `SELECT id, name, created_at, update_at FROM tag ORDER BY created_at DESC`

// TagCountAll 统计所有标签数
const TagCountAll = `SELECT COUNT(*) FROM tag`

// TagSelectPaged 分页查询标签
const TagSelectPaged = `SELECT id, name, created_at, update_at FROM tag ORDER BY created_at DESC LIMIT $1 OFFSET $2`

// TagCountAllByName 按名称关键词统计标签数
const TagCountAllByName = `
SELECT COUNT(*)
FROM tag
WHERE name ILIKE '%' || $1 || '%'
`

// TagSelectPagedByName 按名称关键词分页查询标签
const TagSelectPagedByName = `
SELECT id, name, created_at, update_at
FROM tag
WHERE name ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

// TagDelete 删除标签
const TagDelete = `DELETE FROM tag WHERE id = $1`

// TagDeleteRelate 删除标签关联
const TagDeleteRelate = `DELETE FROM article_tag WHERE tag_id = $1`
