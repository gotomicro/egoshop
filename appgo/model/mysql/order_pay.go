package mysql

import (
	"time"
)

type OrderPay struct {
	Id        int        `gorm:"not null;primary_key"json:"id"` // 主键ID
	CreatedAt time.Time  `gorm:""json:"createdAt"`
	UpdatedAt time.Time  `gorm:""json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index"json:"deletedAt"`
	PaySn     string     `gorm:"not null;"json:"paySn"`
	Uid       int        `gorm:"not null;"json:"uid"`
	PayState  int        `gorm:"not null;"json:"payState"` // 支付状态[0默认未支付 1已支付 只有第三方支付接口通知到时才会更改此状态]
	PartId    int        `gorm:"not null"json:"partid"`    //部门分类
}

func (*OrderPay) TableName() string {
	return "order_pay"
}
