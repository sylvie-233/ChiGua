package main

import (
	"chigua-backend/config"
	"chigua-backend/database"
	"chigua-backend/internal/middleware"
	"chigua-backend/internal/router"
	"chigua-backend/utils/ip2region"
	"chigua-backend/utils/logger"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		// 使用默认日志配置输出错误
		logger.InitLogger("info", "logs/app.log")
		logger.Fatalf("加载配置失败: %v", err)
	}

	// 使用配置初始化日志
	logger.InitLogger(config.AppConfig.Logger.Level, config.AppConfig.Logger.LogFile)
	defer logger.CloseLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.Fatalf("程序崩溃: %v", r)
		}
	}()

	// 初始化数据库连接
	// if err := database.InitDatabase(); err != nil {
	if err := database.InitDatabase(); err != nil {
		logger.Fatalf("初始化数据库失败: %v", err)
	}
	defer database.CloseDatabase()

	// 初始化IP2Region
	ip2region.InitIp2Region()
	defer ip2region.CloseIp2Region()

	// 设置Gin模式
	gin.SetMode(config.AppConfig.Server.Mode)

	r := gin.New()
	// 添加Recovery中间件，捕获panic并恢复
	r.Use(gin.Recovery())
	// 添加自定义日志中间件
	r.Use(middleware.LoggerMiddleware())

	// 路由挂载
	router.InitRouter(r)

	// 使用配置的端口
	serverAddr := fmt.Sprintf(":%s", config.AppConfig.Server.Port)
	logger.Infof("服务器启动在 %s", serverAddr)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	// 在goroutine中启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	// kill (无参数) 默认发送 syscall.SIGTERM
	// kill -2 是 syscall.SIGINT
	// kill -9 是 syscall.SIGKILL 但不能被捕获，所以不需要添加
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("正在关闭服务器...")

	// 设置5秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅地关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("服务器关闭错误: %v", err)
	}

	logger.Info("服务器已关闭")
}
