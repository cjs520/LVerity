package router

import (
	"LVerity/pkg/handler"
	"LVerity/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 使用中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// 健康检查
	r.GET("/health", handler.HealthCheck)

	// 公开路由组
	public := r.Group("/auth")
	{
		public.GET("/captcha", handler.GetCaptcha)      // 获取验证码
		public.POST("/login", handler.Login)            // 用户登录
		public.POST("/refresh", middleware.JWTAuth(), handler.RefreshToken)  // 刷新令牌
	}

	// 用户相关路由组
	user := r.Group("/user")
	user.Use(middleware.JWTAuth())
	{
		user.GET("/profile", handler.GetUserProfile)
		user.POST("/profile", handler.UpdateUserProfile)
		user.POST("/change-password", handler.ChangePassword)
	}

	// 需要认证的API路由组
	api := r.Group("/api")
	api.Use(middleware.JWTAuth())
	{
		// 授权管理
		api.GET("/licenses", handler.ListLicenses)
		api.POST("/licenses", handler.CreateLicense)
		api.GET("/licenses/:id", handler.GetLicense)
		api.PUT("/licenses/:id", handler.UpdateLicense)
		api.DELETE("/licenses/:id", handler.DeleteLicense)
		api.GET("/licenses/stats", handler.GetLicenseStats)

		// 设备管理路由
		devices := api.Group("/devices")
		{
			// 基本设备管理
			devices.GET("", handler.ListDevices)                     // 获取设备列表
			devices.POST("", handler.CreateDevice)                   // 创建设备
			devices.GET("/:id", handler.GetDevice)                   // 获取设备详情
			devices.PUT("/:id", handler.UpdateDevice)                // 更新设备
			devices.DELETE("/:id", handler.DeleteDevice)             // 删除设备
			devices.GET("/:id/status", handler.GetDeviceStatus)      // 获取设备状态
			devices.POST("/:id/block", handler.BlockDevice)          // 封禁设备
			devices.POST("/:id/heartbeat", handler.UpdateDeviceHeartbeat) // 更新心跳
			
			// 设备分组管理
			devices.POST("/groups", handler.CreateGroup)             // 创建分组
			devices.POST("/groups/assign", handler.AssignDevice)     // 分配设备到分组
			devices.GET("/groups/:id/devices", handler.GetDevicesByGroup) // 获取分组内设备

			// 异常行为记录
			devices.GET("/:id/abnormal-behaviors", handler.GetDeviceAbnormalBehaviors) // 获取异常行为
			devices.POST("/abnormal-behaviors", handler.RecordAbnormalBehavior)        // 记录异常行为

			// 设备统计
			devices.GET("/stats", handler.GetDeviceStats)            // 获取设备统计
			devices.GET("/:id/usage", handler.GetDeviceUsage)        // 获取使用情况
			devices.GET("/:id/usage-report", handler.GetDeviceUsageReport) // 获取使用报告
			devices.GET("/:id/info", handler.GetDeviceInfo)          // 获取详细信息
			devices.PUT("/:id/metadata", handler.UpdateDeviceMetadata) // 更新元数据
		}
	}

	return r
}
