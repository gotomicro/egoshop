package mysql

import (
	"time"
)

type Banner struct {
	Id        int        `gorm:"not null;primary_key;comment:'主键ID'"json:"id"`
	CreatedAt time.Time  `gorm:"comment:'注释'"json:"createdAt"`
	UpdatedAt time.Time  `gorm:"comment:'注释'"json:"updatedAt"`
	DeletedAt *time.Time `gorm:"comment:'注释'"json:"deletedAt"`
	CreatedBy int        `gorm:"not null;comment:'注释'"json:"createdBy"`
	UpdatedBy int        `gorm:"not null;comment:'注释'"json:"updatedBy"`
	ResType   int        `gorm:"not null;"json:"resType"`   // 资源类型
	Title     string     `gorm:"not null;"json:"title"`     // 名称
	Link      string     `gorm:"not null;"json:"link"`      // 链接
	Image     string     `gorm:"not null;"json:"image"`     // 图片地址
	Content   string     `gorm:"not null;"json:"content"`   // 内容
	Enable    int        `gorm:"not null;"json:"enable"`    // 是否显示
	StartTime int64      `gorm:"not null;"json:"startTime"` // 开始时间
	EndTime   int64      `gorm:"not null;"json:"endTime"`   // 结束时间

}

func (*Banner) TableName() string {
	return "banner"
}
