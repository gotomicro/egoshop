package mysql

import (
	"time"
)

type Order struct {
	Id                 int        `gorm:"not null;primary_key"json:"id"` // 主键id
	CreatedAt          time.Time  `gorm:""json:"createdAt"`              // 创建时间
	UpdatedAt          time.Time  `gorm:""json:"updatedAt"`              // 更新时间
	DeletedAt          *time.Time `gorm:"index"json:"deletedAt"`         // 删除时间
	OrderType          int        `gorm:"not null;"json:"orderType"`     // 1 为猫咪，2为商品
	Sn                 string     `gorm:"not null;"json:"sn"`            // 订单编号
	PaySn              string     `gorm:"not null;"json:"paySn"`         // 支付单号
	Uid                int        `gorm:"not null;"json:"uid"`           // 买家id
	UserName           string     `gorm:"not null;"json:"userName"`      // 买家姓名
	UserPhone          string     `gorm:"not null;"json:"userPhone"`     // 买家手机号码
	UserEmail          string     `gorm:"not null;"json:"userEmail"`     // 买家电子邮箱
	CreatedBy          int        `gorm:"not null;"json:"createdBy"`
	UpdatedBy          int        `gorm:"not null;"json:"updatedBy"`
	From               int        `gorm:"not null;"json:"from"`   // 1WEB2mobile
	Remark             string     `gorm:"not null;"json:"remark"` // 订单备注
	Amount             float64    `gorm:"not null;"json:"amount"` // 订单总价格
	GoodsAmount        float64    `gorm:"not null;"json:"goodsAmount"`
	PdAmount           float64    `gorm:"not null;"json:"pdAmount"`           // 预存款支付金额
	FreightFee         float64    `gorm:"not null;"json:"freightFee"`         // 实际支付的运费
	FreightUnifiedFee  float64    `gorm:"not null;"json:"freightUnifiedFee"`  // 运费统一运费
	FreightTemplateFee float64    `gorm:"not null;"json:"freightTemplateFee"` // 运费模板运费
	GoodsNum           int        `gorm:"not null;"json:"goodsNum"`           // 商品个数
	ReviseAmount       float64    `gorm:"not null;"json:"reviseAmount"`       // 修改过的订单总价格 大于0起作用
	ReviseFreightFee   float64    `gorm:"not null;"json:"reviseFreightFee"`   // 修改过的实际支付的运费
	PaymentCode        string     `gorm:"not null;"json:"paymentCode"`        // 支付方式名称代码
	PayName            string     `gorm:"not null;"json:"payName"`            // 支付方式
	PaymentTime        int64      `gorm:"not null;"json:"paymentTime"`        // 支付成功的(付款的)时间
	FinishedTime       int64      `gorm:"not null;"json:"finishedTime"`       // 订单完成时间
	State              int        `gorm:"not null;"json:"state"`              // 订单状态[0已取消 10默认未付款 20已付款 30已发货 40已收货]
	RefundAmount       float64    `gorm:"not null;"json:"refundAmount"`       // 退款金额
	RefundState        int        `gorm:"not null;"json:"refundState"`        // 退款状态[0无退款 1部分退款 2全部退款]
	LockState          int        `gorm:"not null;"json:"lockState"`          // 锁定状态[0正常 大于0锁定 默认是0]
	TradeNo            string     `gorm:"not null;"json:"tradeNo"`            // 支付宝交易号OR微信交易号
	EvaluateState      int        `gorm:"not null;"json:"evaluateState"`      // 评价状态 0未评价，1已评价
	PayableTime        int64      `gorm:"not null;"json:"payableTime"`        // 订单可支付时间 下单时间+24小时 时间戳
	OutRequestNo       string     `gorm:"not null;"json:"outRequestNo"`       // 支付宝标识一次退款请求，同一笔交易多次退款需要保证唯一，如需部分退款，则此参数必传。
	PartId             int        `gorm:"not null"json:"partid"`              //部门分类
}

func (*Order) TableName() string {
	return "order"
}
