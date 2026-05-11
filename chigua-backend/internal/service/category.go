package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
	"chigua-backend/internal/sql"
	"errors"
	"time"
)

var ErrCategoryExists = errors.New("分类名称已存在")

func CreateCategory(category model.CategoryCreate) (*model.Category, error) {
	var exists bool
	err := database.DB.QueryRow(sql.CategoryCheckExists, category.Name).Scan(&exists)
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

	err = database.DB.QueryRow(sql.CategoryInsert, newCategory.Name, newCategory.CreatedAt, newCategory.UpdateAt).Scan(&newCategory.ID)
	if err != nil {
		return nil, err
	}

	return &newCategory, nil
}

func GetAllCategories() ([]model.Category, error) {
	var categories []model.Category
	err := database.DB.Select(&categories, sql.CategorySelectAll)
	return categories, err
}

func DeleteCategory(id int64) error {
	_, err := database.DB.Exec(sql.CategoryDelete, id)
	return err
}
