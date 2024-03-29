package envelopes

import (
	"context"
	"errors"
	"fmt"

	acservices "github.com/kakaisaname/account/services"
	"github.com/kakaisaname/goRed/services"
	"github.com/kakaisaname/infra/base"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

const (
	pageSize = 100 //定义每次查出100
)

type ExpiredEnvelopeDomain struct {
	expiredGoods []RedEnvelopeGoods //返回的是一个切片
	offset       int
}

//查询出过期红包,			****
//一次性只查出部分过期红包，而不是查出全部的红包（如果红包多的话）				***
func (e *ExpiredEnvelopeDomain) Next() (ok bool) {

	//在事务方法里面进行		**
	base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		e.expiredGoods = dao.FindExpired(e.offset, pageSize)
		if len(e.expiredGoods) > 0 { //证明查询出了数据
			e.offset += len(e.expiredGoods)
			ok = true
		}
		return nil
	})
	return ok
}

//定义过期方法															***
func (e *ExpiredEnvelopeDomain) Expired() (err error) {
	for e.Next() { //通过循环 Next方法 	**
		for _, g := range e.expiredGoods {
			if g.OrderType == services.OrderTypeSending {

				log.Debugf("过期红包退款开始：%+v", g)
				err := e.ExpiredOne(g)
				if err != nil {
					log.Error(err)
				}
				log.Debugf("过期红包退款结束：%+v", g)
			}
		}
	}
	return err
}

//发起退款流程 		（针对每一个过期红包发起退款流程）										***
func (e *ExpiredEnvelopeDomain) ExpiredOne(goods RedEnvelopeGoods) (err error) {
	//创建一个退款订单
	refund := goods
	refund.OrderType = services.OrderTypeRefund
	//refund.RemainAmount = goods.RemainAmount.Mul(decimal.NewFromFloat(-1)) 			//红包剩余总金额
	//refund.RemainQuantity = -goods.RemainQuantity										//红包剩余总数量
	refund.Status = services.OrderExpired //订单状态
	refund.PayStatus = services.Refunding //支付状态为退款中
	//refund.OriginEnvelopeNo = goods.EnvelopeNo										//原订单号				其实表数据结构中应该加上该字段		**
	refund.EnvelopeNo = ""
	domain := goodsDomain{RedEnvelopeGoods: refund}
	domain.createEnvelopeNo()

	err = base.Tx(func(runner *dbx.TxRunner) error {
		txCtx := base.WithValueContext(context.Background(), runner)
		_, err := domain.Save(txCtx)
		if err != nil {
			return errors.New("创建退款订单失败" + err.Error())
		}
		//修改原订单订单状态																	****
		dao := RedEnvelopeGoodsDao{runner: runner}
		_, err = dao.UpdateOrderStatus(goods.EnvelopeNo, services.OrderExpired)
		if err != nil {
			return errors.New("更新原订单状态失败" + err.Error())
		}
		return nil
	})
	if err != nil {
		return
	}
	//调用资金账户接口进行转账
	systemAccount := base.GetSystemAccount()
	account := acservices.GetAccountService().GetEnvelopeAccountByUserId(goods.UserId)
	if account == nil {
		return errors.New("没有找到该用户的红包资金账户:" + goods.UserId)
	}
	body := acservices.TradeParticipator{
		Username:  systemAccount.Username,
		UserId:    systemAccount.UserId,
		AccountNo: systemAccount.AccountNo,
	}
	target := acservices.TradeParticipator{
		Username:  account.Username,
		UserId:    account.UserId,
		AccountNo: account.AccountNo,
	}
	transfer := acservices.AccountTransferDTO{
		TradeBody:   body,
		TradeTarget: target,
		TradeNo:     domain.RedEnvelopeGoods.EnvelopeNo,
		Amount:      goods.RemainAmount,
		AmountStr:   goods.RemainAmount.String(),
		ChangeType:  acservices.EnvelopExpiredRefund,
		ChangeFlag:  acservices.FlagTransferIn,
		Decs:        "红包过期退款:" + goods.EnvelopeNo,
	}
	fmt.Printf("body: %+v\n", body)
	fmt.Printf("target: %+v\n", target)
	status, err := acservices.GetAccountService().Transfer(transfer)
	if status != acservices.TransferedStatusSuccess {
		return err
	}

	err = base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		//修改原订单状态
		rows, err := dao.UpdateOrderStatus(goods.EnvelopeNo, services.OrderExpiredRefundSuccessful)
		if err != nil || rows == 0 {
			return errors.New("更新原订单状态失败")
		}
		//修改退款订单状态
		rows, err = dao.UpdateOrderStatus(refund.EnvelopeNo, services.OrderExpiredRefundSuccessful)
		if err != nil {
			return errors.New("更新退款订单状态失败" + err.Error())
		}
		if rows == 0 {
			fmt.Println(rows, "更新退款订单状态失败")
		}
		return nil
	})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//把测试和debug作为作业留给同学们，
// 过程中遇到问题和bug可以留言和发起问题讨论
