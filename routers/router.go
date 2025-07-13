package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"hotel-management-system/controllers"
	"hotel-management-system/middlewares"
)

func SetupRouters() *gin.Engine {
	r := gin.Default()

	// 设置跨域
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true // 允许所有域
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的HTTP方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length", "Authorization"},         // 允许暴露的响应头
		AllowCredentials: true,                                                // 允许携带凭证
		MaxAge:           12 * 3600,                                           // 预检请求的缓存时间，单位为秒
	}))

	// 设置路由组
	// 用户认证相关路由
	adminGroup := r.Group("/admin")
	adminGroup.POST("/register", controllers.AdminRegister) // 用户注册
	adminGroup.POST("/login", controllers.AdminLogin)       // 用户登录
	adminGroup.Use(middlewares.AuthMiddleware())            // 使用认证中间件
	{
		adminGroup.GET("/getOne", controllers.GetAdminInfo)     // 获取用户信息
		adminGroup.GET("/list", controllers.GetAdminList)       // 获取用户列表
		adminGroup.POST("/add", controllers.AddAdmin)           // 添加用户
		adminGroup.POST("/delete", controllers.DeleteAdmin)     // 删除用户
		adminGroup.POST("/update", controllers.UpdateAdmin)     // 更新用户信息
		adminGroup.POST("/resetPwd", controllers.ResetAdminPwd) // 重置用户密码
	}
	// 文件相关路由
	uploadsGroup := r.Group("/uploads")
	uploadsGroup.GET("/:filename", controllers.GetImg) // 图片下载
	uploadsGroup.Use(middlewares.AuthMiddleware())     // 使用认证中间件
	{
		uploadsGroup.POST("/img", controllers.UploadImg) // 图片上传
	}
	// 角色相关路由
	roleGroup := r.Group("/role")
	roleGroup.Use(middlewares.AuthMiddleware()) // 使用认证中间件
	{
		roleGroup.POST("/add", controllers.AddRole)       // 添加角色
		roleGroup.POST("/delete", controllers.DeleteRole) // 删除角色
		roleGroup.POST("/update", controllers.UpdateRole) // 更新角色
		roleGroup.GET("/getOne", controllers.GetRole)     // 获取单个角色
		roleGroup.GET("/list", controllers.GetRoleList)   // 获取角色列表
	}
	// 房型相关路由
	roomTypeGroup := r.Group("/roomType")
	roomTypeGroup.Use(middlewares.AuthMiddleware()) // 使用认证中间件
	{
		roomTypeGroup.GET("/list", controllers.GetRoomTypeList)     // 获取房型列表
		roomTypeGroup.POST("/add", controllers.AddRoomType)         // 添加房型
		roomTypeGroup.POST("/delete", controllers.DeleteRoomType)   // 删除房型
		roomTypeGroup.POST("/update", controllers.UpdateRoomType)   // 更新房型信息
		roomTypeGroup.GET("/detail", controllers.GetRoomTypeDetail) // 获取房型详情
	}
	// 房间相关路由
	roomGroup := r.Group("/room")
	roomGroup.Use(middlewares.AuthMiddleware()) // 使用认证中间件
	{
		roomGroup.GET("/list", controllers.GetRoomList)             // 获取房间列表
		roomGroup.POST("/add", controllers.AddRoom)                 // 添加房间
		roomGroup.POST("/delete", controllers.DeleteRoom)           // 删除房间
		roomGroup.POST("/update", controllers.UpdateRoom)           // 更新房间信息
		roomGroup.GET("/detail", controllers.GetRoomDetail)         // 获取房间详情
		roomGroup.GET("/statusList", controllers.GetRoomStatusList) // 获取房间状态列表
	}
	// 客户相关路由
	guestGroup := r.Group("/guestRecord")
	guestGroup.Use(middlewares.AuthMiddleware()) // 使用认证中间件
	{
		guestGroup.GET("/list", controllers.GetGuestList)             // 获取入住信息列表
		guestGroup.POST("/add", controllers.AddGuest)                 // 添加入住信息
		guestGroup.POST("/delete", controllers.DeleteGuest)           // 删除入住信息
		guestGroup.POST("/update", controllers.UpdateGuest)           // 更新入住信息
		guestGroup.GET("/detail", controllers.GetGuestDetail)         // 获取入住信息
		guestGroup.GET("/roomList", controllers.GetGuestRoomList)     // 获取房间列表
		guestGroup.GET("/statusList", controllers.GetGuestStatusList) // 获取入住状态列表
		guestGroup.POST("/checkout", controllers.CheckoutGuest)       // 结账
	}
	return r
}
