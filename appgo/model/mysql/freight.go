package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Freight struct {
	Id        int              `json:"id" form:"id" gorm:"not null;primary_key"` // 主键ID
	CreatedAt time.Time        `gorm:""json:"createdAt"`                         // 创建时间
	UpdatedAt time.Time        `gorm:""json:"updatedAt"`                         // 更新时间
	DeletedAt *time.Time       `gorm:"index"json:"deletedAt"`                    // 删除时间
	Areas     FreightAreasJson `gorm:"not null;type:json"json:"areas"`
	Name      string           `gorm:"not null;"json:"name"`
	PayType   int              `gorm:"not null;"json:"payType"`
	PartId    int              `gorm:"not null"json:"partid"` //部门分类
}

func (*Freight) TableName() string {
	return "freight"
}

// 请在model/mysql/addition.json.go里定义FreightAreasJson结构体
func (c FreightAreasJson) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *FreightAreasJson) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}
