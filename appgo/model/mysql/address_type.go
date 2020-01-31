package mysql

import (
	"time"
)

type AddressType struct {
	Id        int        `gorm:"not null;primary_key;comment:'主键ID'"json:"id"`
	CreatedAt time.Time  `gorm:"comment:'注释'"json:"createdAt"`
	UpdatedAt time.Time  `gorm:"comment:'注释'"json:"updatedAt"`
	DeletedAt *time.Time `gorm:"comment:'注释'"json:"deletedAt"`
	CreatedBy int        `gorm:"not null;comment:'注释'"json:"createdBy"`
	UpdatedBy int        `gorm:"not null;comment:'注释'"json:"updatedBy"`
	Name      string     `gorm:"not null;comment:'注释'"json:"name"`
}

func (*AddressType) TableName() string {
	return "address_type"
}
