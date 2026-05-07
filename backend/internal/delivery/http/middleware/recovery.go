// 浣滅敤锛氬叏灞€寮傚父鎹曡幏涓棿浠讹紝澶勭悊 HTTP 璇锋眰涓殑 panic锛岃繑鍥炵粺涓€ 500 閿欒
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryMiddleware 鎹曡幏 panic锛岄槻姝㈣繘绋嬪穿婧冿紝骞惰褰曞爢鏍堟棩蹇?
func RecoveryMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("HTTP Server Panic Recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  "Internal Server Error",
				})
			}
		}()
		c.Next()
	}
}
