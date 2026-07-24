package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
	"chigua-backend/internal/sql"
	"chigua-backend/utils/logger"
	"errors"
	"time"
)

var ErrCategoryExists = errors.New("分类名称已存在")
var ErrCategoryNotFound = errors.New("分类不存在")

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
		SortOrder: category.SortOrder,
		CreatedAt: now,
		UpdateAt:  now,
	}

	err = database.DB.QueryRow(sql.CategoryInsert, newCategory.Name, newCategory.SortOrder, newCategory.CreatedAt, newCategory.UpdateAt).Scan(&newCategory.ID)
	if err != nil {
		return nil, err
	}

	return &newCategory, nil
}

func UpdateCategory(id int64, name string, sortOrder int) (*model.Category, error) {
	var exists bool
	err := database.DB.QueryRow(sql.CategoryCheckExistsExcludeID, name, id).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCategoryExists
	}

	now := time.Now()
	_, err = database.DB.Exec(sql.CategoryUpdate, name, sortOrder, now, id)
	if err != nil {
		return nil, err
	}

	var category model.Category
	err = database.DB.Get(&category, sql.CategorySelectByID, id)
	if err != nil {
		logger.Error(err)
		return nil, ErrCategoryNotFound
	}

	return &category, nil
}

func GetAllCategories() ([]model.Category, error) {
	var categories []model.Category
	err := database.DB.Select(&categories, sql.CategorySelectAll)
	return categories, err
}

func GetCategoryList(page, pageSize int, keyword string) (*model.CategoryListResponse, error) {
	var total int64
	var categories []model.Category
	var err error

	offset := (page - 1) * pageSize

	if keyword != "" {
		err = database.DB.Get(&total, sql.CategoryCountAllByName, keyword)
		if err != nil {
			return nil, err
		}
		err = database.DB.Select(&categories, sql.CategorySelectPagedByName, keyword, pageSize, offset)
	} else {
		err = database.DB.Get(&total, sql.CategoryCountAll)
		if err != nil {
			return nil, err
		}
		err = database.DB.Select(&categories, sql.CategorySelectPaged, pageSize, offset)
	}

	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page < totalPages
	hasPrev := page > 1

	return &model.CategoryListResponse{
		PageResponse: model.PageResponse{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
		Items: categories,
	}, nil
}

func DeleteCategory(id int64) error {
	_, err := database.DB.Exec(sql.CategoryDelete, id)
	return err
}

// CountCategory 统计分类总数
func CountCategory() (int64, error) {
	var count int64
	err := database.DB.Get(&count, sql.CategoryCountAll)
	if err != nil {
		return 0, err
	}
	return count, nil
}
