// 作用：程序入口，负责初始化配置、日志、依赖注入并启动HTTP/gRPC/WS服务，实现优雅停机
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// @title           Go Clean Architecture API
// @version         1.0
// @description     This is a sample server for Go Clean Architecture.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	// 初始化 App
	app, err := InitializeApp()
	if err != nil {
		fmt.Printf("Failed to initialize app: %v\n", err)
		return
	}
	
	app.Logger.Info("App starting...")

	// 启动 WebSocket Hub
	go app.WSHub.Run()

	// 启动 HTTP 服务
	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.Router,
	}
	go func() {
		app.Logger.Info("HTTP server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatal("http listen error", zap.Error(err))
		}
	}()

	// 启动 gRPC 服务
	go func() {
		if err := app.GRPCSrv.Start(":9090"); err != nil {
			app.Logger.Fatal("grpc listen error", zap.Error(err))
		}
	}()

	// 优雅停机：等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.Logger.Info("Shutdown Server ...")

	// 设置一个 30 秒的超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. 关闭 HTTP 服务
	if err := srv.Shutdown(ctx); err != nil {
		app.Logger.Fatal("Server Shutdown:", zap.Error(err))
	}
	
	// 2. 优雅关闭 gRPC 服务
	app.GRPCSrv.Server.GracefulStop()

	// 3. 关闭 Redis 连接池
	if err := app.Redis.Close(); err != nil {
		app.Logger.Error("Redis Close Error", zap.Error(err))
	}

	// 4. 关闭 Postgres 连接池
	if sqlDB, err := app.DB.DB(); err == nil && sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			app.Logger.Error("Postgres Close Error", zap.Error(err))
		}
	}

	// 5. 将日志缓冲落盘
	_ = app.Logger.Sync()

	fmt.Println("Server exiting")
}
