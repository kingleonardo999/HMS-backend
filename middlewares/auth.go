package middlewares

import (
	"github.com/gin-gonic/gin"
	"hotel-management-system/global"
	"hotel-management-system/models"
	"hotel-management-system/utils"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort() // 中止请求
			return
		}

		// 解析 JWT 令牌
		username, err := utils.ParseJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort() // 中止请求
			return
		}
		c.Set("username", username) // 将用户名存储在上下文中，供后续处理使用
		// 查询角色
		var user models.User
		if err := global.Db.Preload("Role").Model(&models.User{}).Where("login_id = ?", username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort() // 中止请求
			return
		}
		c.Set("role", user.Role.RoleName)
		c.Set("userId", user.ID)
		c.Next() // 继续处理请求
	}
}
