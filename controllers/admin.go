package controllers

import (
	"github.com/gin-gonic/gin"
	"hotel-management-system/config"
	"hotel-management-system/global"
	"hotel-management-system/models"
	"hotel-management-system/utils"
	"log"
	"net/http"
	"strconv"
)

var (
	baseUrl        = "https://imgs.161517.xyz/" // 基础 URL
	defaultPhotoID = uint(1)                    // 默认头像ID
)

func AdminRegister(c *gin.Context) {
	// 获取注册信息
	var registerInfo models.User
	// 绑定 JSON 数据到结构体
	if err := c.ShouldBindJSON(&registerInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}

	// 检查用户是否已存在
	exists, err := utils.CheckUserExists(registerInfo.LoginId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器崩溃，请稍后再试"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "用户已存在"})
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
	user, err := utils.GetUserInfo(loginInfo.LoginId)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "该用户不存在"})
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
	user, err := utils.GetUserInfo(loginId)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "该用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器崩溃，请稍后再试"})
		return
	}
	info := Info{
		LoginId:  user.LoginId,
		Name:     user.Name,
		Phone:    user.Phone,
		Photo:    user.Img.Url,
		RoleName: user.Role.RoleName,
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    info,
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
	var adminList []models.User
	var total int64
	query := global.Db.Model(&models.User{})
	// 预加载角色信息，以便获取角色名称；预加载头像信息，以便获取头像URL
	query = global.Db.Model(&models.User{}).Preload("Role").Preload("Img")
	if roleId != 0 {
		// 如果指定了角色ID，则查询该角色的管理员
		query = query.Where("role_id = ?", roleId)
	}
	query.Count(&total)
	query.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&adminList)

	// list
	var list []Info
	for i := range adminList {
		list = append(list, Info{
			LoginId:  adminList[i].LoginId,
			Name:     adminList[i].Name,
			Phone:    adminList[i].Phone,
			Photo:    adminList[i].Img.Url,
			RoleName: adminList[i].Role.RoleName,
		})
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

	// 检查用户是否已存在
	exists, err := utils.CheckUserExists(newAdmin.LoginId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器崩溃，请稍后再试"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"success": false, "message": "用户已存在"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "添加用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "添加用户成功"})
}

func DeleteAdmin(c *gin.Context) {
	var post struct {
		LoginId string `json:"loginId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	if post.LoginId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "登录ID不能为空"})
		return
	}
	if post.LoginId == "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "不能删除系统账号"})
		return
	}
	// 检查用户是否存在
	exists, err := utils.CheckUserExists(post.LoginId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器崩溃，请稍后再试"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "该用户不存在"})
		return
	}
	// 删除用户
	err = global.Db.Where("login_id = ?", post.LoginId).Delete(&models.User{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除用户失败"})
		return
	}
	if config.Config.DeleteConfig.EnableDelete {
		// 硬删除
		err = global.Db.Unscoped().Where("login_id = ?", post.LoginId).Delete(&models.User{}).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除用户失败"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除用户成功",
	})
}

func UpdateAdmin(c *gin.Context) {
	var updateInfo struct {
		LoginId string `json:"loginId" binding:"required"`
		Name    string `json:"name" binding:"required"`
		Phone   string `json:"phone" binding:"required"`
		ImgId   uint   `json:"imgId"`
		RoleId  uint   `json:"roleId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&updateInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	user, err := utils.GetUserInfo(updateInfo.LoginId)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "该用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器崩溃，请稍后再试"})
		return
	}
	// 更新用户信息
	user.Name = updateInfo.Name
	user.Phone = updateInfo.Phone
	log.Println("更新用户信息:", user.LoginId, user.Name, user.Phone, updateInfo.ImgId, updateInfo.RoleId)
	if updateInfo.ImgId != 0 {
		// 检查提供的头像是否存在
		var img models.Img
		if err := global.Db.First(&img, updateInfo.ImgId).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "头像不存在"})
			return
		}
		// 如果头像ID不为0，则更新头像ID, 注意将关联的 Img 对象设置为空
		user.Img = img                // 清除之前的头像关联
		user.ImgId = updateInfo.ImgId // 更新头像ID
	}
	// 检查角色是否存在
	var role models.Role
	if err := global.Db.First(&role, updateInfo.RoleId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "角色不存在"})
		return
	}
	user.Role = role                // 清除之前的角色关联
	user.RoleId = updateInfo.RoleId // 更新角色ID
	err = global.Db.Save(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新用户信息失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新用户信息成功",
	})
}

func ResetAdminPwd(c *gin.Context) {
	var post struct {
		LoginId     string `json:"loginId" binding:"required"`
		LoginPwd    string `json:"loginPwd" binding:"required"`
		NewLoginPwd string `json:"newLoginPwd" binding:"required"`
	}
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "输入无效"})
		return
	}
	user, err := utils.GetUserInfo(post.LoginId)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "该用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "服务器崩溃，请稍后再试"})
		return
	}
	if !utils.CheckPassword(post.LoginPwd, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "原密码错误"})
		return
	}
	hashedPassword, err := utils.HashPassword(post.NewLoginPwd)
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
