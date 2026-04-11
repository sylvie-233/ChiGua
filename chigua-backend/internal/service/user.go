package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
)

func RegisterUser(user model.UserRegister) (*model.User, error) {
	// 检查用户名是否已存在
	var count int
	err := database.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE username = $1", user.Username)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 密码加密
	hasher := sha256.New()
	hasher.Write([]byte(user.Password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	// 创建用户
	now := time.Now()
	newUser := model.User{
		Username:  user.Username,
		Password:  hashedPassword,
		Nickname:  user.Nickname,
		Role:      "user",
		CreatedAt: now,
		UpdateAt:  now,
	}

	// 插入数据库
	query := `INSERT INTO users (username, password, nickname, role, created_at, update_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = database.DB.QueryRow(query, newUser.Username, newUser.Password, newUser.Nickname, newUser.Role, newUser.CreatedAt, newUser.UpdateAt).Scan(&newUser.ID)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func LoginUser(login model.UserLogin) (*model.User, error) {
	// 密码加密
	hasher := sha256.New()
	hasher.Write([]byte(login.Password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	// 查询用户
	var user model.User
	query := `SELECT id, username, password, nickname, role, created_at, update_at FROM users WHERE username = $1 AND password = $2`
	err := database.DB.Get(&user, query, login.Username, hashedPassword)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	return &user, nil
}

func GetUserByID(id int64) (*model.User, error) {
	var user model.User
	query := `SELECT id, username, password, nickname, role, created_at, update_at FROM users WHERE id = $1`
	err := database.DB.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
