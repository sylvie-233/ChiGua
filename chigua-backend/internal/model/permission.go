package model

import "time"

// Role 角色
type Role struct {
	ID          int64     `json:"id" db:"id"`
	Code        string    `json:"code" db:"code"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdateAt    time.Time `json:"updateAt" db:"update_at"`
}
