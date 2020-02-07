package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type ComSku struct {
	Id            int            `gorm:"not null;comment:'主键id'"json:"id"`
	CreatedAt     time.Time      `gorm:"comment:'创建时间'"json:"createdAt"`
	UpdatedAt     time.Time      `gorm:"comment:'更新时间'"json:"updatedAt"` // 更新时间
	DeletedAt     *time.Time     `gorm:"comment:'注释';index"json:"deletedAt"`
	CreatedBy     int            `gorm:"not null;comment:'注释'"json:"createdBy"`
	UpdatedBy     int            `gorm:"not null;comment:'注释'"json:"updatedBy"`
	ComId         int            `gorm:"not null;"json:"comId"`         // 商品公共表id
	Spec          ComSkuSpecJson `gorm:"not null;type:json"json:"spec"` // 规格json信息
	Price         float64        `gorm:"not null;"json:"price"`         // 商品价格
	Stock         int            `gorm:"not null;"json:"stock"`         // 商品库存
	Code          string         `gorm:"not null;"json:"code"`          // 商品编码
	Cover         string         `gorm:"not null;"json:"cover"`         // 商品主图
	Weight        float64        `gorm:"not null;"json:"weight"`        // 商品重量
	Title         string         `gorm:"not null;"json:"title"`         // 商品名称（+规格名称）
	SaleNum       int            `gorm:"not null;"json:"saleNum"`       // 销售数量
	GroupSaleNum  int            `gorm:"not null;"json:"groupSaleNum"`  // 拼团销量
	//SpecValueSign string         `gorm:"not null;"json:"specValueSign"` // 规格值标识
	//SpecSign      string         `gorm:"not null;"json:"specSign"`      // 规格标识
}

func (*ComSku) TableName() string {
	return "com_sku"
}

type ComSkuSpecJson []ComSkuSpecOneInfo

type ComSkuSpecOneInfo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ValueID   int    `json:"valueId"`
	ValueName string `json:"valueName"`
	ValueImg  string `json:"valueImg"`
}

func (c ComSkuSpecJson) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *ComSkuSpecJson) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}
