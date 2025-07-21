package models

type Order struct {
	Id         uint   `json:"id" gorm:"primary_key"`
	OrderId    string `json:"orderId" gorm:"unique;not null"`
	GuestId    uint   `json:"guestId" gorm:"not null"`
	Guest      Guest
	RoomId     string `json:"roomId" gorm:"not null"`
	Room       Room   `gorm:"foreignKey:RoomId;references:RoomId"`
	OrderDate  string `json:"orderDate" gorm:"not null"`
	LeaveDate  string `json:"leaveDate" gorm:"not null"`
	TotalMoney int32  `json:"totalMoney" gorm:"not null"`
	GuestNum   int32  `json:"guestNum" gorm:"not null"`
}
