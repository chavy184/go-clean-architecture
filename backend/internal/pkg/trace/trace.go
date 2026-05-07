// 浣滅敤锛氭彁渚?Trace ID 鐨勪笂涓嬫枃娉ㄥ叆涓庢彁鍙栧伐鍏凤紝淇濊瘉鍏ㄩ摼璺拷韪兘鍔?
package trace

import (
	"context"

	"github.com/google/uuid"
)

// 浣跨敤绉佹湁绫诲瀷闃叉 context key 鍐茬獊
type traceKey struct{}

// WithTraceID 灏?Trace ID 娉ㄥ叆鍒?context.Context 涓?
func WithTraceID(ctx context.Context, traceID string) context.Context {
	if traceID == "" {
		traceID = uuid.New().String()
	}
	return context.WithValue(ctx, traceKey{}, traceID)
}

// GetTraceID 浠?context.Context 涓彁鍙?Trace ID
func GetTraceID(ctx context.Context) string {
	if val, ok := ctx.Value(traceKey{}).(string); ok {
		return val
	}
	return ""
}
