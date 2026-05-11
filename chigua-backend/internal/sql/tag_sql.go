package sql

// TagSQL 标签相关SQL常量

// TagCheckExists 检查标签名称是否存在
const TagCheckExists = `SELECT EXISTS(SELECT 1 FROM tag WHERE name = $1)`

// TagInsert 插入标签
const TagInsert = `INSERT INTO tag (name, created_at, update_at) VALUES ($1, $2, $3) RETURNING id`

// TagSelectAll 查询所有标签
const TagSelectAll = `SELECT id, name, created_at, update_at FROM tag ORDER BY created_at DESC`

// TagDelete 删除标签
const TagDelete = `DELETE FROM tag WHERE id = $1`

// TagDeleteRelate 删除标签关联
const TagDeleteRelate = `DELETE FROM article_tag WHERE tag_id = $1`
