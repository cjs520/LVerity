package main

import (
	"LVerity/pkg/common"
	"LVerity/pkg/config"
	"LVerity/pkg/database"
	"LVerity/pkg/model"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	dbConfig := &database.Config{
		DBPath: "data/lverity.db",
	}
	if err := database.InitDB(dbConfig); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	db := database.GetDB()

	// 检查管理员角色是否存在
	var adminRole model.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		// 创建管理员角色
		adminRole = model.Role{
			ID:          common.GenerateUUID(),
			Name:        "admin",
			Type:        "system",
			Description: "System Administrator",
		}
		
		if err := db.Create(&adminRole).Error; err != nil {
			log.Printf("Warning: Failed to create admin role: %v", err)
		}
	}

	// 检查管理员用户是否存在
	var adminUser model.User
	if err := db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		// 生成密码哈希
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}

		// 创建管理员用户
		adminUser = model.User{
			ID:       common.GenerateUUID(),
			Username: "admin",
			Password: string(hashedPassword),
			RoleID:   adminRole.ID,
			Status:   "active",
		}

		if err := db.Create(&adminUser).Error; err != nil {
			log.Fatalf("Failed to create admin user: %v", err)
		}
	}

	// 创建管理员权限
	adminPermissions := []model.Permission{
		{ID: common.GenerateUUID(), Resource: "user", Action: "manage"},
		{ID: common.GenerateUUID(), Resource: "role", Action: "manage"},
		{ID: common.GenerateUUID(), Resource: "device", Action: "manage"},
		{ID: common.GenerateUUID(), Resource: "license", Action: "manage"},
	}

	for _, perm := range adminPermissions {
		if err := db.Create(&perm).Error; err != nil {
			log.Printf("Warning: Failed to create permission %s-%s: %v", perm.Resource, perm.Action, err)
			continue
		}
		
		// 关联角色和权限
		rolePermission := &model.RolePermission{
			RoleID:       adminRole.ID,
			PermissionID: perm.ID,
		}
		
		if err := db.Create(rolePermission).Error; err != nil {
			log.Printf("Warning: Failed to associate permission %s-%s with admin role: %v", perm.Resource, perm.Action, err)
		}
	}

	log.Println("Successfully created admin user:")
	log.Println("Username: admin")
	log.Println("Password: admin123")
}
