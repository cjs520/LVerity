package database

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"LVerity/pkg/config"
	"LVerity/pkg/model"
)

var (
	// DB 全局数据库连接
	DB   *gorm.DB
	once sync.Once
)

// InitDB 初始化数据库连接
func InitDB(config *config.DatabaseConfig) error {
	var err error
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.User,
			config.Password,
			config.Host,
			config.Port,
			config.DBName)

		log.Printf("正在连接MySQL数据库: %s:%d/%s", config.Host, config.Port, config.DBName)
		// 连接MySQL数据库
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			SkipDefaultTransaction: true,
			PrepareStmt: true,
		})
		if err != nil {
			err = fmt.Errorf("连接MySQL数据库失败: %v", err)
			return
		}

		// 配置连接池
		sqlDB, err := DB.DB()
		if err != nil {
			err = fmt.Errorf("获取数据库连接池失败: %v", err)
			return
		}

		// 设置连接池参数
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(1000)
		sqlDB.SetConnMaxLifetime(time.Minute * 10)
		sqlDB.SetConnMaxIdleTime(time.Minute * 10)

		// 自动迁移数据库结构
		err = autoMigrate()
		if err != nil {
			err = fmt.Errorf("数据库迁移失败: %v", err)
			return
		}
		log.Println("MySQL数据库连接和迁移成功完成")
	})
	return err
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	if DB == nil {
		panic("数据库连接未初始化")
	}
	return DB
}

// SetDB 设置数据库连接（仅用于测试）
func SetDB(db *gorm.DB) {
	DB = db
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return fmt.Errorf("获取底层数据库连接失败: %v", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("关闭数据库连接失败: %v", err)
		}
		DB = nil
		log.Println("数据库连接已关闭")
	}
	return nil
}

// autoMigrate 自动迁移数据库结构并初始化数据
func autoMigrate() error {
    log.Println("开始数据库迁移...")

    // 先删除外键约束
    if err := DB.Migrator().DropConstraint(&model.User{}, "users_ibfk_1"); err != nil {
        log.Printf("Warning: Failed to drop foreign key constraint: %v", err)
    }

    // 删除已存在的索引
    if err := DB.Migrator().DropIndex(&model.User{}, "idx_username"); err != nil {
        log.Printf("Warning: Failed to drop index: %v", err)
    }

    // 基础表
    if err := DB.AutoMigrate(
        &model.User{},
        &model.Role{},
        &model.Permission{},
        &model.Device{},
        &model.LicenseTag{},
    ); err != nil {
        return fmt.Errorf("迁移基础模型失败: %v", err)
    }

    // 关联表
    if err := DB.AutoMigrate(
        &model.RolePermission{},
        &model.UserRole{},
        &model.License{},
        &model.LicenseUsage{},
    ); err != nil {
        return fmt.Errorf("迁移关联模型失败: %v", err)
    }

    // 其他表
    if err := DB.AutoMigrate(
        &model.DeviceLocation{},
        &model.AbnormalBehavior{},
        &model.BlacklistRule{},
    ); err != nil {
        return fmt.Errorf("迁移其他模型失败: %v", err)
    }

    log.Println("数据库迁移完成")

    // 检查是否需要初始化数据
    var count int64
    if err := DB.Model(&model.Role{}).Count(&count).Error; err != nil {
        return fmt.Errorf("检查角色表失败: %v", err)
    }

    // 如果没有角色数据，则进行初始化
    if count == 0 {
        log.Println("开始初始化数据...")

        // 插入默认角色
        roles := []model.Role{
            {ID: "1", Name: "admin", Description: "系统管理员", Type: "admin"},
            {ID: "2", Name: "operator", Description: "运营人员", Type: "operator"},
            {ID: "3", Name: "viewer", Description: "查看者", Type: "viewer"},
        }
        if err := DB.Create(&roles).Error; err != nil {
            return fmt.Errorf("初始化角色失败: %v", err)
        }

        // 插入默认权限
        permissions := []model.Permission{
            {ID: "1", Resource: "system", Action: "admin"},
            {ID: "2", Resource: "license", Action: "manage"},
            {ID: "3", Resource: "device", Action: "manage"},
            {ID: "4", Resource: "user", Action: "manage"},
            {ID: "5", Resource: "alert", Action: "manage"},
            {ID: "6", Resource: "log", Action: "view"},
        }
        if err := DB.Create(&permissions).Error; err != nil {
            return fmt.Errorf("初始化权限失败: %v", err)
        }

        // 插入角色权限关联
        rolePermissions := []model.RolePermission{
            // 管理员拥有所有权限
            {RoleID: "1", PermissionID: "1"}, {RoleID: "1", PermissionID: "2"},
            {RoleID: "1", PermissionID: "3"}, {RoleID: "1", PermissionID: "4"},
            {RoleID: "1", PermissionID: "5"}, {RoleID: "1", PermissionID: "6"},
            // 运营人员拥有部分权限
            {RoleID: "2", PermissionID: "2"}, {RoleID: "2", PermissionID: "3"},
            {RoleID: "2", PermissionID: "5"},
            // 查看者只有查看权限
            {RoleID: "3", PermissionID: "6"},
        }
        if err := DB.Create(&rolePermissions).Error; err != nil {
            return fmt.Errorf("初始化角色权限关联失败: %v", err)
        }

        // 插入默认管理员用户
        adminUser := model.User{
            ID:       "1",
            Username: "admin",
            Password: "240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9",
            RoleID:   "1",
            Status:   "active",
            Salt:     []byte("123456"),
        }
        if err := DB.Create(&adminUser).Error; err != nil {
            return fmt.Errorf("初始化管理员用户失败: %v", err)
        }

        log.Println("数据初始化完成")
    }

    return nil
}
