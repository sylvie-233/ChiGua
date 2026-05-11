package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
	"chigua-backend/internal/sql"
	"errors"
	"time"
)

var ErrTagExists = errors.New("标签名称已存在")

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

func GetAllTags() ([]model.Tag, error) {
	var tags []model.Tag
	err := database.DB.Select(&tags, sql.TagSelectAll)
	return tags, err
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
