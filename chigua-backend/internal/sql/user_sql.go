package sql

// UserSQL 用户相关SQL常量

// UserCheckExists 检查用户名是否存在
const UserCheckExists = `SELECT COUNT(*) FROM users WHERE username = $1`

// UserInsert 插入用户
const UserInsert = `INSERT INTO users (username, password, nickname, role, created_at, update_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

// UserSelectByUP 按用户名密码查询用户
const UserSelectByUP = `SELECT id, username, password, nickname, role, created_at, update_at FROM users WHERE username = $1 AND password = $2`

// UserSelectByID 按ID查询用户
const UserSelectByID = `SELECT id, username, password, nickname, role, created_at, update_at FROM users WHERE id = $1`

// UserUpdateNickname 更新用户昵称
const UserUpdateNickname = `UPDATE users SET nickname = $1, update_at = $2 WHERE id = $3`

// UserUpdate 更新用户信息
const UserUpdate = `UPDATE users SET nickname = COALESCE($1, nickname), update_at = $2 WHERE id = $3`

// UserCountAll 统计所有用户数量
const UserCountAll = `SELECT COUNT(*) FROM users`

// UserSelectAll 分页查询所有用户
const UserSelectAll = `SELECT id, username, nickname, role, created_at, update_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`

// UserCountAllByKeyword 按关键词（用户名或昵称）统计用户数
const UserCountAllByKeyword = `
SELECT COUNT(*)
FROM users
WHERE username ILIKE '%' || $1 || '%' OR nickname ILIKE '%' || $1 || '%'
`

// UserSelectAllByKeyword 按关键词（用户名或昵称）分页查询用户
const UserSelectAllByKeyword = `
SELECT id, username, nickname, role, created_at, update_at
FROM users
WHERE username ILIKE '%' || $1 || '%' OR nickname ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

// UserDelete 删除用户
const UserDelete = `DELETE FROM users WHERE id = $1`

// UserUpdateRole 更新用户角色
const UserUpdateRole = `UPDATE users SET role = $1, update_at = $2 WHERE id = $3`

// AdminUserUpdate 更新用户信息（管理员）
const AdminUserUpdate = `UPDATE users SET nickname = COALESCE($1, nickname), role = COALESCE($2, role), update_at = $3 WHERE id = $4`

// UserCheckExistsExcludeID 检查用户名是否存在（排除指定ID）
const UserCheckExistsExcludeID = `SELECT COUNT(*) FROM users WHERE username = $1 AND id != $2`
