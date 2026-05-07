// 浣滅敤锛氬疄鐜?pb 鎺ュ彛锛屽唴閮ㄨ皟鐢?application 灞?
package grpc

import (
	appUser "go-clean-architecture/internal/application/user"
)

type UserServiceServer struct {
	// pb.UnimplementedUserServiceServer
	createUC *appUser.CreateUserUseCase
	getUC    *appUser.GetUserUseCase
}

func NewUserServiceServer(createUC *appUser.CreateUserUseCase, getUC *appUser.GetUserUseCase) *UserServiceServer {
	return &UserServiceServer{createUC: createUC, getUC: getUC}
}
