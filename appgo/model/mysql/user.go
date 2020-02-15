package mysql

import (
	"time"
)

type User struct {
	Id            int        `json:"id" form:"id" gorm:"primary_key"` // 主键ID
	CreatedAt     time.Time  `gorm:""json:"createdAt"`                // 创建时间
	UpdatedAt     time.Time  `gorm:""json:"updatedAt"`                // 更新时间
	DeletedAt     *time.Time `gorm:"index"json:"deletedAt"`           // 删除时间
	Name          string     `gorm:"not null;"json:"name"`            // 姓名
	OpenId        int        `gorm:"not null;"json:"openId"`          // 开放平台id
	OpenType      int        `gorm:"not null;"json:"openType"`        // 开放平台类型，微信、支付宝等
	Address       string     `gorm:"not null;"json:"address"`         // 地址
	Password      string     `gorm:"not null;"json:"password"`        // 密码
	Point         int        `gorm:"not null;"json:"point"`           // 积分
	Quota         int        `gorm:"not null;"json:"quota"`           // 免费下载资料的配额
	Balance       int        `gorm:"not null;"json:"balance"`         // 账户余额
	Level         int        `gorm:"not null;"json:"level"`           // 等级
	Account       string     `gorm:"not null;"json:"account"`
	AuthMethod    string     `gorm:"not null;"json:"authMethod"`
	Nickname      string     `gorm:"not null;"json:"nickname"`
	Desc          string     `gorm:"not null;"json:"desc"`
	Email         string     `gorm:"not null;"json:"email"`
	Phone         string     `gorm:"not null;"json:"phone"`
	Avatar        string     `gorm:"not null;"json:"avatar"`
	Role          int        `gorm:"not null;"json:"role"`
	Status        int        `gorm:"not null;"json:"status"`
	LastLoginTime time.Time  `gorm:""json:"lastLoginTime"`
	LastLoginIP   string     `gorm:"not null;"json:"lastLoginIp"` // 最后一次登陆IP
	Wxpay         string     `gorm:"not null;"json:"wxpay"`       //
	Alipay        string     `gorm:"not null;"json:"alipay"`      //
	RoleName      string     `gorm:"-"`
	PartId        int        `gorm:"not null"json:"partid"` //部门分类
}

func (*User) TableName() string {
	return "user"
}

func (m *User) IsAdministrator() bool {
	if m == nil || m.Id <= 0 {
		return false
	}
	return m.Role == 0 || m.Role == 1
}
