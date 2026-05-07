
package grpc

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewUserServiceServer, NewGRPCServer)
