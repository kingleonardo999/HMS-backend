package models

type Room struct {
	ID              uint       `json:"id" gorm:"primary_key"`
	RoomId          string     `json:"roomId" gorm:"unique;not null"` // 房间ID，唯一且不能为空
	RoomTypeId      uint       `json:"roomTypeId" gorm:"not null"`    // 房间类型ID，不能为空
	RoomType        RoomType   `gorm:"foreignKey:RoomTypeId;references:ID"`
	RoomStatusId    uint       `json:"roomStatusId" gorm:"not null"` // 房间状态ID，不能为空
	RoomStatus      RoomStatus `gorm:"foreignKey:RoomStatusId;references:ID"`
	RoomDescription string     `json:"roomDescription"` // 房间描述，允许为空
}

type RoomStatus struct {
	ID         uint   `json:"statusId" gorm:"primary_key"`
	StatusName string `json:"statusName" gorm:"unique;not null"` // 房间状态名称，唯一且不能为空
}
