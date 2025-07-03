package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	RoleName string `json:"roleName" gorm:"unique;not null"` // 角色名称，唯一且不能为空
}
