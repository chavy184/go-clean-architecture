// 浣滅敤锛氱敤鎴蜂粨鍌ㄦ帴鍙ｅ畾涔夛紙涓嶅惈浠讳綍 GORM 鎴栨暟鎹簱閫昏緫锛夛紝瀹氫箟鐢ㄦ埛鑱氬悎鐨勬寔涔呭寲濂戠害
package user

import "context"

type Repository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
}
