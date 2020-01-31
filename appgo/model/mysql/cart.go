package mysql

import (
	"errors"
	"math"
	"time"

	"github.com/thoas/go-funk"
)

type Cart struct {
	Id        int        `gorm:"not null;primary_key;comment:'注释'"json:"id"`
	CreatedAt time.Time  `gorm:"comment:'注释'"json:"createdAt"`
	UpdatedAt time.Time  `gorm:"comment:'注释'"json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index;comment:'注释'"json:"deletedAt"`
	CreatedBy int        `gorm:"not null;comment:'注释'"json:"createdBy"`
	UpdatedBy int        `gorm:"not null;comment:'注释'"json:"updatedBy"`
	IsCheck   int        `gorm:"not null;"json:"isCheck"` // 选中状态 默认1选中 0未选中
	Num       int        `gorm:"not null;"json:"num"`     // 商品个数

	TypeId   int `gorm:"not null;"json:"typeId"` // 类型，1为feed，2为com
	FeedId   int `gorm:"not null;"json:"feedId"`
	ComSkuId int `gorm:"not null;"json:"comSkuId"`

	// 赋值处理的
	Title     string  `gorm:"-"json:"title"`
	SubTitle  string  `gorm:"-"json:"subTitle"`
	Status    int     `gorm:"-"json:"status"`    // 1生效，2失效
	Cover     string  `gorm:"-"json:"cover"`     // goods sku img
	Price     float64 `gorm:"-"json:"price"`     // goods sku price
	PriceType int     `gorm:"-"json:"priceType"` // goods sku price
	PayType   int     `gorm:"-"json:"payType"`   // goods pay type

	Stock           int              `gorm:"-"json:"stock"` // goods sku stock
	ComId           int              `gorm:"-"json:"comId"`
	ComIsOnSale     int              `gorm:"-"json:"comIsOnSale"`
	ComFreightFee   float64          `gorm:"-"json:"comFreightFee"`
	ComFreightId    int              `gorm:"-"json:"comFreightId"`    // goods freight id
	ComSpec         ComSkuSpecJson   `gorm:"-"json:"comSpec"`         // goods sku spec
	ComWeight       float64          `gorm:"-"json:"comWeight"`       // goods sku weight
	ComFreightAreas FreightAreasJson `gorm:"-"json:"comFreightAreas"` // freight areas
}

func (*Cart) TableName() string {
	return "cart"
}

func (c Cart) GetFreightWay() string {
	if c.ComFreightId == 0 {
		return "goods_freight_unified"
	}
	return "goods_freight_template"
}

// 计算运输费用
func (c *Cart) FreightFeeByAddress(address Address) (float64, error) {
	// 运费计算方式
	// 如果模板id为0，说明是freight_unified，直接用指定运费
	if c.ComFreightId == 0 {
		return c.ComFreightFee, nil
	}

	var algorithm FreightAreas
	// 如果模板id大于0，选择对应的oods_freight_template

	var flag = false
	freightAreas := c.ComFreightAreas
	for _, value := range freightAreas {
		if funk.ContainsInt(value.AreaIds, address.AreaId) {
			algorithm = value
			flag = true
			break
		}
	}

	if flag {
		flag = false
		for _, value := range freightAreas {
			if funk.ContainsInt(value.AreaIds, address.CityId) {
				algorithm = value
				flag = true
				break
			}
		}
	}

	// todo 要看这里逻辑
	if flag {
		flag = false
		for _, value := range freightAreas {
			if funk.ContainsInt(value.AreaIds, address.ProvinceId) {
				algorithm = value
				flag = true
				break
			}
		}
	}

	firstAmout := algorithm.FirstAmount
	firstFee := algorithm.FirstFee
	additionalAmount := algorithm.AdditionalAmount
	additionalFee := algorithm.AdditionalFee

	if c.PayType == 2 {
		weight := c.ComWeight * float64(c.Num)
		if weight == 0 {
			return 0, errors.New("重量不可能为0")
		}
		return firstFee + math.Ceil((weight-firstAmout)/additionalAmount)*additionalFee, nil
	}
	if c.Num == 0 {
		return 0, errors.New("个数不可能为0")
	}
	return firstFee + math.Ceil((float64(c.Num)-firstAmout)/additionalAmount)*additionalFee, nil

}
