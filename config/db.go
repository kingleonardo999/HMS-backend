package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"hotel-management-system/global"
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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
}
