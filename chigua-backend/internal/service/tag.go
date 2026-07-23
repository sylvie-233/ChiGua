package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
	"chigua-backend/internal/sql"
	"chigua-backend/utils/logger"
	"errors"
	"time"
)

var ErrTagExists = errors.New("标签名称已存在")
var ErrTagNotFound = errors.New("标签不存在")

func CreateTag(tag model.TagCreate) (*model.Tag, error) {
	var exists bool
	err := database.DB.QueryRow(sql.TagCheckExists, tag.Name).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrTagExists
	}

	now := time.Now()
	newTag := model.Tag{
		Name:      tag.Name,
		CreatedAt: now,
		UpdateAt:  now,
	}

	err = database.DB.QueryRow(sql.TagInsert, newTag.Name, newTag.CreatedAt, newTag.UpdateAt).Scan(&newTag.ID)
	if err != nil {
		return nil, err
	}

	return &newTag, nil
}

func UpdateTag(id int64, name string) (*model.Tag, error) {
	var exists bool
	err := database.DB.QueryRow(sql.TagCheckExistsExcludeID, name, id).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrTagExists
	}

	now := time.Now()
	_, err = database.DB.Exec(sql.TagUpdate, name, now, id)
	if err != nil {
		return nil, err
	}

	var tag model.Tag
	err = database.DB.Get(&tag, sql.TagSelectByID, id)
	if err != nil {
		logger.Error(err)
		return nil, ErrTagNotFound
	}

	return &tag, nil
}

func GetAllTags() ([]model.Tag, error) {
	var tags []model.Tag
	err := database.DB.Select(&tags, sql.TagSelectAll)
	return tags, err
}

func GetTagList(page, pageSize int, keyword string) (*model.TagListResponse, error) {
	var total int64
	var tags []model.Tag
	var err error

	offset := (page - 1) * pageSize

	if keyword != "" {
		err = database.DB.Get(&total, sql.TagCountAllByName, keyword)
		if err != nil {
			return nil, err
		}
		err = database.DB.Select(&tags, sql.TagSelectPagedByName, keyword, pageSize, offset)
	} else {
		err = database.DB.Get(&total, sql.TagCountAll)
		if err != nil {
			return nil, err
		}
		err = database.DB.Select(&tags, sql.TagSelectPaged, pageSize, offset)
	}

	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page < totalPages
	hasPrev := page > 1

	return &model.TagListResponse{
		PageResponse: model.PageResponse{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
		Items: tags,
	}, nil
}

func DeleteTag(id int64) error {
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(sql.TagDeleteRelate, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(sql.TagDelete, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// CountTag 统计标签总数
func CountTag() (int64, error) {
	var count int64
	err := database.DB.Get(&count, sql.TagCountAll)
	if err != nil {
		return 0, err
	}
	return count, nil
}
