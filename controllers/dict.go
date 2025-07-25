package controllers

import (
	"github.com/gin-gonic/gin"
	"hotel-management-system/global"
	"hotel-management-system/models"
	"net/http"
)

func GetDictList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    []string{"room_status", "reside_state", "menu_type"},
	})
}

func GetDictByType(c *gin.Context) {
	dictType := c.Param("dictType")
	switch dictType {
	case "room_status":
		var list []models.RoomStatus
		if err := global.Db.Find(&list).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询失败"})
			return
		}
		var result []gin.H
		for _, item := range list {
			result = append(result, gin.H{"id": item.ID, "name": item.StatusName})
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result,
		})
	case "reside_state":
		var list []models.ResideState
		if err := global.Db.Find(&list).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询失败"})
			return
		}
		var result []gin.H
		for _, item := range list {
			result = append(result, gin.H{"id": item.ID, "name": item.StateName})
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result,
		})
	case "menu_type":
		var list []models.MenuType
		if err := global.Db.Find(&list).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询失败"})
			return
		}
		var result []gin.H
		for _, item := range list {
			result = append(result, gin.H{"id": item.ID, "name": item.Type})
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result,
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
}

func AddDict(c *gin.Context) {
	dictType := c.Param("dictType")
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	switch dictType {
	case "room_status":
		if err := global.Db.First(&models.RoomStatus{}, req.Name).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "已存在"})
			return
		}
		if err := global.Db.Create(&models.RoomStatus{StatusName: req.Name}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "添加失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "添加成功",
		})
	case "reside_state":
		if err := global.Db.First(&models.ResideState{}, req.Name).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "已存在"})
			return
		}
		if err := global.Db.Create(&models.ResideState{StateName: req.Name}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "添加失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "添加成功",
		})
	case "menu_type":
		if err := global.Db.First(&models.MenuType{}, req.Name).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "已存在"})
			return
		}
		if err := global.Db.Create(&models.MenuType{Type: req.Name}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "添加失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "添加成功",
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
}

func DeleteDict(c *gin.Context) {
	dictType := c.Param("dictType")
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	switch dictType {
	case "room_status":
		if err := global.Db.Delete(&models.RoomStatus{}, req.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "删除成功",
		})
	case "reside_state":
		if err := global.Db.Delete(&models.ResideState{}, req.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "删除成功",
		})
	case "menu_type":
		if err := global.Db.Delete(&models.MenuType{}, req.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "删除成功",
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
}

func UpdateDict(c *gin.Context) {
	dictType := c.Param("dictType")
	var req struct {
		ID   uint   `json:"id" binding:"required"`
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	switch dictType {
	case "room_status":
		if err := global.Db.Model(&models.RoomStatus{}).Where("id = ?", req.ID).Update("status_name", req.Name).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "更新成功",
		})
	case "reside_state":
		if err := global.Db.Model(&models.ResideState{}).Where("id = ?", req.ID).Update("state_name", req.Name).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "更新成功",
		})
	case "menu_type":
		if err := global.Db.Model(&models.MenuType{}).Where("id = ?", req.ID).Update("type", req.Name).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "更新成功",
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
}
