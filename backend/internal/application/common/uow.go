// 浣滅敤锛氥€愭牳蹇冩灑绾姐€慤nitOfWork 鎺ュ彛瀹氫箟锛堝簲鐢ㄥ眰濂戠害锛夛紝鐢ㄤ簬淇濊瘉璺ㄥ涓粨鍌ㄦ搷浣滅殑浜嬪姟涓€鑷存€?
package common

import "context"

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
