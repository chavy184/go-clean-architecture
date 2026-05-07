// 浣滅敤锛歐S 鐨?Wire ProviderSet
package websocket

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewHub, NewRouter, NewHandler, NewPusherImpl)
