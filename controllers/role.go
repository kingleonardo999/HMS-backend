package controllers

import (
	"github.com/gin-gonic/gin"
	"hotel-management-system/global"
	"hotel-management-system/models"
	"net/http"
	"strconv"
)

func AddRole(c *gin.Context) {
	var role models.Role
	// 绑定 JSON 数据到结构体
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}

	// 检查角色名是否已存在
	var existingRole models.Role
	if err := global.Db.Where("role_name = ?", role.RoleName).First(&existingRole).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "角色名已存在"})
		return
	}
	// 保存角色到数据库
	if err := global.Db.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "角色添加成功"})
}

func DeleteRole(c *gin.Context) {
	var req struct {
		RoleId uint `json:"roleId" binding:"required"`
	}
	// 绑定 JSON 数据到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	roleId := req.RoleId
	if roleId == 1 {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "无法删除管理员角色"})
		return
	}
	var role models.Role
	if err := global.Db.First(&role, roleId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "角色不存在"})
		return
	}
	if err := global.Db.Delete(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "角色删除成功"})
}

func UpdateRole(c *gin.Context) {
	var role models.Role
	// 绑定 JSON 数据到结构体
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}

	if role.ID == 1 {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "无法修改管理员角色"})
		return
	}

	// 检查角色是否存在
	var existingRole models.Role
	if err := global.Db.First(&existingRole, role.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "角色不存在"})
		return
	}

	// 更新角色信息
	if err := global.Db.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "角色更新成功"})
}

func GetRole(c *gin.Context) {
	roleIdStr := c.Query("roleId")
	roleId, err := strconv.Atoi(roleIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的角色 ID"})
		return
	}

	var role models.Role
	if err := global.Db.First(&role, roleId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "角色不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    role,
	})
}

func GetRoleList(c *gin.Context) {
	var roleList []models.Role
	global.Db.Find(&roleList)

	// 将 id 字段转换为 roleId
	var result []map[string]interface{}
	for _, role := range roleList {
		item := map[string]interface{}{
			"roleId":   role.ID,
			"roleName": role.RoleName,
			// 根据 models.Role 结构体添加其他需要的字段
		}
		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}
