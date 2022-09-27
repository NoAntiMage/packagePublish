package constant

import "errors"

const (
	RpcTryConnecting = "RpcTryConnecting"
	RpcPublished     = "RpcPublished"
	RpcConnected     = "RpcConnected"

	RpcDefaultExpireTime = 900
)

var (
	ErrRpcUnreachable = errors.New("ErrRpcUnreachable")
)
