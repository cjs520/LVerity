package service

import (
	"github.com/sirupsen/logrus"
	"LVerity/pkg/model"
	"time"
)

// InitAdminUser 初始化管理员用户
func InitAdminUser() error {
	logrus.Info("开始初始化管理员用户...")

	// 检查是否已存在管理员用户
	exists, err := CheckUsernameExists("admin")
	if err != nil {
		logrus.WithError(err).Error("检查管理员用户是否存在时发生错误")
		return err
	}
	if exists {
		logrus.Info("管理员用户已存在，跳过初始化")
		return nil
	}

	// 获取管理员角色ID
	adminRoleID := "1"
	logrus.WithField("roleID", adminRoleID).Info("使用管理员角色ID创建用户")

	// 创建管理员用户
	admin, err := CreateUser("admin", "admin123", adminRoleID)
	if err != nil {
		logrus.WithError(err).Error("创建管理员用户失败")
		return err
	}

	// 设置管理员状态为激活
	admin.Status = model.UserStatusActive
	admin.CreateTime = time.Now()
	admin.UpdateTime = time.Now()

	if err := UpdateUser(admin); err != nil {
		logrus.WithError(err).Error("更新管理员用户状态失败")
		return err
	}

	logrus.WithField("userID", admin.ID).Info("成功创建管理员用户")
	return nil
}
