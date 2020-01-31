package mysql

import (
	"time"
)

type ComRelateCate struct {
	Id        int        `json:"id" form:"id" gorm:"primary_key"` // 主键ID
	CreatedAt time.Time  `json:"createdAt"`                       // 创建时间
	UpdatedAt time.Time  `json:"updatedAt"`                       // 更新时间
	DeletedAt *time.Time `json:"deletedAt"gorm:"index"`           // 删除时间
	ComId     int        `gorm:"not null;"json:"comId"`           //
	Cid       int        `gorm:"not null;"json:"cid"`
}

func (*ComRelateCate) TableName() string {
	return "com_relate_cate"
}
