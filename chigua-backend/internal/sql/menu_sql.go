package sql

// PermissionSelectAll 查询所有可见菜单（按排序，排除按钮）
const PermissionSelectAll = `
SELECT id, parent_id, title, path, component, icon, permission_code, sort_order, visible, permission_type, created_at, update_at
FROM permission
WHERE visible = TRUE AND permission_type != 'B'
ORDER BY sort_order, id
`

// PermissionSelectAllAdmin 查询所有菜单（含不可见，管理用）
const PermissionSelectAllAdmin = `
SELECT id, parent_id, title, path, component, icon, permission_code, sort_order, visible, permission_type, created_at, update_at
FROM permission
ORDER BY sort_order, id
`

// PermissionInsert 新增菜单
const PermissionInsert = `
INSERT INTO permission (parent_id, title, path, component, icon, permission_code, sort_order, visible, permission_type, created_at, update_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id
`

// PermissionUpdate 更新菜单
const PermissionUpdate = `
UPDATE permission SET parent_id=$1, title=$2, path=$3, component=$4, icon=$5, permission_code=$6, sort_order=$7, visible=$8, permission_type=$9, update_at=$10
WHERE id=$11
`

// PermissionDelete 删除菜单
const PermissionDelete = `DELETE FROM permission WHERE id = $1`
