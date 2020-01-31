package constx

const (
	OrderNew = 10
	OrderOk  = 20

	OrderStateDefault    = 0
	OrderStateNew        = 1 // 新建订单,未支付订单
	OrderStatePay        = 2 // 已支付
	OrderStateSend       = 3 // 已发货
	OrderStateSuccess    = 4 // 已收货，交易成功
	OrderStateClose      = 5 // 取消订单
	OrderStateCancel     = 6
	OrderStateUnevaluate = 7

	// 未支付订单 OrderPay
	OrderPayStateNew = 1
	// 已支付订单 OrderPay
	OrderPayStatePay = 2

	OrderPayCreate = 1 // 创建order pay
	OrderPayOK     = 2 // 完成order pay

)

var OrderStates = map[string]int{
	"stateClose":      1, // TODO
	"stateNew":        OrderStateNew,
	"statePay":        OrderStatePay,
	"stateSend":       OrderStateSend,
	"stateSuccess":    OrderStateSuccess,
	"stateCancel":     OrderStateCancel,
	"stateRefund":     1,                    // TODO user表handle_state
	"stateUnevaluate": OrderStateUnevaluate, // TODO user表evaluate_state
}

type GoodsType struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Checked bool   `json:"checked"`
}

var GoodsTypeMap = map[int]string{
	1: "is_pay",
	2: "is_star",
	3: "is_create",
	4: "is_share",
	5: "is_read",
}
