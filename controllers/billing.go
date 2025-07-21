package controllers

import (
	"github.com/gin-gonic/gin"
	"hotel-management-system/global"
	"net/http"
)

func GetBillingList(c *gin.Context) {
	// 分为roomTypeName来计算收入
	type Billing struct {
		RoomTypeName string `json:"roomTypeName"`
		TotalMoney   int32  `json:"totalMoney"`
	}
	var results []Billing
	err := global.Db.Table("billings").
		Select("room_type_name, SUM(amount) as total_money").
		Group("room_type_name").
		Order("total_money DESC").
		Scan(&results).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
	})
}

func GetTop3(c *gin.Context) {
	type Billing struct {
		RoomId          string `json:"roomId"`
		TotalMoney      int32  `json:"totalMoney"`
		Num             int    `json:"count"`
		RoomTypeName    string `json:"roomTypeName"`
		RoomTypePrice   int32  `json:"roomTypePrice"`
		RoomDescription string `json:"roomDescription"`
	}

	var results []Billing
	if err := global.Db.Table("billings").
		Select("billings.room_id as room_id, SUM(amount) as total_money, COUNT(billings.room_id) as num, billings.room_type_name as room_type_name, room_types.room_type_price, rooms.room_description").
		Joins("LEFT JOIN rooms ON rooms.room_id = billings.room_id").
		Joins("LEFT JOIN room_types ON rooms.room_type_id = room_types.id").
		Group("billings.room_id, billings.room_type_name, room_types.room_type_price, rooms.room_description").
		Order("num DESC, room_types.room_type_price DESC").
		Limit(3).
		Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
	})
}
