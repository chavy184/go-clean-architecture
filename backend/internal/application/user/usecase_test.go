// 浣滅敤锛欰pplication 灞?UseCase 鐨?TDD 娴嬭瘯
// 娴嬭瘯瀵硅薄锛欳reateUserUseCase, GetUserUseCase
// 娴嬭瘯绛栫暐锛氱函鍗曞厓娴嬭瘯 鈥?閫氳繃 Mock 闅旂涓€鍒囧閮ㄤ緷璧栵紙Repo, UoW锛?
package user_test

import (
	"context"
	"errors"
	"testing"

	appUser "go-clean-architecture/internal/application/user"
	domainUser "go-clean-architecture/internal/domain/user"
	"go-clean-architecture/test/mock"
)

// ==================== CreateUserUseCase ====================

// 馃敶 绾細绌虹敤鎴峰悕搴旇杩斿洖閿欒
func TestCreateUser_EmptyUsername(t *testing.T) {
	mockRepo := mock.NewMockUserRepository()
	mockUoW := mock.NewMockUnitOfWork()
	uc := appUser.NewCreateUserUseCase(mockUoW, mockRepo)

	_, err := uc.Execute(context.Background(), &appUser.CreateUserRequest{Username: ""})

	// 馃煝 缁匡細鏈熸湜鎷垮埌鏍￠獙閿欒
	if err == nil {
		t.Fatal("expected error for empty username, got nil")
	}
	// 鏂█ Repo 娌℃湁琚皟鐢紙鏍￠獙搴旇鍦ㄦ寔涔呭寲涔嬪墠锛?
	if mockRepo.SaveCalled {
		t.Error("Save should NOT be called when username is empty")
	}
}

// 馃敶 绾細nil 璇锋眰搴旇杩斿洖閿欒
func TestCreateUser_NilRequest(t *testing.T) {
	mockRepo := mock.NewMockUserRepository()
	mockUoW := mock.NewMockUnitOfWork()
	uc := appUser.NewCreateUserUseCase(mockUoW, mockRepo)

	_, err := uc.Execute(context.Background(), nil)

	if err == nil {
		t.Fatal("expected error for nil request, got nil")
	}
}

// 馃敶 绾細姝ｅ父鍒涘缓搴旇鎴愬姛锛岃繑鍥炲甫 ID 鐨勫搷搴?
func TestCreateUser_Success(t *testing.T) {
	mockRepo := mock.NewMockUserRepository()
	mockUoW := mock.NewMockUnitOfWork()
	uc := appUser.NewCreateUserUseCase(mockUoW, mockRepo)

	resp, err := uc.Execute(context.Background(), &appUser.CreateUserRequest{Username: "alice"})

	// 馃煝 缁匡細鏈熸湜鏃犻敊璇?
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.ID == "" {
		t.Error("expected non-empty ID in response")
	}
	// 鏂█ Repo.Save 纭疄琚皟鐢ㄤ簡
	if !mockRepo.SaveCalled {
		t.Error("expected Save to be called")
	}
	if mockRepo.SaveCalledWith.Username != "alice" {
		t.Errorf("expected username 'alice', got '%s'", mockRepo.SaveCalledWith.Username)
	}
}

// 馃敶 绾細褰?Repo.Save 杩斿洖閿欒鏃讹紝UseCase 涔熷簲璇ヨ繑鍥為敊璇?
func TestCreateUser_RepoSaveError(t *testing.T) {
	mockRepo := mock.NewMockUserRepository()
	mockRepo.SaveErr = errors.New("db connection lost")
	mockUoW := mock.NewMockUnitOfWork()
	uc := appUser.NewCreateUserUseCase(mockUoW, mockRepo)

	_, err := uc.Execute(context.Background(), &appUser.CreateUserRequest{Username: "alice"})

	// 馃煝 缁匡細鏈熸湜閿欒鍚戜笂浼犳挱
	if err == nil {
		t.Fatal("expected error when repo fails, got nil")
	}
}

// 馃敶 绾細褰?UoW 浜嬪姟澶辫触鏃讹紝UseCase 涔熷簲璇ヨ繑鍥為敊璇?
func TestCreateUser_UoWError(t *testing.T) {
	mockRepo := mock.NewMockUserRepository()
	mockUoW := mock.NewMockUnitOfWork()
	mockUoW.DoErr = errors.New("transaction begin failed")
	uc := appUser.NewCreateUserUseCase(mockUoW, mockRepo)

	_, err := uc.Execute(context.Background(), &appUser.CreateUserRequest{Username: "alice"})

	if err == nil {
		t.Fatal("expected error when UoW fails, got nil")
	}
}

// ==================== GetUserUseCase ====================

// 馃敶 绾細绌?ID 搴旇杩斿洖閿欒
func TestGetUser_EmptyID(t *testing.T) {
	mockRepo := mock.NewMockUserRepository()
	uc := appUser.NewGetUserUseCase(mockRepo)

	_, err := uc.Execute(context.Background(), "")

	if err == nil {
		t.Fatal("expected error for empty id, got nil")
	}
}

// 馃敶 绾細鏌ヨ瀛樺湪鐨勭敤鎴峰簲璇ユ垚鍔?
func TestGetUser_Success(t *testing.T) {
	mockRepo := mock.NewMockUserRepository()
	mockRepo.Users["user-123"] = &domainUser.User{ID: "user-123", Username: "bob"}
	uc := appUser.NewGetUserUseCase(mockRepo)

	resp, err := uc.Execute(context.Background(), "user-123")

	// 馃煝 缁?
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.ID != "user-123" {
		t.Errorf("expected ID 'user-123', got '%s'", resp.ID)
	}
}

// 馃敶 绾細褰?Repo 杩斿洖閿欒鏃跺簲璇ヤ紶鎾?
func TestGetUser_RepoError(t *testing.T) {
	mockRepo := mock.NewMockUserRepository()
	mockRepo.FindByIDErr = errors.New("db timeout")
	uc := appUser.NewGetUserUseCase(mockRepo)

	_, err := uc.Execute(context.Background(), "any-id")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
