package main

import (
	_ "github.com/kakaisaname/account/core/accounts"
	"github.com/kakaisaname/goRed/apis/gorpc"
	_ "github.com/kakaisaname/goRed/apis/web"
	_ "github.com/kakaisaname/goRed/core/envelopes"
	"github.com/kakaisaname/goRed/jobs"
	"github.com/kakaisaname/infra"
	"github.com/kakaisaname/infra/base"
)

//注册我们的starter  启动器
func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	infra.Register(&base.GoRPCStarter{}) //rpc
	infra.Register(&gorpc.GoRpcApiStarter{})
	infra.Register(&jobs.RefundExpiredJobStarter{}) //定时任务	**
	infra.Register(&base.EurekaStarter{})
	infra.Register(&base.IrisServerStarter{})
	infra.Register(&infra.WebApiStarter{})
	//infra.Register(&accounts.AccountClientStarter{})
	infra.Register(&base.HookStarter{})
}
