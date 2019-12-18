package envelopes

import (
	"context"
	"errors"

	acservices "github.com/kakaisaname/account/services"
	"github.com/kakaisaname/infra/base"
	"github.com/tietang/dbx"
	"goRed/services"
	"path"
)

//发红包业务领域代码   					*****
//返回的是发红包的信息
func (d *goodsDomain) SendOut(
	goods services.RedEnvelopeGoodsDTO) (activity *services.RedEnvelopeActivity, err error) {
	//创建红包商品
	d.Create(goods)
	//创建活动							***
	activity = new(services.RedEnvelopeActivity)
	//红包链接，格式：http://域名/v1/envelope/{id}/link/
	link := base.GetEnvelopeActivityLink()
	domain := base.GetEnvelopeDomain()                    //红包域名
	activity.Link = path.Join(domain, link, d.EnvelopeNo) //path.Join		d.EnvelopeNo:红包编号 	**

	err = base.Tx(func(runner *dbx.TxRunner) (err error) {
		ctx := base.WithValueContext(context.Background(), runner)
		//事务逻辑问题：
		//保存红包商品和红包金额的支付必须要保证全部成功或者全部失败
		//保存红包商品
		id, err := d.Save(ctx) //ctx 上下文对象 		**
		if id <= 0 || err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return nil, err
	}
	//红包金额支付
	//1. 需要红包中间商的红包资金账户，定义在配置文件中，事先初始化到资金账户表中
	//2. 从红包发送人的资金账户中扣减红包金额 ，把红包金额从红包发送人的资金账户里扣除
	body := acservices.TradeParticipator{
		AccountNo: goods.AccountNo,
		UserId:    goods.UserId,
		Username:  goods.Username,
	}
	systemAccount := base.GetSystemAccount()
	target := acservices.TradeParticipator{
		AccountNo: systemAccount.AccountNo,
		Username:  systemAccount.Username,
		UserId:    systemAccount.UserId,
	}

	//转账对象 		***
	transfer := acservices.AccountTransferDTO{
		TradeBody:   body,
		TradeTarget: target,
		TradeNo:     d.EnvelopeNo,
		AmountStr:   d.Amount.String(),
		ChangeType:  acservices.EnvelopeOutgoing,
		ChangeFlag:  acservices.FlagTransferOut,
		Decs:        "红包金额支付",
	}

	acsvs := acservices.GetAccountService()

	status, err := acsvs.Transfer(transfer)
	if err != nil {
		return nil, err
	}
	if status != acservices.TransferedStatusSuccess {
		return nil, errors.New("转账失败！")
	}
	//3. 将扣减的红包总金额转入红包中间商的红包资金账户
	//入账
	transfer = acservices.AccountTransferDTO{
		TradeBody:   target,
		TradeTarget: body,
		TradeNo:     d.EnvelopeNo,
		AmountStr:   d.Amount.String(),
		ChangeType:  acservices.EnvelopeIncoming,
		ChangeFlag:  acservices.FlagTransferIn,
		Decs:        "红包金额转入",
	}
	status, err = acsvs.Transfer(transfer)
	if err != nil {
		return nil, err
	}
	if status != acservices.TransferedStatusSuccess {
		return nil, errors.New("转账失败！")
	}
	err = base.Tx(func(runner *dbx.TxRunner) (err error) {
		_, err = d.UpdatePayStatus(d.EnvelopeNo, services.Payed) //更新支付状态    已支付
		return err
	})
	if err != nil {
		return nil, err
	}
	//扣减金额没有问题，返回活动

	activity.RedEnvelopeGoodsDTO = *d.RedEnvelopeGoods.ToDTO()

	return activity, err
}
