package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Com struct {
	Id           int             `gorm:"not null;comment:'主键id'"json:"id"`
	CreatedAt    time.Time       `gorm:"comment:'创建时间'"json:"createdAt"`
	UpdatedAt    time.Time       `gorm:"comment:'更新时间'"json:"updatedAt"` // 更新时间
	DeletedAt    *time.Time      `gorm:"comment:'注释';index"json:"deletedAt"`
	CntView      int             `gorm:"not null;comment:'注释'"json:"cntView"`
	CntStar      int             `gorm:"not null;comment:'注释'"json:"cntStar"`
	CntCollect   int             `gorm:"not null;comment:'注释'"json:"cntCollect"`
	CntShare     int             `gorm:"not null;comment:'注释'"json:"cntShare"`
	CntComment   int             `gorm:"not null;comment:'注释'"json:"cntComment"`
	CntIsPay     int             `gorm:"not null;"json:"cntIsPay"` // 为了兼容usergoods
	CreatedBy    int             `gorm:"not null;comment:'注释'"json:"createdBy"`
	UpdatedBy    int             `gorm:"not null;comment:'注释'"json:"updatedBy"`
	Title        string          `gorm:"not null;"json:"title"`
	SubTitle     string          `gorm:"not null;"json:"subTitle"`
	Cover        string          `gorm:"not null;comment:'注释'"json:"cover"`
	Gallery      ComGalleryJson  `gorm:"not null;type:json;comment:'注释'"json:"gallery"`
	Stock        int             `gorm:"not null;comment:'注释'"json:"stock"` // goods表库存之和
	SaleNum      int             `gorm:"not null;comment:'注释'"json:"saleNum"`
	GroupSaleNum int             `gorm:"not null;comment:'注释'"json:"groupSaleNum"`
	SaleTime     time.Time       `gorm:"comment:'注释'"json:"saleTime"`
	PayType      int             `gorm:"not null;comment:'注释'"json:"payType"`
	FreightFee   float64         `gorm:"not null;comment:'注释'"json:"freightFee"`
	FreightId    int             `gorm:"not null;comment:'注释'"json:"freightId"`
	BaseSaleNum  int             `gorm:"not null;comment:'注释'"json:"baseSaleNum"`
	IsOnSale     int             `gorm:"not null;comment:'注释'"json:"isOnSale"`
	ImageSpecId  int             `gorm:"not null;comment:'注释'"json:"imageSpecId"`
	OriginPrice  float64         `gorm:"not null;"json:"originPrice"` // 原价
	Price        float64         `gorm:"not null;comment:'注释'"json:"price"`
	PreMarkdown  string          `gorm:"-"json:"preMarkdown"`
	PreHtml      string          `gorm:"-"json:"preHtml"`
	Markdown     string          `gorm:"-"json:"markdown"`
	Html         string          `gorm:"-"json:"html"`
	WechatHtml   string          `gorm:"-"json:"wechatHtml"`
	SkuList      []ComSku        `gorm:"-"json:"skuList"`                   // 库存列表信息,冗余字段
	SpecList     ComSpecListJson `gorm:"not null;type:json"json:"specList"` // 规格参数
	Cids         ComCidsJson           `gorm:"not null;type:json"json:"cids"`
}

func (*Com) TableName() string {
	return "com"
}

// 请在model/mysql/addition.json.go里定义ComGalleryJson结构体
func (c ComGalleryJson) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *ComGalleryJson) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

// 请在model/mysql/addition.json.go里定义ComCategoryIdsJson结构体
func (c ComCidsJson) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *ComCidsJson) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

// 请在model/mysql/addition.json.go里定义ComBodyJson结构体
func (c ComBodyJson) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *ComBodyJson) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

// 请在model/mysql/addition.json.go里定义ComImageSpecImagesJson结构体
func (c ComImageSpecImagesJson) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *ComImageSpecImagesJson) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

// 请在model/mysql/addition.json.go里定义ComSpecListJson结构体
func (c ComSpecListJson) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *ComSpecListJson) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}
