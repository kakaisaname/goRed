package services

const (
	DefaultBlessing   = "恭喜发财，鸿富猪到！" //默认的祝福语
	DefaultTimeFormat = "2006-01-02.15:04:05"
)

//订单类型：发布单、退款单
type OrderType int

const (
	OrderTypeSending OrderType = 1 //发送		**
	OrderTypeRefund  OrderType = 2 //退回		**
)

//支付状态：未支付，支付中，已支付，支付失败
//退款：未退款，退款中，已退款，退款失败
type PayStatus int

const (
	PayNothing PayStatus = 1 //未支付
	Paying     PayStatus = 2 //支付中
	Payed      PayStatus = 3 //已支付
	PayFailure PayStatus = 4 //支付失败
	//
	RefundNothing PayStatus = 61 //未退款
	Refunding     PayStatus = 62 //退款中
	Refunded      PayStatus = 63 //已退款
	RefundFailure PayStatus = 64 //退款失败
)

//红包订单状态：创建、发布、过期、失效、过期退款成功，过期退款失败
type OrderStatus int

const (
	OrderCreate                  OrderStatus = 1
	OrderSending                 OrderStatus = 2
	OrderExpired                 OrderStatus = 3
	OrderDisabled                OrderStatus = 4
	OrderExpiredRefundSuccessful OrderStatus = 5
	OrderExpiredRefundFalured    OrderStatus = 6
)

//红包类型：普通红包，碰运气红包
type EnvelopeType int

const (
	GeneralEnvelopeType = 1 //普通红包
	LuckyEnvelopeType   = 2 //碰运气红包
)

var EnvelopeTypes = map[EnvelopeType]string{
	GeneralEnvelopeType: "普通红包",
	LuckyEnvelopeType:   "碰运气红包",
}
