package gorpc

import (
	"github.com/kakaisaname/infra"
	"github.com/kakaisaname/infra/base"
)

type GoRpcApiStarter struct {
	infra.BaseStarter
}

//放在Init阶段，先进行注册												***
func (g *GoRpcApiStarter) Init(ctx infra.StarterContext) {
	base.RpcRegister(new(EnvelopeRpc)) //EnvelopeRpc 这个结构体符合RPC的一些规范		***
}
