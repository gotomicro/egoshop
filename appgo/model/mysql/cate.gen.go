package mysql

import (
	"time"
)

type Cate struct {
	Id        int        `json:"id" form:"id" gorm:"primary_key"`           // 主键ID
	CreatedAt time.Time  `json:"created_at" form:"created_at" `             // 创建时间
	UpdatedAt time.Time  `json:"updated_at" form:"updated_at" `             // 更新时间
	DeletedAt *time.Time `json:"deleted_at" form:"deleted_at" gorm:"index"` // 删除时间
	Name      string     `json:"name" form:"name" `                         // 分类名称
	Icon      string     `json:"icon" form:"icon" `                         // 图标，可能是图片也可能是地址

}

func (*Cate) TableName() string {
	return "cate"
}
