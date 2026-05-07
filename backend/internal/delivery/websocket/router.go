// 浣滅敤锛歐S 鍐呴儴鐨勬秷鎭矾鐢卞櫒 (鎸夋秷鎭被鍨嬪垎鍙?
package websocket

import (
	"go.uber.org/zap"
)

type Router struct {
	handler *Handler
	logger  *zap.Logger
}

func NewRouter(handler *Handler, logger *zap.Logger) *Router {
	return &Router{handler: handler, logger: logger}
}

func (r *Router) Route(client *Client, msg []byte) {
	r.logger.Info("Received WS message")
	r.handler.HandleAction(client, msg)
}
