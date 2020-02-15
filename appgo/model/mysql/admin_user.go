package mysql

import (
	"time"
)

type AdminUser struct {
	Id            int        `json:"id" form:"id" gorm:"primary_key"`           // 主键ID
	CreatedAt     time.Time  `json:"created_at" form:"created_at" `             // 创建时间
	UpdatedAt     time.Time  `json:"updated_at" form:"updated_at" `             // 更新时间
	DeletedAt     *time.Time `json:"deleted_at" form:"deleted_at" gorm:"index"` // 删除时间
	Name          string     `gorm:"not null;type:varchar(100);comment:'姓名'"`
	Address       string     `gorm:"not null;type:varchar(100);comment:'地址'"`
	Password      string     `gorm:"not null;type:varchar(255);comment:'密码'"`
	Nickname      string     `gorm:"not null;"json:"nickname" form:"nickname" `                            //
	Email         string     `gorm:"not null;"json:"email" form:"email" `                                  //
	Phone         string     `gorm:"not null;"json:"phone" form:"phone" `                                  //
	Avatar        string     `gorm:"not null;"json:"avatar" form:"avatar" `                                //
	Role          int        `gorm:"not null;comment:'0 默认，1 普通会员，2 vip会员，3 高级专家'"json:"role" form:"role"` //
	Status        int        `gorm:"not null;"json:"status" form:"status" `                                //
	LastLoginTime time.Time  `json:"last_login_time" form:"last_login_time" `                              //
	LastLoginIP   string     `gorm:"not null;type:varchar(30);comment:'最后一次登陆IP'"`
	PartId        int        `gorm:"not null"json:"partid"` //部门分类
}

func (*AdminUser) TableName() string {
	return "admin_user"
}
