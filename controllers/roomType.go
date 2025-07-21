package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"hotel-management-system/global"
	"hotel-management-system/models"
	"net/http"
)

func GetRoomTypeList(c *gin.Context) {
	roomTypeList := []models.RoomType{}
	global.Db.Select("id, room_type_name, room_type_price, bed_num").Find(&roomTypeList)

	// 将 id 字段转换为 roomTypeId
	result := []map[string]interface{}{}
	for _, roomType := range roomTypeList {
		item := map[string]interface{}{
			"roomTypeId":    roomType.ID,
			"roomTypeName":  roomType.RoomTypeName,
			"roomTypePrice": roomType.RoomTypePrice,
			"bedNum":        roomType.BedNum,
		}
		result = append(result, item)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

func AddRoomType(c *gin.Context) {
	var roomType models.RoomType
	if err := c.ShouldBindJSON(&roomType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	var existingRoomType models.RoomType
	if err := global.Db.Where("room_type_name = ?", roomType.RoomTypeName).First(&existingRoomType).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "房间类型已存在"})
		return
	}
	if err := global.Db.Create(&roomType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "房间类型添加成功"})
}

func DeleteRoomType(c *gin.Context) {
	var req struct {
		RoomTypeId uint `json:"roomTypeId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	var roomType models.RoomType
	if err := global.Db.Where("id = ?", req.RoomTypeId).First(&roomType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "房间类型不存在"})
		return
	}
	if err := global.Db.Delete(&roomType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}

}

func UpdateRoomType(c *gin.Context) {
	newRoomType := models.RoomType{}
	if err := c.ShouldBindJSON(&newRoomType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	if err := global.Db.Save(&newRoomType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "房间类型不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "房间类型更新成功"})
}

func GetRoomTypeDetail(c *gin.Context) {
	roomTypeId := c.Query("roomTypeId")

	var roomType models.RoomType
	if err := global.Db.First(&roomType, roomTypeId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "房间类型不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": roomType})
}

func GetTotalTypePrice(c *gin.Context) {
	type Result struct {
		RoomTypeId    uint   `json:"roomTypeId"`
		RoomTypeName  string `json:"roomTypeName"`
		RoomTypePrice int32  `json:"roomTypePrice"`
		TotalMoney    int64  `json:"totalMoney"`
		Count         int64  `json:"count"`
	}
	var results []Result
	err := global.Db.Table("room_types").
		Select("room_types.id as room_type_id, room_types.room_type_name, room_types.room_type_price, COUNT(rooms.room_type_id) as count, COALESCE(COUNT(rooms.room_type_id) * room_types.room_type_price, 0) as total_money").
		Where("rooms.room_status_id = ?", roomOccupied).
		Joins("LEFT JOIN rooms ON rooms.room_type_id = room_types.id").
		Group("room_types.id, room_types.room_type_name, room_types.room_type_price").
		Scan(&results).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}
