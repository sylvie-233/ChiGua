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

func GetUserList(page, pageSize int, keyword string) (*model.UserListResponse, error) {
	var total int64
	var users []model.UserResponse
	var err error

	offset := (page - 1) * pageSize

	if keyword != "" {
		err = database.DB.Get(&total, sql.UserCountAllByKeyword, keyword)
		if err != nil {
			return nil, err
		}
		err = database.DB.Select(&users, sql.UserSelectAllByKeyword, keyword, pageSize, offset)
	} else {
		err = database.DB.Get(&total, sql.UserCountAll)
		if err != nil {
			return nil, err
		}
		err = database.DB.Select(&users, sql.UserSelectAll, pageSize, offset)
	}

	if err != nil {
		return nil, err
	}

	totalPages := (int(total) + pageSize - 1) / pageSize

	return &model.UserListResponse{
		PageResponse: model.PageResponse{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
		Items: users,
	}, nil
}

func DeleteUserByID(id int64) error {
	_, err := database.DB.Exec(sql.UserDelete, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserRole(id int64, role string) error {
	now := time.Now()
	_, err := database.DB.Exec(sql.UserUpdateRole, role, now, id)
	if err != nil {
		return err
	}
	return nil
}

func CreateUser(username, password, nickname, role string) (*model.User, error) {
	var count int
	err := database.DB.Get(&count, sql.UserCheckExists, username)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrUserExists
	}

	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	now := time.Now()
	newUser := model.User{
		Username:  username,
		Password:  hashedPassword,
		Nickname:  nickname,
		Role:      role,
		CreatedAt: now,
		UpdateAt:  now,
	}

	err = database.DB.QueryRow(sql.UserInsert, newUser.Username, newUser.Password, newUser.Nickname, newUser.Role, newUser.CreatedAt, newUser.UpdateAt).Scan(&newUser.ID)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func AdminUpdateUser(id int64, nickname, role string) error {
	now := time.Now()
	var nicknamePtr *string
	if nickname != "" {
		nicknamePtr = &nickname
	}
	var rolePtr *string
	if role != "" {
		rolePtr = &role
	}

	_, err := database.DB.Exec(sql.AdminUserUpdate, nicknamePtr, rolePtr, now, id)
	if err != nil {
		return err
	}
	return nil
}

// CountUser 统计用户总数
func CountUser() (int64, error) {
	var count int64
	err := database.DB.Get(&count, sql.UserCountAll)
	if err != nil {
		return 0, err
	}
	return count, nil
}
