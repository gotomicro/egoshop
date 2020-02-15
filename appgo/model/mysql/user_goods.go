package mysql

import (
	"time"
)

type UserGoods struct {
	Id        int        `json:"id" form:"id" gorm:"primary_key"` // 主键ID
	CreatedAt time.Time  `json:"createdAt"`                       // 创建时间
	UpdatedAt time.Time  `json:"updatedAt"`                       // 更新时间
	DeletedAt *time.Time `json:"deletedAt"gorm:"index"`           // 删除时间
	Uid       int        `gorm:"not null;index"json:"uid"`        // 用户ID
	GoodsId   int        `gorm:"not null;index"json:"goodsId"`    // 商品ID
	TypeId    int        `gorm:"not null;index"json:"typeId"`     // 类型id
	Name      string     `gorm:"not null;"json:"name"`            // 商品名称
	IsPay     int        `gorm:"not null;"json:"isPay"`           // 是否购买
	IsCollect int        `gorm:"not null;"json:"isCollect"`       // 是否购买
	IsStar    int        `gorm:"not null;"json:"isStar"`          // 是否关注
	IsCreate  int        `gorm:"not null;"json:"isCreate"`        // 是否上传
	IsShare   int        `gorm:"not null;"json:"isShare"`         // 是否分享
	IsPrePay  int        `gorm:"not null;"json:"isPrePay"`        // 是否预购买
	IsRead    int        `gorm:"not null;"json:"isRead"`          // 是否阅读
	PartId    int        `gorm:"not null"json:"partid"`           //部门分类
}

func (*UserGoods) TableName() string {
	return "user_goods"
}
