package service

import (
	"LVerity/pkg/common"
	"LVerity/pkg/config"
	"LVerity/pkg/database"
	"LVerity/pkg/model"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrUserNotFound    = errors.New("用户不存在")
	ErrInvalidPassword = errors.New("密码错误")
	ErrUsernameTaken   = errors.New("用户名已被使用")
	ErrInvalidRoleID   = errors.New("无效的角色ID")
)

// GetUserByID 通过ID获取用户
func GetUserByID(id string) (*model.User, error) {
	var user model.User
	result := database.GetDB().Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

// GetUserByUsername 通过用户名获取用户
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	result := database.GetDB().Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

// Login 用户登录
func Login(username, password string) (string, *model.User, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return "", nil, ErrUserNotFound
	}

	// 验证密码
	if !user.CheckPassword(password) {
		return "", nil, ErrInvalidPassword
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLogin = &now
	err = UpdateUser(user)
	if err != nil {
		return "", nil, err
	}

	// 生成JWT令牌
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RoleID:   user.RoleID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "LVerity",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}

// CreateUser 创建新用户
func CreateUser(username, password, roleID string) (*model.User, error) {
	// 检查用户名是否已存在
	exists, err := CheckUsernameExists(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameTaken
	}

	// 检查角色ID是否有效
	if !IsValidRoleID(roleID) {
		return nil, ErrInvalidRoleID
	}

	// 创建用户
	user := &model.User{
		ID:        common.GenerateUUID(),
		Username:  username,
		RoleID:    roleID,
		Status:    model.UserStatusActive,
		LastLogin: nil, // 使用*time.Time类型，允许为nil
	}

	// 设置密码
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	// 保存到数据库
	if err := database.GetDB().Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(user *model.User) error {
	return database.GetDB().Save(user).Error
}

// UpdateUserProfile 更新用户信息
func UpdateUserProfile(id string, updates map[string]interface{}) (*model.User, error) {
	user, err := GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if username, ok := updates["username"].(string); ok {
		if username != user.Username {
			exists, err := CheckUsernameExists(username)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, ErrUsernameTaken
			}
			user.Username = username
		}
	}

	if roleID, ok := updates["role_id"].(string); ok {
		if !IsValidRoleID(roleID) {
			return nil, ErrInvalidRoleID
		}
		user.RoleID = roleID
	}

	if status, ok := updates["status"].(model.UserStatus); ok {
		user.Status = status
	}

	if err := database.GetDB().Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword 修改密码
func ChangePassword(id string, oldPassword, newPassword string) error {
	user, err := GetUserByID(id)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !user.CheckPassword(oldPassword) {
		return ErrInvalidPassword
	}

	// 设置新密码
	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	// 保存到数据库
	return database.GetDB().Save(user).Error
}

// DeleteUser 删除用户
func DeleteUser(id string) error {
	return database.GetDB().Delete(&model.User{}, "id = ?", id).Error
}

// CheckUsernameExists 检查用户名是否已存在
func CheckUsernameExists(username string) (bool, error) {
	var count int64
	err := database.GetDB().Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsValidRoleID 检查角色ID是否有效
func IsValidRoleID(roleID string) bool {
	// TODO: 实现角色ID验证逻辑
	return true
}

// ListUsers 获取用户列表
func ListUsers(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	tx := database.GetDB().Model(&model.User{})
	err := tx.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = tx.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
