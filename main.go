package main

import (
	"context"
	"errors"
	"fmt"
	"hotel-management-system/config"
	"hotel-management-system/models"
	"hotel-management-system/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// 初始化配置
	config.InitConfig()
	// 迁移数据库表
	models.Init()
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
	fmt.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	fmt.Println("Server exiting gracefully")
}
