package models

type RoomType struct {
	ID              uint   `json:"id" gorm:"primary_key"`
	RoomTypeName    string `json:"roomTypeName" gorm:"unique;not null"` // 房间类型名称，唯一且不能为空
	RoomTypePrice   int32  `json:"roomTypePrice" gorm:"not null"`       // 房间类型价格，不能为空
	TypeDescription string `json:"typeDescription" gorm:"not null"`     // 房间类型描述，不能为空
	BedNum          int32  `json:"bedNum" gorm:"not null"`              // 床位数量，不能为空
}
