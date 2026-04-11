package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
	"time"
)

func CreateTag(tag model.TagCreate) (*model.Tag, error) {
	now := time.Now()
	newTag := model.Tag{
		Name:      tag.Name,
		CreatedAt: now,
		UpdateAt:  now,
	}

	// 插入标签
	query := `INSERT INTO tag (name, created_at, update_at) VALUES ($1, $2, $3) RETURNING id`
	err := database.DB.QueryRow(query, newTag.Name, newTag.CreatedAt, newTag.UpdateAt).Scan(&newTag.ID)
	if err != nil {
		return nil, err
	}

	return &newTag, nil
}

func GetAllTags() ([]model.Tag, error) {
	var tags []model.Tag
	err := database.DB.Select(&tags, "SELECT id, name, created_at, update_at FROM tag ORDER BY created_at DESC")
	return tags, err
}

func DeleteTag(id int64) error {
	// 开始事务
	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 删除标签关联
	_, err = tx.Exec(`DELETE FROM article_tag WHERE tag_id = $1`, id)
	if err != nil {
		return err
	}

	// 删除标签
	_, err = tx.Exec(`DELETE FROM tag WHERE id = $1`, id)
	if err != nil {
		return err
	}

	// 提交事务
	return tx.Commit()
}
