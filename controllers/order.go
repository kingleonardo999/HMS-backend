package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hotel-management-system/global"
	"hotel-management-system/models"
	"math"
	"net/http"
	"strconv"
	"time"
)

func GetOrderList(c *gin.Context) {
	pageIndex, err1 := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, err2 := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	guestName := c.DefaultQuery("guestName", "")
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	var list []struct {
		OrderId      string `json:"orderId"`
		IdentityId   string `json:"identityId"`
		Name         string `json:"guestName"`
		Phone        string `json:"guestPhone"`
		RoomTypeName string `json:"roomTypeName"`
		RoomId       string `json:"roomId"`
		ResideDate   string `json:"resideDate"`
		LeaveDate    string `json:"leaveDate"`
		GuestNum     int    `json:"guestNum"`
		TotalMoney   int    `json:"totalMoney"`
	}
	query := global.Db.Table("orders").
		Select("orders.order_id AS order_id, guests.identity_id AS identity_id, name, phone, room_types.room_type_name AS room_type_name, " +
			" orders.room_id, orders.order_date AS reside_date, leave_date, guest_num, total_money").
		Joins("LEFT JOIN guests ON orders.guest_id = guests.id").
		Joins("LEFT JOIN rooms ON orders.room_id = rooms.room_id").
		Joins("LEFT JOIN room_types ON rooms.room_type_id = room_types.id")
	if guestName != "" {
		query = query.Where("guests.name LIKE ?", "%"+guestName+"%")
	}
	err := query.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Scan(&list).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询失败"})
		return
	}
	var total = len(list)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    list,
		"count":   total,
	})
}

func GetOrderDetail(c *gin.Context) {
	var orderId = c.DefaultQuery("id", "")
	var order struct {
		OrderId      string `json:"orderId"`
		IdentityId   string `json:"identityId"`
		Name         string `json:"guestName"`
		Phone        string `json:"guestPhone"`
		RoomTypeId   uint   `json:"roomTypeId"`
		RoomTypeName string `json:"roomTypeName"`
		RoomId       string `json:"roomId"`
		ResideDate   string `json:"resideDate"`
		LeaveDate    string `json:"leaveDate"`
		GuestNum     int    `json:"guestNum"`
		TotalMoney   int    `json:"totalMoney"`
	}
	if err := global.Db.Table("orders").
		Select("orders.order_id AS order_id, guests.identity_id AS identity_id, guests.name AS name, guests.phone AS phone, room_types.room_type_name AS room_type_name, "+
			"room_types.id AS room_type_id, orders.room_id, orders.order_date AS reside_date, leave_date, guest_num, total_money").
		Joins("LEFT JOIN guests ON orders.guest_id = guests.id").
		Joins("LEFT JOIN rooms ON orders.room_id = rooms.room_id").
		Joins("LEFT JOIN room_types ON rooms.room_type_id = room_types.id").
		Where("orders.order_id = ?", orderId).
		Scan(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    order,
	})
}

func AddOrder(c *gin.Context) {
	var order struct {
		IdentityId string `json:"identityId"`
		Name       string `json:"guestName"`
		Phone      string `json:"guestPhone"`
		RoomTypeId uint   `json:"roomTypeId"`
		RoomId     string `json:"roomId"`
		ResideDate string `json:"resideDate"`
		LeaveDate  string `json:"leaveDate"`
		GuestNum   int    `json:"guestNum"`
		TotalMoney int    `json:"totalMoney"`
	}
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	// 查询客户是否存在，不存在则创建
	var guest models.Guest
	if err := global.Db.Where("identity_id = ?", order.IdentityId).First(&guest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			guest.IdentityId = order.IdentityId
			guest.Name = order.Name
			guest.Phone = order.Phone
			if err := global.Db.Create(&guest).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建客户失败"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询客户失败"})
			return
		}
	} else {
		if guest.Name != order.Name {
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": "客户信息不匹配"})
			return
		}
	}
	var room models.Room
	if err := global.Db.Where("room_id = ?", order.RoomId).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "房间不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询房间失败"})
		return
	} else if room.RoomStatusId != roomFree {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "房间已被占用"})
		return
	} else if room.RoomTypeId != order.RoomTypeId {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	var orderDetail models.Order
	orderDetail.OrderId = fmt.Sprintf("od%d", time.Now().Unix())
	orderDetail.GuestId = guest.ID
	orderDetail.RoomId = order.RoomId
	orderDetail.OrderDate = order.ResideDate
	orderDetail.LeaveDate = order.LeaveDate
	orderDetail.TotalMoney = int32(order.TotalMoney)
	orderDetail.GuestNum = int32(order.GuestNum)
	// 加锁房间, 避免竞争
	tx := global.Db.Begin()
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("room_id = ?", orderDetail.RoomId).First(&room).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "加锁房间失败"})
		return
	}
	if err := tx.Create(&orderDetail).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建入住记录失败"})
		return
	}
	if err := tx.Model(&room).Where("room_id = ?", order.RoomId).Update("room_status_id", roomOrdered).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新房间状态失败"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建订单成功",
	})
}

func UpdateOrder(c *gin.Context) {
	var order struct {
		OrderId    string `json:"orderId" binding:"required"`
		Phone      string `json:"guestPhone" binding:"required"`
		RoomId     string `json:"roomId" binding:"required"`
		LeaveDate  string `json:"leaveDate" binding:"required"`
		GuestNum   int    `json:"guestNum" binding:"required"`
		TotalMoney int    `json:"totalMoney" binding:"required"`
	}
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	var orderDetail models.Order
	if err := global.Db.Preload("Room").Preload("Guest").Where("order_id = ?", order.OrderId).First(&orderDetail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "订单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询订单失败"})
		return
	}
	tx := global.Db.Begin()
	if err := tx.Model(&orderDetail.Guest).Update("phone", order.Phone).Where("id = ?", orderDetail.GuestId).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新入住记录失败"})
		return
	}
	if order.LeaveDate != orderDetail.LeaveDate {
		if err := tx.Model(&orderDetail).Where("order_id = ?", orderDetail.OrderId).Update("leave_date", order.LeaveDate).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新入住记录失败"})
			return
		}
	}
	if order.GuestNum != 0 {
		if err := tx.Model(&orderDetail).Where("order_id = ?", orderDetail.OrderId).Update("guest_num", order.GuestNum).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新入住记录失败"})
			return
		}
	}
	// 更新房间，可以更换同类房间，注意状态改变
	if order.RoomId != orderDetail.RoomId {
		oldRoom := orderDetail.Room
		var newRoom models.Room
		if err := tx.Where("room_id = ?", order.RoomId).First(&newRoom).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询房间失败"})
			return
		}
		oldRoom.RoomStatusId = roomFree
		if err := tx.Save(&oldRoom).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新房间失败"})
			return
		}
		newRoom.RoomStatusId = roomOrdered
		if err := tx.Save(&newRoom).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新房间失败"})
			return
		}
		orderDetail.Room = newRoom
		orderDetail.RoomId = newRoom.RoomId
		if err := tx.Save(&orderDetail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新入住记录失败"})
			return
		}
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新订单成功",
	})
}

func Order2Reside(c *gin.Context) {
	var req struct {
		OrderId    string `json:"id" binding:"required"`
		TotalMoney int32  `json:"totalMoney" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	orderId := req.OrderId
	var orderDetail models.Order
	if err := global.Db.Preload("Room").Preload("Room.RoomType").Preload("Guest").Where("order_id = ?", orderId).First(&orderDetail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "订单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询订单失败"})
		return
	}
	// 根据入住记录计算应付金额
	layout := "2006-01-02T15:04:05.000Z" // 日期格式
	start, err := time.Parse(layout, orderDetail.OrderDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}
	end, err := time.Parse(layout, orderDetail.LeaveDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}
	days := int32(math.Ceil(end.Sub(start).Hours() / 24))
	totalMoney := days * orderDetail.Room.RoomType.RoomTypePrice
	if totalMoney != req.TotalMoney {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "入住记录金额不匹配"})
		return
	}
	var reside models.Reside
	reside.ResideDate = orderDetail.OrderDate
	reside.LeaveDate = orderDetail.LeaveDate
	reside.GuestNum = orderDetail.GuestNum
	reside.TotalMoney = orderDetail.TotalMoney
	reside.RoomId = orderDetail.RoomId
	reside.GuestId = orderDetail.GuestId
	reside.ResideState = "已结账"
	reside.Deposit = 0
	tx := global.Db.Begin()
	if err := tx.Create(&reside).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建入住记录失败"})
		return
	}
	if err := tx.Model(&orderDetail.Room).Update("room_status_id", roomOccupied).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新房间状态失败"})
		return
	}
	if err := tx.Delete(&orderDetail).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除订单失败"})
		return
	}
	if err := tx.Delete(&orderDetail).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除订单失败"})
		return
	}
	var billing models.Billing
	billing.Time = time.Now().Format("2006-01-02 15:04:05")
	billing.Amount = reside.TotalMoney
	billing.GuestId = reside.GuestId
	billing.RoomId = reside.RoomId
	billing.RoomTypeName = orderDetail.Room.RoomType.RoomTypeName
	billing.ResideId = reside.ID
	if err := tx.Create(&billing).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建账单失败"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "成功入住",
	})
}
