package models

type Billing struct {
	ID           uint   `json:"id" gorm:"primary_key"`
	Amount       int32  `json:"amount" gorm:"not null"`
	Time         string `json:"time" gorm:"not null"`
	GuestId      uint   `json:"guestId" gorm:"not null"`
	Guest        Guest
	RoomId       string `json:"roomId" gorm:"not null"`
	Room         Room
	RoomTypeName string `json:"roomTypeName" gorm:"not null"`
	ResideId     uint   `json:"resideId" gorm:"not null"`
	Reside       Reside
}
