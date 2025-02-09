package service

import (
	"fmt"
	"LVerity/pkg/config"
	"github.com/sirupsen/logrus"
)

var jwtSecret []byte

// InitAuth 初始化认证相关配置
func InitAuth() error {
	logrus.Info("开始初始化认证配置...")

	// 验证JWT配置
	if config.GlobalConfig.JWT.Secret == "" {
		return fmt.Errorf("JWT密钥未配置")
	}

	// 初始化JWT密钥
	jwtSecret = []byte(config.GlobalConfig.JWT.Secret)

	logrus.Info("认证配置初始化完成")
	return nil
}