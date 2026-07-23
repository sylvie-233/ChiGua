package model

import "time"

type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username" binding:"required"`
	Password  string    `json:"password" db:"password" binding:"required"`
	Nickname  string    `json:"nickname" db:"nickname"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdateAt  time.Time `json:"updateAt" db:"update_at"`
}

// ToResponse 将 User 转换为 UserResponse，不包含密码字段
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Nickname:  u.Nickname,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}

type UserRegister struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
}

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Nickname  string    `json:"nickname" db:"nickname"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdateAt  time.Time `json:"updateAt" db:"update_at"`
}

type UserListResponse struct {
	PageResponse
	Items []UserResponse `json:"items"`
}
