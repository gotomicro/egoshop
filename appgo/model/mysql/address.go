package mysql

import (
	"time"

	"github.com/goecology/muses/pkg/database/mysql"
)

type Address struct {
	Id         int        `gorm:"not null;primary_key;comment:'主键ID'"json:"id"`
	CreatedAt  time.Time  `gorm:"comment:'注释'"json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"comment:'注释'"json:"updatedAt"`
	DeletedAt  *time.Time `gorm:"comment:'注释'"json:"deletedAt"`
	CreatedBy  int        `gorm:"not null;comment:'注释'"json:"createdBy"`
	UpdatedBy  int        `gorm:"not null;comment:'注释'"json:"updatedBy"`
	Name       string     `gorm:"not null;comment:'注释'"json:"name"`
	ProvinceId int        `gorm:"not null;comment:'注释'"json:"provinceId"`
	CityId     int        `gorm:"not null;comment:'注释'"json:"cityId"`
	AreaId     int        `gorm:"not null;comment:'注释'"json:"areaId"`
	Region     string     `gorm:"not null;"json:"region"`
	Detail     string     `gorm:"not null;comment:'注释'"json:"detail"`
	TelPhone   string     `gorm:"not null;comment:'注释'"json:"telPhone"`
	Mobile     string     `gorm:"not null;comment:'注释'"json:"mobile"`
	ZipCode    string     `gorm:"not null;comment:'注释'"json:"zipCode"`
	IsDefault  int        `gorm:"not null;comment:'注释'"json:"isDefault"`
	TypeId     int        `gorm:"not null;comment:'注释'"json:"typeId"` // 类型 1为家，2为公司,3为学校，4为其他
	StreetId   int        `gorm:"not null;comment:'注释'"json:"streetId"`
	TypeName   string     `gorm:"not null;"json:"typeName"`
}

func (*Address) TableName() string {
	return "address"
}

func (a *Address) WithTypeName() string {
	var info AddressType
	mysql.Caller("egoshop").Select("name").Where("id = ?", a.TypeId).Find(&info)
	a.TypeName = info.Name
	return info.Name
}
