package models

type Reside struct {
	ID          uint `json:"id" gorm:"primary_key"`
	GuestId     uint `json:"guestId" gorm:"not null"`
	Guest       Guest
	RoomId      string `json:"roomId" gorm:"not null"`
	Room        Room   `gorm:"foreignKey:RoomId;references:RoomId"`
	ResideDate  string `json:"resideDate" gorm:"not null"`
	LeaveDate   string `json:"leaveDate"`
	TotalMoney  int32  `json:"totalMoney" gorm:"not null"`
	Deposit     int32  `json:"deposit" gorm:"not null"`
	GuestNum    int32  `json:"guestNum" gorm:"not null"`
	ResideState string `json:"resideState" gorm:"not null"`
}

type ResideState struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	StateName string `json:"stateName" gorm:"unique;not null"`
}
