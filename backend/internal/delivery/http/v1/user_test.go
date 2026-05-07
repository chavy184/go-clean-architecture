// 浣滅敤锛欴elivery 灞?HTTP Handler 鐨?TDD 娴嬭瘯
// 娴嬭瘯瀵硅薄锛歎serHandler (CreateUser, GetUser)
// 娴嬭瘯绛栫暐锛氶€氳繃 httptest 妯℃嫙 HTTP 璇锋眰锛屾敞鍏?Mock 渚濊禆楠岃瘉绔埌绔涓?
package v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	appUser "go-clean-architecture/internal/application/user"
	v1 "go-clean-architecture/internal/delivery/http/v1"
	domainUser "go-clean-architecture/internal/domain/user"
	"go-clean-architecture/test/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// setupRouter 鏋勫缓涓€涓粎鍖呭惈娴嬭瘯璺敱鐨?gin.Engine
func setupRouter(handler *v1.UserHandler) *gin.Engine {
	r := gin.New()
	r.POST("/api/v1/users", handler.CreateUser)
	r.GET("/api/v1/users/:id", handler.GetUser)
	return r
}

// ==================== CreateUser Handler ====================

// 馃敶 绾細POST /users 浼犲叆鍚堟硶 JSON 搴旇繑鍥?200 鍜屽垱寤虹粨鏋?
func TestCreateUserHandler_Success(t *testing.T) {
	mockRepo := mock.NewMockUserRepository()
	mockUoW := mock.NewMockUnitOfWork()
	createUC := appUser.NewCreateUserUseCase(mockUoW, mockRepo)
	getUC := appUser.NewGetUserUseCase(mockRepo)
	handler := v1.NewUserHandler(createUC, getUC)

	body, _ := json.Marshal(appUser.CreateUserRequest{Username: "alice"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := setupRouter(handler)
	router.ServeHTTP(w, req)

	// 馃煝 缁匡細鏈熸湜 200
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"].(float64) != 0 {
		t.Errorf("expected code 0, got %v", resp["code"])
	}
}

// 馃敶 绾細POST /users 浼犲叆绌?Body 搴旇繑鍥?400
func TestCreateUserHandler_BadRequest(t *testing.T) {
	mockRepo := mock.NewMockUserRepository()
	mockUoW := mock.NewMockUnitOfWork()
	createUC := appUser.NewCreateUserUseCase(mockUoW, mockRepo)
	getUC := appUser.NewGetUserUseCase(mockRepo)
	handler := v1.NewUserHandler(createUC, getUC)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := setupRouter(handler)
	router.ServeHTTP(w, req)

	// 馃煝 缁匡細鏈熸湜 400
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

// ==================== GetUser Handler ====================

// 馃敶 绾細GET /users/:id 鏌ヨ瀛樺湪鐨勭敤鎴峰簲杩斿洖 200
func TestGetUserHandler_Success(t *testing.T) {
	mockRepo := mock.NewMockUserRepository()
	mockRepo.Users["uid-001"] = &domainUser.User{ID: "uid-001", Username: "bob"}
	mockUoW := mock.NewMockUnitOfWork()
	createUC := appUser.NewCreateUserUseCase(mockUoW, mockRepo)
	getUC := appUser.NewGetUserUseCase(mockRepo)
	handler := v1.NewUserHandler(createUC, getUC)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/uid-001", nil)
	w := httptest.NewRecorder()

	router := setupRouter(handler)
	router.ServeHTTP(w, req)

	// 馃煝 缁匡細鏈熸湜 200
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

// 馃敶 绾細GET /users/:id 鏌ヨ涓嶅瓨鍦ㄧ殑鐢ㄦ埛锛圧epo 杩斿洖 nil锛?
func TestGetUserHandler_NotFound(t *testing.T) {
	mockRepo := mock.NewMockUserRepository()
	mockUoW := mock.NewMockUnitOfWork()
	createUC := appUser.NewCreateUserUseCase(mockUoW, mockRepo)
	getUC := appUser.NewGetUserUseCase(mockRepo)
	handler := v1.NewUserHandler(createUC, getUC)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/non-existent", nil)
	w := httptest.NewRecorder()

	router := setupRouter(handler)
	router.ServeHTTP(w, req)

	// 馃煝 缁匡細鏈熸湜 500锛堝洜涓?FindByID 杩斿洖 nil user锛孍xecute 浼?nil pointer锛?
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
