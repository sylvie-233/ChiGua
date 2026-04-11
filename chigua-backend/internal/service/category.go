package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
	"time"
)

func CreateCategory(category model.CategoryCreate) (*model.Category, error) {
	now := time.Now()
	newCategory := model.Category{
		Name:      category.Name,
		CreatedAt: now,
		UpdateAt:  now,
	}

	// 插入分类
	query := `INSERT INTO category (name, created_at, update_at) VALUES ($1, $2, $3) RETURNING id`
	err := database.DB.QueryRow(query, newCategory.Name, newCategory.CreatedAt, newCategory.UpdateAt).Scan(&newCategory.ID)
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
