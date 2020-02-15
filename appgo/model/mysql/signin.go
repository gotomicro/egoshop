package mysql

import (
	"time"
)

type Signin struct {
	Id        int       `gorm:"primary_key;not null;comment:'主键ID'"`
	CreatedAt time.Time `gorm:"comment:'创建时间'"`
	UpdatedAt time.Time `gorm:"comment:'更新时间'"`
	Uid       int       `gorm:"not null"`
	Point     int       `gorm:"not null"`
	SigninCnt int       `gorm:"not null"`
	PartId    int       `gorm:"not null"json:"partid"` //部门分类
}

func (*Signin) TableName() string {
	return "signin"
}
