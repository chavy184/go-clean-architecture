// 浣滅敤锛氱敤渚嬩笓鐢ㄨ緭鍏ヨ緭鍑?DTO (Data Transfer Object)锛岀敤浜庨殧绂昏〃绀哄眰鍜屽簲鐢ㄥ眰鐨勬暟鎹粨鏋?
package user

type CreateUserRequest struct {
	Username string
}

type CreateUserResponse struct {
	ID string
}
