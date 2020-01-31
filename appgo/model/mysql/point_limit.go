package mysql

import (
	"time"
)

// 积分限制
type PointLimit struct {
	Id              int       `gorm:"primary_key;not null;comment:'主键ID'"`
	CreatedAt       time.Time `gorm:"comment:'创建时间'"`
	UpdatedAt       time.Time `gorm:"comment:'更新时间'"`
	Uid             int       `gorm:"not null"`
	Limit1          int       `gorm:"not null;comment:'目前有多少积分'"`
	Limit1UpdatedAt time.Time ``
}

func (*PointLimit) TableName() string {
	return "point_limit"
}
