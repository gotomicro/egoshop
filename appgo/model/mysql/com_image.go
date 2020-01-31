package mysql

import "time"

type ComImage struct {
	Id        int        `gorm:"not null;"json:"id"`        //
	ComId     int        `gorm:"not null;"json:"comId"`     // 商品公共内容id
	ColorId   int        `gorm:"not null;"json:"colorId"`   // 颜色规格值id
	Img       string     `gorm:"not null;"json:"img"`       // 商品图片
	Sort      int        `gorm:"not null;"json:"sort"`      // 排序
	IsDefault int        `gorm:"not null;"json:"isDefault"` // 默认主题，1是，0否
	DeletedAt *time.Time `gorm:"index"json:"deletedAt"`
	CreatedAt time.Time  `gorm:""json:"createdAt"`
}

func (*ComImage) TableName() string {
	return "com_image"
}
