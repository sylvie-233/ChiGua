package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
	"chigua-backend/internal/sql"
	"time"
)

// GetUserPermissions 查询用户拥有的所有权限码（通过角色→角色菜单→菜单）
func GetUserPermissions(userID int64) ([]string, error) {
	var codes []string
	err := database.DB.Select(&codes, sql.UserPermissionCodes, userID)
	if err != nil {
		return nil, err
	}
	if codes == nil {
		codes = []string{}
	}
	return codes, nil
}

// GetAdminPermissions admin 角色获取全部权限码
func GetAdminPermissions() ([]string, error) {
	var codes []string
	err := database.DB.Select(&codes, `SELECT DISTINCT permission_code FROM permission WHERE permission_code != ''`)
	if err != nil {
		return nil, err
	}
	return codes, nil
}

// GetUserRoles 查询用户拥有的角色code
func GetUserRoles(userID int64) ([]string, error) {
	var codes []string
	err := database.DB.Select(&codes, sql.UserRoleSelectByUser, userID)
	if err != nil {
		return nil, err
	}
	if codes == nil {
		codes = []string{}
	}
	return codes, nil
}

// GetAllRoles 查询所有角色
// GetRoleByCode 按 code 查询角色
func GetRoleByCode(code string) (*model.Role, error) {
	var role model.Role
	err := database.DB.Get(&role, sql.RoleSelectByCode, code)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func GetAllRoles() ([]model.Role, error) {
	var roles []model.Role
	err := database.DB.Select(&roles, sql.RoleSelectAll)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRoleList 分页查询角色
func GetRoleList(page, pageSize int, keyword string) ([]model.Role, int64, error) {
	var total int64
	var roles []model.Role
	var err error

	offset := (page - 1) * pageSize

	if keyword != "" {
		err = database.DB.Get(&total, sql.RoleCountByName, keyword)
		if err != nil {
			return nil, 0, err
		}
		err = database.DB.Select(&roles, sql.RoleSelectPagedByName, keyword, pageSize, offset)
	} else {
		err = database.DB.Get(&total, sql.RoleCountAll)
		if err != nil {
			return nil, 0, err
		}
		err = database.DB.Select(&roles, sql.RoleSelectPaged, pageSize, offset)
	}

	return roles, total, err
}

// CreateRole 新增角色
func CreateRole(code, name, description string) (*model.Role, error) {
	now := time.Now()
	var role model.Role
	err := database.DB.QueryRow(sql.RoleInsert, code, name, description, now, now).Scan(&role.ID)
	if err != nil {
		return nil, err
	}
	role.Code = code
	role.Name = name
	role.Description = description
	role.CreatedAt = now
	role.UpdateAt = now
	return &role, nil
}

// UpdateRoleInfo 更新角色信息
func UpdateRoleInfo(id int64, name, description string) error {
	now := time.Now()
	_, err := database.DB.Exec(sql.RoleUpdate, name, description, now, id)
	return err
}

// DeleteRole 删除角色
func DeleteRole(id int64) error {
	_, err := database.DB.Exec(sql.RoleDelete, id)
	return err
}

// GetRoleMenuIDs 查询角色拥有的菜单ID列表
func GetRoleMenuIDs(roleID int64) ([]int64, error) {
	var ids []int64
	err := database.DB.Select(&ids, sql.RoleMenuSelect, roleID)
	if err != nil {
		return nil, err
	}
	if ids == nil {
		ids = []int64{}
	}
	return ids, nil
}

// UpdateRoleMenus 更新角色菜单权限（先删后插）
func UpdateRoleMenus(roleID int64, menuIDs []int64) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(sql.RoleMenuDeleteAll, roleID)
	if err != nil {
		return err
	}

	for _, menuID := range menuIDs {
		_, err = tx.Exec(sql.RoleMenuInsert, roleID, menuID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// SetUserRoles 设置用户角色（先删后插）
// AssignUserRole 给用户分配单个角色
func AssignUserRole(userID, roleID int64) error {
	_, err := database.DB.Exec(sql.UserRoleInsert, userID, roleID)
	return err
}

func SetUserRoles(userID int64, roleIDs []int64) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(sql.UserRoleDeleteAll, userID)
	if err != nil {
		return err
	}

	for _, roleID := range roleIDs {
		_, err = tx.Exec(sql.UserRoleInsert, userID, roleID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// HasPermission 检查用户是否拥有指定权限
func HasPermission(permissions []string, code string) bool {
	for _, p := range permissions {
		if p == code {
			return true
		}
	}
	return false
}

// HasAnyPermission 检查用户是否拥有任一权限
func HasAnyPermission(permissions []string, codes ...string) bool {
	for _, code := range codes {
		if HasPermission(permissions, code) {
			return true
		}
	}
	return false
}
