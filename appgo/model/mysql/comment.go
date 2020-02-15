package mysql

import (
	"time"
)

type Comment struct {
	Id        int        `json:"id" form:"id" gorm:"primary_key"`  // 主键ID
	CreatedAt time.Time  `json:"createdAt" `                       // 创建时间
	UpdatedAt time.Time  `json:"updatedAt"`                        // 更新时间
	DeletedAt *time.Time `json:"deletedAt"gorm:"index"`            // 删除时间
	GoodsId   int        `gorm:"not null"json:"goodsId"`           // 商品id
	TypeId    int        `gorm:"not null"json:"typeId"`            // 类型id
	Content   string     `gorm:"not null;type:text"json:"content"` // 评论内容
	Uid       int        `gorm:"not null"json:"uid"`               // 用户id
	Score     int        `gorm:"not null"json:"score"`             // 评分，存储为int，前端做小数转换
	Nickname  string     `gorm:"not null"json:"nickname"`          //
	Avatar    string     `gorm:"not null;"json:"avatar" `          //
	PartId    int        `gorm:"not null"json:"partid"`            //部门分类
}

func (*Comment) TableName() string {
	return "comment"
}
