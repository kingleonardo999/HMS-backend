package controllers

import (
	"errors"
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

func GetGuestList(c *gin.Context) {
	// 分页查询
	pageIndex, err1 := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, err2 := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	resideState, err3 := strconv.Atoi(c.DefaultQuery("resideState", "0"))
	guestName := c.DefaultQuery("guestName", "")
	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	var list []struct {
		Id            uint   `json:"id"`
		IdentityId    string `json:"identityId"`
		Name          string `json:"guestName"`
		Phone         string `json:"guestPhone"`
		RoomTypeName  string `json:"roomTypeName"`
		RoomTypePrice int32  `json:"roomTypePrice"`
		RoomId        string `json:"roomId"`
		RoomStatus    string `json:"roomStatus"`
		ResideDate    string `json:"resideDate"`
		LeaveDate     string `json:"leaveDate"`
		Deposit       int    `json:"deposit"`
		GuestNum      int    `json:"guestNum"`
		ResideState   string `json:"resideState"`
	}
	query := global.Db.Table("resides").
		Select("resides.id AS id, guests.identity_id AS identity_id, name, phone, room_types.room_type_name AS room_type_name, " +
			"room_types.room_type_price AS room_type_price, resides.room_id, room_statuses.status_name AS room_status, " +
			"reside_date, leave_date, deposit, guest_num, reside_state").
		Joins("LEFT JOIN guests ON resides.guest_id = guests.id").
		Joins("LEFT JOIN rooms ON resides.room_id = rooms.room_id").
		Joins("LEFT JOIN room_types ON rooms.room_type_id = room_types.id").
		Joins("LEFT JOIN room_statuses ON rooms.room_status_id = room_statuses.id")
	if resideState != 0 {
		query = query.Where("reside_state = ?", resideState)
	}
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

func AddGuest(c *gin.Context) {
	var reside struct {
		IdentityId string `json:"identityId" binding:"required"`
		Name       string `json:"guestName" binding:"required"`
		Phone      string `json:"guestPhone" binding:"required"`
		RoomTypeId uint   `json:"roomTypeId" binding:"required"`
		RoomId     string `json:"roomId" binding:"required"`
		ResideDate string `json:"resideDate" binding:"required"`
		Deposit    int    `json:"deposit" binding:"required"`
		GuestNum   int    `json:"guestNum" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reside); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	// 查询客户是否存在，不存在则创建
	var guest models.Guest
	if err := global.Db.Where("identity_id = ?", reside.IdentityId).First(&guest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			guest.IdentityId = reside.IdentityId
			guest.Name = reside.Name
			guest.Phone = reside.Phone
			if err := global.Db.Create(&guest).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建客户失败"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询客户失败"})
			return
		}
	} else {
		if guest.Name != reside.Name {
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": "客户信息不匹配"})
			return
		}
	}
	var room models.Room
	if err := global.Db.Where("room_id = ?", reside.RoomId).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "房间不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询房间失败"})
		return
	} else if room.RoomStatusId != roomFree {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "房间已入住"})
		return
	} else if room.RoomTypeId != reside.RoomTypeId {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	// 创建入住记录
	var record models.Reside
	record.GuestId = guest.ID
	record.RoomId = reside.RoomId
	record.ResideDate = reside.ResideDate
	record.Deposit = int32(reside.Deposit)
	record.GuestNum = int32(reside.GuestNum)
	record.ResideState = resideState
	record.TotalMoney = 0
	// 加锁房间, 避免竞争
	tx := global.Db.Begin()
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("room_id = ?", reside.RoomId).First(&room).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "加锁房间失败"})
		return
	}
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建入住记录失败"})
		return
	}
	if err := tx.Model(&room).Where("room_id = ?", reside.RoomId).Update("room_status_id", roomOccupied).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新房间状态失败"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "入住成功",
	})
}

func UpdateGuest(c *gin.Context) {
	var req struct {
		Id         uint   `json:"id" binding:"required"`
		Phone      string `json:"guestPhone" binding:"required"`
		RoomTypeId uint   `json:"roomTypeId" binding:"required"`
		RoomId     string `json:"roomId" binding:"required"`
		LeaveDate  string `json:"leaveDate"`
		GuestNum   int    `json:"guestNum" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	var reside models.Reside
	if err := global.Db.Preload("Room").Preload("Room.RoomType").Preload("Guest").Where("id = ?", req.Id).First(&reside).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "入住记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询入住记录失败"})
		return
	}
	tx := global.Db.Begin()
	if err := tx.Model(&reside.Guest).Update("phone", req.Phone).Where("id = ?", reside.GuestId).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新入住记录失败"})
		return
	}
	if req.RoomTypeId != reside.Room.RoomTypeId {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	if req.LeaveDate != "" {
		if err := tx.Model(&reside).Where("id = ?", req.Id).Update("leave_date", req.LeaveDate).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新入住记录失败"})
			return
		}
	}
	if req.GuestNum != 0 {
		if err := tx.Model(&reside).Where("id = ?", req.Id).Update("guest_num", req.GuestNum).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新入住记录失败"})
			return
		}
	}
	// 更新房间，可以更换同类房间，注意状态改变
	if req.RoomId != reside.RoomId {
		oldRoom := reside.Room
		var newRoom models.Room
		if err := tx.Where("room_id = ?", req.RoomId).First(&newRoom).Error; err != nil {
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
		newRoom.RoomStatusId = roomOccupied
		if err := tx.Save(&newRoom).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新房间失败"})
			return
		}
		reside.Room = newRoom
		reside.RoomId = newRoom.RoomId
		if err := tx.Save(&reside).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新入住记录失败"})
			return
		}
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新入住记录成功",
	})
}

func DeleteGuest(c *gin.Context) {
	var req struct {
		Id uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	// 查询状态是否为已结账
	var state string
	if err := global.Db.Model(&models.Reside{}).Select("reside_state").Where("id = ?", req.Id).First(&state).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "入住记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询入住记录失败"})
		return
	}
	if state != "已结账" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "入住记录未结账"})
		return
	}
	if err := global.Db.Where("id = ?", req.Id).Delete(&models.Reside{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "入住记录不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除入住记录成功",
	})
}

func GetGuestDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	var info struct {
		Id            uint   `json:"id"`
		IdentityId    string `json:"identityId"`
		Name          string `json:"guestName"`
		Phone         string `json:"guestPhone"`
		RoomTypeId    uint   `json:"roomTypeId"`
		RoomTypePrice int32  `json:"roomTypePrice"`
		RoomId        string `json:"roomId"`
		RoomStatus    string `json:"roomStatus"`
		ResideDate    string `json:"resideDate"`
		LeaveDate     string `json:"leaveDate"`
		Deposit       int    `json:"deposit"`
		GuestNum      int    `json:"guestNum"`
		ResideState   string `json:"resideState"`
		TotalMoney    int    `json:"totalMoney"`
	}
	if err := global.Db.Table("resides").
		Select("resides.id AS id, guests.identity_id AS identity_id, name, phone, room_types.id AS room_type_id, "+
			"room_types.room_type_price AS room_type_price, resides.room_id, room_statuses.status_name AS room_status, "+
			"reside_date, leave_date, deposit, guest_num, reside_state, total_money").
		Joins("LEFT JOIN guests ON resides.guest_id = guests.id").
		Joins("LEFT JOIN rooms ON resides.room_id = rooms.room_id").
		Joins("LEFT JOIN room_types ON rooms.room_type_id = room_types.id").
		Joins("LEFT JOIN room_statuses ON rooms.room_status_id = room_statuses.id").
		Where("resides.id = ?", id).
		Scan(&info).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "入住信息不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    info,
	})
}

func GetGuestRoomList(c *gin.Context) {
	roomTypeId, err := strconv.Atoi(c.Query("roomTypeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	var list []struct {
		RoomId string `json:"roomId"`
	}
	if err := global.Db.Table("rooms").Select("room_id").Where("room_type_id = ? AND room_status_id = ?", roomTypeId, roomFree).Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    list,
	})
}

func GetGuestStatusList(c *gin.Context) {
	var list []models.ResideState
	if err := global.Db.Select("id, state_name").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    list,
	})
}

func CheckoutGuest(c *gin.Context) {
	var req struct {
		Id         uint  `json:"id"`
		TotalMoney int32 `json:"totalMoney"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	var reside models.Reside
	if err := global.Db.Preload("Room").Preload("Room.RoomType").Where("id = ?", req.Id).First(&reside).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "入住记录不存在"})
		return
	}
	// 根据入住记录计算应付金额
	layout := "2006-01-02T15:04:05.000Z" // 日期格式
	start, err := time.Parse(layout, reside.ResideDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}
	end, err := time.Parse(layout, reside.LeaveDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器错误"})
		return
	}
	days := int32(math.Ceil(end.Sub(start).Hours() / 24))
	totalMoney := days * reside.Room.RoomType.RoomTypePrice
	if totalMoney != req.TotalMoney {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "入住记录金额不匹配"})
		return
	}
	tx := global.Db.Begin()
	if err := tx.Where("id = ?", req.Id).Updates(&models.Reside{ResideState: "已结账", TotalMoney: totalMoney}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新入住记录失败"})
		return
	}
	if err := tx.Where("room_id = ?", reside.Room.RoomId).Updates(&models.Room{RoomStatusId: roomFree}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新房间状态失败"})
		return
	}
	var billing models.Billing
	billing.ResideId = reside.ID
	billing.Time = time.Now().Format("2006-01-02 15:04:05")
	billing.Amount = totalMoney
	billing.GuestId = reside.GuestId
	billing.RoomId = reside.Room.RoomId
	billing.RoomTypeName = reside.Room.RoomType.RoomTypeName
	if err := tx.Create(&billing).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建账单失败"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "结账成功",
	})
}
