// 浣滅敤锛氱湡姝ｅ彲浠ヨ法寰湇鍔″叡浜殑绾妧鏈寘锛屽锛氬簳灞傚瓧绗︿覆澶勭悊宸ュ叿锛屼笉鍖呭惈浠讳綍涓氬姟閫昏緫
package strutil

import "strings"

func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
