package models

import "hotel-management-system/global"

func Init() {
	if global.Db == nil {
		panic("数据库未初始化")
	} else {
		_ = global.Db.AutoMigrate(
			&Role{},
			&User{},
			&Img{},
			&RoomType{},
			&Room{},
			&RoomStatus{},
			&Guest{},
			&Reside{},
			&ResideState{},
			&Order{},
			&Billing{},
		)
	}
}
