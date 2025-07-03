package main

import (
	"context"
	"errors"
	"hotel-management-system/config"
	"hotel-management-system/routers"
	"hotel-management-system/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// 初始化配置
	config.InitConfig()
	if config.Config.ImgConfig.EnableImgUpload {
		// 启动图片上传协程
		go utils.Upload2cloud()
	}
	// global.Db.AutoMigrate(&models.Img{}, &models.User{}, &models.Role{}, &models.RoomType{}, &models.Room{})
	r := routers.SetupRouters()
	port := config.Config.App.Port // 获取配置文件中的端口号
	if port == "" {
		port = "8080" // 设置默认端口为8080
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
	// 关闭上传协程
	if config.Config.ImgConfig.EnableImgUpload {
		close(utils.ImgUploadChan)
		log.Println("Image upload channel closed")
	}
	// 删除临时图片
	if err := utils.DeleteTempImages(); err != nil {
		log.Println("Error deleting temporary images:", err)
	} else {
		log.Println("Temporary images deleted successfully")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting gracefully")
}
