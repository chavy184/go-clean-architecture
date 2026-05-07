// 浣滅敤锛氥€愬叧閿€戝疄鐜?application.common.Notifier锛屽弽鍚戣皟鐢?Hub
package websocket

import "go-clean-architecture/internal/application/common"

type PusherImpl struct {
	hub *Hub
}

func NewPusherImpl(hub *Hub) common.Notifier {
	return &PusherImpl{hub: hub}
}

func (p *PusherImpl) PushToUser(userID string, message []byte) error {
	p.hub.SendToUser(userID, message)
	return nil
}

func (p *PusherImpl) Broadcast(message []byte) error {
	p.hub.broadcast <- message
	return nil
}
