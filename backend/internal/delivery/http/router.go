// 浣滅敤锛欸in 璺敱娉ㄥ唽锛岀粺涓€绠＄悊 HTTP 鎺ュ彛鍜屼腑闂翠欢
package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/google/wire"
	"go.uber.org/zap"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/gin-contrib/pprof"

	"go-clean-architecture/internal/delivery/http/middleware"
	v1 "go-clean-architecture/internal/delivery/http/v1"
	
	_ "go-clean-architecture/docs" // 瀵煎叆鐢熸垚鐨?swagger docs
)

var ProviderSet = wire.NewSet(RegisterRoutes)

// RegisterRoutes 娉ㄥ唽鎵€鏈?HTTP 璺敱鍜屼腑闂翠欢
// 浼犻€?logger 鐢ㄤ簬涓棿浠?
func RegisterRoutes(userHandler *v1.UserHandler, logger *zap.Logger) *gin.Engine {
	router := gin.Default()

	// 娉ㄥ唽鍏ㄥ眬寮傚父鎹曡幏
	router.Use(middleware.RecoveryMiddleware(logger))

	// 娉ㄥ唽鍏ㄥ眬閾捐矾杩借釜涓棿浠?
	router.Use(middleware.TraceMiddleware())

	// 娉ㄥ唽鎬ц兘鎸囨爣缁熻涓棿浠?
	router.Use(middleware.MetricsMiddleware())

	// 鏆撮湶 Prometheus Metrics 鎶撳彇绔偣
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 娉ㄥ唽 pprof 鎬ц兘鍒嗘瀽璺敱锛岄粯璁ゅ墠缂€ /debug/pprof (鍙厤鍚堥壌鏉冩満鍒朵娇鐢?
	pprof.Register(router)

	// 娉ㄥ唽 Swagger 璺敱
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 璺敱缁?
	api := router.Group("/api/v1")
	{
		api.POST("/users", userHandler.CreateUser)
	}

	return router
}
