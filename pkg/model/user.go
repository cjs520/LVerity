package model

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"

	"gorm.io/gorm"
)

// UserStatus 用户状态
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"   // 启用
	UserStatusInactive UserStatus = "inactive" // 禁用
	UserStatusBlocked  UserStatus = "blocked"  // 封禁
)

// User 用户模型
type User struct {
	ID         string     `json:"id" gorm:"primaryKey"`
	Username   string     `json:"username" gorm:"uniqueIndex:idx_username,length:191;type:varchar(191)"`
	Password   string     `json:"-" gorm:"not null;type:varchar(191)"` // 密码不返回给前端
	Salt       []byte     `json:"-"`                                   // 密码盐值
	RoleID     string     `json:"role_id" gorm:"type:varchar(191)"`    // 关联角色ID
	Status     UserStatus `json:"status" gorm:"type:varchar(20)"`      // true: 启用, false: 禁用
	LastLogin  *time.Time  `json:"last_login"`
	CreateTime time.Time  `json:"create_time"`
	UpdateTime time.Time  `json:"update_time"`
	MFASecret  string     `json:"mfa_secret" gorm:"type:varchar(255)"`
	MFAEnabled bool       `json:"mfa_enabled" gorm:"default:false"`
}

// SetPassword 设置密码
func (u *User) SetPassword(password string) error {
	// 生成随机盐值
	u.Salt = make([]byte, 16)
	_, err := rand.Read(u.Salt)
	if err != nil {
		return err
	}

	// 使用SHA256计算密码哈希
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write(u.Salt)
	u.Password = hex.EncodeToString(hash.Sum(nil))
	return nil
}

// CheckPassword 检查密码
func (u *User) CheckPassword(password string) bool {
	// 添加日志记录，输出密码验证的详细信息
	log.Printf("正在验证用户 %s 的密码", u.Username)
	log.Printf("数据库中存储的密码哈希值: %s", u.Password)

	// 首先尝试直接使用SHA256验证（用于处理旧格式的密码）
	hash := sha256.New()
	hash.Write([]byte(password))
	directHash := hex.EncodeToString(hash.Sum(nil))

	if directHash == u.Password {
		log.Printf("使用旧格式验证密码成功")
		return true
	}

	// 如果直接验证失败，且存在盐值，则使用盐值进行验证
	if len(u.Salt) > 0 {
		hash := sha256.New()
		hash.Write([]byte(password))
		hash.Write(u.Salt)
		calculatedHash := hex.EncodeToString(hash.Sum(nil))

		if calculatedHash == u.Password {
			log.Printf("使用盐值验证密码成功")
			return true
		}
	}

	log.Printf("密码验证失败: 哈希值不匹配")
	return false
}

// BeforeCreate GORM 创建钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Status == "" {
		u.Status = UserStatusActive
	}
	if u.CreateTime.IsZero() {
		u.CreateTime = time.Now()
	}
	if u.UpdateTime.IsZero() {
		u.UpdateTime = time.Now()
	}
	return nil
}

// BeforeUpdate GORM 更新钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdateTime = time.Now()
	return nil
}
