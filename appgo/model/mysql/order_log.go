package mysql

import (
	"time"
)

type OrderLog struct {
	Id         int        `gorm:"not null;primary_key"json:"id"`
	CreatedAt  time.Time  `gorm:""json:"createdAt"`
	UpdatedAt  time.Time  `gorm:""json:"updatedAt"`
	DeletedAt  *time.Time `gorm:""json:"deletedAt"`
	OrderId    int        `gorm:"not null;"json:"orderId"`    // 订单ID
	Msg        string     `gorm:"not null;"json:"msg"`        // 文字描述
	Role       string     `gorm:"not null;"json:"role"`       // 操作角色
	UpdatedBy  int        `gorm:"not null;"json:"updatedBy"`  // 操作人
	OrderState int        `gorm:"not null;"json:"orderState"` // 订单状态[0已取消 10默认未付款 20已付款 30已发货 40已收货]

}

func (*OrderLog) TableName() string {
	return "order_log"
}
