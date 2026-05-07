// 浣滅敤锛氫簨浠舵€荤嚎鎺ュ彛瀹氫箟锛屾彁渚涘簲鐢ㄥ眰鍜岄鍩熷眰鐨勪簨浠跺彂甯冧笌璁㈤槄鑳藉姏
package common

type EventBus interface {
	Publish(event interface{}) error
}
