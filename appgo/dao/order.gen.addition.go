package dao

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/i2eco/egoshop/appgo/model/common"
	"github.com/i2eco/egoshop/appgo/model/constx"
	"github.com/i2eco/egoshop/appgo/model/mysql"
	"github.com/i2eco/egoshop/appgo/model/trans"
	"github.com/jinzhu/gorm"
	"github.com/thoas/go-funk"
)

const ComplainTimeLimit = 24 * 60 * 60 * 7

const ORDER_EVALUATE_TIME = 1296000 // todo 需要可修改

type OrderFull struct {
	Order          mysql.Order       `json:"order"`
	StateDesc      string            `json:"stateDesc"`
	GroupStateDesc string            `json:"groupStateDesc"`
	PaymentName    string            `json:"paymentName"`
	UserOpen       mysql.UserOpen    `json:"userOpen"`
	OrderExtend    mysql.OrderExtend `json:"orderExtend"`
	OrderGoods     []OrderGoodsJson  `json:"orderGoods"`
	// 新增信息
	IsCancel       bool `json:"isCancel"`
	IsPay          bool `json:"isPay"`
	IsRefundCancel bool `json:"isRefundCancel"`
	IsComplain     bool `json:"isComplain"`
	IsReceive      bool `json:"isReceive"`
	IsLock         bool `json:"isLock"`
	IsDeliver      bool `json:"isDeliver"`
	IsEvaluate     bool `json:"isEvaluate"`
}

type OrderGoodsJson struct {
	mysql.OrderGoods
	IsRefund bool `json:"isRefund"`
}

func (o *order) Pay(c *gin.Context, paySn, paymentCode, tradeNO string) (err error) {
	tx := o.db.Begin()

	_, err = OrderPay.InfoX(c, mysql.Conds{
		"pay_sn":    paySn,
		"pay_state": constx.OrderPayCreate,
	})
	if err != nil && err == gorm.ErrRecordNotFound {
		err = errors.New("订单支付信息不存在")
		tx.Rollback()
		return
	}

	order, err := o.InfoX(c, mysql.Conds{
		"pay_sn": paySn,
	})
	if err != nil && err == gorm.ErrRecordNotFound {
		err = errors.New("该订单不存在")
		tx.Rollback()
		return
	}

	err = OrderPay.UpdateX(c, tx, mysql.Conds{
		"pay_sn": paySn,
	}, mysql.Ups{
		"pay_state": constx.OrderPayOK,
	})
	if err != nil {
		err = errors.New("更新订单支付状态失败")
		tx.Rollback()
		return
	}

	// 修改订单
	orderUpdate := map[string]interface{}{
		"state":          constx.OrderStatePay,
		"payment_time":   time.Now().Unix(),
		"payment_code":   paymentCode,
		"trade_no":       tradeNO,
		"out_request_no": fmt.Sprintf("HZ01RF00%d", order.Id),
	}

	err = o.UpdateX(c, tx, mysql.Conds{
		"pay_sn": paySn,
		"state":  constx.OrderStateNew,
	}, orderUpdate)
	if err != nil {
		err = errors.New("更新订单状态失败")
		tx.Rollback()
		return
	}

	err = OrderLog.Create(c, tx, &mysql.OrderLog{
		OrderId:    order.Id,
		Role:       "buyer",
		Msg:        fmt.Sprintf("支付成功，支付平台交易号 : %s", tradeNO),
		OrderState: constx.OrderOk,
	})
	if err != nil {
		err = errors.New("记录订单日志出现错误")
		tx.Rollback()
		return
	}

	// 修改用户和该商品的关系

	err = UserGoods.CreateOrUpdate(tx, order.Uid, common.Goods{
		Gid:    order.Id,
		TypeId: constx.TypeCom,
	}, "is_pay")
	if err != nil {
		err = errors.New("更新失敗")
		tx.Rollback()
		return
	}
	// todo 统计数据
	//err = Com.UpdateX(c, tx, mysql.Conds{
	//	"id": order.Id,
	//}, mysql.Ups{
	//	"pay_cnt": gorm.Expr("pay_cnt+?", 1),
	//})
	//
	//if err != nil {
	//	err = errors.New("更新失敗")
	//	tx.Rollback()
	//	return
	//}
	tx.Commit()
	return
}

// GetOrderFullsPage 改写getOrderList
func (f *order) GetOrderFullsPage(c *gin.Context, reqList *trans.ReqPage, conds mysql.Conds, extend map[string]bool) (int, map[int]*OrderFull) {
	total, orders := Order.ListPage(c, conds, reqList)
	tmpMap := make(map[int]*OrderFull, 0)
	orderIds := make([]int, 0, total)
	uids := make([]int, 0, total)
	uidstmp := make(map[int]bool, 0)
	for _, v := range orders {
		ret := &OrderFull{Order: v}
		// 1默认2拼团商品3限时折扣商品4组合套装5赠品

		ret.StateDesc = v.GetOrderState()
		ret.PaymentName = v.GetOrderPaymentName()
		tmpMap[v.Id] = ret

		orderIds = append(orderIds, v.Id)
		if _, ok := uidstmp[v.Uid]; !ok {
			uids = append(uids, v.Uid)
		}
	}

	if _, ok := extend["user"]; ok {
		users, _ := UserOpen.ListMap(c, mysql.Conds{
			"id": mysql.Cond{"in", uids},
		})
		for oid := range tmpMap {
			ret := tmpMap[oid]
			ret.UserOpen = users[ret.Order.Uid]
			tmpMap[oid] = ret
		}
	}

	if _, ok := extend["order_extend"]; ok {
		extendOrderGoods, _ := OrderExtend.ListMap(c, mysql.Conds{
			"id": mysql.Cond{"in", orderIds},
		})
		for oid := range tmpMap {
			ret := tmpMap[oid]
			ret.OrderExtend = extendOrderGoods[oid]
			tmpMap[oid] = ret
		}
	}

	if _, ok := extend["order_goods"]; ok {
		// 根据订单id，找到关联的商品
		orderGoods, _ := OrderGoods.ListMap(c, mysql.Conds{
			"order_id": mysql.Cond{"in", orderIds},
		})
		for oid := range tmpMap {
			ret := tmpMap[oid]
			if ret.OrderGoods == nil {
				ret.OrderGoods = make([]OrderGoodsJson, 0)
			}
			// 遍历所有sku里的goods，如果满足order id，添加到对应的订单表里
			for _, goods := range orderGoods {
				if goods.OrderId == oid {
					ret.OrderGoods = append(ret.OrderGoods, OrderGoodsJson{OrderGoods: goods})
				}
			}
			tmpMap[oid] = ret
		}
	}

	return total, tmpMap
}

func (g *order) GetStateCnt(uid int, state int) (cnt int) {
	g.db.Model(mysql.Order{}).Where("uid = ? AND state = ?", uid, state).Count(&cnt)
	return
}

type OrderCondition struct {
	// 新增信息
	IsCancel       bool `json:"isCancel"`
	IsPay          bool `json:"isPay"`
	IsRefundCancel bool `json:"isRefundCancel"`
	IsComplain     bool `json:"isComplain"`
	IsReceive      bool `json:"isReceive"`
	IsLock         bool `json:"isLock"`
	IsDeliver      bool `json:"isDeliver"`
	IsEvaluate     bool `json:"isEvaluate"`
}

type RespOrderInfo struct {
	Info              mysql.Order        `gorm:"not null;"json:"info"`
	OrderCondition    OrderCondition     `gorm:"not null;"json:"orderCondition"`
	OrderLog          []mysql.OrderLog   `gorm:"not null;"json:"orderLog"`
	ExtendOrderExtend mysql.OrderExtend  `gorm:"not null;"json:"extendOrderExtend"`
	ExtendOrderCom    []mysql.OrderGoods `gorm:"not null;"json:"extendOrderCom"`
}

// 取单条订单信息
func (o *order) GetOrderInfo(c *gin.Context, id int, uid int) (resp RespOrderInfo, err error) {
	var orderInfo mysql.Order
	orderInfo, err = o.InfoX(c, mysql.Conds{
		"id":  id,
		"uid": uid,
	})
	if err != nil {
		return
	}

	_, orderLog := OrderLog.ListPage(c, mysql.Conds{
		"order_id": orderInfo.Id,
	}, &trans.ReqPage{
		Current:  0,
		PageSize: 1000,
	})

	extendInfo, _ := OrderExtend.Info(c, orderInfo.Id)

	//追加返回商品信息

	//取商品列表
	extendOrderGoodsOrigin, _ := OrderGoods.List(c, mysql.Conds{"id": orderInfo.Id}, "id desc")

	extendOrderCom := make([]mysql.OrderGoods, 0)
	//
	for _, value := range extendOrderGoodsOrigin {
		// 退款平台处理状态 默认0处理中(未处理) 10拒绝(驳回) 20同意 30成功(已完成) 50取消(用户主动撤销) 51取消(用户主动收货)
		// 不可退款

		if value.RefundId > 0 && funk.ContainsInt([]int{20, 30, 51}, value.RefundHandleState) {
			value.IsRefund = 0
		} else {
			value.IsRefund = 1
		}

		//refund_state 0不显示申请退款按钮 1显示申请退款按钮 2显示退款中按钮 3显示退款完成
		if orderInfo.RefundState <= 10 {
			value.RefundState = 0
		}

		if orderInfo.RefundState > 10 && orderInfo.LockState == 0 && value.RefundId == 0 {
			value.RefundState = 1

		}

		if orderInfo.RefundState > 10 && (orderInfo.LockState != 0 || value.RefundId != 0) && value.RefundHandleState == 30 {
			value.RefundState = 3
		}

		if orderInfo.RefundState > 10 && (orderInfo.LockState != 0 || value.RefundId != 0) && value.RefundHandleState != 30 {
			value.RefundState = 2
		}

		extendOrderCom = append(extendOrderCom, value)

	}

	orderCondition := OrderCondition{}
	// 显示取消订单
	orderCondition.IsCancel = o.GetOrderOperateState("user_cancel", orderInfo)
	// 显示是否需能支付（todo 计算后台过期时间）
	orderCondition.IsPay = o.GetOrderOperateState("user_pay", orderInfo)
	// 显示退款取消订单
	orderCondition.IsRefundCancel = o.GetOrderOperateState("refund_cancel", orderInfo)
	// 显示投诉
	orderCondition.IsComplain = o.GetOrderOperateState("complain", orderInfo)
	// 显示收货
	orderCondition.IsReceive = o.GetOrderOperateState("receive", orderInfo)
	// 显示锁定中
	orderCondition.IsLock = o.GetOrderOperateState("lock", orderInfo)
	// 显示物流跟踪
	orderCondition.IsDeliver = o.GetOrderOperateState("deliver", orderInfo)
	// 显示评价
	orderCondition.IsEvaluate = o.GetOrderOperateState("evaluate", orderInfo)

	resp = RespOrderInfo{
		Info:              orderInfo,
		OrderCondition:    orderCondition,
		OrderLog:          orderLog,
		ExtendOrderExtend: extendInfo,
		ExtendOrderCom:    extendOrderCom, // refundState
	}
	return
}

//  返回是否允许某些操作
func (*order) GetOrderOperateState(operateState string, orderInfo mysql.Order) bool {
	var state bool
	switch operateState {
	//买家支付
	case "user_pay":
		state = orderInfo.State == constx.OrderStateNew && time.Now().Unix() < orderInfo.PayableTime
	//买家取消订单
	case "user_cancel":
		state = orderInfo.State == constx.OrderStateNew || orderInfo.PaymentCode == "offline" && orderInfo.State == constx.OrderStatePay
	//买家取消订单
	case "refund_cancel":
		state = orderInfo.RefundState == 1 && orderInfo.LockState != 1
		//商家取消订单
	case "cancel":
		state = orderInfo.State == constx.OrderStateNew || orderInfo.PaymentCode == "offline" && funk.ContainsInt([]int{constx.OrderStatePay, constx.OrderStateSend}, orderInfo.State)
		//平台取消订单
	case "system_cancel":
		state = orderInfo.State == constx.OrderStateNew || orderInfo.PaymentCode == "offline" && orderInfo.State == constx.OrderStatePay
		//平台收款
	case "system_receive_pay":
		state = orderInfo.State == constx.OrderStateNew && orderInfo.PaymentCode == "online"
		//买家投诉
	case "complain":
		state = funk.ContainsInt([]int{constx.OrderStateNew, constx.OrderStateSend}, orderInfo.State) || orderInfo.FinishedTime > (time.Now().Unix()-ComplainTimeLimit)
		//调整运费
	case "modify_price":
		state = orderInfo.State == constx.OrderStateNew || orderInfo.PaymentCode == "offline" && orderInfo.State == constx.OrderStatePay
		//发货
	case "send":
		state = orderInfo.LockState == 0 && orderInfo.State == constx.OrderStatePay
		//收货
	case "receive":
		state = orderInfo.State == constx.OrderStateSend
		//评价
	case "evaluate":
		state = orderInfo.LockState == 0 && orderInfo.EvaluateState == 0 && orderInfo.State == constx.OrderStateSuccess && time.Now().Unix()-orderInfo.FinishedTime < ORDER_EVALUATE_TIME
		// 子商品是否可评价
	case "evaluate_goods":
		state = orderInfo.State == constx.OrderStateSuccess
		//锁定
	case "lock":
		state = orderInfo.LockState == 1
		//快递跟踪
	case "deliver":
		//分享
	case "share":
		state = orderInfo.State == constx.OrderStateSuccess
	}
	return state
}
