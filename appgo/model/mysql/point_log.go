package mysql

import (
	"time"
)

type PointLog struct {
	Id        int       `json:"id" form:"id" gorm:"primary_key"` // 主键ID
	CreatedAt time.Time `json:"created_at" form:"created_at" `   // 创建时间
	Uid       int       `gorm:"not null"`
	TypeId    int       `gorm:"not null"`
	Point     int       `gorm:"not null"`
	PartId    int       `gorm:"not null"json:"partid"` //部门分类
}

func (*PointLog) TableName() string {
	return "point_log"
}
