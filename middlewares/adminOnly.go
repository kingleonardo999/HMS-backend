package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			if c.Request.Method != "GET" {
				c.JSON(http.StatusForbidden, gin.H{"message": "当前用户无权限"})
				c.Abort() // 中止请求
				return
			}
		}
		c.Next() // 继续处理请求
	}
}
