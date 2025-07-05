package models

type Role struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	RoleName string `json:"roleName" gorm:"unique;not null"` // 角色名称，唯一且不能为空
}
