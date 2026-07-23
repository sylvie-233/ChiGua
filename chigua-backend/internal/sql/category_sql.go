package sql

// CategorySQL 分类相关SQL常量

// CategoryCheckExists 检查分类名称是否存在
const CategoryCheckExists = `SELECT EXISTS(SELECT 1 FROM category WHERE name = $1)`

// CategoryCheckExistsExcludeID 检查分类名称是否存在（排除指定ID）
const CategoryCheckExistsExcludeID = `SELECT EXISTS(SELECT 1 FROM category WHERE name = $1 AND id != $2)`

// CategorySelectByID 按ID查询分类
const CategorySelectByID = `SELECT id, name, created_at, update_at FROM category WHERE id = $1`

// CategoryInsert 插入分类
const CategoryInsert = `INSERT INTO category (name, created_at, update_at) VALUES ($1, $2, $3) RETURNING id`

// CategorySelectAll 查询所有分类
const CategorySelectAll = `SELECT id, name, created_at, update_at FROM category ORDER BY created_at DESC`

// CategoryCountAll 统计所有分类数
const CategoryCountAll = `SELECT COUNT(*) FROM category`

// CategorySelectPaged 分页查询分类
const CategorySelectPaged = `SELECT id, name, created_at, update_at FROM category ORDER BY created_at DESC LIMIT $1 OFFSET $2`

// CategoryCountAllByName 按名称关键词统计分类数
const CategoryCountAllByName = `
SELECT COUNT(*)
FROM category
WHERE name ILIKE '%' || $1 || '%'
`

// CategorySelectPagedByName 按名称关键词分页查询分类
const CategorySelectPagedByName = `
SELECT id, name, created_at, update_at
FROM category
WHERE name ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

// CategoryDelete 删除分类
const CategoryDelete = `DELETE FROM category WHERE id = $1`

// CategoryUpdate 更新分类
const CategoryUpdate = `UPDATE category SET name = $1, update_at = $2 WHERE id = $3`
