// 浣滅敤锛欻TTP 涓棿浠讹細澶勭悊閾捐矾杩借釜 (Trace ID)
package middleware

import (
	"go-clean-architecture/internal/pkg/trace"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HeaderTraceID = "X-Trace-Id"

// TraceMiddleware 纭繚姣忎釜璇锋眰閮芥湁 Trace ID锛屽苟灏嗗叾娉ㄥ叆鍒?context 涓?
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 灏濊瘯浠庡墠绔紶鍏ョ殑 Header 鑾峰彇
		traceID := c.GetHeader(HeaderTraceID)
		if traceID == "" {
			// 濡傛灉鍓嶇娌℃湁浼犻€掞紝鍒欑敱鍚庣鑷姩鐢熸垚
			traceID = uuid.New().String()
		}

		// 2. 灏?Trace ID 鍐欏洖 Response Header锛屾柟渚垮鎴风鎺掓煡闂
		c.Header(HeaderTraceID, traceID)

		// 3. 灏?Trace ID 娉ㄥ叆鍒板師鐢熺殑 context.Context 涓?
		// 骞朵笖鏇存柊 Gin 涓簳灞傚師鐢熺殑 Request Context
		ctx := trace.WithTraceID(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)

		// 娉ㄦ剰锛氬悗缁殑 Handler 鎴?UseCase 搴斿綋鐩存帴閫氳繃 c.Request.Context() 灏?context 寰€鍚庝紶
		c.Next()
	}
}
