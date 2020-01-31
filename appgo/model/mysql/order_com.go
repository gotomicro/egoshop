package mysql

import (
	"time"
)

type OrderGoods struct {
	Id        int        `gorm:"not null;primary_key"json:"id"` // 主键ID
	CreatedAt time.Time  `gorm:""json:"createdAt"`              // 创建时间
	UpdatedAt time.Time  `gorm:""json:"updatedAt"`              // 更新时间
	DeletedAt *time.Time `gorm:"index"json:"deletedAt"`         // 删除时间
	Uid       int        `gorm:"not null;"json:"uid"`           // 买家ID
	OrderId   int        `gorm:"not null;"json:"orderId"`       // 订单ID
	TypeId    int        `gorm:"not null;"json:"typeId"`
	FeedId    int        `gorm:"not null;"json:"feedId"`
	Title     string     `gorm:"not null;"json:"title"`     // 商品名称
	Price     float64    `gorm:"not null;"json:"price"`     // 商品价格 just价格
	PriceType int        `gorm:"not null;"json:"priceType"` // 1 为总款，2 为预付款
	PayPrice  float64    `gorm:"not null;"json:"payPrice"`  // 商品实际支付价格(拼团商品适用)
	Cover     string     `gorm:"not null;"json:"cover"`     // 商品图片
	Num       int        `gorm:"not null;"json:"num"`       // 商品数量

	ComId    int `gorm:"not null;"json:"comId"`    // 商品主表ID
	ComSkuId int `gorm:"not null;"json:"comSkuId"` // 商品ID

	ComType           int     `gorm:"not null;"json:"comType"`           // 类型[1默认 2团购商品 3限时折扣商品 4组合套装 5赠品]
	ComFreightWay     string  `gorm:"not null;"json:"comFreightWay"`     // 商品运费方式
	ComFreightFee     float64 `gorm:"not null;"json:"comFreightFee"`     // 商品的运费
	EvaluateState     int     `gorm:"not null;"json:"evaluateState"`     // 评价状态[0未评价 1已评价 2已追评]
	EvaluateTime      int64   `gorm:"not null;"json:"evaluateTime"`      // 评价时间
	CouponId          int     `gorm:"not null;"json:"couponId"`          // 线上卡券,大于0线上卡券,微信卡券表表ID,一个规格的商品对应一张微信卡券
	CouponCardId      string  `gorm:"not null;"json:"couponCardId"`      // 线上卡券,微信卡券表微信卡券ID
	LockState         int     `gorm:"not null;"json:"lockState"`         // 锁定状态[1退款中]
	RefundHandleState int     `gorm:"not null;"json:"refundHandleState"` // 退款平台处理状态[默认0处理中(未处理) 10拒绝(驳回) 20同意 30成功(已完成) 50取消(用户主动撤销) 51取消(用户主动收货)]
	RefundId          int     `gorm:"not null;"json:"refundId"`          // 退款ID
	GroupPrice        float64 `gorm:"not null;"json:"groupPrice"`        // 拼团价格 just价格
	CaptainPrice      float64 `gorm:"not null;"json:"captainPrice"`      // 团长价格 just价格
	ComRevisePrice    float64 `gorm:"not null;"json:"comRevisePrice"`    // 修改过的商品实际支付费用 大于0起作用
	RefundState       int     `gorm:"-"json:"refundState"`               // 退款状态
	IsRefund          int     `gorm:"-"json:"isRefund"`                  // 是否退款
}

func (*OrderGoods) TableName() string {
	return "order_goods"
}
