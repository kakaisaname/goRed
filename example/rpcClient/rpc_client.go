package main

import (
	"errors"
	"github.com/kakaisaname/infra/lb"
	"github.com/kakaisaname/goRed/services"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/rpc"
	"strings"
)

func main() {
	//conf := ini.NewIniFileConfigSource("ec_test.ini")
	//client := eureka.NewClient(conf)
	//client.Start()
	//client.Applications, _ = client.GetApplications()
	//apps := &lb.Apps{Client: client}
	//c := &GoRpcClient{apps: apps}
	//cs := &EnvelopeClientService{client: c, serviceId: "resk"}
	//
	//in := services.RedEnvelopeSendingDTO{
	//	Amount:       decimal.NewFromFloat(1),
	//	UserId:       "47692588035919872",
	//	Username:     "测试用户",
	//	EnvelopeType: services.GeneralEnvelopeType,
	//	Quantity:     2,
	//	Blessing:     "",
	//}
	//out, err := cs.SendOut(in)
	//if err != nil {
	//	logrus.Panic(err)
	//}
	//logrus.Infof("%+v", out)
	c, err := rpc.Dial("tcp", ":18082")
	if err != nil {
		logrus.Panic(err)
	}
	sendout(c)
	//receive(c)

}

func receive(c *rpc.Client) {
	in := services.RedEnvelopeReceiveDTO{
		EnvelopeNo:   "",
		RecvUserId:   "",
		RecvUsername: "",
		AccountNo:    "",
	}
	out := &services.RedEnvelopeItemDTO{}
	err := c.Call("Envelope.Receive", in, out)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Infof("%+v", out)
}

func sendout(c *rpc.Client) {
	in := services.RedEnvelopeSendingDTO{
		Amount:       decimal.NewFromFloat(1),
		UserId:       "1UuUyPbc9cdx3YFhy7MdotQtpbJ",
		Username:     "周五了",
		EnvelopeType: services.GeneralEnvelopeType,
		Quantity:     2,
		Blessing:     "恭喜发财，发红包了",
	}
	out := &services.RedEnvelopeActivity{}
	err := c.Call("EnvelopeRpc.SendOut", in, &out)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Infof("%+v", out)
}



//下面是微服务						***
type GoRpcClient struct {
	apps *lb.Apps
}

func (g *GoRpcClient) Call(serviceId, serviceMethod string, in interface{}, out interface{}) error {

	//通过微服务名称从本地服务注册表中查询应用和应用实例列表
	app := g.apps.Get(strings.ToUpper(serviceId))
	if app == nil {
		return errors.New("没有可用的微服务应用，应用名称：" + serviceId + ",请求：" + serviceMethod)
	}

	//通过负载均衡算法从应用实例列表中选择一个实例
	ins := app.Get(serviceMethod)
	if ins == nil {
		return errors.New("没有可用的应用实例，应用名称：" + serviceId + ",请求：" + serviceMethod)
	}
	//选择的实例IP和端口
	address := ins.Address

	c, err := rpc.Dial("tcp", address)
	if err != nil {
		logrus.Error(err)
	}
	err = c.Call(serviceId, in, &out)
	if err != nil {
		logrus.Error(err)
	}
	defer c.Close()
	return err
}

type EnvelopeClientService struct {
	client    *GoRpcClient
	serviceId string
}

func (e *EnvelopeClientService) SendOut(dto services.RedEnvelopeSendingDTO) (activity *services.RedEnvelopeActivity, err error) {
	activity = &services.RedEnvelopeActivity{}
	err = e.client.Call(e.serviceId, "EnvelopeRpc.SendOut", dto, activity)
	return
}

func (e *EnvelopeClientService) Receive(dto services.RedEnvelopeReceiveDTO) (item *services.RedEnvelopeItemDTO, err error) {
	item = &services.RedEnvelopeItemDTO{}
	err = e.client.Call(e.serviceId, "EnvelopeRpc.Receive", dto, item)
	return
}
