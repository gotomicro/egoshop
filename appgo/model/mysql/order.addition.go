package mysql

import (
	"strings"

	"github.com/goecology/egoshop/appgo/model/constx"
)

/**
 * 取得订单状态文字输出形式
 */

func (orderInfo Order) GetOrderState() string {
	switch orderInfo.State {
	case constx.OrderStateClose:
		return "已取消"
	case constx.OrderStateNew:
		return "待付款"
	case constx.OrderStatePay:
		return "待发货"
	case constx.OrderStateSend:
		return "待收货"
	case constx.OrderStateSuccess:
		return "交易完成"
	}
	return "系统错误"
}

func (orderInfo Order) GetOrderPaymentName() string {
	if strings.Index(orderInfo.PaymentCode, "offline") > 0 {
		return strings.ReplaceAll(orderInfo.PaymentCode, "offline", "货到付款")
	} else if strings.Index(orderInfo.PaymentCode, "online") > 0 {
		return strings.ReplaceAll(orderInfo.PaymentCode, "online", "在线付款")
	} else if strings.Index(orderInfo.PaymentCode, "alipay") > 0 {
		return strings.ReplaceAll(orderInfo.PaymentCode, "alipay", "支付宝")
	} else if strings.Index(orderInfo.PaymentCode, "tenpay") > 0 {
		return strings.ReplaceAll(orderInfo.PaymentCode, "tenpay", "财付通")
	} else if strings.Index(orderInfo.PaymentCode, "chinabank") > 0 {
		return strings.ReplaceAll(orderInfo.PaymentCode, "chinabank", "网银在线")
	} else if strings.Index(orderInfo.PaymentCode, "predeposit") > 0 {
		return strings.ReplaceAll(orderInfo.PaymentCode, "predeposit", "预存款")
	}
	return "系统错误"
}
