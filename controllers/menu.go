package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"hotel-management-system/global"
	"hotel-management-system/models"
	"net/http"
	"strconv"
)

func GetMenuList(c *gin.Context) {
	pageIndex, err1 := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, err2 := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	typeId, err3 := strconv.Atoi(c.DefaultQuery("typeId", "0"))
	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	var total int64
	query := global.Db.Table("menus").
		Select("menus.id AS id, menus.name, menus.price, imgs.url AS img").
		Joins("LEFT JOIN imgs ON menus.img_id = imgs.id")
	if typeId != 0 {
		query = query.Where("type_id = ?", typeId)
	}
	query.Count(&total)
	var ret []struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Price int32  `json:"price"`
		Img   string `json:"img"`
	}
	if err := query.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Scan(&ret).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ret,
		"total":   total,
	})
}

func AddMenu(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required"`
		TypeId uint   `json:"typeId" binding:"required"`
		Price  int32  `json:"price" binding:"required"`
		ImgId  uint   `json:"imgId" binding:"required"`
		Desc   string `json:"desc" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	if err := global.Db.Create(&models.Menu{
		Name:   req.Name,
		TypeId: req.TypeId,
		Price:  req.Price,
		ImgId:  req.ImgId,
		Desc:   req.Desc,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "添加失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "添加成功",
	})
}

func UpdateMenu(c *gin.Context) {
	var req struct {
		ID     uint   `json:"id" binding:"required"`
		Name   string `json:"name" binding:"required"`
		TypeId uint   `json:"typeId" binding:"required"`
		Price  int32  `json:"price" binding:"required"`
		ImgId  uint   `json:"imgId"`
		Desc   string `json:"desc" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	var menu models.Menu
	if err := global.Db.Where("id = ?", req.ID).First(&menu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "菜单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询菜单失败"})
		return
	}
	menu.Name = req.Name
	menu.TypeId = req.TypeId
	menu.Price = req.Price
	if req.ImgId != 0 {
		menu.ImgId = req.ImgId
	}
	menu.Desc = req.Desc
	if err := global.Db.Save(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
	})
}

func DeleteMenu(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	if err := global.Db.Delete(&models.Menu{}, req.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}

func GetMenuDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	var menu models.Menu
	if err := global.Db.Preload("Img").Where("id = ?", id).First(&menu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "菜单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询菜单失败"})
		return
	}
	ret := struct {
		ID     uint   `json:"id"`
		Name   string `json:"name"`
		TypeId uint   `json:"typeId"`
		Price  int32  `json:"price"`
		Img    string `json:"img"`
		Desc   string `json:"desc"`
	}{
		ID:     menu.ID,
		Name:   menu.Name,
		TypeId: menu.TypeId,
		Price:  menu.Price,
		Img:    menu.Img.Url,
		Desc:   menu.Desc,
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ret,
	})
}

func GetMenuTypeList(c *gin.Context) {
	var list []models.MenuType
	if err := global.Db.Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    list,
	})
}
