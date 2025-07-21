package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"hotel-management-system/global"
	"hotel-management-system/models"
	"net/http"
	"strconv"
	"strings"
)

type RoomInfo struct {
	RoomId          string  `json:"roomId"`
	RoomTypeId      uint    `json:"roomTypeId"`
	RoomTypeName    string  `json:"roomTypeName"`
	RoomTypePrice   float32 `json:"roomTypePrice"`
	BedNum          uint    `json:"bedNum"`
	RoomStatusId    uint    `json:"roomStatusId"`
	RoomStatus      string  `json:"roomStatusName"`
	RoomDescription string  `json:"roomDescription"`
}

func GetRoomList(c *gin.Context) {
	pageIndex, err1 := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, err2 := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	roomTypeId, err3 := strconv.Atoi(c.DefaultQuery("roomTypeId", "0"))
	roomStatusId, err4 := strconv.Atoi(c.DefaultQuery("roomStatusId", "0"))
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	var list []RoomInfo
	var total int64
	query := global.Db.Table("rooms").
		Select("rooms.room_id, room_types.room_type_name, room_types.room_type_price, room_types.bed_num, room_statuses.status_name AS room_status, rooms.room_description").
		Joins("LEFT JOIN room_types ON rooms.room_type_id = room_types.id").
		Joins("LEFT JOIN room_statuses ON rooms.room_status_id = room_statuses.id")
	if roomTypeId != 0 {
		query = query.Where("rooms.room_type_id = ?", roomTypeId)
	}
	if roomStatusId != 0 {
		query = query.Where("rooms.room_status_id = ?", roomStatusId)
	}
	query.Count(&total)
	err := query.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Scan(&list).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    list,
		"count":   total,
	})
}

func AddRoom(c *gin.Context) {
	var newRoom models.Room
	if err := c.ShouldBindJSON(&newRoom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	if err := global.Db.Create(&newRoom).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "uni") {
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": "房间已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "添加房间成功",
	})
}

func DeleteRoom(c *gin.Context) {
	var req struct {
		RoomId string `json:"roomId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	if err := global.Db.Where("room_id = ?", req.RoomId).Delete(&models.Room{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "房间不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除房间成功",
	})
}

func UpdateRoom(c *gin.Context) {
	var newRoom models.Room
	if err := c.ShouldBindJSON(&newRoom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	if newRoom.RoomStatusId == roomOccupied || newRoom.RoomStatusId == roomOrdered {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "房间状态不可更改"})
		return
	}
	var oldRoom models.Room
	if err := global.Db.Where("room_id = ?", newRoom.RoomId).First(&oldRoom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "房间不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}
	if err := global.Db.Model(&oldRoom).Updates(&newRoom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "房间不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新房间成功",
	})
}

func GetRoomDetail(c *gin.Context) {
	var roomInfo RoomInfo
	roomId := c.Query("roomId")
	err := global.Db.Table("rooms").
		Select("rooms.room_id, rooms.room_type_id, rooms.room_status_id, room_types.room_type_name, room_types.room_type_price, room_types.bed_num, room_statuses.status_name AS room_status, rooms.room_description").
		Joins("LEFT JOIN room_types ON rooms.room_type_id = room_types.id").
		Joins("LEFT JOIN room_statuses ON rooms.room_status_id = room_statuses.id").
		Where("rooms.room_id = ?", roomId).
		Scan(&roomInfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "房间不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    roomInfo,
	})
}

func GetRoomStatusList(c *gin.Context) {
	var list []models.RoomStatus
	if err := global.Db.Select("id, status_name").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    list,
	})
}
