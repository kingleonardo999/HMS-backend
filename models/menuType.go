package models

type MenuType struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Type string `json:"type" gorm:"unique;not null"`
}
