// 浣滅敤锛氫笟鍔＄敤渚嬶細鍒涘缓鐢ㄦ埛锛屽寘鍚?CreateUserUseCase锛屾敞鍏?UnitOfWork 涓?Repository 鍗忚皟瀹屾垚涓氬姟閫昏緫
package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/google/wire"

	"go-clean-architecture/internal/application/common"
	domainUser "go-clean-architecture/internal/domain/user"
)

var ProviderSet = wire.NewSet(NewCreateUserUseCase, NewGetUserUseCase)

type CreateUserUseCase struct {
	uow  common.UnitOfWork
	repo domainUser.Repository
}

func NewCreateUserUseCase(uow common.UnitOfWork, repo domainUser.Repository) *CreateUserUseCase {
	return &CreateUserUseCase{uow: uow, repo: repo}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	// 1. 鍏ュ弬鏍￠獙
	if req == nil || req.Username == "" {
		return nil, errors.New("username is required")
	}

	// 2. 鏋勫缓棰嗗煙瀹炰綋
	u := &domainUser.User{
		ID:       uuid.New().String(),
		Username: req.Username,
	}

	// 3. 鍦?UoW 浜嬪姟涓墽琛屾寔涔呭寲
	err := uc.uow.Do(ctx, func(txCtx context.Context) error {
		return uc.repo.Save(txCtx, u)
	})
	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{ID: u.ID}, nil
}
