package controllers

import (
	"github.com/gin-gonic/gin"
	"hotel-management-system/global"
	"hotel-management-system/models"
	"hotel-management-system/utils"
	"net/http"
	"strconv"
)

var (
	defaultPhotoID = uint(1) // 默认头像ID
)

func AdminRegister(c *gin.Context) {
	// 获取注册信息
	var registerInfo models.User
	// 绑定 JSON 数据到结构体
	if err := c.ShouldBindJSON(&registerInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}

	// 哈希密码
	hashedPassword, err := utils.HashPassword(registerInfo.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器崩溃，请稍后再试"})
		return
	}

	// 设置哈希后的密码
	registerInfo.Password = hashedPassword

	// 保存用户信息到数据库
	err = global.Db.Create(&registerInfo).Error
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.login_id" {
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": "用户已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "注册失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "注册成功"})
}

func AdminLogin(c *gin.Context) {
	// 获取登录信息
	var loginInfo struct {
		LoginId  string `json:"loginId" binding:"required"`
		LoginPwd string `json:"loginPwd" binding:"required"`
	}
	// 绑定 JSON 数据到结构体
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	// 从数据库查找用户信息
	var user models.User
	if err := global.Db.Where("login_id = ?", loginInfo.LoginId).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器崩溃，请稍后再试"})
		return
	}

	// 验证密码
	if !utils.CheckPassword(loginInfo.LoginPwd, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "密码错误"})
		return
	}

	// 生成 JWT 令牌
	token, err := utils.GenerateJWT(user.LoginId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器崩溃，请稍后再试"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"token":   token,
	})
}

type Info struct {
	LoginId  string `json:"loginId"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Photo    string `json:"photo"`
	RoleName string `json:"roleName"`
}

func GetAdminInfo(c *gin.Context) {
	loginId := c.Query("loginId")
	var userInfo Info
	err := global.Db.Table("users").
		Select("users.login_id, users.name, users.phone, imgs.url AS photo, roles.role_name AS role_name").
		Joins("LEFT JOIN roles ON users.role_id = roles.id"). // 明确连接两个表
		Joins("LEFT JOIN imgs ON users.img_id = imgs.id").
		Where("users.login_id = ?", loginId).
		Scan(&userInfo). // 使用 Scan 将结果直接映射到 Info 结构体
		Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器崩溃，请稍后再试"})
		return
	}
	if userInfo.LoginId == "" {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    userInfo,
	})
}

func GetAdminList(c *gin.Context) {
	pageIndex, err1 := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, err2 := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	roleId, err3 := strconv.Atoi(c.DefaultQuery("roleId", "0"))
	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}

	var list []Info
	var total int64

	query := global.Db.Table("users").
		Select("users.login_id, users.name, users.phone, imgs.url AS photo, roles.role_name AS role_name").
		Joins("LEFT JOIN roles ON users.role_id = roles.id").
		Joins("LEFT JOIN imgs ON users.img_id = imgs.id")

	if roleId != 0 {
		query = query.Where("users.role_id = ?", roleId)
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

func AddAdmin(c *gin.Context) {
	var newAdmin models.User
	// 绑定 JSON 数据到结构体
	if err := c.ShouldBindJSON(&newAdmin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}

	// 哈希密码
	hashedPassword, err := utils.HashPassword(newAdmin.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器崩溃，请稍后再试"})
		return
	}
	newAdmin.Password = hashedPassword

	// 设置默认头像
	if newAdmin.ImgId == 0 {
		newAdmin.ImgId = defaultPhotoID
	} else {
		// 检查提供的头像是否存在
		var img models.Img
		if err := global.Db.First(&img, newAdmin.ImgId).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "头像不存在"})
			return
		}
	}
	// 保存信息到数据库
	err = global.Db.Create(&newAdmin).Error
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.login_id" {
			c.JSON(http.StatusConflict, gin.H{"success": false, "message": "用户已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "添加用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "添加用户成功"})
}

func DeleteAdmin(c *gin.Context) {
	var req struct {
		LoginId string `json:"loginId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	if req.LoginId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "账号不能为空"})
		return
	}
	if req.LoginId == "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "不能删除系统账号"})
		return
	}
	// 删除用户
	err := global.Db.Where("login_id = ?", req.LoginId).Delete(&models.User{}).Error
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除用户失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除用户成功",
	})
}

type updateInfo struct {
	LoginId string `json:"loginId" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	ImgId   uint   `json:"imgId"`
	RoleId  uint   `json:"roleId" binding:"required"`
}

func UpdateAdmin(c *gin.Context) {
	var req updateInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}

	var user models.User
	if err := global.Db.Where("login_id = ?", req.LoginId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不存在"})
		return
	}

	user.Name = req.Name
	user.Phone = req.Phone

	// 检查头像
	if req.ImgId != 0 {
		var img models.Img
		if err := global.Db.First(&img, req.ImgId).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "头像不存在"})
			return
		}
		user.Img = img
		user.ImgId = req.ImgId
	}

	// 检查角色
	var role models.Role
	if err := global.Db.First(&role, req.RoleId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "角色不存在"})
		return
	}
	user.Role = role
	user.RoleId = req.RoleId

	if err := global.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新用户信息成功",
	})
}

func ResetAdminPwd(c *gin.Context) {
	var req struct {
		LoginId     string `json:"loginId" binding:"required"`
		LoginPwd    string `json:"loginPwd" binding:"required"`
		NewLoginPwd string `json:"newLoginPwd" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	var user models.User
	if err := global.Db.Where("login_id = ?", req.LoginId).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器崩溃，请稍后再试"})
		return
	}
	if !utils.CheckPassword(req.LoginPwd, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "原密码错误"})
		return
	}
	hashedPassword, err := utils.HashPassword(req.NewLoginPwd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器崩溃，请稍后再试"})
		return
	}
	user.Password = hashedPassword // 更新密码
	err = global.Db.Save(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "重置密码失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "重置密码成功",
	})
}
