package mysql

import (
	"time"
)

// 图片
type Image struct {
	Id        int        `json:"id" form:"id" gorm:"primary_key"` // 主键ID
	CreatedAt time.Time  `json:"createdAt" `                      // 创建时间
	UpdatedAt time.Time  `json:"updatedat" `                      // 更新时间
	DeletedAt *time.Time `json:"deletedAt"`                       // 删除时间
	CreatedBy int        `gorm:"not null;"json:"createdBy"`       // 创建者
	UpdatedBy int        `gorm:"not null;"json:"updatedBy"`       // 更新者
	TypeId    int        `gorm:"not null;"json:"typeId"`
	Name      string     `gorm:"not null;"json:"name" form:"name" ` //
	Size      float64    `gorm:"not null;"json:"size" form:"size" ` //
	Type      string     `gorm:"not null;"json:"type" form:"type" ` //
	Url       string     `gorm:"not null;"json:"url" form:"url" `   //
	PartId    int        `gorm:"not null"json:"partid"`             //部门分类

}

func (*Image) TableName() string {
	return "image"
}
