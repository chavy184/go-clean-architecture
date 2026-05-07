// 浣滅敤锛氬綋鍓嶅井鏈嶅姟涓撶敤鐨勫唴閮ㄥ伐鍏峰簱锛氱粺涓€鐨?HTTP JSON 鍝嶅簲鏍煎紡灏佽
package response

type BaseResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}
