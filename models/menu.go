package models

type Menu struct {
	ID     uint     `json:"id" gorm:"primary_key"`
	Name   string   `json:"name" gorm:"not null"`
	TypeId uint     `json:"typeId" gorm:"not null"`
	Type   MenuType `gorm:"foreignKey:TypeId;references:ID"`
	Price  int32    `json:"price" gorm:"not null"`
	ImgId  uint     `json:"imgId" gorm:"not null"`
	Img    Img      `gorm:"foreignKey:ImgId;references:ID"`
	Desc   string   `json:"desc" gorm:"not null"`
}
