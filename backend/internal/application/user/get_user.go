// 浣滅敤锛氭煡璇㈢敤鎴风敤渚嬶紝璐熻矗澶勭悊鏌ヨ鐢ㄦ埛鐨勪笟鍔￠€昏緫
package user

import (
	"context"
	"errors"

	domainUser "go-clean-architecture/internal/domain/user"
)

type GetUserUseCase struct {
	repo domainUser.Repository
}

func NewGetUserUseCase(repo domainUser.Repository) *GetUserUseCase {
	return &GetUserUseCase{repo: repo}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, id string) (*CreateUserResponse, error) {
	if id == "" {
		return nil, errors.New("user id is required")
	}

	u, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	return &CreateUserResponse{ID: u.ID}, nil
}
