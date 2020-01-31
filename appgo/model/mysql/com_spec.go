package mysql

import "time"

type ComSpec struct {
	Id        int        `gorm:"not null;primary_key"json:"id"` // 规格id
	Name      string     `gorm:"not null;"json:"name"`          // 规格名称
	Sort      int        `gorm:"not null;"json:"sort"`          // 规格排序
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index"json:"deletedAt"` // 软删除时间
}

func (*ComSpec) TableName() string {
	return "com_spec"
}
