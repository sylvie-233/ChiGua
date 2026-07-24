package sql

// --- 用户权限（通过角色菜单） ---

// UserPermissionCodes 查询用户权限码
const UserPermissionCodes = `
SELECT DISTINCT m.permission_code
FROM permission m
JOIN role_permission rm ON m.id = rm.menu_id
JOIN user_role ur ON rm.role_id = ur.role_id
WHERE ur.user_id = $1 AND m.permission_code != ''
ORDER BY m.permission_code
`

// --- 角色 ---

// RoleSelectAll 查询所有角色
const RoleSelectAll = `SELECT id, code, name, description, created_at, update_at FROM role ORDER BY id`

// RoleCountAll 统计角色数
const RoleCountAll = `SELECT COUNT(*) FROM role`

// RoleSelectPaged 分页查询角色
const RoleSelectPaged = `SELECT id, code, name, description, created_at, update_at FROM role ORDER BY id LIMIT $1 OFFSET $2`

// RoleCountByName 按名称搜索角色数
const RoleCountByName = `SELECT COUNT(*) FROM role WHERE name ILIKE '%' || $1 || '%' OR code ILIKE '%' || $1 || '%'`

// RoleSelectPagedByName 按名称搜索角色分页
const RoleSelectPagedByName = `SELECT id, code, name, description, created_at, update_at FROM role WHERE name ILIKE '%' || $1 || '%' OR code ILIKE '%' || $1 || '%' ORDER BY id LIMIT $2 OFFSET $3`

// RoleInsert 新增角色
const RoleInsert = `INSERT INTO role (code, name, description, created_at, update_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

// RoleUpdate 更新角色
const RoleUpdate = `UPDATE role SET name = $1, description = $2, update_at = $3 WHERE id = $4`

// RoleDelete 删除角色
const RoleDelete = `DELETE FROM role WHERE id = $1`

// RoleSelectByCode 按code查询角色
const RoleSelectByCode = `SELECT id, code, name, description, created_at, update_at FROM role WHERE code = $1`

// --- 用户角色 ---

// UserRoleSelectByUser 查询用户角色code
const UserRoleSelectByUser = `
SELECT r.code FROM role r JOIN user_role ur ON r.id = ur.role_id
WHERE ur.user_id = $1 ORDER BY r.id
`

// UserRoleInsert 分配角色
const UserRoleInsert = `INSERT INTO user_role (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

// UserRoleDeleteAll 移除所有角色
const UserRoleDeleteAll = `DELETE FROM user_role WHERE user_id = $1`

// --- 角色菜单管理 ---

// RoleMenuSelect 查询角色菜单ID
const RoleMenuSelect = `SELECT menu_id FROM role_permission WHERE role_id = $1`

// RoleMenuDeleteAll 删除角色所有菜单
const RoleMenuDeleteAll = `DELETE FROM role_permission WHERE role_id = $1`

// RoleMenuInsert 给角色分配菜单
const RoleMenuInsert = `INSERT INTO role_permission (role_id, menu_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
