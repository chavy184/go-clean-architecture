// 浣滅敤锛氳鍗曚粨鍌ㄦ帴鍙ｅ畾涔夛紙涓嶅惈浠讳綍 GORM 鎴栨暟鎹簱閫昏緫锛夛紝閫氳繃渚濊禆鍊掔疆鎻愪緵鏁版嵁鎸佷箙鍖栫殑鑳藉姏
package order

import "context"

type Repository interface {
	Save(ctx context.Context, order *Order) error
}
