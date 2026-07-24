package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
	"chigua-backend/internal/sql"
	"errors"
	"strings"
	"time"
)

var ErrPermissionCodeRequired = errors.New("权限标识不能为空")
var ErrPermissionCodeExists = errors.New("权限标识已存在")
var ErrPermissionTitleRequired = errors.New("名称不能为空")

// GetPermissionTree 获取菜单树（按权限过滤）
func GetPermissionTree(permissions []string) ([]model.Menu, error) {
	var menus []model.Menu
	err := database.DB.Select(&menus, sql.PermissionSelectAll)
	if err != nil {
		return nil, err
	}

	tree := model.BuildPermissionTree(menus, 0)
	if len(permissions) > 0 {
		tree = model.FilterPermissionByPermissions(tree, permissions)
	}
	return tree, nil
}

// GetAllPermissions 获取所有菜单（管理用，含不可见）
func GetAllPermissions() ([]model.Menu, error) {
	var menus []model.Menu
	err := database.DB.Select(&menus, sql.PermissionSelectAllAdmin)
	if err != nil {
		return nil, err
	}
	return menus, nil
}

// CreateMenu 创建菜单
func CreateMenu(menu model.Menu) (*model.Menu, error) {
	if strings.TrimSpace(menu.Title) == "" {
		return nil, ErrPermissionTitleRequired
	}
	if strings.TrimSpace(menu.PermissionCode) == "" {
		return nil, ErrPermissionCodeRequired
	}

	// 检查唯一性
	var count int
	err := database.DB.Get(&count, `SELECT COUNT(*) FROM permission WHERE permission_code = $1`, menu.PermissionCode)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrPermissionCodeExists
	}

	now := time.Now()
	err = database.DB.QueryRow(sql.PermissionInsert,
		menu.ParentID, menu.Title, menu.Path, menu.Component,
		menu.Icon, menu.PermissionCode, menu.SortOrder,
		menu.Visible, menu.PermissionType, now, now,
	).Scan(&menu.ID)
	if err != nil {
		return nil, err
	}
	menu.CreatedAt = now
	menu.UpdateAt = now
	return &menu, nil
}

// UpdateMenu 更新菜单
func UpdateMenu(id int64, menu model.Menu) error {
	if strings.TrimSpace(menu.Title) == "" {
		return ErrPermissionTitleRequired
	}
	if strings.TrimSpace(menu.PermissionCode) == "" {
		return ErrPermissionCodeRequired
	}

	// 检查唯一性（排除自身）
	var count int
	err := database.DB.Get(&count, `SELECT COUNT(*) FROM permission WHERE permission_code = $1 AND id != $2`, menu.PermissionCode, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrPermissionCodeExists
	}

	now := time.Now()
	_, err = database.DB.Exec(sql.PermissionUpdate,
		menu.ParentID, menu.Title, menu.Path, menu.Component,
		menu.Icon, menu.PermissionCode, menu.SortOrder,
		menu.Visible, menu.PermissionType, now, id,
	)
	return err
}

// DeleteMenu 删除菜单
func DeleteMenu(id int64) error {
	_, err := database.DB.Exec(sql.PermissionDelete, id)
	return err
}
