package gorpc

import (
	"github.com/kakaisaname/goRed/services"
)

type EnvelopeRpc struct {
}

//Go内置的RPC接口有一些规范：
//										1. 入参和出参都要作为方法参数
//										2. 方法必须有2个参数，并且是可导出类型
//										3. 第二个参数（返回值）必须是指针类型			(返回参数必须为指针，入参可为或不可为指针)
//										4. 方法返回值要返回error类型
//										5. 方法必须是可导出的							（第一个字母是大写）
func (e *EnvelopeRpc) SendOut(
	in services.RedEnvelopeSendingDTO,
	out *services.RedEnvelopeActivity) error {
	s := services.GetRedEnvelopeService() //调唯一的暴露点
	a, err := s.SendOut(in)               //调发红包方法
	if err != nil {
		return err
	}
	a.CopyTo(out) //将返回的 a 拷贝到 out中，不能修改出参（out）的引用
	return err    //完成后需要将接口注册到rpc  server中去						***
}

func (e *EnvelopeRpc) Receive(
	in services.RedEnvelopeReceiveDTO,
	out *services.RedEnvelopeItemDTO) error {
	s := services.GetRedEnvelopeService()
	a, err := s.Receive(in)
	if err != nil {
		return err
	}
	a.CopeTo(out)
	return err
}

//感兴趣的同学，可以尝试使用thrift来适配用户接口
