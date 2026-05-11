package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
	"errors"
	"time"
)

var ErrCategoryExists = errors.New("分类名称已存在")

func CreateCategory(category model.CategoryCreate) (*model.Category, error) {
	// 检查名称是否已存在
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM category WHERE name = $1)`
	err := database.DB.QueryRow(checkQuery, category.Name).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCategoryExists
	}

	now := time.Now()
	newCategory := model.Category{
		Name:      category.Name,
		CreatedAt: now,
		UpdateAt:  now,
	}

	// 插入分类
	query := `INSERT INTO category (name, created_at, update_at) VALUES ($1, $2, $3) RETURNING id`
	err = database.DB.QueryRow(query, newCategory.Name, newCategory.CreatedAt, newCategory.UpdateAt).Scan(&newCategory.ID)
	if err != nil {
		return nil, err
	}

	return &newCategory, nil
}

func GetAllCategories() ([]model.Category, error) {
	var categories []model.Category
	err := database.DB.Select(&categories, "SELECT id, name, created_at, update_at FROM category ORDER BY created_at DESC")
	return categories, err
}

func DeleteCategory(id int64) error {
	// 删除分类
	_, err := database.DB.Exec(`DELETE FROM category WHERE id = $1`, id)
	return err
}
