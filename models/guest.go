package models

type Guest struct {
	ID         uint   `json:"id" gorm:"primary_key"`
	IdentityId string `json:"identityId" gorm:"unique;not null"`
	Name       string `json:"guestName" gorm:"not null"`
	Phone      string `json:"guestPhone" gorm:"not null"`
}
