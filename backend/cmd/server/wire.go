//go:build wireinject
// +build wireinject

// 作用：顶层 Wire 依赖注入组装文件，定义依赖图
package main

import (
	"github.com/google/wire"
	"go-clean-architecture/config"
	"go-clean-architecture/internal/application/user"
	httpDelivery "go-clean-architecture/internal/delivery/http"
	"go-clean-architecture/internal/delivery/http/v1"
	grpcDelivery "go-clean-architecture/internal/delivery/grpc"
	wsDelivery "go-clean-architecture/internal/delivery/websocket"
	"go-clean-architecture/internal/infrastructure/persistence/postgres"
	"go-clean-architecture/internal/infrastructure/persistence/redis"
	"go-clean-architecture/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	goRedis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"go.uber.org/zap"
)

func InitializeApp() (*App, error) {
	// 组装所有分层的 ProviderSet
	wire.Build(
		config.ProviderSet,      // 配置解析层
		logger.ProviderSet,      // 日志等基础组件层
		redis.ProviderSet,       // 缓存层
		postgres.ProviderSet,    // 基础设施持久化层
		user.ProviderSet,        // 应用层(用户用例)
		v1.ProviderSet,
		httpDelivery.ProviderSet,
		grpcDelivery.ProviderSet,
		wsDelivery.ProviderSet,  // 交付层
		NewApp,                  // 顶层应用实例
	)
	return &App{}, nil
}

type App struct {
	GRPCSrv *grpcDelivery.GRPCServer
	WSHub   *wsDelivery.Hub
	Router  *gin.Engine
	DB      *gorm.DB
	Redis   *goRedis.Client
	Logger  *zap.Logger
}

func NewApp(grpcSrv *grpcDelivery.GRPCServer, wsHub *wsDelivery.Hub, router *gin.Engine, db *gorm.DB, rdb *goRedis.Client, logger *zap.Logger) *App {
	return &App{
		GRPCSrv: grpcSrv,
		WSHub:   wsHub,
		Router:  router,
		DB:      db,
		Redis:   rdb,
		Logger:  logger,
	}
}
