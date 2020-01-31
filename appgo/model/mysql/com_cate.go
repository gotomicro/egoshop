package mysql

import (
	"time"
)

type ComCate struct {
	Id        int        `json:"id" form:"id" gorm:"primary_key"` // 主键ID
	CreatedAt time.Time  `json:"createdAt"`                       // 创建时间
	UpdatedAt time.Time  `json:"updatedAt"`                       // 更新时间
	DeletedAt *time.Time `json:"deletedAt"gorm:"index"`           // 删除时间
	CreatedBy int        `gorm:"not null;"json:"createdBy"`       // 创建者
	UpdatedBy int        `gorm:"not null;"json:"updatedBy"`       // 更新者
	Name      string     `gorm:"not null"json:"name"`
	Sequence  int        `gorm:"not null"json:"sequence"`
	Pid       int        `gorm:"not null"json:"pid"`
	Icon      string     `gorm:"not null"json:"icon" form:"icon" ` //
	Sort      int        `gorm:"not null"json:"sort" form:"sort" ` //
	Status    int        `gorm:"not null"json:"status"`
	Remark    string     `gorm:"not null"json:"remark"`
}

func (*ComCate) TableName() string {
	return "com_cate"
}
