// 浣滅敤锛氭帴鏀?WS 娑堟伅锛岃皟鐢?application 灞?
package websocket

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleAction(client *Client, msg []byte) {
	// call usecase
}
