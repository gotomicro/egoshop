package mysql

import "time"

type ComSpecValue struct {
	Id        int        `gorm:"not null;primary_key"json:"id"` // 规格值id
	SpecId    int        `gorm:"not null;index"json:"specId"`   // 规格id
	Name      string     `gorm:"not null;"json:"name"`          // 规格值名称
	Sort      int        `gorm:"not null;"json:"sort"`          // 排序
	Color     string     `gorm:"not null;"json:"color"`         // 色彩值
	Img       string     `gorm:"not null;"json:"img"`           // 图片
	DeletedAt *time.Time `gorm:"index"json:"deletedAt"`         // 软删除时间

}

func (*ComSpecValue) TableName() string {
	return "com_spec_value"
}
