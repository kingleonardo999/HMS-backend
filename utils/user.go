package utils

import (
	"errors"
	"gorm.io/gorm"
	"hotel-management-system/global"
	"hotel-management-system/models"
)

func GetUserInfo(loginId string) (models.User, error) {
	var userInfo models.User
	// 预加载 Role 和 Img
	err := global.Db.Preload("Role").Preload("Img").Where("login_id = ?", loginId).First(&userInfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userInfo, errors.New("user not found")
		}
		return userInfo, err
	}
	return userInfo, nil
}

func CheckUserExists(loginId string) (bool, error) {
	var count int64
	err := global.Db.Model(&models.User{}).Where("login_id = ?", loginId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
