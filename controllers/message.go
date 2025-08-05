package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"hotel-management-system/global"
	"hotel-management-system/models"
	"net/http"
)

func GetMessageList(c *gin.Context) {
	id, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取用户名失败"})
		return
	}
	var list []struct {
		Id       uint   `json:"id"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		CreateAt string `json:"createAt"`
	}
	query := global.Db.Model(&models.UserMessage{}).
		Select("user_messages.id as id, messages.title as title, messages.content as content, messages.created_at as create_at").
		Joins("join messages on user_messages.message_id = messages.id").
		Where("user_messages.user_id = ?", id).
		Order("messages.created_at desc")
	if err := query.Scan(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取消息列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    list,
	})
}

func AddMessage(c *gin.Context) {
	var req struct {
		LoginId string `json:"loginId" binding:"required"`
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	var user models.User
	if err := global.Db.Where("login_id = ?", req.LoginId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询用户失败"})
		return
	}
	message := models.Message{
		Title:   req.Title,
		Content: req.Content,
		AdminId: user.ID,
	}
	if err := global.Db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "添加失败"})
		return
	}
	// 插入用户消息
	var list []uint
	if err := global.Db.Model(&models.User{}).Select("id").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取用户列表失败"})
		return
	}
	for _, v := range list {
		if err := global.Db.Create(&models.UserMessage{
			UserId:    v,
			MessageId: message.ID,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "添加失败"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "添加成功",
	})
}

func DeleteMessage(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	var userMessage models.UserMessage
	if err := global.Db.Where("id = ?", req.ID).First(&userMessage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "消息不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询消息失败"})
		return
	}
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取用户ID失败"})
		return
	}
	if userMessage.UserId != userId.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "没有权限"})
		return
	}
	if err := global.Db.Delete(&userMessage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}
