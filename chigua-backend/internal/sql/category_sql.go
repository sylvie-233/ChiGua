package sql

// CategorySQL 分类相关SQL常量

// CategoryCheckExists 检查分类名称是否存在
const CategoryCheckExists = `SELECT EXISTS(SELECT 1 FROM category WHERE name = $1)`

// CategoryInsert 插入分类
const CategoryInsert = `INSERT INTO category (name, created_at, update_at) VALUES ($1, $2, $3) RETURNING id`

// CategorySelectAll 查询所有分类
const CategorySelectAll = `SELECT id, name, created_at, update_at FROM category ORDER BY created_at DESC`

// CategoryDelete 删除分类
const CategoryDelete = `DELETE FROM category WHERE id = $1`
