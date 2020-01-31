package mysql

import (
	"time"
)

type SigninLog struct {
	Id        int       `json:"id" form:"id" gorm:"primary_key"` // 主键ID
	CreatedAt time.Time `json:"created_at" form:"created_at" `   // 创建时间
	Uid       int       `gorm:"not null"`
	Point     int       `gorm:"not null"`
	SigninCnt int       `gorm:"not null"`
}

func (*SigninLog) TableName() string {
	return "signin_log"
}
