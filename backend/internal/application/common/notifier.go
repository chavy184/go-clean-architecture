// 浣滅敤锛氭秷鎭帹閫佸櫒鎺ュ彛锛岀敤浜庡悗绔悜鍓嶇鎺ㄩ€佸疄鏃舵秷鎭紙闅旂鍏蜂綋鐨?WebSocket 绛夊疄鐜帮級
package common

type Notifier interface {
	PushToUser(userID string, message []byte) error
	Broadcast(message []byte) error
}
