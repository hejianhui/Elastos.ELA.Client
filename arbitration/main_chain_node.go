package arbitration

import "ELAClient/rpc"

type MainChainNode interface {

	GetCurrentHeight() (uint32, error)
	GetBlockByHeight(height uint32) (rpc.BlockInfo)
}

