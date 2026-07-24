package model

import "time"

// Menu 菜单
type Menu struct {
	ID             int64     `json:"id" db:"id"`
	ParentID       int64     `json:"parentId" db:"parent_id"`
	Title          string    `json:"title" db:"title"`
	Path           string    `json:"path" db:"path"`
	Component      string    `json:"component" db:"component"`
	Icon           string    `json:"icon" db:"icon"`
	PermissionCode string    `json:"permissionCode" db:"permission_code"`
	SortOrder      int       `json:"sortOrder" db:"sort_order"`
	Visible        bool      `json:"visible" db:"visible"`
	PermissionType       string    `json:"menuType" db:"permission_type"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdateAt       time.Time `json:"updateAt" db:"update_at"`
	Children       []Menu    `json:"children,omitempty"`
}

// PermissionTree 构建菜单树
func BuildPermissionTree(menus []Menu, parentID int64) []Menu {
	var tree []Menu
	for _, m := range menus {
		if m.ParentID == parentID {
			m.Children = BuildPermissionTree(menus, m.ID)
			tree = append(tree, m)
		}
	}
	if tree == nil {
		tree = []Menu{}
	}
	return tree
}

// FilterPermissionByPermissions 根据权限过滤菜单树
func FilterPermissionByPermissions(tree []Menu, permissions []string) []Menu {
	permSet := make(map[string]bool)
	for _, p := range permissions {
		permSet[p] = true
	}

	var filtered []Menu
	for _, m := range tree {
		if !m.Visible {
			continue
		}
		// 先递归过滤子菜单
		m.Children = FilterPermissionByPermissions(m.Children, permissions)
		// 有子菜单的父级：子菜单非空就保留
		hasChildren := len(m.Children) > 0
		if hasChildren {
			filtered = append(filtered, m)
			continue
		}
		// 无子菜单的叶子节点：检查自身权限
		if m.PermissionCode == "" || permSet[m.PermissionCode] {
			filtered = append(filtered, m)
		}
	}
	if filtered == nil {
		filtered = []Menu{}
	}
	return filtered
}
