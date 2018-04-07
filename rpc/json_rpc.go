package rpc

import (
	core_rpc "github.com/elastos/Elastos.ELA.Client.Core/rpc"
	"github.com/elastos/Elastos.ELA.Client/common/config"
)

func init() {
	core_rpc.Url = config.Params().Host
}
