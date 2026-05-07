// 浣滅敤锛歡RPC Server 鐨勫垵濮嬪寲涓庢嫤鎴櫒閰嶇疆
package grpc

import (
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCServer 灏佽鍘熺敓 grpc.Server
type GRPCServer struct {
	Server      *grpc.Server
	UserService *UserServiceServer
	logger      *zap.Logger
}

func NewGRPCServer(userService *UserServiceServer, logger *zap.Logger) *GRPCServer {
	// 閰嶇疆 Recovery
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Error("gRPC Server Panic Recovered", zap.Any("panic", p))
			return status.Errorf(codes.Internal, "Internal Server Error")
		}),
	}

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOpts...),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(recoveryOpts...),
		),
	)
	
	// 鍦ㄦ澶勬敞鍐?protobuf 鐨?server 瀹炵幇
	// pb.RegisterUserServiceServer(srv, userService)

	return &GRPCServer{
		Server:      srv,
		UserService: userService,
		logger:      logger,
	}
}

func (s *GRPCServer) Start(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	s.logger.Info("gRPC server listening on " + port)
	return s.Server.Serve(lis)
}
