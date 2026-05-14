package service

import (
	"chigua-backend/database"
	"chigua-backend/internal/model"
	"chigua-backend/internal/sql"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
)

var ErrUserExists = errors.New("用户名已存在")

func RegisterUser(user model.UserRegister) (*model.User, error) {
	var count int
	err := database.DB.Get(&count, sql.UserCheckExists, user.Username)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrUserExists
	}

	hasher := sha256.New()
	hasher.Write([]byte(user.Password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	now := time.Now()
	newUser := model.User{
		Username:  user.Username,
		Password:  hashedPassword,
		Nickname:  user.Nickname,
		Role:      "user",
		CreatedAt: now,
		UpdateAt:  now,
	}

	err = database.DB.QueryRow(sql.UserInsert, newUser.Username, newUser.Password, newUser.Nickname, newUser.Role, newUser.CreatedAt, newUser.UpdateAt).Scan(&newUser.ID)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func LoginUser(login model.UserLogin) (*model.UserResponse, error) {
	hasher := sha256.New()
	hasher.Write([]byte(login.Password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	var user model.User
	err := database.DB.Get(&user, sql.UserSelectByUP, login.Username, hashedPassword)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	return user.ToResponse(), nil
}

func GetUserByID(id int64) (*model.User, error) {
	var user model.User
	err := database.DB.Get(&user, sql.UserSelectByID, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUserNickname(userID int64, nickname string) error {
	now := time.Now()
	_, err := database.DB.Exec(sql.UserUpdateNickname, nickname, now, userID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(userID int64, updates map[string]interface{}) error {
	now := time.Now()

	var nickname *string
	if val, ok := updates["nickname"]; ok && val != nil {
		nick := val.(string)
		nickname = &nick
	}

	_, err := database.DB.Exec(sql.UserUpdate, nickname, now, userID)
	if err != nil {
		return err
	}
	return nil
}
