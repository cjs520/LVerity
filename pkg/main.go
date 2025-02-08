// main.go
package main

import (
	"fmt"
	"log"
	"LVerity/pkg/config"
	"LVerity/pkg/database"
	"LVerity/pkg/router"
	"LVerity/pkg/service"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	dbConfig := config.GetConfig().Database
	if err := database.InitDB(&dbConfig); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化默认权限
	if err := service.InitDefaultPermissions(); err != nil {
		log.Printf("Warning: Failed to initialize default permissions: %v", err)
	}

	// 初始化管理员账户
	if err := service.InitAdminUser(); err != nil {
		log.Printf("Warning: Failed to initialize admin user: %v", err)
	}

	// 创建路由
	r := router.SetupRouter()

	// 启动服务器
	port := config.GetConfig().Server.Port
	if port == 0 {
		port = 8080
	}
	addr := fmt.Sprintf("%s:%d", config.GetConfig().Server.Host, port)
	log.Printf("Server is running on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
