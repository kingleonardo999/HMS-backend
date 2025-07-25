package config

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"hotel-management-system/global"
	"hotel-management-system/models"
	"hotel-management-system/utils"
	"log"
	"time"
)

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func initDB() {
	DB := Config.Database

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DB.Username, DB.Password, DB.Host, DB.Port, DB.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束, 使用逻辑外键
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to set database connection pool: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)           // 设置连接池中空闲连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接的最大生命周期为无限制
	global.Db = db                      // 将数据库连接赋值给全局变量

	if global.Db == nil {
		panic("数据库未初始化")
	} else {
		_ = global.Db.AutoMigrate(
			&models.Role{},
			&models.User{},
			&models.Img{},
			&models.RoomType{},
			&models.Room{},
			&models.RoomStatus{},
			&models.Guest{},
			&models.Reside{},
			&models.ResideState{},
			&models.Order{},
			&models.Billing{},
			&models.Menu{},
			&models.MenuType{},
		)
		// 插入系统角色和管理员
		if err := global.Db.First(&models.User{}, "login_id = ?", "admin").Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				global.Db.Create(&models.Role{RoleName: "admin"})
				hashedPassword, err := utils.HashPassword("admin")
				if err != nil {
					panic(err)
				}
				global.Db.Create(&models.User{Name: "admin", LoginId: "admin", Password: hashedPassword, RoleId: 1})
			}
		}
	}
}
