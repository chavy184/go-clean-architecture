// 浣滅敤锛氬崟涓繛鎺ョ殑鐘舵€佹満 (鍖呭惈 ReadPump 鍜?WritePump)
package websocket

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	UserID string
	send   chan []byte
}

func (c *Client) ReadPump(router *Router) {
	defer func() {
		if r := recover(); r != nil {
			c.Hub.logger.Error("WebSocket ReadPump Panic Recovered", zap.Any("panic", r))
		}
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		router.Route(c, msg)
	}
}

func (c *Client) WritePump() {
	defer func() {
		if r := recover(); r != nil {
			c.Hub.logger.Error("WebSocket WritePump Panic Recovered", zap.Any("panic", r))
		}
		c.Conn.Close()
	}()
	for message := range c.send {
		c.Conn.WriteMessage(websocket.TextMessage, message)
	}
	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}
