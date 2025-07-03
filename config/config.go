package config

import (
	"github.com/spf13/viper"
	"log"
)

type AppConfig struct {
	App struct {
		Port string
		Name string
	}
	Database struct {
		Host     string
		Port     string
		Username string
		Password string
		DBName   string
	}
	ImgConfig struct {
		EnableImgUpload bool
	}
	DeleteConfig struct {
		EnableDelete bool
	}
}

var Config *AppConfig

func InitConfig() {
	// 1. 路径
	viper.AddConfigPath("./config")
	// 2. 文件名
	viper.SetConfigName("config")
	// 3. 文件类型
	viper.SetConfigType("yaml")
	// 4. 读取配置
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	// 5. 绑定配置到结构体变量
	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	// 6. 设置全局变量
	Config = &config

	// 初始化数据库
	initDB()
}
