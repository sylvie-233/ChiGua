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
