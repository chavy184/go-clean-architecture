// 浣滅敤锛欻TTP Handler锛岃礋璐ｈВ鏋?HTTP 璇锋眰锛堝 JSON銆丵uery锛夊苟璋冪敤 Application 灞傜殑瀵瑰簲閫昏緫锛屾渶鍚庣粍瑁?HTTP 鍝嶅簲
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	appUser "go-clean-architecture/internal/application/user"
)

var ProviderSet = wire.NewSet(NewUserHandler)

type UserHandler struct {
	createUC *appUser.CreateUserUseCase
	getUC    *appUser.GetUserUseCase
}

func NewUserHandler(createUC *appUser.CreateUserUseCase, getUC *appUser.GetUserUseCase) *UserHandler {
	return &UserHandler{createUC: createUC, getUC: getUC}
}

// CreateUser godoc
// @Summary      鍒涘缓鐢ㄦ埛
// @Description  鍒涘缓涓€涓柊鐨勭敤鎴?// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      appUser.CreateUserRequest  true  "Create User Request"
// @Success      200      {object}  appUser.CreateUserResponse
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req appUser.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	resp, err := h.createUC.Execute(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": resp})
}

// GetUser godoc
// @Summary      鏌ヨ鐢ㄦ埛
// @Description  鏍规嵁 ID 鏌ヨ鐢ㄦ埛
// @Tags         users
// @Produce      json
// @Param        id   path  string  true  "User ID"
// @Success      200  {object}  appUser.CreateUserResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.getUC.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": resp})
}
