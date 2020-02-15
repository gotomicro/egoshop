package mysql

import (
	"time"
)

type Attachment struct {
	Id        int        `gorm:"not null;primary_key;comment:'注释'"`
	CreatedAt time.Time  `gorm:"comment:'注释'"json:"createdAt"`
	UpdatedAt time.Time  `gorm:"comment:'注释'"json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index;comment:'注释'"json:"deletedAt"`
	GoodsId   int        `gorm:"not null;comment:'注释'"json:"goodsId"`
	GoodsType int        `gorm:"not null;comment:'注释'"json:"goodsType"`
	FileName  string     `gorm:"not null;comment:'注释'"json:"fileName"`
	FilePath  string     `gorm:"not null;comment:'注释'"json:"filePath"`
	FileSize  float64    `gorm:"not null;comment:'注释'"json:"fileSize"`
	HttpPath  string     `gorm:"not null;comment:'注释'"json:"httpPath"`
	FileExt   string     `gorm:"not null;comment:'注释'"json:"fileExt"`
	CreatedBy int        `gorm:"not null;comment:'注释'"json:"createdBy"`
	PartId    int        `gorm:"not null"json:"partid"` //部门分类
}

func (*Attachment) TableName() string {
	return "attachment"
}
