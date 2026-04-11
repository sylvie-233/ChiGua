package model

import "time"

type Category struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" binding:"required"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdateAt  time.Time `json:"updateAt" db:"update_at"`
}

type CategoryCreate struct {
	Name string `json:"name" binding:"required"`
}
